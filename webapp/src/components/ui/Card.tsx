import React from 'react';

interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
  hoverable?: boolean;
}

export const Card = React.forwardRef<HTMLDivElement, CardProps>(
  ({ children, hoverable = true, className = '', ...props }, ref) => {
    const finalClass = `card ${className}`.trim();
    return (
      <div ref={ref} className={finalClass} {...props}>
        {children}
      </div>
    );
  }
);

Card.displayName = 'Card';

interface CardHeaderProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
}

export const CardHeader = React.forwardRef<HTMLDivElement, CardHeaderProps>(
  ({ children, className = '', ...props }, ref) => (
    <div ref={ref} className={`card-header ${className}`.trim()} {...props}>
      {children}
    </div>
  )
);

CardHeader.displayName = 'CardHeader';

interface CardTitleProps extends React.HTMLAttributes<HTMLHeadingElement> {
  children: React.ReactNode;
  level?: 'h1' | 'h2' | 'h3';
}

export const CardTitle = React.forwardRef<
  HTMLHeadingElement,
  CardTitleProps
>(({ children, level = 'h2', className = '', ...props }, ref) => {
  const Component = level;
  return React.createElement(Component, {
    ref,
    className: `card-title ${className}`.trim(),
    children,
    ...props,
  });
});

CardTitle.displayName = 'CardTitle';

interface CardSubtitleProps extends React.HTMLAttributes<HTMLParagraphElement> {
  children: React.ReactNode;
}

export const CardSubtitle = React.forwardRef<
  HTMLParagraphElement,
  CardSubtitleProps
>(({ children, className = '', ...props }, ref) => (
  <p ref={ref} className={`card-subtitle ${className}`.trim()} {...props}>
    {children}
  </p>
));

CardSubtitle.displayName = 'CardSubtitle';

interface CardBodyProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
}

export const CardBody = React.forwardRef<HTMLDivElement, CardBodyProps>(
  ({ children, className = '', ...props }, ref) => (
    <div ref={ref} className={`card-body ${className}`.trim()} {...props}>
      {children}
    </div>
  )
);

CardBody.displayName = 'CardBody';

interface CardFooterProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
}

export const CardFooter = React.forwardRef<HTMLDivElement, CardFooterProps>(
  ({ children, className = '', ...props }, ref) => (
    <div ref={ref} className={`card-footer ${className}`.trim()} {...props}>
      {children}
    </div>
  )
);

CardFooter.displayName = 'CardFooter';
