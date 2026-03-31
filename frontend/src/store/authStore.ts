/**
 * Store de Autenticación – Gestiona el estado del usuario y token JWT
 * Persistido en localStorage bajo la clave 'cloudmart-auth'
 */
import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User } from '../types';

export interface AuthState {
  user: User | null;
  token: string | null;
  setAuth: (user: User, token: string) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      setAuth: (user, token) => set({ user, token }),
      logout: () => set({ user: null, token: null }),
    }),
    { name: 'cloudmart-auth' },
  ),
);
