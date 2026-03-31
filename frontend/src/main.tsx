/**
 * ═══════════════════════════════════════════════════════════════
 * CloudMart – Punto de entrada principal del frontend
 * Configura proveedores: React Query, Router, Notificaciones
 * ═══════════════════════════════════════════════════════════════
 */
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter } from 'react-router-dom';
import { Toaster } from 'react-hot-toast';
import App from './App';
import './index.css';

// Configuración del cliente de React Query con políticas de cache y reintentos
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // Los datos se consideran frescos por 5 minutos
      retry: 1,                   // Un solo reintento en caso de error
      refetchOnWindowFocus: false, // No refrescar al cambiar de pestaña
    },
  },
});

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <App />
        <Toaster
          position="top-right"
          toastOptions={{
            duration: 3000,
            style: {
              borderRadius: '12px',
              background: '#1e293b',
              color: '#f8fafc',
              fontSize: '14px',
            },
          }}
        />
      </BrowserRouter>
    </QueryClientProvider>
  </StrictMode>,
);
