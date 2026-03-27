import { useQuery, useQueryClient } from '@tanstack/react-query';
import { orderService } from '../services/orderService';
import { useAuthStore, type AuthState } from '../store/authStore';
import { Navigate } from 'react-router-dom';
import { Package, Clock, ShoppingBag, ChevronRight, XCircle, CheckCircle2, Truck, Timer, AlertCircle } from 'lucide-react';
import type { Order } from '../types';
import Badge from '../components/ui/Badge';
import { OrderSkeleton } from '../components/ui/Skeleton';
import EmptyState from '../components/ui/EmptyState';
import clsx from 'clsx';
import toast from 'react-hot-toast';

const statusConfig: Record<string, { color: 'default' | 'success' | 'warning' | 'danger' | 'info' | 'accent'; icon: React.ElementType }> = {
  pending: { color: 'warning', icon: Timer },
  confirmed: { color: 'info', icon: CheckCircle2 },
  processing: { color: 'accent', icon: Package },
  shipped: { color: 'info', icon: Truck },
  delivered: { color: 'success', icon: CheckCircle2 },
  cancelled: { color: 'danger', icon: XCircle },
};

export default function OrdersPage() {
  const user = useAuthStore((s: AuthState) => s.user);
  const queryClient = useQueryClient();
  if (!user) return <Navigate to="/login" />;

  const { data, isLoading } = useQuery({
    queryKey: ['orders'],
    queryFn: () => orderService.list({ limit: 50 }),
  });

  const handleCancel = async (orderId: string) => {
    try {
      await orderService.cancel(orderId);
      toast.success('Order cancelled');
      queryClient.invalidateQueries({ queryKey: ['orders'] });
    } catch {
      toast.error('Failed to cancel order');
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50/50 to-white">
      <div className="max-w-4xl mx-auto px-4 py-8">
        {/* Header */}
        <div className="flex items-center gap-3 mb-8 animate-fade-in-up">
          <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-primary-500 to-accent-500 flex items-center justify-center">
            <Package className="w-6 h-6 text-white" />
          </div>
          <div>
            <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">My Orders</h1>
            {data?.data && (
              <p className="text-sm text-gray-500 mt-0.5">{data.data.length} order{data.data.length !== 1 ? 's' : ''}</p>
            )}
          </div>
        </div>

        {isLoading ? (
          <div className="space-y-4">
            {[...Array(3)].map((_, i) => <OrderSkeleton key={i} />)}
          </div>
        ) : !data?.data?.length ? (
          <div className="py-12">
            <EmptyState
              icon={ShoppingBag}
              title="No orders yet"
              description="When you place orders, they'll appear here so you can track them."
              actionLabel="Start Shopping"
              actionTo="/products"
            />
          </div>
        ) : (
          <div className="space-y-4">
            {data.data.map((order: Order, index: number) => {
              const config = statusConfig[order.status] ?? { color: 'default' as const, icon: AlertCircle };
              const StatusIcon = config.icon;
              return (
                <div
                  key={order.id}
                  className="bg-white rounded-2xl border border-gray-100 hover:shadow-lg hover:shadow-gray-100/80 transition-all overflow-hidden animate-fade-in-up group"
                  style={{ animationDelay: `${index * 60}ms` }}
                >
                  <div className="p-5 sm:p-6">
                    <div className="flex items-start sm:items-center justify-between gap-3 flex-col sm:flex-row">
                      <div className="flex items-center gap-3">
                        <div className={clsx(
                          'w-10 h-10 rounded-xl flex items-center justify-center shrink-0',
                          config.color === 'warning' && 'bg-amber-50 text-amber-600',
                          config.color === 'info' && 'bg-blue-50 text-blue-600',
                          config.color === 'accent' && 'bg-purple-50 text-purple-600',
                          config.color === 'success' && 'bg-green-50 text-green-600',
                          config.color === 'danger' && 'bg-red-50 text-red-600',
                          config.color === 'default' && 'bg-gray-50 text-gray-600',
                        )}>
                          <StatusIcon className="w-5 h-5" />
                        </div>
                        <div>
                          <div className="flex items-center gap-2.5 flex-wrap">
                            <span className="font-mono text-sm font-semibold text-gray-900">{order.order_number}</span>
                            <Badge variant={config.color}>{order.status}</Badge>
                          </div>
                          <div className="flex items-center gap-1.5 text-xs text-gray-400 mt-1">
                            <Clock className="w-3.5 h-3.5" />
                            {new Date(order.created_at).toLocaleDateString('en-US', {
                              year: 'numeric', month: 'short', day: 'numeric'
                            })}
                          </div>
                        </div>
                      </div>

                      <div className="flex items-center gap-4 w-full sm:w-auto justify-between sm:justify-end">
                        <div className="text-right">
                          <p className="text-xs text-gray-400">{order.items?.length ?? 0} items</p>
                          <p className="text-lg font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
                            ${order.total?.toFixed(2)}
                          </p>
                        </div>
                        <ChevronRight className="w-5 h-5 text-gray-300 group-hover:text-primary-500 group-hover:translate-x-0.5 transition-all" />
                      </div>
                    </div>

                    {order.status === 'pending' && (
                      <div className="mt-4 pt-4 border-t border-gray-50 flex items-center justify-between">
                        <span className="text-xs text-amber-600 font-medium">Awaiting confirmation</span>
                        <button
                          onClick={() => handleCancel(order.id)}
                          className="flex items-center gap-1.5 text-sm text-red-500 hover:text-red-600 font-medium px-3 py-1.5 hover:bg-red-50 rounded-lg transition-all"
                        >
                          <XCircle className="w-4 h-4" />
                          Cancel Order
                        </button>
                      </div>
                    )}
                  </div>

                  {/* Progress bar for active orders */}
                  {['confirmed', 'processing', 'shipped'].includes(order.status) && (
                    <div className="h-1 bg-gray-100">
                      <div
                        className={clsx(
                          'h-full bg-gradient-to-r from-primary-500 to-accent-500 transition-all duration-500',
                          order.status === 'confirmed' && 'w-1/3',
                          order.status === 'processing' && 'w-2/3',
                          order.status === 'shipped' && 'w-full',
                        )}
                      />
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
