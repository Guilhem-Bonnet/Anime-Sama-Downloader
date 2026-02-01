import React from 'react';

interface BadgeProps extends React.HTMLAttributes<HTMLSpanElement> {
  variant?: 'primary' | 'secondary' | 'success' | 'warning' | 'error' | 'info';
  children: React.ReactNode;
}

export const Badge = React.forwardRef<HTMLSpanElement, BadgeProps>(
  ({ variant = 'primary', children, className = '', ...props }, ref) => {
    const variantClass = `badge-${variant}`;
    const finalClass = `badge ${variantClass} ${className}`.trim();

    return (
      <span ref={ref} className={finalClass} {...props}>
        {children}
      </span>
    );
  }
);

Badge.displayName = 'Badge';
