import api from '../lib/api';
import type { Order, PaginatedResponse } from '../types';

export interface CreateOrderRequest {
  items: { product_id: string; quantity: number; unit_price: number }[];
  shipping_address: {
    street: string;
    city: string;
    state: string;
    zip_code: string;
    country: string;
  };
}

export const orderService = {
  create: (data: CreateOrderRequest) =>
    api.post<Order>('/orders', data).then((r) => r.data),

  list: (params?: { page?: number; limit?: number }) =>
    api.get<PaginatedResponse<Order>>('/orders', { params }).then((r) => r.data),

  getById: (id: string) => api.get<Order>(`/orders/${id}`).then((r) => r.data),

  cancel: (id: string) => api.post(`/orders/${id}/cancel`).then((r) => r.data),
};
