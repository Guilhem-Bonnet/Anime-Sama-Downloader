from __future__ import annotations

import threading
import time
import uuid
from collections import deque
from concurrent.futures import Future, ThreadPoolExecutor
from dataclasses import dataclass, field
from typing import Callable, Deque, Dict, Optional


JobRunner = Callable[[threading.Event], Optional[str]]
JobEventCallback = Callable[["DownloadJob", str], None]


@dataclass
class DownloadJob:
    label: str
    run: JobRunner
    job_id: str = field(default_factory=lambda: uuid.uuid4().hex)
    cancel_event: threading.Event = field(default_factory=threading.Event)
    status: str = "PENDING"  # PENDING|RUNNING|SUCCESS|FAILED|CANCELLED
    result_path: Optional[str] = None
    error: Optional[str] = None
    created_at: float = field(default_factory=time.time)
    started_at: Optional[float] = None
    finished_at: Optional[float] = None

    # Optional structured progress (for UIs).
    progress_percent: Optional[float] = None
    progress_downloaded: Optional[int] = None
    progress_total: Optional[int] = None
    progress_speed_bps: Optional[float] = None
    progress_eta_seconds: Optional[float] = None
    progress_stage: Optional[str] = None
    progress_message: Optional[str] = None


class DownloadManager:
    """Simple job queue with bounded parallelism.

    Designed to be used from both CLI (blocking wait) and Textual TUI (callbacks).
    """

    def __init__(
        self,
        max_parallel: int = 10,
        on_event: JobEventCallback | None = None,
        executor_name: str = "download",
    ):
        self.max_parallel = max(1, min(int(max_parallel or 1), 10))
        self._on_event = on_event
        self._lock = threading.RLock()
        self._pending: Deque[DownloadJob] = deque()
        self._running: Dict[str, Future] = {}
        self._jobs: Dict[str, DownloadJob] = {}
        self._executor = ThreadPoolExecutor(max_workers=self.max_parallel, thread_name_prefix=executor_name)

    def _emit(self, job: DownloadJob, event: str) -> None:
        cb = self._on_event
        if cb is None:
            return
        try:
            cb(job, event)
        except Exception:
            # Never crash the downloader on UI/log callback failures.
            return

    def set_max_parallel(self, max_parallel: int) -> None:
        # ThreadPoolExecutor can't be resized; keep the cap here for future versions.
        with self._lock:
            self.max_parallel = max(1, min(int(max_parallel or 1), 10))

    def enqueue(self, job: DownloadJob) -> str:
        with self._lock:
            self._jobs[job.job_id] = job
            self._pending.append(job)
            self._emit(job, "queued")
            self._drain_locked()
            return job.job_id

    def list_jobs(self) -> list[DownloadJob]:
        with self._lock:
            return list(self._jobs.values())

    def update_job_progress(
        self,
        job_id: str,
        *,
        percent: float | None = None,
        downloaded: int | None = None,
        total: int | None = None,
        speed_bps: float | None = None,
        eta_seconds: float | None = None,
        stage: str | None = None,
        message: str | None = None,
    ) -> bool:
        """Update structured progress for a job.

        Returns True if job exists.
        """
        with self._lock:
            job = self._jobs.get(job_id)
            if job is None:
                return False
            if percent is not None:
                job.progress_percent = float(percent)
            if downloaded is not None:
                job.progress_downloaded = int(downloaded)
            if total is not None:
                job.progress_total = int(total)
            if speed_bps is not None:
                job.progress_speed_bps = float(speed_bps)
            if eta_seconds is not None:
                job.progress_eta_seconds = float(eta_seconds)
            if stage is not None:
                job.progress_stage = str(stage)
            if message is not None:
                job.progress_message = str(message)
            return True

    def clear_finished(self) -> int:
        """Remove finished jobs (SUCCESS/FAILED/CANCELLED) from history."""
        with self._lock:
            finished_ids = [
                jid
                for jid, job in self._jobs.items()
                if job.status in {"SUCCESS", "FAILED", "CANCELLED"} and jid not in self._running
            ]
            for jid in finished_ids:
                self._jobs.pop(jid, None)
            return len(finished_ids)

    def retry(self, job_id: str) -> str | None:
        """Retry a job by cloning it as a new job (new job_id)."""
        with self._lock:
            job = self._jobs.get(job_id)
            if job is None:
                return None
            new_job = DownloadJob(label=job.label, run=job.run)
        # enqueue outside the lock (enqueue acquires lock internally)
        self.enqueue(new_job)
        return new_job.job_id

    def get_job(self, job_id: str) -> Optional[DownloadJob]:
        with self._lock:
            return self._jobs.get(job_id)

    def pending_count(self) -> int:
        with self._lock:
            return len(self._pending)

    def running_count(self) -> int:
        with self._lock:
            return len(self._running)

    def cancel(self, job_id: str) -> bool:
        with self._lock:
            job = self._jobs.get(job_id)
            if not job:
                return False

            job.cancel_event.set()

            # If it's pending, remove it immediately.
            for i, pending_job in enumerate(list(self._pending)):
                if pending_job.job_id == job_id:
                    try:
                        self._pending.remove(pending_job)
                    except ValueError:
                        pass
                    pending_job.status = "CANCELLED"
                    pending_job.finished_at = time.time()
                    self._emit(pending_job, "cancelled")
                    return True

            # If running, cooperative cancel.
            fut = self._running.get(job_id)
            if fut is not None:
                fut.cancel()
            self._emit(job, "cancelling")
            return True

    def clear_pending(self) -> int:
        with self._lock:
            cleared = 0
            while self._pending:
                job = self._pending.popleft()
                job.cancel_event.set()
                job.status = "CANCELLED"
                job.finished_at = time.time()
                self._emit(job, "cancelled")
                cleared += 1
            return cleared

    def cancel_all(self) -> None:
        with self._lock:
            for job in list(self._jobs.values()):
                job.cancel_event.set()
            self.clear_pending()
            for job_id, fut in list(self._running.items()):
                fut.cancel()
                job = self._jobs.get(job_id)
                if job is not None:
                    self._emit(job, "cancelling")

    def _drain_locked(self) -> None:
        # Start as many pending jobs as we can.
        while self._pending and len(self._running) < self.max_parallel:
            job = self._pending.popleft()
            if job.cancel_event.is_set():
                job.status = "CANCELLED"
                job.finished_at = time.time()
                self._emit(job, "cancelled")
                continue

            job.status = "RUNNING"
            job.started_at = time.time()
            self._emit(job, "started")

            fut = self._executor.submit(self._run_job, job)
            self._running[job.job_id] = fut
            fut.add_done_callback(lambda f, jid=job.job_id: self._on_done(jid, f))

    def _run_job(self, job: DownloadJob) -> Optional[str]:
        if job.cancel_event.is_set():
            return None
        return job.run(job.cancel_event)

    def _on_done(self, job_id: str, fut: Future) -> None:
        with self._lock:
            self._running.pop(job_id, None)
            job = self._jobs.get(job_id)
            if job is None:
                self._drain_locked()
                return

            job.finished_at = time.time()

            if job.cancel_event.is_set():
                job.status = "CANCELLED"
                job.error = None
                job.result_path = None
                self._emit(job, "cancelled")
                self._drain_locked()
                return

            try:
                res = fut.result()
                if isinstance(res, str) and res:
                    job.status = "SUCCESS"
                    job.result_path = res
                    job.error = None
                    self._emit(job, "success")
                else:
                    job.status = "FAILED"
                    job.result_path = None
                    if not job.error:
                        job.error = "failed"
                    self._emit(job, "failed")
            except Exception as e:
                job.status = "FAILED"
                job.result_path = None
                job.error = str(e)
                self._emit(job, "failed")

            self._drain_locked()

    def wait(self) -> None:
        """Block until no pending/running jobs remain."""
        while True:
            with self._lock:
                if not self._pending and not self._running:
                    return
            time.sleep(0.1)

    def shutdown(self) -> None:
        try:
            self.cancel_all()
        except Exception:
            pass
        try:
            self._executor.shutdown(wait=False, cancel_futures=True)
        except Exception:
            pass
