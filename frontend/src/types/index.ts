export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: string;
  avatar_url?: string;
  created_at: string;
}

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  expires_at: string;
}

export interface Product {
  id: string;
  name: string;
  slug?: string;
  description: string;
  price: number;
  discount_price?: number;
  category_id: string;
  images: string[];
  tags: string[];
  rating: number;
  review_count: number;
  is_active?: boolean;
  stock?: number;
  category_name?: string;
  created_at?: string;
}

export interface Category {
  id: string;
  name: string;
  slug?: string;
  description?: string;
  parent_id?: string;
}

export interface CartItem {
  product: Product;
  quantity: number;
}

export interface OrderItem {
  product_id: string;
  product_name: string;
  quantity: number;
  unit_price: number;
  total: number;
}

export interface Order {
  id: string;
  order_number: string;
  user_id: string;
  status: string;
  items: OrderItem[];
  subtotal: number;
  tax: number;
  total: number;
  created_at: string;
}

export interface Review {
  id: string;
  product_id: string;
  user_id: string;
  rating: number;
  title: string;
  comment: string;
  created_at: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
}
