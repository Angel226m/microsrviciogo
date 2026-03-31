/**
 * ═══════════════════════════════════════════════════════════════
 * Página de Perfil – Información del usuario y opciones de cuenta
 * ═══════════════════════════════════════════════════════════════
 */
import { useAuthStore } from '../store/authStore';
import { Navigate, Link } from 'react-router-dom';
import { Mail, Calendar, ShoppingBag, Heart, Settings, ChevronRight, Shield, LogOut } from 'lucide-react';
import { useQuery } from '@tanstack/react-query';
import { orderService } from '../services/orderService';
import Button from '../components/ui/Button';

export default function ProfilePage() {
  const user = useAuthStore((s) => s.user);
  const logout = useAuthStore((s) => s.logout);
  if (!user) return <Navigate to="/login" />;

  const { data: ordersData } = useQuery({
    queryKey: ['orders'],
    queryFn: () => orderService.list({ limit: 50 }),
  });

  const orderCount = ordersData?.data?.length ?? 0;
  const initial = (user.first_name?.[0] ?? user.email[0]).toUpperCase();

  const menuItems = [
    { icon: ShoppingBag, label: 'Mis Pedidos', desc: `${orderCount} pedidos`, to: '/orders' },
    { icon: Heart, label: 'Lista de Deseos', desc: 'Artículos guardados', to: '#' },
    { icon: Settings, label: 'Configuración', desc: 'Preferencias de cuenta', to: '#' },
    { icon: Shield, label: 'Seguridad', desc: 'Contraseña y privacidad', to: '#' },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50/50 to-white">
      <div className="max-w-2xl mx-auto px-4 py-8">
        {/* Profile card */}
        <div className="bg-white rounded-3xl border border-gray-100 shadow-sm overflow-hidden animate-fade-in-up">
          {/* Gradient header */}
          <div className="h-32 bg-gradient-to-r from-primary-500 via-primary-600 to-accent-600 relative">
            <div className="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAiIGhlaWdodD0iNjAiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGRlZnM+PHBhdHRlcm4gaWQ9ImciIHdpZHRoPSI2MCIgaGVpZ2h0PSI2MCIgcGF0dGVyblVuaXRzPSJ1c2VyU3BhY2VPblVzZSI+PGNpcmNsZSBjeD0iMzAiIGN5PSIzMCIgcj0iMSIgZmlsbD0icmdiYSgyNTUsMjU1LDI1NSwwLjEpIi8+PC9wYXR0ZXJuPjwvZGVmcz48cmVjdCBmaWxsPSJ1cmwoI2cpIiB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIi8+PC9zdmc+')] opacity-60" />
          </div>

          {/* Avatar & info */}
          <div className="px-6 sm:px-8 pb-8 -mt-14 relative">
            <div className="w-28 h-28 rounded-2xl bg-gradient-to-br from-primary-400 to-accent-500 flex items-center justify-center text-white text-4xl font-bold border-4 border-white shadow-lg">
              {initial}
            </div>

            <div className="mt-5">
              <h1 className="text-2xl font-bold text-gray-900">{user.first_name} {user.last_name}</h1>
              <span className="inline-flex items-center gap-1.5 mt-1.5 px-3 py-1 bg-primary-50 text-primary-700 rounded-full text-xs font-semibold uppercase tracking-wider">
                <Shield className="w-3 h-3" />
                {user.role}
              </span>
            </div>

            <div className="mt-6 grid grid-cols-2 gap-4">
              <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-xl">
                <Mail className="w-5 h-5 text-gray-400" />
                <div className="min-w-0">
                  <p className="text-xs text-gray-400 font-medium">Correo</p>
                  <p className="text-sm text-gray-900 truncate">{user.email}</p>
                </div>
              </div>
              <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-xl">
                <Calendar className="w-5 h-5 text-gray-400" />
                <div>
                  <p className="text-xs text-gray-400 font-medium">Miembro desde</p>
                  <p className="text-sm text-gray-900">
                    {new Date(user.created_at).toLocaleDateString('es-MX', { month: 'short', year: 'numeric' })}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Stats cards */}
        <div className="grid grid-cols-3 gap-3 mt-6 animate-fade-in-up" style={{ animationDelay: '100ms' }}>
          {[
            { label: 'Pedidos', value: orderCount, color: 'from-primary-500 to-primary-600' },
            { label: 'Deseos', value: 0, color: 'from-pink-500 to-rose-500' },
            { label: 'Reseñas', value: 0, color: 'from-amber-500 to-orange-500' },
          ].map((stat) => (
            <div key={stat.label} className="bg-white rounded-2xl border border-gray-100 p-4 text-center">
              <p className={`text-2xl font-bold bg-gradient-to-r ${stat.color} bg-clip-text text-transparent`}>
                {stat.value}
              </p>
              <p className="text-xs text-gray-500 font-medium mt-0.5">{stat.label}</p>
            </div>
          ))}
        </div>

        {/* Menu items */}
        <div className="mt-6 bg-white rounded-2xl border border-gray-100 overflow-hidden animate-fade-in-up" style={{ animationDelay: '200ms' }}>
          {menuItems.map((item, i) => (
            <Link
              key={item.label}
              to={item.to}
              className={`flex items-center gap-4 px-6 py-4 hover:bg-gray-50 transition-colors ${
                i !== menuItems.length - 1 ? 'border-b border-gray-50' : ''
              }`}
            >
              <div className="w-10 h-10 rounded-xl bg-gray-100 flex items-center justify-center">
                <item.icon className="w-5 h-5 text-gray-600" />
              </div>
              <div className="flex-1">
                <p className="font-medium text-gray-900 text-sm">{item.label}</p>
                <p className="text-xs text-gray-400">{item.desc}</p>
              </div>
              <ChevronRight className="w-4 h-4 text-gray-300" />
            </Link>
          ))}
        </div>

        {/* Logout */}
        <div className="mt-6 animate-fade-in-up" style={{ animationDelay: '300ms' }}>
          <Button
            variant="outline"
            className="w-full"
            icon={<LogOut className="w-4 h-4" />}
            onClick={logout}
          >
            Cerrar Sesión
          </Button>
        </div>
      </div>
    </div>
  );
}
