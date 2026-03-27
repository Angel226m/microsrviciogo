import type { Product } from '../types';
import { Link } from 'react-router-dom';
import { ShoppingCart, Star, Eye, Package, Heart } from 'lucide-react';
import { useCartStore } from '../store/cartStore';
import toast from 'react-hot-toast';
import { useState } from 'react';
import clsx from 'clsx';

interface Props {
  product: Product;
}

export default function ProductCard({ product }: Props) {
  const addItem = useCartStore((s) => s.addItem);
  const [liked, setLiked] = useState(false);
  const [imageLoaded, setImageLoaded] = useState(false);

  const handleAdd = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    addItem(product);
    toast.success(`${product.name} added to cart`);
  };

  const price = product.discount_price ?? product.price;
  const hasDiscount = product.discount_price != null && product.discount_price < product.price;
  const discountPercent = hasDiscount
    ? Math.round(((product.price - product.discount_price!) / product.price) * 100)
    : 0;

  return (
    <div className="group relative bg-white rounded-2xl border border-gray-100 overflow-hidden hover:shadow-xl hover:shadow-primary-100/50 hover:-translate-y-1 transition-all duration-300">
      {/* Image section */}
      <Link to={`/products/${product.id}`} className="block">
        <div className="aspect-[4/3] bg-gradient-to-br from-gray-50 to-gray-100 relative overflow-hidden">
          {product.images?.[0] ? (
            <>
              {!imageLoaded && (
                <div className="absolute inset-0 skeleton" />
              )}
              <img
                src={product.images[0]}
                alt={product.name}
                onLoad={() => setImageLoaded(true)}
                className={clsx(
                  'w-full h-full object-cover group-hover:scale-110 transition-transform duration-500',
                  !imageLoaded && 'opacity-0'
                )}
              />
            </>
          ) : (
            <div className="w-full h-full flex items-center justify-center">
              <div className="w-20 h-20 rounded-2xl bg-gradient-to-br from-primary-100 to-primary-50 flex items-center justify-center">
                <Package className="w-10 h-10 text-primary-400" />
              </div>
            </div>
          )}

          {/* Overlay actions */}
          <div className="absolute inset-0 bg-gradient-to-t from-black/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />

          {/* Quick actions */}
          <div className="absolute top-3 right-3 flex flex-col gap-2 translate-x-12 group-hover:translate-x-0 transition-transform duration-300">
            <button
              onClick={(e) => { e.preventDefault(); e.stopPropagation(); setLiked(!liked); }}
              className={clsx(
                'w-9 h-9 rounded-full flex items-center justify-center backdrop-blur-md transition-all shadow-lg',
                liked ? 'bg-red-500 text-white' : 'bg-white/90 text-gray-600 hover:bg-white hover:text-red-500'
              )}
            >
              <Heart className={clsx('w-4 h-4', liked && 'fill-current')} />
            </button>
            <Link
              to={`/products/${product.id}`}
              className="w-9 h-9 rounded-full bg-white/90 text-gray-600 flex items-center justify-center backdrop-blur-md hover:bg-white hover:text-primary-600 transition-all shadow-lg"
              onClick={(e) => e.stopPropagation()}
            >
              <Eye className="w-4 h-4" />
            </Link>
          </div>

          {/* Badges */}
          <div className="absolute top-3 left-3 flex flex-col gap-1.5">
            {hasDiscount && (
              <span className="bg-gradient-to-r from-red-500 to-orange-500 text-white text-[11px] font-bold px-2.5 py-1 rounded-full shadow-lg">
                -{discountPercent}%
              </span>
            )}
            {product.stock != null && product.stock <= 5 && product.stock > 0 && (
              <span className="bg-gradient-to-r from-amber-500 to-orange-500 text-white text-[11px] font-bold px-2.5 py-1 rounded-full shadow-lg">
                Only {product.stock} left
              </span>
            )}
          </div>
        </div>
      </Link>

      {/* Content */}
      <div className="p-4">
        <Link to={`/products/${product.id}`} className="block group/title">
          <h3 className="font-semibold text-gray-900 truncate group-hover/title:text-primary-600 transition-colors text-[15px]">
            {product.name}
          </h3>
        </Link>

        {product.category_name && (
          <span className="text-xs text-gray-400 font-medium uppercase tracking-wider mt-0.5 block">
            {product.category_name}
          </span>
        )}

        <div className="flex items-center gap-1.5 mt-2">
          <div className="flex items-center">
            {[...Array(5)].map((_, i) => (
              <Star
                key={i}
                className={clsx(
                  'w-3.5 h-3.5',
                  i < Math.round(product.rating ?? 0) ? 'fill-amber-400 text-amber-400' : 'text-gray-200'
                )}
              />
            ))}
          </div>
          <span className="text-xs text-gray-500">
            ({product.review_count ?? 0})
          </span>
        </div>

        <div className="flex items-center justify-between mt-3 pt-3 border-t border-gray-50">
          <div className="flex items-baseline gap-2">
            <span className="text-xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
              ${price.toFixed(2)}
            </span>
            {hasDiscount && (
              <span className="text-sm text-gray-400 line-through">${product.price.toFixed(2)}</span>
            )}
          </div>
          <button
            onClick={handleAdd}
            className="p-2.5 bg-gradient-to-r from-primary-500 to-primary-600 text-white rounded-xl hover:from-primary-600 hover:to-primary-700 transition-all shadow-md shadow-primary-200 hover:shadow-lg hover:shadow-primary-300 active:scale-95"
          >
            <ShoppingCart className="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  );
}
