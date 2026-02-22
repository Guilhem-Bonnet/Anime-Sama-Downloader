import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { apiAnimeSamaResolve, apiAnimeSamaEnqueue } from '../api';
import { useJobsStore } from '../stores/jobs.store';

interface Episode {
  number: number;
  title: string;
  season_number: number;
  url: string;
}

interface Season {
  number: number;
  name: string;
  episodes: Episode[];
}

interface AnimeDetail {
  id: string;
  title: string;
  thumbnail_url: string;
  synopsis: string;
  year: number;
  status: string;
  genres: string[];
  episode_count: number;
  seasons: Season[];
}

export function AnimeDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const loadJobs = useJobsStore((s) => s.loadJobs);
  const [anime, setAnime] = useState<AnimeDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedEpisodes, setSelectedEpisodes] = useState<Set<string>>(new Set());
  const [activeSeason, setActiveSeason] = useState(1);
  const [downloading, setDownloading] = useState(false);

  useEffect(() => {
    if (!id) return;

    setLoading(true);
    setError(null);

    fetch(`/api/v1/anime/${id}`)
      .then((res) => {
        if (!res.ok) {
          if (res.status === 404) {
            throw new Error('Anime non trouvé');
          }
          throw new Error(`Erreur HTTP ${res.status}`);
        }
        return res.json();
      })
      .then((data: AnimeDetail) => {
        setAnime(data);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message || 'Erreur lors du chargement');
        setLoading(false);
      });
  }, [id]);

  const handleEpisodeToggle = (seasonNumber: number, episodeNumber: number) => {
    const key = `${seasonNumber}-${episodeNumber}`;
    setSelectedEpisodes((prev) => {
      const next = new Set(prev);
      if (next.has(key)) {
        next.delete(key);
      } else {
        next.add(key);
      }
      return next;
    });
  };

  const handleSelectAll = (season: Season) => {
    setSelectedEpisodes((prev) => {
      const next = new Set(prev);
      season.episodes.forEach((ep) => {
        next.add(`${season.number}-${ep.number}`);
      });
      return next;
    });
  };

  const handleDeselectAll = (season: Season) => {
    setSelectedEpisodes((prev) => {
      const next = new Set(prev);
      season.episodes.forEach((ep) => {
        next.delete(`${season.number}-${ep.number}`);
      });
      return next;
    });
  };

  const handleDownload = async () => {
    if (selectedEpisodes.size === 0 || !anime) return;
    setDownloading(true);

    try {
      // Groupe les épisodes sélectionnés par saison
      const bySeason = new Map<number, number[]>();
      Array.from(selectedEpisodes).forEach((key) => {
        const [seasonStr, epStr] = key.split('-');
        const season = Number(seasonStr);
        const ep = Number(epStr);
        if (!bySeason.has(season)) bySeason.set(season, []);
        bySeason.get(season)!.push(ep);
      });

      let totalEnqueued = 0;
      const errors: string[] = [];

      // 1. Résolution du catalogueUrl (une seule fois pour tout l'anime)
      let catalogueUrl: string | null = null;
      try {
        const resolveResp = await apiAnimeSamaResolve({
          titles: [anime.title],
          maxCandidates: 3,
        });
        catalogueUrl = resolveResp.candidates[0]?.catalogueUrl ?? null;
      } catch (e: any) {
        errors.push(`Résolution impossible : ${e.message ?? 'erreur réseau'}`);
      }

      if (!catalogueUrl) {
        errors.push(`"${anime.title}" introuvable sur anime-sama — ajoute-le via un abonnement manuel`);
      } else {
        // 2. Enqueue par saison — essaie vostfr puis vf en fallback
        //    (pas de scan préalable : on tente directement et on gère les erreurs)
        const LANGS = ['vostfr', 'vf'];
        const catalogueBase = catalogueUrl.replace(/\/?$/, '/');

        for (const [season, episodes] of bySeason) {
          let enqueued = false;

          for (const lang of LANGS) {
            const baseUrl = `${catalogueBase}saison${season}/${lang}/`;
            try {
              const resp = await apiAnimeSamaEnqueue({
                baseUrl,
                label: anime.title,
                episodes,
              });

              // episodes.js trouvé — même si tous les éps sont skipped c'est ok,
              // on a la bonne source
              totalEnqueued += resp.enqueuedEpisodes.length;
              resp.skipped.forEach((s) => {
                errors.push(`Épisode ${s.episode} ignoré : ${s.reason}`);
              });
              enqueued = true;
              break; // source trouvée, pas besoin d'essayer la langue suivante
            } catch (_) {
              // 502 = episodes.js 404 → on essaie la langue suivante
              continue;
            }
          }

          if (!enqueued) {
            errors.push(`Saison ${season} : aucune source disponible (vostfr/vf) — épisodes peut-être pas encore uploadés sur anime-sama`);
          }
        }
      }

      // Si rien n'a été enqueué, montrer une erreur claire — NE PAS naviguer
      if (totalEnqueued === 0) {
        const msg = errors.length > 0
          ? `Aucun épisode lancé :\n\n${errors.join('\n')}`
          : `Aucun épisode disponible pour "${anime.title}"`;
        alert(msg);
        return;
      }

      // Au moins un job créé → rafraîchir le store et naviguer
      await loadJobs();
      setSelectedEpisodes(new Set());
      navigate('/downloads');

      // Avertissements partiels (certains ont quand même été lancés)
      if (errors.length > 0) {
        console.warn('[AnimeDetailPage] avertissements partiels:', errors);
      }
    } catch (err: any) {
      alert(err.message || 'Échec de la création des téléchargements');
    } finally {
      setDownloading(false);
    }
  };

  if (loading) {
    return (
      <div className="container">
        <div style={{ textAlign: 'center', padding: '48px' }}>
          <div style={{ fontSize: '32px', marginBottom: '16px' }}>⏳</div>
          <p className="muted">Chargement des détails...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container">
        <div className="card" style={{ padding: '32px', textAlign: 'center' }}>
          <div style={{ fontSize: '48px', marginBottom: '16px' }}>❌</div>
          <h2>Erreur</h2>
          <p className="muted" style={{ marginTop: '8px' }}>{error}</p>
          <button className="btn primary" onClick={() => navigate(-1)} style={{ marginTop: '24px' }}>
            ← Retour
          </button>
        </div>
      </div>
    );
  }

  if (!anime) {
    return null;
  }

  const currentSeason = anime.seasons.find((s) => s.number === activeSeason) || anime.seasons[0];

  return (
    <div className="container">
      {/* Header */}
      <div className="topbar">
        <button className="btn sm" onClick={() => navigate(-1)}>
          ← Retour
        </button>
      </div>

      {/* Anime Info Card */}
      <div className="card" style={{ marginTop: '24px' }}>
        <div style={{ display: 'flex', gap: '24px', flexWrap: 'wrap' }}>
          {/* Thumbnail */}
          <img
            src={anime.thumbnail_url}
            alt={anime.title}
            onError={(e) => {
              const target = e.currentTarget;
              if (target.src !== '/assets/cover-placeholder.svg') {
                target.src = '/assets/cover-placeholder.svg';
              }
            }}
            style={{
              width: '200px',
              height: '280px',
              objectFit: 'cover',
              borderRadius: '12px',
              border: '1px solid var(--border)',
            }}
          />

          {/* Info */}
          <div style={{ flex: 1, minWidth: '300px' }}>
            <h1 className="h1" style={{ marginBottom: '12px' }}>{anime.title}</h1>
            
            <div style={{ display: 'flex', gap: '8px', marginBottom: '16px', flexWrap: 'wrap' }}>
              <span className="badge">{anime.year}</span>
              <span className={`badge ${anime.status === 'completed' ? 'ok' : anime.status === 'ongoing' ? 'run' : ''}`}>
                {anime.status}
              </span>
              <span className="badge">{anime.episode_count} épisodes</span>
            </div>

            {/* Genres */}
            <div style={{ display: 'flex', gap: '6px', marginBottom: '16px', flexWrap: 'wrap' }}>
              {anime.genres.map((genre) => (
                <span key={genre} className="pill info">
                  {genre}
                </span>
              ))}
            </div>

            {/* Synopsis */}
            <div>
              <h3 style={{ fontSize: '14px', fontWeight: 600, marginBottom: '8px' }}>Synopsis</h3>
              <p className="muted" style={{ maxHeight: '200px', overflow: 'auto', lineHeight: '1.6' }}>
                {anime.synopsis}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Episodes Section */}
      <div className="card" style={{ marginTop: '24px' }}>
        <div className="cardTitle">Épisodes</div>

        {/* Season Tabs (if multiple seasons) */}
        {anime.seasons.length > 1 && (
          <div style={{ display: 'flex', gap: '8px', marginTop: '16px', marginBottom: '16px', flexWrap: 'wrap' }}>
            {anime.seasons.map((season) => (
              <button
                key={season.number}
                className={`btn ${activeSeason === season.number ? 'primary' : ''}`}
                onClick={() => setActiveSeason(season.number)}
              >
                {season.name}
              </button>
            ))}
          </div>
        )}

        {/* Select All / Deselect All */}
        <div style={{ display: 'flex', gap: '8px', marginTop: '16px', marginBottom: '16px' }}>
          <button className="btn sm" onClick={() => handleSelectAll(currentSeason)}>
            ✓ Tout sélectionner
          </button>
          <button className="btn sm" onClick={() => handleDeselectAll(currentSeason)}>
            ✗ Tout désélectionner
          </button>
        </div>

        {/* Episode Grid */}
        <div className="epgrid">
          {currentSeason.episodes.map((episode) => {
            const key = `${currentSeason.number}-${episode.number}`;
            const isSelected = selectedEpisodes.has(key);

            return (
              <label key={key} className="ep" style={{ cursor: 'pointer' }}>
                <input
                  type="checkbox"
                  checked={isSelected}
                  onChange={() => handleEpisodeToggle(currentSeason.number, episode.number)}
                />
                <span title={episode.title || `Épisode ${episode.number}`}>
                  {episode.number}
                </span>
              </label>
            );
          })}
        </div>

        {/* Download Button */}
        <button
          className="btn primary"
          disabled={selectedEpisodes.size === 0 || downloading}
          onClick={handleDownload}
          style={{ marginTop: '24px', width: '100%' }}
        >
          {downloading
            ? 'Création des jobs…'
            : selectedEpisodes.size === 0
            ? 'Sélectionnez des épisodes'
            : `Télécharger ${selectedEpisodes.size} épisode${selectedEpisodes.size > 1 ? 's' : ''}`}
        </button>
      </div>
    </div>
  );
}
