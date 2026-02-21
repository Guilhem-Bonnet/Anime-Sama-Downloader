import React from 'react';

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  hint?: string;
  error?: string;
}

export const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ label, hint, error, className = '', ...props }, ref) => {
    return (
      <div className="flex flex-col gap-2">
        {label && <label className="input-label">{label}</label>}
        <input
          ref={ref}
          className={`input ${className}`.trim()}
          {...props}
        />
        {hint && !error && <p className="input-hint">{hint}</p>}
        {error && <p className="input-hint" style={{ color: 'var(--night-error-text)' }}>{error}</p>}
      </div>
    );
  }
);

Input.displayName = 'Input';

interface TextAreaProps
  extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string;
  hint?: string;
  error?: string;
}

export const TextArea = React.forwardRef<HTMLTextAreaElement, TextAreaProps>(
  ({ label, hint, error, className = '', ...props }, ref) => {
    return (
      <div className="flex flex-col gap-2">
        {label && <label className="input-label">{label}</label>}
        <textarea
          ref={ref}
          className={`input ${className}`.trim()}
          style={{ minHeight: '120px', resize: 'vertical' }}
          {...props}
        />
        {hint && !error && <p className="input-hint">{hint}</p>}
        {error && <p className="input-hint" style={{ color: 'var(--night-error-text)' }}>{error}</p>}
      </div>
    );
  }
);

TextArea.displayName = 'TextArea';

interface SelectProps extends React.SelectHTMLAttributes<HTMLSelectElement> {
  label?: string;
  hint?: string;
  error?: string;
  options: Array<{ value: string; label: string }>;
}

export const Select = React.forwardRef<HTMLSelectElement, SelectProps>(
  ({ label, hint, error, options, className = '', ...props }, ref) => {
    return (
      <div className="flex flex-col gap-2">
        {label && <label className="input-label">{label}</label>}
        <select
          ref={ref}
          className={`input ${className}`.trim()}
          {...props}
        >
          {options.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </select>
        {hint && !error && <p className="input-hint">{hint}</p>}
        {error && <p className="input-hint" style={{ color: 'var(--night-error-text)' }}>{error}</p>}
      </div>
    );
  }
);

Select.displayName = 'Select';
