import React from 'react';

export const EmptySearchIllustration: React.FC = () => {
  return (
    <div className="sakura-illustration">
      <svg
        width="320"
        height="200"
        viewBox="0 0 320 200"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        className="sakura-float"
        role="img"
        aria-label="Illustration de recherche"
      >
        <defs>
          <linearGradient id="sakura-search-gradient" x1="0" y1="0" x2="1" y2="1">
            <stop offset="0%" stopColor="#d946ef" />
            <stop offset="50%" stopColor="#fb6f8a" />
            <stop offset="100%" stopColor="#06b6d4" />
          </linearGradient>
          <radialGradient id="sakura-search-glow" cx="0" cy="0" r="1" gradientUnits="userSpaceOnUse" gradientTransform="translate(160 160) rotate(90) scale(90 140)">
            <stop offset="0%" stopColor="#d946ef" stopOpacity="0.28" />
            <stop offset="100%" stopColor="#06b6d4" stopOpacity="0" />
          </radialGradient>
          <linearGradient id="sakura-search-panel" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="#111827" />
            <stop offset="100%" stopColor="#0b1220" />
          </linearGradient>
        </defs>

        {/* Glow */}
        <ellipse cx="160" cy="165" rx="130" ry="28" fill="url(#sakura-search-glow)" className="sakura-pulse" />

        {/* Background stars */}
        <circle cx="52" cy="38" r="2" fill="#94a3b8" opacity="0.6" />
        <circle cx="282" cy="42" r="1.5" fill="#cbd5f5" opacity="0.7" />
        <circle cx="268" cy="150" r="1.5" fill="#94a3b8" opacity="0.6" />

        {/* Card */}
        <rect x="54" y="48" width="212" height="118" rx="18" fill="url(#sakura-search-panel)" opacity="0.95" />
        <rect x="60" y="54" width="200" height="106" rx="16" stroke="#1f2937" strokeWidth="1.5" />

        {/* Top bar */}
        <rect x="75" y="66" width="170" height="16" rx="8" fill="#1f2937" />
        <rect x="83" y="70" width="52" height="8" rx="4" fill="#374151" />
        <rect x="142" y="70" width="36" height="8" rx="4" fill="#334155" />
        <rect x="184" y="70" width="32" height="8" rx="4" fill="#334155" />

        {/* Content lines */}
        <rect x="75" y="92" width="130" height="10" rx="5" fill="#374151" />
        <rect x="75" y="108" width="150" height="10" rx="5" fill="#2b3648" />
        <rect x="75" y="124" width="90" height="10" rx="5" fill="#374151" />

        {/* Lens */}
        <circle cx="210" cy="95" r="18" stroke="url(#sakura-search-gradient)" strokeWidth="4" />
        <line x1="222" y1="107" x2="238" y2="123" stroke="url(#sakura-search-gradient)" strokeWidth="4" strokeLinecap="round" />

        {/* Sakura petals */}
        <path d="M105 40C110 32 122 32 127 40C130 46 124 54 116 54C108 54 102 46 105 40Z" fill="#fb6f8a" opacity="0.8" className="sakura-drift" />
        <path d="M240 45C244 39 254 39 258 45C261 50 256 57 249 57C242 57 237 50 240 45Z" fill="#06b6d4" opacity="0.7" className="sakura-drift" />
        <path d="M78 140C82 134 92 134 96 140C99 145 94 152 87 152C80 152 75 145 78 140Z" fill="#d946ef" opacity="0.65" className="sakura-drift" />
        <path d="M258 130C262 124 272 124 276 130C279 135 274 142 267 142C260 142 255 135 258 130Z" fill="#fb6f8a" opacity="0.6" className="sakura-drift" />
      </svg>
    </div>
  );
};

export const EmptyDownloadsIllustration: React.FC = () => {
  return (
    <div className="sakura-illustration">
      <svg
        width="320"
        height="200"
        viewBox="0 0 320 200"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        className="sakura-float"
        role="img"
        aria-label="Illustration de téléchargements"
      >
        <defs>
          <linearGradient id="sakura-download-gradient" x1="0" y1="0" x2="1" y2="1">
            <stop offset="0%" stopColor="#06b6d4" />
            <stop offset="50%" stopColor="#d946ef" />
            <stop offset="100%" stopColor="#fb6f8a" />
          </linearGradient>
          <linearGradient id="sakura-box" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="#0f172a" />
            <stop offset="100%" stopColor="#0b1220" />
          </linearGradient>
        </defs>

        {/* Base glow */}
        <ellipse cx="160" cy="170" rx="120" ry="24" fill="#06b6d4" opacity="0.12" className="sakura-pulse" />

        {/* Floating chips */}
        <rect x="60" y="48" width="40" height="18" rx="9" fill="#1f2937" opacity="0.8" />
        <rect x="220" y="52" width="48" height="18" rx="9" fill="#1f2937" opacity="0.8" />

        {/* Box */}
        <rect x="90" y="58" width="140" height="96" rx="16" fill="url(#sakura-box)" stroke="#1f2937" strokeWidth="2" />
        <rect x="100" y="68" width="120" height="14" rx="7" fill="#1f2937" />
        <rect x="100" y="90" width="90" height="10" rx="5" fill="#334155" />
        <rect x="100" y="108" width="110" height="10" rx="5" fill="#2b3648" />
        <rect x="100" y="126" width="70" height="10" rx="5" fill="#334155" />

        {/* Arrow */}
        <path d="M160 40V80" stroke="url(#sakura-download-gradient)" strokeWidth="6" strokeLinecap="round" />
        <path d="M146 68L160 86L174 68" stroke="url(#sakura-download-gradient)" strokeWidth="6" strokeLinecap="round" strokeLinejoin="round" />
        <circle cx="160" cy="34" r="6" fill="#06b6d4" opacity="0.6" className="sakura-pulse" />

        {/* Petals */}
        <path d="M70 90C74 84 84 84 88 90C91 95 86 102 79 102C72 102 67 95 70 90Z" fill="#d946ef" opacity="0.8" className="sakura-drift" />
        <path d="M250 88C254 82 264 82 268 88C271 93 266 100 259 100C252 100 247 93 250 88Z" fill="#06b6d4" opacity="0.8" className="sakura-drift" />
        <path d="M230 132C234 126 244 126 248 132C251 137 246 144 239 144C232 144 227 137 230 132Z" fill="#fb6f8a" opacity="0.7" className="sakura-drift" />
        <path d="M86 132C90 126 100 126 104 132C107 137 102 144 95 144C88 144 83 137 86 132Z" fill="#06b6d4" opacity="0.7" className="sakura-drift" />
      </svg>
    </div>
  );
};

export const SuccessIllustration: React.FC = () => {
  return (
    <div className="sakura-illustration">
      <svg
        width="240"
        height="180"
        viewBox="0 0 240 180"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        className="sakura-float"
        role="img"
        aria-label="Illustration de succès"
      >
        <defs>
          <linearGradient id="sakura-success-gradient" x1="0" y1="0" x2="1" y2="1">
            <stop offset="0%" stopColor="#22d3ee" />
            <stop offset="100%" stopColor="#4ade80" />
          </linearGradient>
        </defs>

        <circle cx="120" cy="90" r="52" fill="#0f172a" stroke="url(#sakura-success-gradient)" strokeWidth="6" />
        <path d="M92 92L112 112L148 74" stroke="url(#sakura-success-gradient)" strokeWidth="8" strokeLinecap="round" strokeLinejoin="round" />
        <circle cx="120" cy="155" r="30" fill="#22d3ee" opacity="0.12" className="sakura-pulse" />
      </svg>
    </div>
  );
};
