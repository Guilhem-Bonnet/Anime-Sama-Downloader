import React from 'react';
import { useNavigate } from 'react-router-dom';
import { JobRow } from '../anime';

interface JobsRecentProps {
  jobs: any[];
}

const JobsRecent: React.FC<JobsRecentProps> = ({ jobs }) => {
  const navigate = useNavigate();

  if (!jobs || jobs.length === 0) {
    return (
      <div className="empty-state">
        <div className="empty-state__icon">📥</div>
        <h3 className="empty-state__title">Aucun téléchargement récent</h3>
        <p className="empty-state__description">
          Cherchez un anime pour commencer vos téléchargements
        </p>
        <button 
          className="empty-state__cta"
          onClick={() => navigate('/search')}
        >
          Rechercher un anime
        </button>
      </div>
    );
  }

  return (
    <div className="jobs-recent">
      <div className="jobs-recent__list">
        {jobs.slice(0, 5).map((job) => (
          <JobRow
            key={job.id}
            id={job.id}
            animeTitle={job.animeTitle}
            episode={job.episode}
            progress={job.progress || 0}
            eta={job.eta}
            speed={job.speed}
            status={job.status}
            onPause={() => {}}
            onResume={() => {}}
            onCancel={() => {}}
            onRetry={() => {}}
          />
        ))}
      </div>
      {jobs.length > 5 && (
        <button 
          className="jobs-recent__view-all"
          onClick={() => navigate('/downloads')}
        >
          Voir tous les téléchargements
        </button>
      )}
    </div>
  );
};

JobsRecent.displayName = 'JobsRecent';

export default JobsRecent;
