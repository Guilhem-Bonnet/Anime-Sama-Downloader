
import os
import tempfile
import threading
import time
from concurrent.futures import ThreadPoolExecutor, as_completed
from typing import Callable

import requests
try:
    from tqdm import tqdm
except Exception:  # pragma: no cover
    class _DummyTqdm:
        def __init__(self, iterable=None, total=None, **kwargs):
            self.iterable = iterable
            self.total = total

        def __iter__(self):
            return iter(self.iterable or [])

        def __enter__(self):
            return self

        def __exit__(self, exc_type, exc, tb):
            return False

        def update(self, n=1):
            return None

    def tqdm(iterable=None, **kwargs):  # type: ignore
        return _DummyTqdm(iterable=iterable, **kwargs)

from utils.parsers import parse_ts_segments
from utils.var import print_status
from utils.http_pool import get_http_pool


http_pool = get_http_pool()

def download_video(
    video_url,
    save_path,
    use_ts_threading=False,
    url='',
    automatic_mp4=False,
    threaded_mp4=False,
    log_callback: Callable[[str, str], None] | None = None,
    progress_callback: Callable[[dict], None] | None = None,
    use_tqdm: bool = True,
    use_mp4_threading: bool = False,
    mp4_workers: int = 4,
    ts_workers: int = 10,
    cancel_event: threading.Event | None = None,
):
    def _status(message: str, status_type: str = "info") -> None:
        if log_callback is not None:
            try:
                log_callback(message, status_type)
                return
            except Exception:
                # Fallback to stdout if callback fails
                pass
        print_status(message, status_type)

    def _cancelled() -> bool:
        return cancel_event is not None and cancel_event.is_set()

    progress_last_t = 0.0
    progress_last_bytes = 0
    progress_start_t = time.time()

    def _progress(
        *,
        stage: str,
        downloaded: int | None = None,
        total: int | None = None,
        percent: float | None = None,
        message: str | None = None,
        force: bool = False,
    ) -> None:
        nonlocal progress_last_t, progress_last_bytes
        if progress_callback is None:
            return
        now = time.time()
        # Throttle to avoid flooding SSE/UI.
        if not force and (now - progress_last_t) < 0.35:
            return

        speed_bps = None
        eta_seconds = None
        if downloaded is not None:
            dt = max(1e-6, now - progress_last_t) if progress_last_t else max(1e-6, now - progress_start_t)
            db = downloaded - progress_last_bytes if progress_last_t else downloaded
            if db >= 0 and dt > 0:
                speed_bps = float(db) / float(dt)
            if total is not None and speed_bps and speed_bps > 1e-6:
                remaining = max(0, int(total) - int(downloaded))
                eta_seconds = float(remaining) / float(speed_bps)

        payload = {
            "stage": stage,
            "downloaded": downloaded,
            "total": total,
            "percent": percent,
            "speed_bps": speed_bps,
            "eta_seconds": eta_seconds,
            "message": message,
            "ts": now,
        }
        try:
            progress_callback(payload)
        except Exception:
            return
        progress_last_t = now
        if downloaded is not None:
            progress_last_bytes = int(downloaded)

    def _cleanup_partial(path: str) -> None:
        try:
            if path and os.path.exists(path):
                os.remove(path)
        except Exception:
            pass

    if _cancelled():
        _status("Download cancelled before start.", "warning")
        return False, None

    _status(f"Starting download: {os.path.basename(save_path)}", "loading")
    headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Gecko/20100101 Firefox/108.0',
        'Accept': 'video/webm,video/mp4,video/*;q=0.9,*/*;q=0.8',
        'Accept-Language': 'en-US,en;q=0.5',
        'Referer': 'https://movearnpre.com/' if 'movearnpre.com' in url or 'ovaltinecdn.com' in url else
                  'https://vidmoly.net/' if 'vidmoly.net'in url else 'https://vidmoly.net/' if 'vidmoly.to' in url else
                  'https://oneupload.net/' if 'oneupload.net' in url else 
                  'https://sendvid.com/' if 'sendvid.com' in url else 
                  'https://mivalyo.com/' if 'mivalyo.com' in url else
                  'https://video.sibnet.ru/'
    }

    def _range_probe(target_url: str) -> tuple[str, int, bool]:
        """Return (final_url, content_length, supports_ranges)."""
        try:
            head = http_pool.session.head(target_url, headers=headers, allow_redirects=True, timeout=15)
            try:
                final_url = head.url or target_url
                content_length = int(head.headers.get("Content-Length", 0) or 0)
                accept_ranges = str(head.headers.get("Accept-Ranges", "")).lower()
            finally:
                try:
                    head.close()
                except Exception:
                    pass
            if accept_ranges == "bytes" and content_length > 0:
                return final_url, content_length, True
            # Some hosts don't expose Accept-Ranges on HEAD; probe with a tiny Range GET
            probe_headers = dict(headers)
            probe_headers["Range"] = "bytes=0-0"
            probe = http_pool.get(final_url, headers=probe_headers, stream=True, timeout=15)
            try:
                if probe.status_code == 206:
                    try:
                        cr = probe.headers.get("Content-Range", "")
                        # format: bytes 0-0/12345
                        if "/" in cr:
                            content_length = int(cr.split("/")[-1])
                    except Exception:
                        pass
                    return final_url, content_length, True
            finally:
                try:
                    probe.close()
                except Exception:
                    pass
        except Exception:
            pass
        return target_url, 0, False

    def _download_mp4_ranged(target_url: str, total_size: int) -> tuple[bool, str | None]:
        """Download an MP4 using concurrent Range requests when supported."""
        workers = max(2, int(mp4_workers or 4))
        workers = min(workers, 16)
        part_size = max(5 * 1024 * 1024, total_size // workers)  # at least 5MB
        ranges: list[tuple[int, int]] = []
        start = 0
        while start < total_size:
            end = min(total_size - 1, start + part_size - 1)
            ranges.append((start, end))
            start = end + 1

        save_dirname = os.path.dirname(save_path)
        if save_dirname:
            os.makedirs(save_dirname, exist_ok=True)

        # Write directly into a preallocated temp file to avoid a costly merge step.
        tmp_path = save_path + ".part"
        try:
            with open(tmp_path, "wb") as f:
                f.truncate(total_size)
        except Exception:
            # If truncate fails, fall back to old behavior by signalling failure.
            _cleanup_partial(tmp_path)
            return False, None

        downloaded = 0
        downloaded_lock = threading.Lock()
        last_report_t = 0.0
        last_report_pct = -1.0
        last_report_downloaded = -1

        def _report(delta: int) -> None:
            nonlocal downloaded, last_report_t, last_report_pct, last_report_downloaded
            with downloaded_lock:
                downloaded += delta
                if total_size <= 0:
                    return
                pct = (downloaded * 100.0) / float(total_size)
                now = time.time()

                # Never show 100% until we have safely renamed the file.
                if pct >= 100.0:
                    pct = 99.9

                pct_step = 0.2 if pct >= 90 else 5.0
                time_step = 2.0 if pct >= 90 else 5.0
                should = (
                    (pct - last_report_pct) >= pct_step
                    or (now - last_report_t) >= time_step
                )

                if should and downloaded != last_report_downloaded:
                    pct_txt = f"{pct:.1f}%" if pct >= 90 else f"{int(pct)}%"
                    _status(f"Progress: {pct_txt} ({downloaded}/{total_size} bytes)", "loading")
                    _progress(
                        stage="mp4_range",
                        downloaded=downloaded,
                        total=total_size,
                        percent=float(pct),
                        message=None,
                    )
                    last_report_pct = pct
                    last_report_downloaded = downloaded
                    last_report_t = now

        def _fetch_part(i: int, byte_range: tuple[int, int]) -> bool:
            if _cancelled():
                return False

            start_b, end_b = byte_range
            expected = (end_b - start_b) + 1
            h = dict(headers)
            h["Range"] = f"bytes={start_b}-{end_b}"

            for attempt in range(3):
                if _cancelled():
                    return False
                try:
                    r = http_pool.get(target_url, headers=h, stream=True, timeout=60)
                    try:
                        if r.status_code != 206:
                            continue

                        wrote = 0
                        with open(tmp_path, "r+b") as out:
                            out.seek(start_b)
                            for chunk in r.iter_content(chunk_size=1024 * 1024):
                                if _cancelled():
                                    return False
                                if not chunk:
                                    continue
                                out.write(chunk)
                                wrote += len(chunk)
                                if not use_tqdm:
                                    _report(len(chunk))

                        if wrote != expected:
                            continue
                        return True
                    finally:
                        try:
                            r.close()
                        except Exception:
                            pass
                except Exception:
                    if attempt < 2:
                        time.sleep(1.5)
                    continue

            return False

        if use_tqdm:
            _status(f"MP4 multi-part download: {len(ranges)} parts x ~{part_size // (1024*1024)}MB", "info")

        with ThreadPoolExecutor(max_workers=workers) as executor:
            futures = [executor.submit(_fetch_part, i, r) for i, r in enumerate(ranges)]
            for fut in as_completed(futures):
                if _cancelled():
                    _cleanup_partial(tmp_path)
                    return False, None
                ok = fut.result()
                if not ok:
                    _cleanup_partial(tmp_path)
                    return False, None

        if _cancelled():
            _cleanup_partial(tmp_path)
            return False, None

        _status("Finalisation: écriture terminée, renommage…", "loading")
        _progress(stage="finalize", downloaded=total_size, total=total_size, percent=99.9, message="rename", force=True)
        os.replace(tmp_path, save_path)
        _status(f"Progress: 100% ({total_size}/{total_size} bytes)", "loading")
        _progress(stage="done", downloaded=total_size, total=total_size, percent=100.0, message="ok", force=True)
        _status("Finalisation: OK", "info")
        return True, save_path

    try:
        if _cancelled():
            _status("Download cancelled.", "warning")
            return False, None

        if 'm3u8' in video_url:
            _progress(stage="m3u8", downloaded=0, total=None, percent=0.0, message="playlist")
            if _cancelled():
                _status("Download cancelled.", "warning")
                return False, None
            response = http_pool.get(video_url, headers=headers, timeout=10)
            try:
                response.raise_for_status()
                segments = parse_ts_segments(response.text)
            finally:
                try:
                    response.close()
                except Exception:
                    pass
            if not segments:
                _status("No .ts segments found in M3U8 playlist", "error")
                return False, None
            
            save_dirname = os.path.dirname(save_path)
            if save_dirname:
                os.makedirs(save_dirname, exist_ok=True)
            temp_ts_path = save_path.replace('.mp4', '.ts')
            random_string = os.path.basename(save_path).replace('.mp4', '.ts')

            # Never prompt here; the caller (main.py) decides.
            use_threads = use_ts_threading
            
            if use_threads:
                segment_data = []
                
                def download_segment(segment_url, index):
                    for attempt in range(3):
                        if _cancelled():
                            return index, None
                        try:
                            with http_pool.get(segment_url, headers=headers, stream=True, timeout=10) as seg_response:
                                seg_response.raise_for_status()
                                return index, seg_response.content
                        except requests.RequestException as e:
                            if attempt < 2:
                                time.sleep(2)
                            else:
                                print_status(f"Failed to download segment {index+1}: {str(e)}", "error")
                                return index, None
                    return index, None

                workers = max(1, min(int(ts_workers or 10), 32))
                with ThreadPoolExecutor(max_workers=workers) as executor:
                    future_to_segment = {executor.submit(download_segment, url, i): i for i, url in enumerate(segments)}
                    total = len(segments)
                    done = 0
                    if use_tqdm:
                        with tqdm(total=total, desc=f"📥 {random_string}", unit="segment") as pbar:
                            for future in as_completed(future_to_segment):
                                if _cancelled():
                                    _status("Download cancelled.", "warning")
                                    return False, None
                                index, content = future.result()
                                if content is None:
                                    _status(f"Aborting download due to failure in segment {index+1}", "error")
                                    return False, None
                                segment_data.append((index, content))
                                done += 1
                                pbar.update(1)
                                pct = (done * 100.0) / float(total) if total > 0 else None
                                _progress(stage="m3u8", downloaded=done, total=total, percent=pct, message="segments")
                    else:
                        for future in as_completed(future_to_segment):
                            if _cancelled():
                                _status("Download cancelled.", "warning")
                                return False, None
                            index, content = future.result()
                            if content is None:
                                _status(f"Aborting download due to failure in segment {index+1}", "error")
                                return False, None
                            segment_data.append((index, content))
                            done += 1
                            if done == total or done % 10 == 0:
                                _status(f"Segments: {done}/{total}", "loading")
                                pct = (done * 100.0) / float(total) if total > 0 else None
                                _progress(stage="m3u8", downloaded=done, total=total, percent=pct, message="segments")

                segment_data.sort(key=lambda x: x[0])
                
                with open(temp_ts_path, 'wb') as f:
                    for _, content in segment_data:
                        f.write(content)
            else:
                with open(temp_ts_path, 'wb') as f:
                    seq = enumerate(segments)
                    if use_tqdm:
                        seq = enumerate(tqdm(segments, desc=f"📥 {random_string}", unit="segment"))
                    total = len(segments)
                    for i, segment_url in seq:
                        if _cancelled():
                            _status("Download cancelled.", "warning")
                            try:
                                f.flush()
                            except Exception:
                                pass
                            _cleanup_partial(temp_ts_path)
                            return False, None
                        for attempt in range(3):
                            if _cancelled():
                                _status("Download cancelled.", "warning")
                                _cleanup_partial(temp_ts_path)
                                return False, None
                            try:
                                with http_pool.get(segment_url, headers=headers, stream=True, timeout=10) as seg_response:
                                    seg_response.raise_for_status()
                                    for chunk in seg_response.iter_content(chunk_size=1024 * 512):
                                        if _cancelled():
                                            _status("Download cancelled.", "warning")
                                            _cleanup_partial(temp_ts_path)
                                            return False, None
                                        if chunk:
                                            f.write(chunk)
                                break
                            except requests.RequestException as e:
                                if attempt < 2:
                                    time.sleep(2)
                                else:
                                    _status(f"Failed to download segment {i+1}: {str(e)}", "error")
                                    return False, None
                        if not use_tqdm and ((i + 1) == total or (i + 1) % 10 == 0):
                            _status(f"Segments: {i+1}/{total}", "loading")
                            pct = ((i + 1) * 100.0) / float(total) if total > 0 else None
                            _progress(stage="m3u8", downloaded=(i + 1), total=total, percent=pct, message="segments")
            
            _status(f"Combined {len(segments)} segments into {temp_ts_path}", "success")
            _progress(stage="done", downloaded=len(segments), total=len(segments), percent=100.0, message="ts_ready", force=True)
            return True, temp_ts_path
        else:
            if _cancelled():
                _status("Download cancelled.", "warning")
                return False, None
            final_url, total_size, supports_ranges = _range_probe(video_url)
            _progress(stage="mp4", downloaded=0, total=total_size or None, percent=0.0, message="start")
            if use_mp4_threading and supports_ranges and total_size >= 25 * 1024 * 1024:
                _status("Using multi-part MP4 download (Range)", "info")
                ok, out = _download_mp4_ranged(final_url, total_size)
                if ok:
                    _status("Download completed successfully!", "success")
                    return True, out
                _status("Multi-part download failed; falling back to single stream", "warning")

            with http_pool.get(final_url, stream=True, headers=headers, timeout=60) as response:
                # content-length may be missing if the server doesn't expose it
                if total_size <= 0:
                    try:
                        total_size = int(response.headers.get('content-length', 0) or 0)
                    except Exception:
                        total_size = 0

                if response.status_code != 200:
                    _status(f"Download failed with status code: {response.status_code}", "error")
                    return False, None

                save_dirname = os.path.dirname(save_path)
                if save_dirname:
                    os.makedirs(save_dirname, exist_ok=True)

                with open(save_path, 'wb') as f:
                    if use_tqdm:
                        with tqdm(
                            total=total_size,
                            unit='B',
                            unit_scale=True,
                            desc=f"📥 {os.path.basename(save_path)}",
                            bar_format='{l_bar}{bar}| {n_fmt}/{total_fmt} [{elapsed}<{remaining}, {rate_fmt}]'
                        ) as pbar:
                            downloaded = 0
                            last_emit_t = 0.0
                            for chunk in response.iter_content(chunk_size=1024 * 1024):
                                if _cancelled():
                                    _status("Download cancelled.", "warning")
                                    try:
                                        f.flush()
                                    except Exception:
                                        pass
                                    _cleanup_partial(save_path)
                                    return False, None
                                if chunk:
                                    f.write(chunk)
                                    downloaded += len(chunk)
                                    pbar.update(len(chunk))
                                    # emit progress occasionally
                                    now = time.time()
                                    if (now - last_emit_t) >= 0.5:
                                        pct = (downloaded * 100.0) / float(total_size) if total_size > 0 else None
                                        _progress(stage="mp4", downloaded=downloaded, total=total_size or None, percent=pct, message=None)
                                        last_emit_t = now
                    else:
                        downloaded = 0
                        last_report_t = 0.0
                        last_report_pct = -1.0
                        last_report_downloaded = -1
                        for chunk in response.iter_content(chunk_size=1024 * 1024):
                            if _cancelled():
                                _status("Download cancelled.", "warning")
                                try:
                                    f.flush()
                                except Exception:
                                    pass
                                _cleanup_partial(save_path)
                                return False, None
                            if not chunk:
                                continue
                            f.write(chunk)
                            downloaded += len(chunk)
                            if total_size > 0:
                                pct = (downloaded * 100.0) / float(total_size)
                                now = time.time()
                                pct_step = 0.2 if pct >= 90 else 5.0
                                time_step = 2.0 if pct >= 90 else 5.0
                                should = (
                                    pct >= 100
                                    or (pct - last_report_pct) >= pct_step
                                    or (now - last_report_t) >= time_step
                                )
                                if should and downloaded != last_report_downloaded:
                                    pct_txt = f"{pct:.1f}%" if pct >= 90 else f"{int(pct)}%"
                                    _status(f"Progress: {pct_txt} ({downloaded}/{total_size} bytes)", "loading")
                                    _progress(stage="mp4", downloaded=downloaded, total=total_size, percent=float(pct), message=None)
                                    last_report_pct = pct
                                    last_report_downloaded = downloaded
                                    last_report_t = now
            
            _status("Download completed successfully!", "success")
            if total_size > 0:
                _progress(stage="done", downloaded=total_size, total=total_size, percent=100.0, message="ok", force=True)
            else:
                _progress(stage="done", downloaded=None, total=None, percent=None, message="ok", force=True)
            return True, save_path
    except Exception as e:
        if _cancelled():
            _status("Download cancelled.", "warning")
            _cleanup_partial(save_path)
            return False, None
        _status(f"Download failed: {str(e)}", "error")
        return False, None
    