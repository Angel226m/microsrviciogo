import { Link } from 'react-router-dom';
import Button from './Button';
import type { LucideIcon } from 'lucide-react';

interface EmptyStateProps {
  icon: LucideIcon | React.ReactNode;
  title: string;
  description: string;
  actionLabel?: string;
  actionTo?: string;
  actionHref?: string;
}

export default function EmptyState({ icon: IconOrNode, title, description, actionLabel, actionTo, actionHref }: EmptyStateProps) {
  const href = actionTo ?? actionHref;
  const isComponent = typeof IconOrNode === 'function';
  return (
    <div className="flex flex-col items-center justify-center py-20 px-4 animate-fade-in-up">
      <div className="w-20 h-20 bg-gradient-to-br from-gray-100 to-gray-50 rounded-2xl flex items-center justify-center text-gray-300 mb-6">
        {isComponent ? <IconOrNode className="w-10 h-10" /> : IconOrNode}
      </div>
      <h2 className="text-xl font-bold text-gray-900">{title}</h2>
      <p className="mt-2 text-gray-500 text-center max-w-sm">{description}</p>
      {actionLabel && href && (
        <Link to={href} className="mt-6">
          <Button size="lg">{actionLabel}</Button>
        </Link>
      )}
    </div>
  );
}
