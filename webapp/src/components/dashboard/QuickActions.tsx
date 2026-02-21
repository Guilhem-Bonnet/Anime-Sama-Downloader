import React from 'react';
import { useNavigate } from 'react-router-dom';

interface QuickActionsProps {
  jobsCount: number;
  subscriptionsCount: number;
  newEpisodesCount: number;
}

const QuickActions: React.FC<QuickActionsProps> = ({
  jobsCount,
  subscriptionsCount,
  newEpisodesCount,
}) => {
  const navigate = useNavigate();

  const actions = [
    {
      id: 'search',
      title: 'Rechercher anime',
      icon: '🔍',
      onClick: () => navigate('/search'),
    },
    {
      id: 'jobs',
      title: 'Téléchargements',
      icon: '⬇️',
      count: jobsCount,
      onClick: () => navigate('/downloads'),
      countLabel: jobsCount > 0 ? `${jobsCount} en cours` : '',
    },
    {
      id: 'subscriptions',
      title: 'Mes abonnements',
      icon: '🔔',
      count: newEpisodesCount,
      onClick: () => navigate('/search'),
      countLabel: newEpisodesCount > 0 ? `${newEpisodesCount} nouveaux` : '',
    },
  ];

  return (
    <div className="quick-actions">
      {actions.map((action) => (
        <button
          key={action.id}
          className="quick-action-card"
          onClick={action.onClick}
        >
          <span className="quick-action-card__icon">{action.icon}</span>
          <h3 className="quick-action-card__title">{action.title}</h3>
          {action.countLabel && (
            <span className="quick-action-card__badge">{action.countLabel}</span>
          )}
        </button>
      ))}
    </div>
  );
};

QuickActions.displayName = 'QuickActions';

export default QuickActions;
