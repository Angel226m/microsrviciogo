/**
 * Servicio de Autenticación – Comunicación con la API de usuarios
 * Endpoints: login, registro, perfil, actualización de perfil
 */
import api from '../lib/api';
import type { User, AuthTokens } from '../types';

export const authService = {
  login: (email: string, password: string) =>
    api.post<{ user: User; tokens: AuthTokens }>('/auth/login', { email, password }).then((r) => r.data),

  register: (data: { email: string; password: string; first_name: string; last_name: string }) =>
    api.post<{ user: User; tokens: AuthTokens }>('/auth/register', data).then((r) => r.data),

  getProfile: () => api.get<User>('/users/profile').then((r) => r.data),

  updateProfile: (data: Partial<User>) =>
    api.put<User>('/users/profile', data).then((r) => r.data),
};
