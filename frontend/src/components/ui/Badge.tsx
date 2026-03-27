import clsx from 'clsx';

type Variant = 'default' | 'success' | 'warning' | 'danger' | 'info' | 'accent';

const variants: Record<Variant, string> = {
  default: 'bg-gray-100 text-gray-700',
  success: 'bg-emerald-50 text-emerald-700',
  warning: 'bg-amber-50 text-amber-700',
  danger: 'bg-red-50 text-red-700',
  info: 'bg-blue-50 text-blue-700',
  accent: 'bg-primary-50 text-primary-700',
};

interface BadgeProps {
  variant?: Variant;
  children: React.ReactNode;
  className?: string;
}

export default function Badge({ variant = 'default', children, className }: BadgeProps) {
  return (
    <span
      className={clsx(
        'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold tracking-wide',
        variants[variant],
        className,
      )}
    >
      {children}
    </span>
  );
}
