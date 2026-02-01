import React from 'react';
import { X } from 'lucide-react';

export interface SearchFilters {
  genres: string[];
  status: string;
  yearMin: number;
  yearMax: number;
}

interface FilterPanelProps {
  filters: SearchFilters;
  onFiltersChange: (filters: SearchFilters) => void;
}

const AVAILABLE_GENRES = [
  'Action',
  'Adventure',
  'Comedy',
  'Drama',
  'Fantasy',
  'Horror',
  'Mecha',
  'Mystery',
  'Psychological',
  'Romance',
  'Sci-Fi',
  'Shonen',
  'Slice of Life',
  'Sports',
  'Supernatural',
  'Superhero',
  'Thriller',
];

const STATUS_OPTIONS = [
  { value: '', label: 'All Status' },
  { value: 'ongoing', label: 'Ongoing' },
  { value: 'completed', label: 'Completed' },
  { value: 'planning', label: 'Planning' },
];

export const FilterPanel: React.FC<FilterPanelProps> = ({ filters, onFiltersChange }) => {
  const [isExpanded, setIsExpanded] = React.useState(false);

  const toggleGenre = (genre: string) => {
    const newGenres = filters.genres.includes(genre)
      ? filters.genres.filter((g) => g !== genre)
      : [...filters.genres, genre];
    onFiltersChange({ ...filters, genres: newGenres });
  };

  const removeGenre = (genre: string) => {
    onFiltersChange({
      ...filters,
      genres: filters.genres.filter((g) => g !== genre),
    });
  };

  const setStatus = (status: string) => {
    onFiltersChange({ ...filters, status });
  };

  const setYearMin = (year: number) => {
    onFiltersChange({ ...filters, yearMin: year });
  };

  const setYearMax = (year: number) => {
    onFiltersChange({ ...filters, yearMax: year });
  };

  const clearAllFilters = () => {
    onFiltersChange({
      genres: [],
      status: '',
      yearMin: 0,
      yearMax: 0,
    });
  };

  const hasActiveFilters =
    filters.genres.length > 0 ||
    filters.status !== '' ||
    filters.yearMin > 0 ||
    filters.yearMax > 0;

  return (
    <div className="w-full max-w-2xl mx-auto mb-6">
      {/* Filter toggle button */}
      <button
        onClick={() => setIsExpanded(!isExpanded)}
        className="flex items-center justify-between w-full px-4 py-2 bg-gray-100 dark:bg-gray-800 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
      >
        <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
          {isExpanded ? 'Hide Filters' : 'Show Filters'}
          {hasActiveFilters && (
            <span className="ml-2 text-xs text-cyan-500">
              ({filters.genres.length +
                (filters.status ? 1 : 0) +
                (filters.yearMin || filters.yearMax ? 1 : 0)}{' '}
              active)
            </span>
          )}
        </span>
        <svg
          className={`w-5 h-5 transition-transform ${isExpanded ? 'rotate-180' : ''}`}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
      </button>

      {/* Filter panel */}
      {isExpanded && (
        <div className="mt-4 p-4 bg-gray-50 dark:bg-gray-900 rounded-lg space-y-4">
          {/* Genre filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Genres
            </label>
            <div className="flex flex-wrap gap-2">
              {AVAILABLE_GENRES.map((genre) => (
                <button
                  key={genre}
                  onClick={() => toggleGenre(genre)}
                  className={`px-3 py-1 text-sm rounded-full transition-colors ${
                    filters.genres.includes(genre)
                      ? 'bg-cyan-500 text-white'
                      : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600'
                  }`}
                >
                  {genre}
                </button>
              ))}
            </div>
          </div>

          {/* Status filter */}
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Status
            </label>
            <select
              value={filters.status}
              onChange={(e) => setStatus(e.target.value)}
              className="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-cyan-500"
            >
              {STATUS_OPTIONS.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>

          {/* Year range filter */}
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Year Min
              </label>
              <input
                type="number"
                value={filters.yearMin || ''}
                onChange={(e) => setYearMin(parseInt(e.target.value) || 0)}
                placeholder="1990"
                className="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-cyan-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Year Max
              </label>
              <input
                type="number"
                value={filters.yearMax || ''}
                onChange={(e) => setYearMax(parseInt(e.target.value) || 0)}
                placeholder="2024"
                className="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-cyan-500"
              />
            </div>
          </div>

          {/* Clear filters button */}
          {hasActiveFilters && (
            <button
              onClick={clearAllFilters}
              className="w-full px-4 py-2 text-sm text-red-500 hover:text-red-600 dark:text-red-400 dark:hover:text-red-300 font-medium"
            >
              Clear All Filters
            </button>
          )}
        </div>
      )}

      {/* Active filter chips */}
      {hasActiveFilters && (
        <div className="mt-3 flex flex-wrap gap-2">
          {filters.genres.map((genre) => (
            <div
              key={genre}
              className="flex items-center gap-1 px-3 py-1 bg-cyan-100 dark:bg-cyan-900 text-cyan-800 dark:text-cyan-200 rounded-full text-sm"
            >
              <span>{genre}</span>
              <button
                onClick={() => removeGenre(genre)}
                className="hover:bg-cyan-200 dark:hover:bg-cyan-800 rounded-full p-0.5"
              >
                <X className="w-3 h-3" />
              </button>
            </div>
          ))}
          {filters.status && (
            <div className="flex items-center gap-1 px-3 py-1 bg-purple-100 dark:bg-purple-900 text-purple-800 dark:text-purple-200 rounded-full text-sm">
              <span>Status: {filters.status}</span>
              <button
                onClick={() => setStatus('')}
                className="hover:bg-purple-200 dark:hover:bg-purple-800 rounded-full p-0.5"
              >
                <X className="w-3 h-3" />
              </button>
            </div>
          )}
          {(filters.yearMin > 0 || filters.yearMax > 0) && (
            <div className="flex items-center gap-1 px-3 py-1 bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200 rounded-full text-sm">
              <span>
                Year:{' '}
                {filters.yearMin > 0 && filters.yearMax > 0
                  ? `${filters.yearMin}-${filters.yearMax}`
                  : filters.yearMin > 0
                  ? `${filters.yearMin}+`
                  : `≤${filters.yearMax}`}
              </span>
              <button
                onClick={() => onFiltersChange({ ...filters, yearMin: 0, yearMax: 0 })}
                className="hover:bg-green-200 dark:hover:bg-green-800 rounded-full p-0.5"
              >
                <X className="w-3 h-3" />
              </button>
            </div>
          )}
        </div>
      )}
    </div>
  );
};
