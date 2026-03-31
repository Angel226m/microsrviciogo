/**
 * ═══════════════════════════════════════════════════════════════
 * Página de Inicio de Sesión – Autenticación de usuarios
 * Incluye panel lateral profesional y formulario de acceso
 * ═══════════════════════════════════════════════════════════════
 */
import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { authService } from '../services/authService';
import { useAuthStore, type AuthState } from '../store/authStore';
import toast from 'react-hot-toast';
import {
  Mail, Lock, ArrowRight, ShoppingBag, Sparkles, Shield, Zap,
  Eye, EyeOff, Users, Globe,
} from 'lucide-react';
import Button from '../components/ui/Button';
import Input from '../components/ui/Input';

export default function LoginPage() {
  const navigate = useNavigate();
  const setAuth = useAuthStore((s: AuthState) => s.setAuth);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      const { user, tokens } = await authService.login(email, password);
      setAuth(user, tokens.access_token);
      toast.success('¡Bienvenido de vuelta!');
      navigate('/');
    } catch {
      toast.error('Credenciales inválidas');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex bg-gray-50">
      {/* Left side — Dark professional branding */}
      <div className="hidden lg:flex lg:w-[55%] bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 relative overflow-hidden">
        <div className="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAiIGhlaWdodD0iNjAiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+PGRlZnM+PHBhdHRlcm4gaWQ9ImciIHdpZHRoPSI2MCIgaGVpZ2h0PSI2MCIgcGF0dGVyblVuaXRzPSJ1c2VyU3BhY2VPblVzZSI+PGNpcmNsZSBjeD0iMzAiIGN5PSIzMCIgcj0iMSIgZmlsbD0icmdiYSgyNTUsMjU1LDI1NSwwLjAzKSIvPjwvcGF0dGVybj48L2RlZnM+PHJlY3QgZmlsbD0idXJsKCNnKSIgd2lkdGg9IjEwMCUiIGhlaWdodD0iMTAwJSIvPjwvc3ZnPg==')] opacity-100" />
        <div className="absolute top-20 left-10 w-80 h-80 bg-primary-500/10 rounded-full blur-[100px]" />
        <div className="absolute bottom-20 right-10 w-96 h-96 bg-accent-500/8 rounded-full blur-[100px]" />

        <div className="relative z-10 flex flex-col justify-between p-12 xl:p-16 w-full">
          <Link to="/" className="flex items-center gap-3 group">
            <div className="w-11 h-11 bg-primary-500 rounded-xl flex items-center justify-center group-hover:scale-105 transition-transform">
              <ShoppingBag className="w-5 h-5 text-white" />
            </div>
            <span className="text-xl font-bold text-white">CloudMart</span>
          </Link>

          <div className="max-w-lg">
            <div className="inline-flex items-center gap-2 bg-primary-500/15 border border-primary-500/20 rounded-full px-4 py-1.5 text-sm font-medium text-primary-300 mb-6">
              <Sparkles className="w-3.5 h-3.5" />
              Más de 100K clientes confían en nosotros
            </div>
            <h2 className="text-4xl xl:text-5xl font-bold text-white leading-tight tracking-tight">
              Bienvenido de nuevo a
              <span className="block bg-gradient-to-r from-primary-400 to-accent-400 bg-clip-text text-transparent">
                CloudMart
              </span>
            </h2>
            <p className="mt-5 text-gray-400 text-lg leading-relaxed">
              Accede a tu panel personalizado, rastrea pedidos en tiempo real y descubre ofertas exclusivas para ti.
            </p>

            <div className="mt-10 grid grid-cols-1 gap-4">
              {[
                { icon: Shield, title: 'Seguro y Privado', desc: 'Cifrado de extremo a extremo en todas las transacciones' },
                { icon: Zap, title: 'Pago Relámpago', desc: 'Compra en un clic con preferencias guardadas' },
              ].map(({ icon: Icon, title, desc }) => (
                <div key={title} className="flex items-start gap-4 p-4 rounded-xl bg-white/5 border border-white/5">
                  <div className="w-10 h-10 rounded-lg bg-primary-500/15 flex items-center justify-center shrink-0">
                    <Icon className="w-5 h-5 text-primary-400" />
                  </div>
                  <div>
                    <p className="text-sm font-semibold text-white">{title}</p>
                    <p className="text-sm text-gray-500 mt-0.5">{desc}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="flex items-center gap-8 text-sm text-gray-500">
            <span className="flex items-center gap-2"><Users className="w-4 h-4" /> 100K+ usuarios</span>
            <span className="flex items-center gap-2"><Globe className="w-4 h-4" /> 50+ países</span>
            <span className="flex items-center gap-2"><Shield className="w-4 h-4" /> Certificado SOC 2</span>
          </div>
        </div>
      </div>

      {/* Right side — Form */}
      <div className="flex-1 flex items-center justify-center px-4 sm:px-8 py-12">
        <div className="w-full max-w-[420px] animate-fade-in-up">
          <div className="text-center mb-8 lg:hidden">
            <Link to="/" className="inline-flex items-center gap-2.5">
              <div className="w-11 h-11 bg-gradient-to-br from-primary-500 to-primary-600 rounded-xl flex items-center justify-center shadow-lg shadow-primary-500/20">
                <ShoppingBag className="w-6 h-6 text-white" />
              </div>
              <span className="text-2xl font-bold text-gray-900">CloudMart</span>
            </Link>
          </div>

          <div>
            <div className="mb-8">
              <h2 className="text-2xl sm:text-3xl font-bold text-gray-900 tracking-tight">Iniciar Sesión</h2>
              <p className="text-gray-500 mt-2">Ingresa tus credenciales para acceder a tu cuenta</p>
            </div>

            <div className="mb-6 p-4 bg-primary-50 border border-primary-100 rounded-xl">
              <p className="text-xs font-semibold text-primary-700 uppercase tracking-wider mb-1">Credenciales de demostración</p>
              <p className="text-sm text-primary-600">
                <span className="font-mono bg-primary-100 px-1.5 py-0.5 rounded text-xs">demo@cloudmart.com</span> / <span className="font-mono bg-primary-100 px-1.5 py-0.5 rounded text-xs">demo123</span>
              </p>
            </div>

            <form onSubmit={handleSubmit} className="space-y-5">
              <Input
                label="Correo electrónico"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                placeholder="tu@ejemplo.com"
                icon={<Mail className="w-5 h-5" />}
              />
              <div className="relative">
                <Input
                  label="Contraseña"
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                  placeholder="Ingresa tu contraseña"
                  icon={<Lock className="w-5 h-5" />}
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3.5 top-[38px] text-gray-400 hover:text-gray-600 transition-colors"
                >
                  {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                </button>
              </div>

              <div className="flex items-center justify-between">
                <label className="flex items-center gap-2 cursor-pointer select-none">
                  <input type="checkbox" className="w-4 h-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500" />
                  <span className="text-sm text-gray-600">Recordarme</span>
                </label>
                <button type="button" className="text-sm text-primary-600 font-medium hover:text-primary-700 transition-colors">
                  ¿Olvidaste tu contraseña?
                </button>
              </div>

              <Button
                type="submit"
                loading={loading}
                className="w-full"
                size="lg"
                icon={<ArrowRight className="w-5 h-5" />}
              >
                Sign In
              >Iniciar Sesión</Button>
            </form>

            <div className="mt-6 relative">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-gray-200" />
              </div>
              <div className="relative flex justify-center text-xs">
                <span className="px-3 bg-gray-50 text-gray-400 uppercase tracking-wider font-medium">o</span>
              </div>
            </div>

            <div className="mt-6 grid grid-cols-2 gap-3">
              <button className="flex items-center justify-center gap-2 px-4 py-3 bg-white border border-gray-200 rounded-xl text-sm font-medium text-gray-700 hover:bg-gray-50 hover:border-gray-300 transition-all">
                <svg className="w-4 h-4" viewBox="0 0 24 24"><path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 01-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" fill="#4285F4"/><path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/><path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/><path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/></svg>
                Google
              </button>
              <button className="flex items-center justify-center gap-2 px-4 py-3 bg-white border border-gray-200 rounded-xl text-sm font-medium text-gray-700 hover:bg-gray-50 hover:border-gray-300 transition-all">
                <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/></svg>
                GitHub
              </button>
            </div>

            <p className="mt-8 text-center text-sm text-gray-500">
              ¿No tienes una cuenta?{' '}
              <Link to="/register" className="text-primary-600 font-semibold hover:text-primary-700 transition-colors">
                Crea una gratis
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
