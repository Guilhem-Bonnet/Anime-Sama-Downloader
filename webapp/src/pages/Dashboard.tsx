import React, { useEffect } from 'react';
import { useJobsStore } from '../stores/jobs.store';
import { useSubscriptionsStore } from '../stores/subscriptions.store';
import type { Job, Subscription } from '../api';
import HeroBanner from '../components/dashboard/HeroBanner';
import QuickActions from '../components/dashboard/QuickActions';
import JobsRecent from '../components/dashboard/JobsRecent';
import SubscriptionsSection from '../components/dashboard/SubscriptionsSection';
import '../styles/dashboard.css';

const Dashboard: React.FC = () => {
  const { jobs, loadJobs } = useJobsStore();
  const { subscriptions, loadSubscriptions } = useSubscriptionsStore();
  const [trendingAnime, setTrendingAnime] = React.useState<any>(null);

  useEffect(() => {
    loadJobs();
    loadSubscriptions();

    // Fetch trending anime from AniList (best-effort, cached in state)
    fetch('https://graphql.anilist.co', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        query: `{
          Page(page: 1, perPage: 1) {
            media(sort: TRENDING_DESC, type: ANIME) {
              id
              title { romaji }
              coverImage { large }
              seasonYear
              season
            }
          }
        }`
      })
    })
      .then(r => r.json())
      .then(data => {
        if (data.data?.Page?.media?.[0]) {
          setTrendingAnime(data.data.Page.media[0]);
        }
      })
      .catch(() => { /* AniList is best-effort */ });
  }, []);

  const mappedJobs = jobs.slice(0, 5).map(mapApiJob);
  const activeJobsCount = jobs.filter(j => j.state === 'running' || j.state === 'queued' || j.state === 'muxing').length;
  const newEpisodesCount = subscriptions.reduce((acc, sub) => {
    const diff = (sub.lastAvailableEpisode || 0) - (sub.lastDownloadedEpisode || 0);
    return acc + Math.max(0, diff);
  }, 0);

  return (
    <main className="dashboard" aria-label="Tableau de bord principal">
      <section className="dashboard__hero">
        <HeroBanner 
          jobs={mappedJobs}
          subscriptions={subscriptions}
          trendingAnime={trendingAnime}
        />
      </section>

      <section className="dashboard__actions">
        <QuickActions 
          jobsCount={activeJobsCount}
          subscriptionsCount={subscriptions.length}
          newEpisodesCount={newEpisodesCount}
        />
      </section>

      <section className="dashboard__jobs">
        <h2>Téléchargements Récents</h2>
        <JobsRecent jobs={mappedJobs} />
      </section>

      <section className="dashboard__subscriptions">
        <h2>Mes Abonnements</h2>
        <SubscriptionsSection subscriptions={subscriptions} />
      </section>
    </main>
  );
};

Dashboard.displayName = 'Dashboard';

export default Dashboard;

function mapApiJob(job: any) {
  const params = job.params || {};
  const status = mapJobState(job.state);
  return {
    id: job.id,
    animeTitle: params.animeTitle || params.animeId || 'Anime',
    episode: params.episodeNumber || 1,
    progress: Math.round((job.progress || 0) * 100),
    eta: undefined,
    speed: undefined,
    status,
  };
}

function mapJobState(state: string) {
  switch (state) {
    case 'queued':
      return 'queued';
    case 'running':
    case 'muxing':
      return 'downloading';
    case 'completed':
      return 'completed';
    case 'failed':
    case 'canceled':
      return 'failed';
    default:
      return 'queued';
  }
}
