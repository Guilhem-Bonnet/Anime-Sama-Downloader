/**
 * Domain types - Shared TypeScript enums and types
 * Used across all anime components to ensure consistency
 */

// ============================================
// Job & Download Status
// ============================================

export type JobStatus = 
  | 'queued'      // En attente dans la queue
  | 'downloading' // Téléchargement en cours
  | 'paused'      // Mis en pause par l'utilisateur
  | 'completed'   // Téléchargé avec succès
  | 'failed';     // Échec après retries

// ============================================
// Anime Metadata
// ============================================

export type Language = 
  | 'VOSTFR'  // Version originale sous-titrée français
  | 'VF';     // Version française doublée

export type AnimeStatus = 
  | 'ongoing'   // En cours de diffusion
  | 'completed' // Série terminée
  | 'upcoming'; // À venir

export type Quality = 
  | '1080p'
  | '720p'
  | '480p';

// ============================================
// Subscription Status
// ============================================

export type SubscriptionStatus = 
  | 'active'  // Abonnement actif (sync auto)
  | 'paused'; // Abonnement en pause (pas de sync)

// ============================================
// Calendar Event Status
// ============================================

export type CalendarEventStatus = 
  | 'scheduled'  // Prévu dans le futur
  | 'available'  // Disponible pour téléchargement
  | 'downloaded'; // Déjà téléchargé

// ============================================
// Notification Types
// ============================================

export type NotificationType = 
  | 'success'
  | 'error'
  | 'warning'
  | 'info';

// ============================================
// Badge Variants (from 2-12)
// ============================================

export type BadgeVariant = 
  | 'neutral'  // Gris (queued, default)
  | 'primary'  // Rouge accent (active, downloading)
  | 'success'  // Vert (completed, available)
  | 'warning'  // Orange (paused, upcoming)
  | 'error';   // Rouge foncé (failed, error)

// ============================================
// Common Interfaces
// ============================================

export interface Anime {
  id: string;
  title: string;
  coverUrl: string;
  season?: string;
  language: Language;
  status: AnimeStatus;
  anilistId?: number;
  malId?: number;
}

export interface Episode {
  number: number;
  title?: string;
  duration?: string;
  releaseDate?: string;
}

export interface Job {
  id: string;
  animeTitle: string;
  episode: number;
  status: JobStatus;
  progress?: number; // 0-100
  eta?: string;
  speed?: string;
  errorMessage?: string;
  createdAt: string;
}

export interface Subscription {
  id: string;
  animeId: string;
  animeTitle: string;
  season: string;
  coverUrl: string;
  language: Language;
  updateFrequency: 'weekly' | 'daily';
  lastDownloaded: number;
  nextAvailable?: {
    episode: number;
    eta: string;
  };
  newEpisodesCount: number;
  status: SubscriptionStatus;
}

export interface CalendarEvent {
  id: string;
  animeTitle: string;
  episode: number;
  releaseTime: string; // ISO datetime
  coverUrl: string;
  status: CalendarEventStatus;
}

// ============================================
// API Response Types
// ============================================

export interface SearchResult {
  results: Anime[];
  total: number;
  page: number;
  pageSize: number;
}

export interface JobsResponse {
  jobs: Job[];
  stats: {
    active: number;
    queued: number;
    completedToday: number;
  };
}

export interface SubscriptionsResponse {
  subscriptions: Subscription[];
  total: number;
}

export interface CalendarResponse {
  events: CalendarEvent[];
  startDate: string;
  endDate: string;
}
