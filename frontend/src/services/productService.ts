/**
 * Servicio de Productos – Comunicación con la API de productos
 * Incluye fallback a datos de demostración cuando la API no está disponible
 */
import api from '../lib/api';
import type { Product, PaginatedResponse, Category, Review } from '../types';
import { mockProducts, mockCategories } from '../data/mockProducts';

function getMockPaginated(params?: { page?: number; limit?: number; category?: string; search?: string }): PaginatedResponse<Product> {
  const page = params?.page ?? 1;
  const limit = params?.limit ?? 12;
  let filtered = [...mockProducts];

  if (params?.category) {
    filtered = filtered.filter((p) => p.category_id === params.category);
  }
  if (params?.search) {
    const q = params.search.toLowerCase();
    filtered = filtered.filter(
      (p) => p.name.toLowerCase().includes(q) || p.description.toLowerCase().includes(q) || p.tags.some((t) => t.toLowerCase().includes(q)),
    );
  }

  const start = (page - 1) * limit;
  const data = filtered.slice(start, start + limit);
  return { data, total: filtered.length, page, limit };
}

export const productService = {
  list: async (params?: { page?: number; limit?: number; category?: string; search?: string }) => {
    try {
      return await api.get<PaginatedResponse<Product>>('/products', { params }).then((r) => r.data);
    } catch {
      console.warn('[CloudMart] API no disponible — cargando productos de demostración');
      return getMockPaginated(params);
    }
  },

  getById: async (id: string) => {
    try {
      return await api.get<Product>(`/products/${id}`).then((r) => r.data);
    } catch {
      const mock = mockProducts.find((p) => p.id === id);
      if (mock) return mock;
      throw new Error('Product not found');
    }
  },

  getCategories: async () => {
    try {
      return await api.get<Category[]>('/products/categories').then((r) => r.data);
    } catch {
      return mockCategories;
    }
  },

  getReviews: (productId: string) =>
    api.get<Review[]>(`/products/${productId}/reviews`).then((r) => r.data).catch(() => []),

  createReview: (productId: string, data: { rating: number; title: string; comment: string }) =>
    api.post<Review>(`/products/${productId}/reviews`, data).then((r) => r.data),
};
