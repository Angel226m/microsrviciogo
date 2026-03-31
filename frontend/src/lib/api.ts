/**
 * Configuración del cliente HTTP (Axios)
 * - Interceptor de solicitud: inyecta token JWT automáticamente
 * - Interceptor de respuesta: redirige a login en caso de 401 (no autorizado)
 * - URL base: /api/v1 (proxy al API Gateway)
 */
import axios from 'axios';
import { useAuthStore } from '../store/authStore';

const api = axios.create({
  baseURL: '/api/v1',
  headers: { 'Content-Type': 'application/json' },
});

api.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      useAuthStore.getState().logout();
      window.location.href = '/login';
    }
    return Promise.reject(err);
  },
);

export default api;
