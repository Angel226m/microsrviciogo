# ═══════════════════════════════════════════════════════════════
# CloudMart Frontend – Build multi-etapa optimizado
# Seguridad: usuario no-root, headers OWASP, compresión
# ═══════════════════════════════════════════════════════════════
FROM node:22-alpine AS builder

WORKDIR /app

# Copiar solo archivos de dependencias para aprovechar cache de Docker
COPY package.json package-lock.json* ./
RUN npm ci --no-audit --no-fund

# Copiar código fuente y construir
COPY . .
RUN npm run build

# ── Imagen de producción con Nginx hardened ──────────────────
FROM nginx:1.27-alpine

# Eliminar archivos de configuración por defecto y contenido innecesario
RUN rm -rf /etc/nginx/conf.d/default.conf /usr/share/nginx/html/* \
    && apk add --no-cache curl \
    && apk upgrade --no-cache

# Copiar artefactos de build y configuración
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Usuario no-root para cumplimiento OWASP/ISO 27001
RUN chown -R nginx:nginx /usr/share/nginx/html \
    && chown -R nginx:nginx /var/cache/nginx \
    && chown -R nginx:nginx /var/log/nginx \
    && touch /var/run/nginx.pid && chown nginx:nginx /var/run/nginx.pid

EXPOSE 3000

HEALTHCHECK --interval=15s --timeout=3s --retries=3 \
    CMD curl -f http://localhost:3000/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
