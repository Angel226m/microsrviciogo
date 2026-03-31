/**
 * ═══════════════════════════════════════════════════════════════
 * Página de Inicio – Landing page principal de CloudMart
 * Secciones: Hero, barra de confianza, categorías, promociones, newsletter
 * ═══════════════════════════════════════════════════════════════
 */
import { Link } from 'react-router-dom';
import {
  ArrowRight, Truck, Shield, Headphones, Star, Sparkles,
  Laptop, Shirt, Home as HomeIcon, Dumbbell, BookOpen, Gamepad2,
  ChevronRight, CreditCard, RotateCcw, Package,
} from 'lucide-react';
import Button from '../components/ui/Button';

// Categorías principales de la tienda
const categorias = [
  { nombre: 'Electrónica', icon: Laptop, color: 'from-blue-500 to-cyan-500', href: '/products?category=electronics' },
  { nombre: 'Moda', icon: Shirt, color: 'from-teal-500 to-emerald-500', href: '/products?category=fashion' },
  { nombre: 'Hogar y Vida', icon: HomeIcon, color: 'from-amber-500 to-yellow-400', href: '/products?category=home' },
  { nombre: 'Deportes', icon: Dumbbell, color: 'from-green-500 to-emerald-400', href: '/products?category=sports' },
  { nombre: 'Libros', icon: BookOpen, color: 'from-sky-500 to-blue-500', href: '/products?category=books' },
  { nombre: 'Gaming', icon: Gamepad2, color: 'from-slate-600 to-gray-500', href: '/products?category=gaming' },
];

// Características de valor de la tienda
const caracteristicas = [
  { icon: Truck, titulo: 'Envío Gratis', desc: 'En pedidos mayores a $50', color: 'bg-blue-50 text-blue-600' },
  { icon: Shield, titulo: 'Pagos Seguros', desc: 'Cifrado de nivel empresarial', color: 'bg-emerald-50 text-emerald-600' },
  { icon: RotateCcw, titulo: 'Devoluciones Fáciles', desc: 'Política de 30 días', color: 'bg-amber-50 text-amber-600' },
  { icon: Headphones, titulo: 'Soporte 24/7', desc: 'Siempre aquí para ayudarte', color: 'bg-sky-50 text-sky-600' },
];

// Estadísticas de la plataforma
const estadisticas = [
  { valor: '50K+', etiqueta: 'Productos' },
  { valor: '100K+', etiqueta: 'Clientes' },
  { valor: '99.9%', etiqueta: 'Disponibilidad' },
  { valor: '4.9', etiqueta: 'Valoración', icon: Star },
];

export default function HomePage() {
  return (
    <div className="overflow-hidden">
      {/* ═══ Hero ═══ */}
      <section className="relative bg-gradient-to-br from-primary-700 via-primary-600 to-accent-600 animate-gradient text-white overflow-hidden">
        {/* Decorative blobs */}
        <div className="absolute top-0 left-0 w-96 h-96 bg-white/5 rounded-full blur-3xl -translate-x-1/2 -translate-y-1/2" />
        <div className="absolute bottom-0 right-0 w-[500px] h-[500px] bg-accent-500/10 rounded-full blur-3xl translate-x-1/3 translate-y-1/3" />
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-white/5 rounded-full blur-3xl" />

        <div className="max-w-7xl mx-auto px-4 py-20 sm:py-28 lg:py-36 relative">
          <div className="grid lg:grid-cols-2 gap-12 items-center">
            {/* Left */}
            <div className="animate-fade-in-up">
              <div className="inline-flex items-center gap-2 bg-white/10 backdrop-blur-sm border border-white/20 rounded-full px-4 py-1.5 text-sm font-medium mb-6">
                <Sparkles className="w-4 h-4 text-yellow-300" />
                Nueva colección disponible ahora
              </div>
              <h1 className="text-4xl sm:text-5xl lg:text-6xl font-extrabold tracking-tight leading-[1.1]">
                Tu Experiencia de
                <span className="block mt-2 bg-gradient-to-r from-emerald-200 to-cyan-200 bg-clip-text text-transparent">
                  Compras Cloud-Native
                </span>
              </h1>
              <p className="mt-6 text-lg sm:text-xl text-primary-100 max-w-lg leading-relaxed">
                Descubre productos premium impulsados por microservicios modernos. Ultra rápido, infinitamente confiable.
              </p>
              <div className="mt-8 flex flex-wrap gap-4">
                <Link to="/products">
                  <Button size="lg" className="bg-white text-primary-700 hover:bg-primary-50 shadow-2xl shadow-black/20">
                    <span>Comprar Ahora</span>
                    <ArrowRight className="w-5 h-5" />
                  </Button>
                </Link>
                <Link to="/register">
                  <Button size="lg" variant="outline" className="border-white/30 text-white hover:bg-white/10 hover:border-white/50">
                    Crear Cuenta
                  </Button>
                </Link>
              </div>

              {/* Fila de estadísticas */}
              <div className="mt-12 flex flex-wrap gap-8">
                {estadisticas.map(({ valor, etiqueta, icon: Icon }) => (
                  <div key={etiqueta} className="animate-fade-in-up">
                    <p className="text-2xl sm:text-3xl font-extrabold flex items-center gap-1">
                      {valor}
                      {Icon && <Icon className="w-5 h-5 fill-yellow-300 text-yellow-300" />}
                    </p>
                    <p className="text-sm text-primary-200">{etiqueta}</p>
                  </div>
                ))}
              </div>
            </div>

            {/* Right — floating cards illustration */}
            <div className="hidden lg:flex justify-center relative">
              <div className="relative w-80 h-80">
                <div className="absolute top-0 right-0 w-48 h-60 bg-white/10 backdrop-blur-md rounded-3xl border border-white/20 animate-float shadow-2xl p-6 flex flex-col justify-between">
                  <Package className="w-10 h-10 text-white/80" />
                  <div>
                    <p className="text-2xl font-bold">$49.99</p>
                    <p className="text-sm text-white/70">Producto Premium</p>
                  </div>
                </div>
                <div className="absolute bottom-4 left-0 w-52 h-36 bg-white/10 backdrop-blur-md rounded-2xl border border-white/20 animate-float shadow-2xl p-5 flex items-center gap-4" style={{ animationDelay: '2s' }}>
                  <div className="w-14 h-14 rounded-xl bg-gradient-to-br from-yellow-400 to-orange-400 flex items-center justify-center shrink-0">
                    <CreditCard className="w-7 h-7 text-white" />
                  </div>
                  <div>
                    <p className="text-sm font-semibold">Pago Seguro</p>
                    <p className="text-xs text-white/60 mt-0.5">Cifrado de 256 bits</p>
                  </div>
                </div>
                <div className="absolute top-1/3 left-8 w-14 h-14 bg-gradient-to-br from-emerald-400 to-cyan-400 rounded-2xl flex items-center justify-center animate-float shadow-lg" style={{ animationDelay: '4s' }}>
                  <Star className="w-7 h-7 text-white fill-white" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* ═══ Barra de confianza ═══ */}
      <section className="border-b border-gray-100 bg-white">
        <div className="max-w-7xl mx-auto px-4 py-6">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
            {caracteristicas.map(({ icon: Icon, titulo, desc, color }) => (
              <div key={titulo} className="flex items-center gap-3">
                <div className={`w-11 h-11 rounded-xl flex items-center justify-center shrink-0 ${color}`}>
                  <Icon className="w-5 h-5" />
                </div>
                <div>
                  <p className="text-sm font-semibold text-gray-900">{titulo}</p>
                  <p className="text-xs text-gray-500">{desc}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ═══ Categorías ═══ */}
      <section className="max-w-7xl mx-auto px-4 py-16 sm:py-20">
        <div className="text-center mb-12 animate-fade-in-up">
          <h2 className="text-3xl sm:text-4xl font-extrabold text-gray-900">
            Comprar por Categoría
          </h2>
          <p className="mt-3 text-gray-500 max-w-md mx-auto">
            Explora nuestras colecciones seleccionadas para cada estilo de vida
          </p>
        </div>

        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-4">
          {categorias.map(({ nombre, icon: Icon, color, href }, i) => (
            <Link
              key={nombre}
              to={href}
              className={`group relative flex flex-col items-center gap-3 p-6 rounded-2xl bg-white border border-gray-100 hover:border-gray-200 hover:shadow-xl hover:shadow-gray-200/50 transition-all duration-300 animate-fade-in-up animate-delay-${i * 100}`}
            >
              <div className={`w-14 h-14 rounded-2xl bg-gradient-to-br ${color} flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform duration-300`}>
                <Icon className="w-7 h-7 text-white" />
              </div>
              <span className="text-sm font-semibold text-gray-800 group-hover:text-primary-600 transition-colors">
                {nombre}
              </span>
            </Link>
          ))}
        </div>
      </section>

      {/* ═══ Promo banners ═══ */}
      <section className="max-w-7xl mx-auto px-4 pb-16 sm:pb-20">
        <div className="grid md:grid-cols-2 gap-6">
          <div className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-primary-600 to-primary-800 p-8 sm:p-10 text-white group">
            <div className="absolute -right-8 -bottom-8 w-40 h-40 bg-white/10 rounded-full group-hover:scale-150 transition-transform duration-500" />
            <div className="relative">
              <span className="text-sm font-semibold text-primary-200 uppercase tracking-wider">Tiempo limitado</span>
              <h3 className="mt-2 text-2xl sm:text-3xl font-extrabold">Hasta 40% de descuento</h3>
              <p className="mt-2 text-primary-100">En electrónica y gadgets seleccionados</p>
              <Link to="/products?category=electronics" className="inline-flex items-center gap-1.5 mt-5 text-sm font-semibold hover:gap-3 transition-all">
                Ver ofertas <ChevronRight className="w-4 h-4" />
              </Link>
            </div>
          </div>

          <div className="relative overflow-hidden rounded-3xl bg-gradient-to-br from-gray-900 to-gray-800 p-8 sm:p-10 text-white group">
            <div className="absolute -right-8 -bottom-8 w-40 h-40 bg-white/5 rounded-full group-hover:scale-150 transition-transform duration-500" />
            <div className="relative">
              <span className="text-sm font-semibold text-gray-400 uppercase tracking-wider">Nueva colección</span>
              <h3 className="mt-2 text-2xl sm:text-3xl font-extrabold">Primavera 2026</h3>
              <p className="mt-2 text-gray-400">Estilos frescos para la nueva temporada</p>
              <Link to="/products?category=fashion" className="inline-flex items-center gap-1.5 mt-5 text-sm font-semibold hover:gap-3 transition-all">
                Descubrir ahora <ChevronRight className="w-4 h-4" />
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* ═══ CTA Newsletter ═══ */}
      <section className="bg-gradient-to-r from-primary-50 via-gray-50 to-accent-50">
        <div className="max-w-7xl mx-auto px-4 py-16 sm:py-20">
          <div className="max-w-2xl mx-auto text-center animate-fade-in-up">
            <h2 className="text-3xl sm:text-4xl font-extrabold text-gray-900">
              Mantente al día
            </h2>
            <p className="mt-3 text-gray-500">
              Recibe las últimas ofertas, novedades y promociones exclusivas en tu correo.
            </p>
            <form className="mt-8 flex flex-col sm:flex-row gap-3 max-w-md mx-auto" onSubmit={(e) => e.preventDefault()}>
              <input
                type="email"
                placeholder="Ingresa tu correo electrónico"
                className="flex-1 px-5 py-3.5 rounded-xl border border-gray-200 bg-white text-sm focus:ring-2 focus:ring-primary-500/30 focus:border-primary-500 outline-none"
              />
              <Button size="lg" type="submit">Suscribirse</Button>
            </form>
            <p className="mt-3 text-xs text-gray-400">Sin spam, cancela cuando quieras.</p>
          </div>
        </div>
      </section>
    </div>
  );
}
