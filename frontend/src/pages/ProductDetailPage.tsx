import { useParams, Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { productService } from '../services/productService';
import { useCartStore } from '../store/cartStore';
import { Star, ShoppingCart, Minus, Plus, ChevronRight, Heart, Share2, Shield, Truck, RotateCcw, Package, Check } from 'lucide-react';
import { useState } from 'react';
import toast from 'react-hot-toast';
import clsx from 'clsx';
import Button from '../components/ui/Button';
import Badge from '../components/ui/Badge';

export default function ProductDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [qty, setQty] = useState(1);
  const [selectedImage, setSelectedImage] = useState(0);
  const [liked, setLiked] = useState(false);
  const addItem = useCartStore((s) => s.addItem);

  const { data: product, isLoading } = useQuery({
    queryKey: ['product', id],
    queryFn: () => productService.getById(id!),
    enabled: !!id,
  });

  const { data: reviews } = useQuery({
    queryKey: ['reviews', id],
    queryFn: () => productService.getReviews(id!),
    enabled: !!id,
  });

  if (isLoading) {
    return (
      <div className="max-w-7xl mx-auto px-4 py-12">
        {/* Breadcrumb skeleton */}
        <div className="h-5 w-48 skeleton rounded mb-8" />
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12">
          <div className="space-y-4">
            <div className="aspect-square skeleton rounded-3xl" />
            <div className="flex gap-3">
              {[...Array(4)].map((_, i) => (
                <div key={i} className="w-20 h-20 skeleton rounded-xl" />
              ))}
            </div>
          </div>
          <div className="space-y-6">
            <div className="h-10 skeleton rounded-xl w-3/4" />
            <div className="h-6 skeleton rounded w-1/3" />
            <div className="h-8 skeleton rounded w-1/4" />
            <div className="space-y-2">
              <div className="h-4 skeleton rounded w-full" />
              <div className="h-4 skeleton rounded w-5/6" />
              <div className="h-4 skeleton rounded w-4/6" />
            </div>
            <div className="h-14 skeleton rounded-2xl w-full" />
          </div>
        </div>
      </div>
    );
  }

  if (!product) {
    return (
      <div className="max-w-7xl mx-auto px-4 py-20 text-center animate-fade-in-up">
        <div className="w-24 h-24 mx-auto rounded-2xl bg-gradient-to-br from-gray-100 to-gray-50 flex items-center justify-center mb-6">
          <Package className="w-12 h-12 text-gray-300" />
        </div>
        <h2 className="text-2xl font-bold text-gray-900">Product not found</h2>
        <p className="text-gray-500 mt-2">The product you're looking for doesn't exist.</p>
        <Link
          to="/products"
          className="inline-block mt-6 px-6 py-3 bg-primary-50 text-primary-600 rounded-xl font-semibold hover:bg-primary-100 transition-colors"
        >
          Browse Products
        </Link>
      </div>
    );
  }

  const price = product.discount_price ?? product.price;
  const hasDiscount = product.discount_price != null && product.discount_price < product.price;
  const discountPercent = hasDiscount
    ? Math.round(((product.price - product.discount_price!) / product.price) * 100)
    : 0;
  const images = product.images?.length ? product.images : [];
  const inStock = product.stock == null || product.stock > 0;

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50/50 to-white">
      <div className="max-w-7xl mx-auto px-4 py-8">
        {/* Breadcrumbs */}
        <nav className="flex items-center gap-1.5 text-sm text-gray-500 mb-8 animate-fade-in-up">
          <Link to="/" className="hover:text-primary-600 transition-colors">Home</Link>
          <ChevronRight className="w-3.5 h-3.5" />
          <Link to="/products" className="hover:text-primary-600 transition-colors">Products</Link>
          <ChevronRight className="w-3.5 h-3.5" />
          <span className="text-gray-900 font-medium truncate max-w-[200px]">{product.name}</span>
        </nav>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 animate-fade-in-up">
          {/* Image gallery */}
          <div className="space-y-4">
            <div className="aspect-square bg-white rounded-3xl overflow-hidden border border-gray-100 relative group shadow-sm">
              {images[selectedImage] ? (
                <img
                  src={images[selectedImage]}
                  alt={product.name}
                  className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                />
              ) : (
                <div className="w-full h-full flex items-center justify-center bg-gradient-to-br from-gray-50 to-gray-100">
                  <div className="text-center">
                    <Package className="w-20 h-20 text-gray-300 mx-auto" />
                    <p className="text-gray-400 text-sm mt-2">No image available</p>
                  </div>
                </div>
              )}

              {/* Badges */}
              <div className="absolute top-4 left-4 flex flex-col gap-2">
                {hasDiscount && (
                  <span className="bg-gradient-to-r from-red-500 to-pink-500 text-white text-sm font-bold px-3 py-1.5 rounded-full shadow-lg">
                    -{discountPercent}% OFF
                  </span>
                )}
                {!inStock && (
                  <span className="bg-gray-900/80 text-white text-sm font-bold px-3 py-1.5 rounded-full backdrop-blur-md">
                    Out of Stock
                  </span>
                )}
              </div>

              {/* Actions overlay */}
              <div className="absolute top-4 right-4 flex flex-col gap-2">
                <button
                  onClick={() => setLiked(!liked)}
                  className={clsx(
                    'w-11 h-11 rounded-full flex items-center justify-center backdrop-blur-md shadow-lg transition-all',
                    liked ? 'bg-red-500 text-white' : 'bg-white/90 text-gray-600 hover:text-red-500'
                  )}
                >
                  <Heart className={clsx('w-5 h-5', liked && 'fill-current')} />
                </button>
                <button className="w-11 h-11 rounded-full bg-white/90 text-gray-600 flex items-center justify-center backdrop-blur-md shadow-lg hover:text-primary-600 transition-all">
                  <Share2 className="w-5 h-5" />
                </button>
              </div>
            </div>

            {/* Thumbnails */}
            {images.length > 1 && (
              <div className="flex gap-3 overflow-x-auto pb-2">
                {images.map((img, i) => (
                  <button
                    key={i}
                    onClick={() => setSelectedImage(i)}
                    className={clsx(
                      'w-20 h-20 rounded-xl overflow-hidden border-2 shrink-0 transition-all',
                      selectedImage === i
                        ? 'border-primary-500 ring-2 ring-primary-200'
                        : 'border-gray-100 hover:border-gray-300'
                    )}
                  >
                    <img src={img} alt="" className="w-full h-full object-cover" />
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* Product details */}
          <div className="flex flex-col">
            {product.category_name && (
              <span className="text-sm text-primary-600 font-semibold uppercase tracking-wider mb-2">
                {product.category_name}
              </span>
            )}

            <h1 className="text-3xl sm:text-4xl font-bold text-gray-900 leading-tight">{product.name}</h1>

            {/* Rating */}
            <div className="flex items-center gap-3 mt-4">
              <div className="flex items-center gap-0.5">
                {[...Array(5)].map((_, i) => (
                  <Star
                    key={i}
                    className={clsx(
                      'w-5 h-5',
                      i < Math.round(product.rating ?? 0) ? 'fill-amber-400 text-amber-400' : 'text-gray-200'
                    )}
                  />
                ))}
              </div>
              <span className="text-sm text-gray-500 font-medium">
                {product.rating?.toFixed(1)} ({product.review_count} reviews)
              </span>
            </div>

            {/* Price */}
            <div className="mt-6 flex items-baseline gap-3">
              <span className="text-4xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
                ${price.toFixed(2)}
              </span>
              {hasDiscount && (
                <>
                  <span className="text-xl text-gray-400 line-through">${product.price.toFixed(2)}</span>
                  <Badge variant="danger">Save {discountPercent}%</Badge>
                </>
              )}
            </div>

            {/* Description */}
            <p className="mt-6 text-gray-600 leading-relaxed text-[15px]">{product.description}</p>

            {/* Tags */}
            {product.tags?.length > 0 && (
              <div className="mt-5 flex flex-wrap gap-2">
                {product.tags.map((tag) => (
                  <span key={tag} className="px-3 py-1.5 bg-gray-100 text-gray-600 rounded-full text-xs font-medium hover:bg-gray-200 transition-colors cursor-default">
                    {tag}
                  </span>
                ))}
              </div>
            )}

            {/* Stock status */}
            <div className="mt-6 flex items-center gap-2">
              {inStock ? (
                <>
                  <div className="w-2.5 h-2.5 rounded-full bg-green-500 animate-pulse" />
                  <span className="text-sm font-medium text-green-700">
                    {product.stock != null ? `${product.stock} in stock` : 'In stock'}
                  </span>
                </>
              ) : (
                <>
                  <div className="w-2.5 h-2.5 rounded-full bg-red-500" />
                  <span className="text-sm font-medium text-red-600">Out of stock</span>
                </>
              )}
            </div>

            {/* Quantity selector + Add to cart */}
            <div className="mt-8 flex items-center gap-4">
              <div className="flex items-center bg-gray-100 rounded-xl">
                <button
                  onClick={() => setQty(Math.max(1, qty - 1))}
                  className="p-3 hover:bg-gray-200 rounded-l-xl transition-colors"
                >
                  <Minus className="w-4 h-4 text-gray-600" />
                </button>
                <span className="px-5 font-semibold text-gray-900 min-w-[3rem] text-center">{qty}</span>
                <button
                  onClick={() => setQty(qty + 1)}
                  className="p-3 hover:bg-gray-200 rounded-r-xl transition-colors"
                >
                  <Plus className="w-4 h-4 text-gray-600" />
                </button>
              </div>
              <Button
                size="lg"
                icon={<ShoppingCart className="w-5 h-5" />}
                className="flex-1"
                disabled={!inStock}
                onClick={() => {
                  addItem(product, qty);
                  toast.success('Added to cart!');
                }}
              >
                Add to Cart — ${(price * qty).toFixed(2)}
              </Button>
            </div>

            {/* Trust features */}
            <div className="mt-8 grid grid-cols-3 gap-4">
              {[
                { icon: Truck, label: 'Free Shipping', desc: 'Orders over $50' },
                { icon: Shield, label: 'Secure Payment', desc: '256-bit SSL' },
                { icon: RotateCcw, label: 'Easy Returns', desc: '30-day policy' },
              ].map(({ icon: Icon, label, desc }) => (
                <div key={label} className="text-center p-3 rounded-xl bg-gray-50">
                  <Icon className="w-5 h-5 text-primary-600 mx-auto" />
                  <p className="text-xs font-semibold text-gray-900 mt-1.5">{label}</p>
                  <p className="text-[11px] text-gray-500">{desc}</p>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Reviews section */}
        {reviews && reviews.length > 0 && (
          <section className="mt-20 animate-fade-in-up">
            <div className="flex items-center justify-between mb-8">
              <div>
                <h2 className="text-2xl font-bold text-gray-900">Customer Reviews</h2>
                <p className="text-gray-500 mt-1">{reviews.length} reviews for this product</p>
              </div>
              <div className="flex items-center gap-2 bg-amber-50 px-4 py-2.5 rounded-xl">
                <Star className="w-5 h-5 fill-amber-400 text-amber-400" />
                <span className="text-lg font-bold text-amber-700">{product.rating?.toFixed(1)}</span>
                <span className="text-sm text-amber-600">/ 5</span>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
              {reviews.map((review) => (
                <div
                  key={review.id}
                  className="bg-white p-6 rounded-2xl border border-gray-100 hover:shadow-md transition-shadow"
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-0.5">
                      {[...Array(5)].map((_, i) => (
                        <Star
                          key={i}
                          className={clsx(
                            'w-4 h-4',
                            i < review.rating ? 'fill-amber-400 text-amber-400' : 'text-gray-200'
                          )}
                        />
                      ))}
                    </div>
                    <div className="flex items-center gap-1.5 text-xs text-gray-400">
                      <Check className="w-3.5 h-3.5 text-green-500" />
                      Verified
                    </div>
                  </div>
                  <h4 className="mt-3 font-semibold text-gray-900">{review.title}</h4>
                  <p className="mt-2 text-sm text-gray-600 leading-relaxed">{review.comment}</p>
                </div>
              ))}
            </div>
          </section>
        )}
      </div>
    </div>
  );
}
