import { useQuery } from '@tanstack/react-query';
import { useSearchParams } from 'react-router-dom';
import { productService } from '../services/productService';
import ProductCard from '../components/ProductCard';
import { ProductCardSkeleton } from '../components/ui/Skeleton';
import { Search, SlidersHorizontal, X, ChevronLeft, ChevronRight, Sparkles, Grid3X3, LayoutGrid } from 'lucide-react';
import type { Category, Product } from '../types';
import { useState } from 'react';
import clsx from 'clsx';

export default function ProductsPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [search, setSearch] = useState(searchParams.get('search') ?? '');
  const [filtersOpen, setFiltersOpen] = useState(false);
  const [gridCols, setGridCols] = useState<2 | 3>(3);
  const page = Number(searchParams.get('page') ?? '1');
  const category = searchParams.get('category') ?? undefined;
  const searchQuery = searchParams.get('search') || undefined;

  const { data, isLoading } = useQuery({
    queryKey: ['products', page, category, searchQuery],
    queryFn: () => productService.list({ page, limit: 12, category, search: searchQuery }),
  });

  const { data: categories } = useQuery({
    queryKey: ['categories'],
    queryFn: productService.getCategories,
  });

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSearchParams((p) => {
      if (search) p.set('search', search);
      else p.delete('search');
      p.set('page', '1');
      return p;
    });
  };

  const clearFilters = () => {
    setSearch('');
    setSearchParams({});
  };

  const totalPages = data ? Math.ceil(data.total / 12) : 0;
  const hasActiveFilters = !!category || !!searchQuery;

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50/80 to-white">
      {/* Page header */}
      <div className="bg-white border-b border-gray-100">
        <div className="max-w-7xl mx-auto px-4 py-6 sm:py-8">
          <div className="flex items-center gap-3 mb-1">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-primary-500 to-accent-500 flex items-center justify-center">
              <Sparkles className="w-5 h-5 text-white" />
            </div>
            <div>
              <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">Products</h1>
              {data && (
                <p className="text-sm text-gray-500 mt-0.5">
                  Showing {data.data?.length ?? 0} of {data.total} products
                </p>
              )}
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 py-6">
        {/* Search & toolbar */}
        <div className="flex flex-col sm:flex-row gap-3 mb-6">
          <form onSubmit={handleSearch} className="flex-1 relative">
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              type="text"
              placeholder="Search by name, category, or keyword..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-full pl-12 pr-12 py-3.5 bg-white border border-gray-200 rounded-2xl focus:ring-2 focus:ring-primary-500/20 focus:border-primary-400 outline-none transition-all shadow-sm placeholder:text-gray-400"
            />
            {search && (
              <button
                type="button"
                onClick={() => { setSearch(''); setSearchParams((p) => { p.delete('search'); return p; }); }}
                className="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
              >
                <X className="w-4 h-4" />
              </button>
            )}
          </form>

          <div className="flex items-center gap-2">
            <button
              onClick={() => setFiltersOpen(!filtersOpen)}
              className={clsx(
                'flex items-center gap-2 px-4 py-3.5 rounded-2xl border font-medium text-sm transition-all',
                filtersOpen
                  ? 'bg-primary-50 border-primary-200 text-primary-700'
                  : 'bg-white border-gray-200 text-gray-600 hover:border-gray-300'
              )}
            >
              <SlidersHorizontal className="w-4 h-4" />
              <span className="hidden sm:inline">Filters</span>
              {hasActiveFilters && (
                <span className="w-2 h-2 rounded-full bg-primary-500" />
              )}
            </button>

            {/* Grid toggle - desktop */}
            <div className="hidden sm:flex items-center bg-white border border-gray-200 rounded-2xl p-1">
              <button
                onClick={() => setGridCols(2)}
                className={clsx('p-2.5 rounded-xl transition-all', gridCols === 2 ? 'bg-primary-50 text-primary-600' : 'text-gray-400 hover:text-gray-600')}
              >
                <Grid3X3 className="w-4 h-4" />
              </button>
              <button
                onClick={() => setGridCols(3)}
                className={clsx('p-2.5 rounded-xl transition-all', gridCols === 3 ? 'bg-primary-50 text-primary-600' : 'text-gray-400 hover:text-gray-600')}
              >
                <LayoutGrid className="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>

        {/* Active filters pills */}
        {hasActiveFilters && (
          <div className="flex items-center gap-2 mb-5 flex-wrap animate-fade-in-up">
            {searchQuery && (
              <span className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-primary-50 text-primary-700 rounded-full text-sm font-medium">
                Search: "{searchQuery}"
                <button onClick={() => { setSearch(''); setSearchParams((p) => { p.delete('search'); p.set('page', '1'); return p; }); }}>
                  <X className="w-3.5 h-3.5" />
                </button>
              </span>
            )}
            {category && (
              <span className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-accent-50 text-accent-700 rounded-full text-sm font-medium">
                {categories?.find((c: Category) => c.id === category)?.name ?? 'Category'}
                <button onClick={() => setSearchParams((p) => { p.delete('category'); p.set('page', '1'); return p; })}>
                  <X className="w-3.5 h-3.5" />
                </button>
              </span>
            )}
            <button onClick={clearFilters} className="text-sm text-gray-500 hover:text-gray-700 font-medium ml-1">
              Clear all
            </button>
          </div>
        )}

        <div className="flex gap-6">
          {/* Sidebar filters */}
          <aside className={clsx(
            'shrink-0 transition-all duration-300',
            filtersOpen ? 'w-64 opacity-100' : 'w-0 opacity-0 overflow-hidden'
          )}>
            <div className="bg-white rounded-2xl p-5 border border-gray-100 shadow-sm sticky top-24">
              <div className="flex items-center justify-between mb-4">
                <h3 className="font-semibold text-gray-900 text-sm uppercase tracking-wider">Categories</h3>
                {category && (
                  <button onClick={() => setSearchParams((p) => { p.delete('category'); return p; })} className="text-xs text-primary-600 font-medium hover:underline">
                    Reset
                  </button>
                )}
              </div>
              <div className="space-y-1">
                <button
                  onClick={() => setSearchParams((p) => { p.delete('category'); p.set('page', '1'); return p; })}
                  className={clsx(
                    'w-full text-left px-3 py-2.5 rounded-xl text-sm font-medium transition-all',
                    !category
                      ? 'bg-gradient-to-r from-primary-500 to-primary-600 text-white shadow-md shadow-primary-200'
                      : 'text-gray-600 hover:bg-gray-50'
                  )}
                >
                  All Products
                </button>
                {categories?.map((cat: Category) => (
                  <button
                    key={cat.id}
                    onClick={() => setSearchParams((p) => { p.set('category', cat.id); p.set('page', '1'); return p; })}
                    className={clsx(
                      'w-full text-left px-3 py-2.5 rounded-xl text-sm font-medium transition-all',
                      category === cat.id
                        ? 'bg-gradient-to-r from-primary-500 to-primary-600 text-white shadow-md shadow-primary-200'
                        : 'text-gray-600 hover:bg-gray-50'
                    )}
                  >
                    {cat.name}
                  </button>
                ))}
              </div>
            </div>
          </aside>

          {/* Product grid */}
          <div className="flex-1 min-w-0">
            {isLoading ? (
              <div className={clsx(
                'grid gap-5',
                gridCols === 2
                  ? 'grid-cols-1 sm:grid-cols-2'
                  : 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3'
              )}>
                {[...Array(6)].map((_, i) => (
                  <ProductCardSkeleton key={i} />
                ))}
              </div>
            ) : !data?.data?.length ? (
              <div className="text-center py-20 animate-fade-in-up">
                <div className="w-20 h-20 mx-auto rounded-2xl bg-gradient-to-br from-gray-100 to-gray-50 flex items-center justify-center mb-4">
                  <Search className="w-10 h-10 text-gray-300" />
                </div>
                <h3 className="text-xl font-semibold text-gray-900">No products found</h3>
                <p className="text-gray-500 mt-2 max-w-sm mx-auto">
                  Try adjusting your search or filter to find what you're looking for.
                </p>
                <button
                  onClick={clearFilters}
                  className="mt-4 px-5 py-2.5 bg-primary-50 text-primary-600 rounded-xl font-medium text-sm hover:bg-primary-100 transition-colors"
                >
                  Clear all filters
                </button>
              </div>
            ) : (
              <>
                <div className={clsx(
                  'grid gap-5',
                  gridCols === 2
                    ? 'grid-cols-1 sm:grid-cols-2'
                    : 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3'
                )}>
                  {data.data.map((product: Product, i: number) => (
                    <div key={product.id} className="animate-fade-in-up" style={{ animationDelay: `${i * 50}ms` }}>
                      <ProductCard product={product} />
                    </div>
                  ))}
                </div>

                {/* Pagination */}
                {totalPages > 1 && (
                  <div className="mt-10 flex items-center justify-center gap-1">
                    <button
                      onClick={() => setSearchParams((p) => { p.set('page', String(Math.max(1, page - 1))); return p; })}
                      disabled={page === 1}
                      className="p-2.5 rounded-xl bg-white border border-gray-200 text-gray-600 hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed transition-all"
                    >
                      <ChevronLeft className="w-4 h-4" />
                    </button>
                    {Array.from({ length: totalPages }, (_, i) => {
                      const pageNum = i + 1;
                      // Show first, last, and pages around current
                      if (pageNum === 1 || pageNum === totalPages || (pageNum >= page - 1 && pageNum <= page + 1)) {
                        return (
                          <button
                            key={i}
                            onClick={() => setSearchParams((p) => { p.set('page', String(pageNum)); return p; })}
                            className={clsx(
                              'w-10 h-10 rounded-xl text-sm font-semibold transition-all',
                              page === pageNum
                                ? 'bg-gradient-to-r from-primary-500 to-primary-600 text-white shadow-md shadow-primary-200'
                                : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50'
                            )}
                          >
                            {pageNum}
                          </button>
                        );
                      }
                      if (pageNum === page - 2 || pageNum === page + 2) {
                        return <span key={i} className="px-1 text-gray-400">...</span>;
                      }
                      return null;
                    })}
                    <button
                      onClick={() => setSearchParams((p) => { p.set('page', String(Math.min(totalPages, page + 1))); return p; })}
                      disabled={page === totalPages}
                      className="p-2.5 rounded-xl bg-white border border-gray-200 text-gray-600 hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed transition-all"
                    >
                      <ChevronRight className="w-4 h-4" />
                    </button>
                  </div>
                )}
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
