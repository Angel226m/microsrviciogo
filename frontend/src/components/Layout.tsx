/**
 * ═══════════════════════════════════════════════════════════════
 * Layout Principal – Estructura base de la aplicación
 * Incluye: barra de navegación, menú móvil y pie de página
 * Implementa diseño responsivo profesional
 * ═══════════════════════════════════════════════════════════════
 */
import { Outlet, Link, NavLink, useLocation } from 'react-router-dom';
import { ShoppingCart, Package, LogOut, Menu, X, Zap, Shield, Truck, Heart, Award, Headphones } from 'lucide-react';
import { useAuthStore } from '../store/authStore';
import { useCartStore, type CartState } from '../store/cartStore';
import { useState, useEffect } from 'react';
import clsx from 'clsx';
import Button from './ui/Button';

export default function Layout() {
  const { user, logout } = useAuthStore();
  const cartCount = useCartStore((s: CartState) => s.count());
  const [mobileOpen, setMobileOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);
  const location = useLocation();

  // Cerrar menú móvil al cambiar de ruta
  useEffect(() => { setMobileOpen(false); }, [location.pathname]);

  // Detectar scroll para efecto glass en el header
  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 10);
    window.addEventListener('scroll', onScroll, { passive: true });
    return () => window.removeEventListener('scroll', onScroll);
  }, []);

  const navLinkClass = ({ isActive }: { isActive: boolean }) =>
    clsx(
      'relative py-1 text-sm font-medium transition-colors duration-200',
      isActive ? 'text-primary-600' : 'text-gray-600 hover:text-primary-600',
    );

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-white">
      {/* Banner promocional */}
      <div className="bg-gradient-to-r from-gray-900 via-gray-800 to-gray-900 text-gray-300 text-center py-2.5 text-xs font-medium tracking-wide">
        <Zap className="w-3.5 h-3.5 inline mr-1 -mt-0.5 text-primary-400" />
        Envío gratis en pedidos mayores a $50 &mdash; Plataforma cloud-native con microservicios
      </div>

      {/* Header */}
      <header
        className={clsx(
          'sticky top-0 z-50 transition-all duration-300',
          scrolled ? 'glass shadow-lg shadow-gray-200/50 border-b border-white/50' : 'bg-white/95 border-b border-gray-100',
        )}
      >
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            {/* Logo */}
            <Link to="/" className="flex items-center gap-2.5 group">
              <div className="w-9 h-9 bg-gradient-to-br from-primary-500 to-accent-600 rounded-xl flex items-center justify-center shadow-lg shadow-primary-500/25 group-hover:shadow-primary-500/40 transition-shadow">
                <Package className="w-5 h-5 text-white" />
              </div>
              <span className="text-xl font-extrabold bg-gradient-to-r from-primary-600 to-accent-600 bg-clip-text text-transparent">
                CloudMart
              </span>
            </Link>

            {/* Desktop nav */}
            <nav className="hidden md:flex items-center gap-8">
              <NavLink to="/products" className={navLinkClass}>Productos</NavLink>
              {user && <NavLink to="/orders" className={navLinkClass}>Mis Pedidos</NavLink>}
            </nav>

            {/* Actions */}
            <div className="flex items-center gap-3">
              {/* Cart */}
              <Link
                to="/cart"
                className="relative p-2.5 rounded-xl text-gray-500 hover:text-primary-600 hover:bg-primary-50 transition-all duration-200"
              >
                <ShoppingCart className="w-5 h-5" />
                {cartCount > 0 && (
                  <span className="absolute -top-0.5 -right-0.5 bg-gradient-to-r from-primary-500 to-accent-500 text-white text-[10px] font-bold rounded-full min-w-[20px] h-5 flex items-center justify-center px-1 shadow-lg shadow-primary-500/30 animate-bounce">
                    {cartCount}
                  </span>
                )}
              </Link>

              {user ? (
                <div className="hidden sm:flex items-center gap-2">
                  <Link
                    to="/profile"
                    className="flex items-center gap-2 px-3 py-2 rounded-xl text-gray-600 hover:text-primary-600 hover:bg-primary-50 transition-all duration-200"
                  >
                    <div className="w-7 h-7 bg-gradient-to-br from-primary-400 to-accent-400 rounded-lg flex items-center justify-center">
                      <span className="text-white text-xs font-bold">{user.first_name[0]}</span>
                    </div>
                    <span className="text-sm font-medium">{user.first_name}</span>
                  </Link>
                  <button
                    onClick={logout}
                    className="p-2 rounded-xl text-gray-400 hover:text-red-500 hover:bg-red-50 transition-all duration-200"
                    title="Cerrar sesión"
                  >
                    <LogOut className="w-4.5 h-4.5" />
                  </button>
                </div>
              ) : (
                <Link to="/login" className="hidden sm:block">
                  <Button size="sm">Iniciar Sesión</Button>
                </Link>
              )}

              {/* Mobile menu toggle */}
              <button
                onClick={() => setMobileOpen(!mobileOpen)}
                className="md:hidden p-2 rounded-xl text-gray-600 hover:bg-gray-100 transition-colors"
              >
                {mobileOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
              </button>
            </div>
          </div>
        </div>

        {/* Mobile menu */}
        <div
          className={clsx(
            'md:hidden overflow-hidden transition-all duration-300 ease-in-out',
            mobileOpen ? 'max-h-80 border-t border-gray-100' : 'max-h-0',
          )}
        >
          <div className="px-4 py-4 space-y-1 bg-white">
            <MobileLink to="/products">Productos</MobileLink>
            {user && <MobileLink to="/orders">Mis Pedidos</MobileLink>}
            {user ? (
              <>
                <MobileLink to="/profile">Mi Perfil</MobileLink>
                <button onClick={logout} className="w-full text-left px-4 py-3 rounded-xl text-red-500 hover:bg-red-50 font-medium text-sm">
                  Cerrar Sesión
                </button>
              </>
            ) : (
              <>
                <MobileLink to="/login">Iniciar Sesión</MobileLink>
                <MobileLink to="/register">Crear Cuenta</MobileLink>
              </>
            )}
          </div>
        </div>
      </header>

      {/* Main content */}
      <main className="min-h-[60vh]">
        <Outlet />
      </main>

      {/* Footer */}
      <footer className="bg-gray-900 text-gray-400 mt-20">
        <div className="max-w-7xl mx-auto px-4 py-16">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-10">
            {/* Brand */}
            <div className="md:col-span-1">
              <Link to="/" className="flex items-center gap-2.5">
                <div className="w-9 h-9 bg-gradient-to-br from-primary-500 to-accent-500 rounded-xl flex items-center justify-center">
                  <Package className="w-5 h-5 text-white" />
                </div>
                <span className="text-xl font-extrabold text-white">CloudMart</span>
              </Link>
              <p className="mt-4 text-sm leading-relaxed">
                Plataforma de e-commerce cloud-native premium. Construida con microservicios Go, React y Kubernetes.
              </p>
            </div>

            {/* Links */}
            <div>
              <h4 className="text-white font-semibold text-sm mb-4">Tienda</h4>
              <ul className="space-y-2.5 text-sm">
                <li><Link to="/products" className="hover:text-white transition-colors">Todos los Productos</Link></li>
                <li><Link to="/products?category=new" className="hover:text-white transition-colors">Novedades</Link></li>
                <li><Link to="/products?category=sale" className="hover:text-white transition-colors">En Oferta</Link></li>
              </ul>
            </div>

            <div>
              <h4 className="text-white font-semibold text-sm mb-4">Mi Cuenta</h4>
              <ul className="space-y-2.5 text-sm">
                <li><Link to="/orders" className="hover:text-white transition-colors">Mis Pedidos</Link></li>
                <li><Link to="/profile" className="hover:text-white transition-colors">Perfil</Link></li>
                <li><Link to="/cart" className="hover:text-white transition-colors">Carrito</Link></li>
              </ul>
            </div>

            {/* Features */}
            <div>
              <h4 className="text-white font-semibold text-sm mb-4">¿Por qué CloudMart?</h4>
              <ul className="space-y-3 text-sm">
                <li className="flex items-center gap-2"><Truck className="w-4 h-4 text-primary-400" /> Envío Gratuito</li>
                <li className="flex items-center gap-2"><Shield className="w-4 h-4 text-primary-400" /> Pagos Seguros</li>
                <li className="flex items-center gap-2"><Award className="w-4 h-4 text-primary-400" /> Calidad Premium</li>
                <li className="flex items-center gap-2"><Headphones className="w-4 h-4 text-primary-400" /> Soporte 24/7</li>
              </ul>
            </div>
          </div>

          <div className="mt-12 pt-8 border-t border-gray-800 flex flex-col sm:flex-row justify-between items-center gap-4">
            <p className="text-xs">&copy; {new Date().getFullYear()} CloudMart. Todos los derechos reservados.</p>
            <p className="text-xs flex items-center gap-1.5">
              Hecho con <Heart className="w-3 h-3 text-red-500 fill-red-500" /> usando Go + React + Kubernetes
            </p>
          </div>
        </div>
      </footer>
    </div>
  );
}

function MobileLink({ to, children }: { to: string; children: React.ReactNode }) {
  return (
    <NavLink
      to={to}
      className={({ isActive }) =>
        clsx(
          'block px-4 py-3 rounded-xl font-medium text-sm transition-colors',
          isActive ? 'bg-primary-50 text-primary-700' : 'text-gray-700 hover:bg-gray-50',
        )
      }
    >
      {children}
    </NavLink>
  );
}
