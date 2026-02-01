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
            <stop offset="0%" stopColor="#8b7a5a" />
            <stop offset="55%" stopColor="#6b7280" />
            <stop offset="100%" stopColor="#4b5563" />
          </linearGradient>
          <radialGradient id="sakura-search-glow" cx="0" cy="0" r="1" gradientUnits="userSpaceOnUse" gradientTransform="translate(160 160) rotate(90) scale(90 140)">
            <stop offset="0%" stopColor="#8b7a5a" stopOpacity="0.14" />
            <stop offset="100%" stopColor="#0f172a" stopOpacity="0" />
          </radialGradient>
          <linearGradient id="sakura-search-panel" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="#111827" />
            <stop offset="100%" stopColor="#0b0f19" />
          </linearGradient>
          <linearGradient id="sakura-ink-stroke" x1="0" y1="0" x2="1" y2="0">
            <stop offset="0%" stopColor="#3f4a34" />
            <stop offset="100%" stopColor="#5c6b4a" />
          </linearGradient>
          <linearGradient id="sakura-paper" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="#0f172a" stopOpacity="0.92" />
            <stop offset="100%" stopColor="#0b0f19" stopOpacity="0.98" />
          </linearGradient>
        </defs>

        {/* Ornate frame */}
        <rect x="10" y="10" width="300" height="180" rx="18" stroke="#6e4f2e" strokeWidth="2" opacity="0.7" />
        <rect x="18" y="18" width="284" height="164" rx="14" stroke="#a8926b" strokeWidth="1" opacity="0.5" />
        <path d="M18 34C26 22 40 18 56 18" stroke="#6e4f2e" strokeWidth="2" strokeLinecap="round" opacity="0.6" />
        <path d="M302 34C294 22 280 18 264 18" stroke="#6e4f2e" strokeWidth="2" strokeLinecap="round" opacity="0.6" />
        <path d="M18 166C26 178 40 182 56 182" stroke="#6e4f2e" strokeWidth="2" strokeLinecap="round" opacity="0.6" />
        <path d="M302 166C294 178 280 182 264 182" stroke="#6e4f2e" strokeWidth="2" strokeLinecap="round" opacity="0.6" />

        {/* Subtle base */}
        <ellipse cx="160" cy="165" rx="130" ry="24" fill="url(#sakura-search-glow)" />

        {/* Brushy sky */}
        <path d="M28 54C58 34 98 30 136 42C164 50 182 60 206 58C236 56 260 42 292 44" stroke="#334155" strokeWidth="2" strokeLinecap="round" opacity="0.5" />
        <path d="M24 68C68 52 112 50 150 60C178 68 202 70 234 66C258 62 282 52 304 54" stroke="#475569" strokeWidth="1.6" strokeLinecap="round" opacity="0.45" />

        {/* Card */}
        <rect x="54" y="52" width="212" height="112" rx="18" fill="url(#sakura-paper)" opacity="0.96" />
        <rect x="60" y="58" width="200" height="100" rx="16" stroke="#1f2937" strokeWidth="1.2" />

        {/* Top bar */}
        <rect x="75" y="66" width="170" height="16" rx="8" fill="#1f2937" />
        <rect x="83" y="70" width="52" height="8" rx="4" fill="#374151" />
        <rect x="142" y="70" width="36" height="8" rx="4" fill="#334155" />
        <rect x="184" y="70" width="32" height="8" rx="4" fill="#334155" />

        {/* Content lines */}
        <rect x="75" y="92" width="130" height="10" rx="5" fill="#374151" />
        <rect x="75" y="108" width="150" height="10" rx="5" fill="#2b3648" />
        <rect x="75" y="124" width="90" height="10" rx="5" fill="#374151" />

        {/* Thumbnails */}
        <rect x="75" y="140" width="24" height="12" rx="3" fill="#1f2937" />
        <rect x="104" y="140" width="24" height="12" rx="3" fill="#1f2937" />
        <rect x="133" y="140" width="24" height="12" rx="3" fill="#1f2937" />

        {/* Lens */}
        <circle cx="210" cy="95" r="18" stroke="url(#sakura-search-gradient)" strokeWidth="3" />
        <line x1="222" y1="107" x2="238" y2="123" stroke="url(#sakura-search-gradient)" strokeWidth="3" strokeLinecap="round" />

        {/* Ink mountains */}
        <path d="M58 156L96 118L132 152L164 124L214 160H58Z" fill="#0f172a" opacity="0.85" />
        <path d="M86 156L114 130L142 156H86Z" fill="#1f2937" opacity="0.9" />
        <path d="M168 160L196 134L234 160H168Z" fill="#1f2937" opacity="0.85" />

        {/* Fog layer */}
        <path d="M62 150C86 144 122 144 150 148C182 152 204 154 246 150" stroke="#334155" strokeWidth="2" strokeLinecap="round" opacity="0.25" />

        {/* Paper fibers */}
        <path d="M72 96C88 92 110 92 128 98" stroke="url(#sakura-ink-stroke)" strokeWidth="1.5" strokeLinecap="round" opacity="0.35" />
        <path d="M88 112C106 108 122 110 140 116" stroke="#5c6b4a" strokeWidth="1.2" strokeLinecap="round" opacity="0.3" />
        <path d="M170 112C186 108 204 110 222 116" stroke="#475569" strokeWidth="1.1" strokeLinecap="round" opacity="0.25" />
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
            <stop offset="0%" stopColor="#8b7a5a" />
            <stop offset="50%" stopColor="#6b7280" />
            <stop offset="100%" stopColor="#4b5563" />
          </linearGradient>
          <linearGradient id="sakura-box" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="#0f172a" />
            <stop offset="100%" stopColor="#0b1220" />
          </linearGradient>
          <linearGradient id="sakura-wind" x1="0" y1="0" x2="1" y2="0">
            <stop offset="0%" stopColor="#3f4a34" />
            <stop offset="100%" stopColor="#6b7b4e" />
          </linearGradient>
          <linearGradient id="sakura-wood" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="#2b2f3a" />
            <stop offset="100%" stopColor="#1c2230" />
          </linearGradient>
        </defs>

        {/* Ornate frame */}
        <rect x="10" y="10" width="300" height="180" rx="18" stroke="#6e4f2e" strokeWidth="2" opacity="0.7" />
        <rect x="18" y="18" width="284" height="164" rx="14" stroke="#a8926b" strokeWidth="1" opacity="0.5" />
        <path d="M18 34C26 22 40 18 56 18" stroke="#6e4f2e" strokeWidth="2" strokeLinecap="round" opacity="0.6" />
        <path d="M302 34C294 22 280 18 264 18" stroke="#6e4f2e" strokeWidth="2" strokeLinecap="round" opacity="0.6" />
        <path d="M18 166C26 178 40 182 56 182" stroke="#6e4f2e" strokeWidth="2" strokeLinecap="round" opacity="0.6" />
        <path d="M302 166C294 178 280 182 264 182" stroke="#6e4f2e" strokeWidth="2" strokeLinecap="round" opacity="0.6" />

        {/* Base shade */}
        <ellipse cx="160" cy="170" rx="120" ry="22" fill="#1f2937" opacity="0.2" />

        {/* Wind lines */}
        <path d="M24 86C52 78 86 78 112 84" stroke="url(#sakura-wind)" strokeWidth="2" strokeLinecap="round" opacity="0.45" />
        <path d="M208 78C232 72 264 72 294 80" stroke="#6b7b4e" strokeWidth="1.8" strokeLinecap="round" opacity="0.4" />

        {/* Floating chips */}
        <rect x="60" y="48" width="40" height="18" rx="9" fill="#1f2937" opacity="0.8" />
        <rect x="220" y="52" width="48" height="18" rx="9" fill="#1f2937" opacity="0.8" />

        {/* Box */}
        <rect x="90" y="58" width="140" height="96" rx="16" fill="url(#sakura-box)" stroke="#1f2937" strokeWidth="2" />
        <rect x="96" y="64" width="128" height="8" rx="4" fill="url(#sakura-wood)" opacity="0.9" />
        <rect x="100" y="68" width="120" height="14" rx="7" fill="#1f2937" />
        <rect x="100" y="90" width="90" height="10" rx="5" fill="#334155" />
        <rect x="100" y="108" width="110" height="10" rx="5" fill="#2b3648" />
        <rect x="100" y="126" width="70" height="10" rx="5" fill="#334155" />

        {/* Arrow */}
        <path d="M160 40V80" stroke="url(#sakura-download-gradient)" strokeWidth="5" strokeLinecap="round" />
        <path d="M146 68L160 86L174 68" stroke="url(#sakura-download-gradient)" strokeWidth="5" strokeLinecap="round" strokeLinejoin="round" />
        <circle cx="160" cy="34" r="6" fill="#8b7a5a" opacity="0.45" />

        {/* Leaves */}
        <path d="M72 104C78 96 90 96 96 104C99 110 92 116 84 116C76 116 69 110 72 104Z" fill="#5c6b4a" opacity="0.55" className="sakura-drift" />
        <path d="M246 102C252 94 264 94 270 102C273 108 266 114 258 114C250 114 243 108 246 102Z" fill="#6b7b4e" opacity="0.55" className="sakura-drift" />
        <path d="M232 138C238 130 250 130 256 138C259 144 252 150 244 150C236 150 229 144 232 138Z" fill="#4b5563" opacity="0.5" className="sakura-drift" />
        <path d="M86 138C92 130 104 130 110 138C113 144 106 150 98 150C90 150 83 144 86 138Z" fill="#64748b" opacity="0.5" className="sakura-drift" />
      </svg>
    </div>
  );
};

export const HeroLandscapeIllustration: React.FC = () => {
  return (
    <div className="sakura-illustration">
      <svg
        width="720"
        height="240"
        viewBox="0 0 720 240"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
        className="sakura-float"
        role="img"
        aria-label="Illustration paysage"
      >
        <defs>
          <linearGradient id="hero-sky" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="#0f172a" />
            <stop offset="100%" stopColor="#111827" />
          </linearGradient>
          <linearGradient id="hero-ridge" x1="0" y1="0" x2="1" y2="0">
            <stop offset="0%" stopColor="#1f2937" />
            <stop offset="100%" stopColor="#0b1220" />
          </linearGradient>
          <linearGradient id="hero-ink" x1="0" y1="0" x2="1" y2="0">
            <stop offset="0%" stopColor="#6b7b4e" />
            <stop offset="100%" stopColor="#8b7a5a" />
          </linearGradient>
        </defs>

        {/* Ornate frame */}
        <rect x="10" y="10" width="700" height="220" rx="22" stroke="#6e4f2e" strokeWidth="2" opacity="0.7" />
        <rect x="18" y="18" width="684" height="204" rx="18" stroke="#a8926b" strokeWidth="1" opacity="0.5" />
        <path d="M18 36C30 22 50 18 76 18" stroke="#6e4f2e" strokeWidth="2" strokeLinecap="round" opacity="0.6" />
        <path d="M702 36C690 22 670 18 644 18" stroke="#6e4f2e" strokeWidth="2" strokeLinecap="round" opacity="0.6" />
        <path d="M18 214C30 228 50 222 76 222" stroke="#6e4f2e" strokeWidth="2" strokeLinecap="round" opacity="0.6" />
        <path d="M702 214C690 228 670 222 644 222" stroke="#6e4f2e" strokeWidth="2" strokeLinecap="round" opacity="0.6" />

        <rect x="0" y="0" width="720" height="240" rx="20" fill="url(#hero-sky)" />

        {/* Wind strokes */}
        <path d="M40 60C110 40 190 44 260 60" stroke="#475569" strokeWidth="2" strokeLinecap="round" opacity="0.4" />
        <path d="M300 54C360 40 430 42 500 56" stroke="#475569" strokeWidth="2" strokeLinecap="round" opacity="0.35" />
        <path d="M520 62C580 48 650 50 700 62" stroke="#334155" strokeWidth="2" strokeLinecap="round" opacity="0.3" />

        {/* Mountains */}
        <path d="M0 200L120 120L240 190L320 130L420 200H0Z" fill="url(#hero-ridge)" />
        <path d="M220 210L320 150L420 210H220Z" fill="#1f2937" opacity="0.9" />
        <path d="M420 210L520 140L620 210H420Z" fill="#0f172a" opacity="0.95" />

        {/* Shrine gate */}
        <path d="M520 140H620" stroke="url(#hero-ink)" strokeWidth="6" strokeLinecap="round" />
        <path d="M540 140V200" stroke="url(#hero-ink)" strokeWidth="6" strokeLinecap="round" />
        <path d="M600 140V200" stroke="url(#hero-ink)" strokeWidth="6" strokeLinecap="round" />

        {/* Foreground grass */}
        <path d="M0 220C80 205 140 210 220 220C320 235 420 235 520 220C600 210 660 210 720 220V240H0Z" fill="#0b1220" />
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
            <stop offset="0%" stopColor="#8b7a5a" />
            <stop offset="100%" stopColor="#6b7b4e" />
          </linearGradient>
        </defs>

        <circle cx="120" cy="90" r="52" fill="#0f172a" stroke="url(#sakura-success-gradient)" strokeWidth="6" />
        <path d="M92 92L112 112L148 74" stroke="url(#sakura-success-gradient)" strokeWidth="7" strokeLinecap="round" strokeLinejoin="round" />
        <circle cx="120" cy="155" r="28" fill="#6b7b4e" opacity="0.14" />
      </svg>
    </div>
  );
};
