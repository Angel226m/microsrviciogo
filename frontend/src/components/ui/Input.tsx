import clsx from 'clsx';
import { forwardRef, type InputHTMLAttributes } from 'react';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  icon?: React.ReactNode;
}

const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, icon, className, id, ...props }, ref) => {
    const inputId = id ?? label?.toLowerCase().replace(/\s+/g, '-');
    return (
      <div className="space-y-1.5">
        {label && (
          <label htmlFor={inputId} className="block text-sm font-medium text-gray-700">
            {label}
          </label>
        )}
        <div className="relative">
          {icon && (
            <div className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400">
              {icon}
            </div>
          )}
          <input
            ref={ref}
            id={inputId}
            className={clsx(
              'w-full rounded-xl border bg-white px-4 py-3 text-sm text-gray-900 placeholder:text-gray-400 outline-none transition-all duration-200',
              'focus:ring-2 focus:ring-primary-500/30 focus:border-primary-500',
              error ? 'border-red-300 focus:ring-red-500/30 focus:border-red-500' : 'border-gray-200',
              icon && 'pl-11',
              className,
            )}
            {...props}
          />
        </div>
        {error && <p className="text-xs text-red-500 font-medium">{error}</p>}
      </div>
    );
  },
);

Input.displayName = 'Input';
export default Input;
