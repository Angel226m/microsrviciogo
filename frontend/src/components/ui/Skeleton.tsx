import clsx from 'clsx';

interface SkeletonProps {
  className?: string;
  rounded?: 'md' | 'lg' | 'xl' | '2xl' | 'full';
}

export default function Skeleton({ className, rounded = 'xl' }: SkeletonProps) {
  return (
    <div
      className={clsx('skeleton', `rounded-${rounded}`, className)}
      aria-hidden="true"
    />
  );
}

export function ProductCardSkeleton() {
  return (
    <div className="bg-white rounded-2xl overflow-hidden border border-gray-100">
      <Skeleton className="aspect-square w-full" rounded="md" />
      <div className="p-5 space-y-3">
        <Skeleton className="h-5 w-3/4" />
        <Skeleton className="h-4 w-1/2" />
        <div className="flex justify-between items-center pt-2">
          <Skeleton className="h-6 w-20" />
          <Skeleton className="h-10 w-10 rounded-xl" />
        </div>
      </div>
    </div>
  );
}

export function OrderSkeleton() {
  return (
    <div className="bg-white rounded-2xl p-6 border border-gray-100 space-y-4">
      <div className="flex justify-between items-center">
        <Skeleton className="h-5 w-32" />
        <Skeleton className="h-6 w-24 rounded-full" />
      </div>
      <Skeleton className="h-4 w-48" />
      <div className="flex justify-between items-center">
        <Skeleton className="h-7 w-28" />
        <Skeleton className="h-4 w-20" />
      </div>
    </div>
  );
}
