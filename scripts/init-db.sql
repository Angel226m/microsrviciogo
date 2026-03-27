-- CloudMart Database Initialization
-- Creates all schemas for microservices (each service owns its schema)

-- ═══════════════════════════════════════════════════════════════
-- USER SERVICE
-- ═══════════════════════════════════════════════════════════════
CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE users.accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    avatar_url TEXT,
    role VARCHAR(20) DEFAULT 'customer' CHECK (role IN ('customer', 'admin', 'seller')),
    is_active BOOLEAN DEFAULT true,
    email_verified BOOLEAN DEFAULT false,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE users.addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users.accounts(id) ON DELETE CASCADE,
    label VARCHAR(50) DEFAULT 'home',
    street VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    zip_code VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL DEFAULT 'MX',
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users.accounts(email);
CREATE INDEX idx_addresses_user ON users.addresses(user_id);

-- ═══════════════════════════════════════════════════════════════
-- PRODUCT SERVICE
-- ═══════════════════════════════════════════════════════════════
CREATE SCHEMA IF NOT EXISTS products;

CREATE TABLE products.categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    image_url TEXT,
    parent_id UUID REFERENCES products.categories(id),
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE products.items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    short_description VARCHAR(500),
    price DECIMAL(12,2) NOT NULL CHECK (price >= 0),
    compare_at_price DECIMAL(12,2),
    cost DECIMAL(12,2),
    category_id UUID REFERENCES products.categories(id),
    brand VARCHAR(100),
    tags TEXT[],
    images TEXT[],
    thumbnail_url TEXT,
    weight DECIMAL(8,2),
    dimensions JSONB,
    attributes JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    is_featured BOOLEAN DEFAULT false,
    rating_avg DECIMAL(3,2) DEFAULT 0,
    rating_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE products.reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products.items(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    title VARCHAR(200),
    body TEXT,
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_products_category ON products.items(category_id);
CREATE INDEX idx_products_sku ON products.items(sku);
CREATE INDEX idx_products_slug ON products.items(slug);
CREATE INDEX idx_products_featured ON products.items(is_featured) WHERE is_featured = true;
CREATE INDEX idx_reviews_product ON products.reviews(product_id);

-- ═══════════════════════════════════════════════════════════════
-- ORDER SERVICE
-- ═══════════════════════════════════════════════════════════════
CREATE SCHEMA IF NOT EXISTS orders;

CREATE TABLE orders.orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(20) UNIQUE NOT NULL,
    user_id UUID NOT NULL,
    status VARCHAR(30) DEFAULT 'pending' CHECK (status IN ('pending','confirmed','processing','shipped','delivered','cancelled','refunded')),
    subtotal DECIMAL(12,2) NOT NULL,
    tax DECIMAL(12,2) DEFAULT 0,
    shipping_cost DECIMAL(12,2) DEFAULT 0,
    discount DECIMAL(12,2) DEFAULT 0,
    total DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'MXN',
    shipping_address JSONB NOT NULL,
    billing_address JSONB,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE orders.order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders.orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    product_sku VARCHAR(50) NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(12,2) NOT NULL,
    total_price DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_orders_user ON orders.orders(user_id);
CREATE INDEX idx_orders_status ON orders.orders(status);
CREATE INDEX idx_order_items_order ON orders.order_items(order_id);

-- ═══════════════════════════════════════════════════════════════
-- PAYMENT SERVICE
-- ═══════════════════════════════════════════════════════════════
CREATE SCHEMA IF NOT EXISTS payments;

CREATE TABLE payments.transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    user_id UUID NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'MXN',
    method VARCHAR(30) NOT NULL CHECK (method IN ('credit_card','debit_card','paypal','bank_transfer','cash')),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending','processing','completed','failed','refunded')),
    provider VARCHAR(50),
    provider_tx_id VARCHAR(255),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_payments_order ON payments.transactions(order_id);
CREATE INDEX idx_payments_user ON payments.transactions(user_id);

-- ═══════════════════════════════════════════════════════════════
-- INVENTORY SERVICE
-- ═══════════════════════════════════════════════════════════════
CREATE SCHEMA IF NOT EXISTS inventory;

CREATE TABLE inventory.stock (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID UNIQUE NOT NULL,
    sku VARCHAR(50) NOT NULL,
    quantity INT NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    reserved INT NOT NULL DEFAULT 0 CHECK (reserved >= 0),
    warehouse VARCHAR(100) DEFAULT 'main',
    reorder_level INT DEFAULT 10,
    reorder_quantity INT DEFAULT 50,
    last_restocked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE inventory.movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('inbound','outbound','reservation','release','adjustment')),
    quantity INT NOT NULL,
    reference_id UUID,
    reference_type VARCHAR(30),
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_stock_product ON inventory.stock(product_id);
CREATE INDEX idx_stock_sku ON inventory.stock(sku);
CREATE INDEX idx_movements_product ON inventory.movements(product_id);

-- ═══════════════════════════════════════════════════════════════
-- SEED DATA
-- ═══════════════════════════════════════════════════════════════
INSERT INTO users.accounts (id, email, password_hash, first_name, last_name, role) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'admin@cloudmart.dev', '$2a$12$LJ3m4ys3ez0RVFEahwXkLejOJRYgvFsm6G5XYLO.FjWHPc/zvpjuW', 'Admin', 'CloudMart', 'admin'),
    ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'customer@cloudmart.dev', '$2a$12$LJ3m4ys3ez0RVFEahwXkLejOJRYgvFsm6G5XYLO.FjWHPc/zvpjuW', 'María', 'García', 'customer');

INSERT INTO products.categories (id, name, slug, description) VALUES
    ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'Electronics', 'electronics', 'Gadgets, devices & accessories'),
    ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'Clothing', 'clothing', 'Fashion & apparel'),
    ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'Home & Garden', 'home-garden', 'Home improvement & garden supplies');

INSERT INTO products.items (id, sku, name, slug, description, short_description, price, compare_at_price, category_id, brand, tags, images, is_featured) VALUES
    ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'ELEC-001', 'Wireless Noise-Cancelling Headphones', 'wireless-noise-cancelling-headphones', 'Premium wireless headphones with active noise cancellation, 40h battery life, and Hi-Res Audio support.', 'Premium ANC headphones with 40h battery', 2499.99, 3299.99, 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'SoundPro', ARRAY['headphones','wireless','anc'], ARRAY['https://picsum.photos/seed/hp1/800/800','https://picsum.photos/seed/hp2/800/800'], true),
    ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'ELEC-002', 'Smart Watch Ultra', 'smart-watch-ultra', 'Advanced smartwatch with GPS, heart rate monitor, SpO2 sensor, and 14-day battery life.', 'Advanced smartwatch with 14-day battery', 4999.99, 5999.99, 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'TechWear', ARRAY['smartwatch','fitness','gps'], ARRAY['https://picsum.photos/seed/sw1/800/800'], true),
    ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'CLTH-001', 'Premium Cotton T-Shirt', 'premium-cotton-tshirt', 'Ultra-soft organic cotton t-shirt, pre-shrunk with reinforced seams.', 'Organic cotton, ultra-soft feel', 599.99, NULL, 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'EcoWear', ARRAY['tshirt','cotton','organic'], ARRAY['https://picsum.photos/seed/ts1/800/800'], false),
    ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'HOME-001', 'Smart LED Desk Lamp', 'smart-led-desk-lamp', 'Wi-Fi enabled desk lamp with color temperature control, brightness adjustment, and app control.', 'Smart desk lamp with Wi-Fi control', 1299.99, 1599.99, 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'LumiHome', ARRAY['lamp','smart-home','led'], ARRAY['https://picsum.photos/seed/dl1/800/800'], true),
    ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a05', 'ELEC-003', 'Mechanical Gaming Keyboard', 'mechanical-gaming-keyboard', 'RGB mechanical keyboard with hot-swappable switches, PBT keycaps, and wireless connectivity.', 'Hot-swappable mechanical keyboard', 1899.99, 2199.99, 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'KeyForge', ARRAY['keyboard','mechanical','gaming'], ARRAY['https://picsum.photos/seed/kb1/800/800'], true);

INSERT INTO inventory.stock (product_id, sku, quantity, reorder_level) VALUES
    ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'ELEC-001', 150, 20),
    ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'ELEC-002', 75, 15),
    ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'CLTH-001', 500, 50),
    ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'HOME-001', 200, 25),
    ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a05', 'ELEC-003', 100, 15);
