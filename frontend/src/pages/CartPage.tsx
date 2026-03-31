/**
 * ═══════════════════════════════════════════════════════════════
 * Página del Carrito – Gestión de artículos y resumen del pedido
 * ═══════════════════════════════════════════════════════════════
 */
import { Link } from 'react-router-dom';
import { useCartStore } from '../store/cartStore';
import { Trash2, Minus, Plus, ShoppingBag, ArrowRight, Shield, Truck, Tag, ChevronLeft } from 'lucide-react';
import type { CartItem } from '../types';
import Button from '../components/ui/Button';
import EmptyState from '../components/ui/EmptyState';

export default function CartPage() {
  const { items, removeItem, updateQuantity, total, clearCart, count: itemCount } = useCartStore();

  if (items.length === 0) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center px-4">
        <EmptyState
          icon={ShoppingBag}
          title="Tu carrito está vacío"
          description="Parece que aún no has agregado nada a tu carrito. ¡Empieza a explorar nuestros increíbles productos!"
          actionLabel="Explorar Productos"
          actionTo="/products"
        />
      </div>
    );
  }

  const subtotal = total();
  const shipping = subtotal >= 50 ? 0 : 9.99;
  const tax = subtotal * 0.16;
  const grandTotal = subtotal + tax + shipping;
  const count = itemCount();

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50/50 to-white">
      <div className="max-w-7xl mx-auto px-4 py-8">
        {/* Header */}
        <div className="flex items-center justify-between mb-8 animate-fade-in-up">
          <div>
            <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">Carrito de Compras</h1>
            <p className="text-gray-500 mt-1">{count} {count === 1 ? 'artículo' : 'artículos'} en tu carrito</p>
          </div>
          <Link
            to="/products"
            className="flex items-center gap-2 text-sm text-primary-600 font-semibold hover:text-primary-700 transition-colors"
          >
            <ChevronLeft className="w-4 h-4" />
            Seguir Comprando
          </Link>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Cart items */}
          <div className="lg:col-span-2 space-y-4 animate-fade-in-up">
            {items.map(({ product, quantity }: CartItem, index: number) => {
              const price = product.discount_price ?? product.price;
              const hasDiscount = product.discount_price != null && product.discount_price < product.price;
              return (
                <div
                  key={product.id}
                  className="bg-white rounded-2xl p-5 sm:p-6 border border-gray-100 flex gap-4 sm:gap-6 hover:shadow-md transition-shadow group"
                  style={{ animationDelay: `${index * 50}ms` }}
                >
                  {/* Image */}
                  <Link
                    to={`/products/${product.id}`}
                    className="w-24 h-24 sm:w-28 sm:h-28 bg-gradient-to-br from-gray-50 to-gray-100 rounded-xl overflow-hidden shrink-0 group-hover:shadow-md transition-shadow"
                  >
                    {product.images?.[0] ? (
                      <img src={product.images[0]} alt={product.name} className="w-full h-full object-cover" />
                    ) : (
                      <div className="w-full h-full flex items-center justify-center">
                        <ShoppingBag className="w-8 h-8 text-gray-300" />
                      </div>
                    )}
                  </Link>

                  {/* Info */}
                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between gap-2">
                      <div className="min-w-0">
                        <Link
                          to={`/products/${product.id}`}
                          className="font-semibold text-gray-900 hover:text-primary-600 transition-colors truncate block text-[15px]"
                        >
                          {product.name}
                        </Link>
                        {hasDiscount && (
                          <span className="text-xs text-red-500 font-medium mt-0.5 inline-block">
                            {Math.round(((product.price - product.discount_price!) / product.price) * 100)}% off
                          </span>
                        )}
                      </div>
                      <button
                        onClick={() => removeItem(product.id)}
                        className="p-2 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-all shrink-0"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>

                    <div className="flex items-end justify-between mt-3 sm:mt-4">
                      {/* Quantity */}
                      <div className="flex items-center bg-gray-100 rounded-xl">
                        <button
                          onClick={() => updateQuantity(product.id, quantity - 1)}
                          className="p-2 hover:bg-gray-200 rounded-l-xl transition-colors"
                        >
                          <Minus className="w-3.5 h-3.5 text-gray-600" />
                        </button>
                        <span className="px-4 text-sm font-semibold text-gray-900 min-w-[2.5rem] text-center">
                          {quantity}
                        </span>
                        <button
                          onClick={() => updateQuantity(product.id, quantity + 1)}
                          className="p-2 hover:bg-gray-200 rounded-r-xl transition-colors"
                        >
                          <Plus className="w-3.5 h-3.5 text-gray-600" />
                        </button>
                      </div>

                      {/* Price */}
                      <div className="text-right">
                        <p className="text-lg font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
                          ${(price * quantity).toFixed(2)}
                        </p>
                        {quantity > 1 && (
                          <p className="text-xs text-gray-400">${price.toFixed(2)} each</p>
                        )}
                      </div>
                    </div>
                  </div>
                </div>
              );
            })}

            <button
              onClick={clearCart}
              className="flex items-center gap-2 text-sm text-red-500 hover:text-red-600 font-medium px-2 py-1 hover:bg-red-50 rounded-lg transition-all"
            >
              <Trash2 className="w-3.5 h-3.5" />
              Vaciar Carrito
            </button>
          </div>

          {/* Order summary */}
          <div className="animate-fade-in-up" style={{ animationDelay: '100ms' }}>
            <div className="bg-white rounded-2xl border border-gray-100 shadow-sm sticky top-24 overflow-hidden">
              <div className="p-6">
                <h3 className="text-lg font-bold text-gray-900 mb-5">Resumen del Pedido</h3>

                <div className="space-y-3.5">
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-500">Subtotal ({count} artículos)</span>
                    <span className="font-medium text-gray-900">${subtotal.toFixed(2)}</span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-500">Envío</span>
                    {shipping === 0 ? (
                      <span className="font-medium text-green-600">Gratis</span>
                    ) : (
                      <span className="font-medium text-gray-900">${shipping.toFixed(2)}</span>
                    )}
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-500">Impuesto (16% IVA)</span>
                    <span className="font-medium text-gray-900">${tax.toFixed(2)}</span>
                  </div>

                  {subtotal < 50 && (
                    <div className="bg-amber-50 text-amber-700 text-xs p-3 rounded-xl">
                      <Truck className="w-4 h-4 inline mr-1.5" />
                      Agrega <span className="font-bold">${(50 - subtotal).toFixed(2)}</span> más para envío gratis!
                    </div>
                  )}

                  <hr className="border-gray-100" />

                  <div className="flex justify-between items-baseline pt-1">
                    <span className="text-base font-bold text-gray-900">Total</span>
                    <span className="text-2xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
                      ${grandTotal.toFixed(2)}
                    </span>
                  </div>
                </div>
              </div>

              <div className="px-6 pb-6">
                <Link to="/login">
                  <Button className="w-full" size="lg" icon={<ArrowRight className="w-5 h-5" />}>
                    Proceder al Pago
                  </Button>
                </Link>
              </div>

              {/* Trust badges */}
              <div className="bg-gray-50 px-6 py-4 flex items-center justify-center gap-6 text-xs text-gray-500">
                <span className="flex items-center gap-1.5">
                  <Shield className="w-3.5 h-3.5 text-green-500" />
                  Seguro
                </span>
                <span className="flex items-center gap-1.5">
                  <Truck className="w-3.5 h-3.5 text-primary-500" />
                  Envío Rápido
                </span>
                <span className="flex items-center gap-1.5">
                  <Tag className="w-3.5 h-3.5 text-amber-500" />
                  Mejor Precio
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
