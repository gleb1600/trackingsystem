CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,                            -- Уникальный ID заказа
    status VARCHAR(20) NOT NULL,  -- Статус заказа
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),   -- Дата создания
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()    -- Дата обновления
);

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,                          -- Уникальный ID товара
    name VARCHAR(255) NOT NULL,                     -- Название товара
    description TEXT,                               -- Описание
    quantity INT NOT NULL DEFAULT 0,                -- Количество на складе
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),   -- Дата создания
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()    -- Дата обновления
);

CREATE TABLE IF NOT EXISTS order_items (
    order_id UUID REFERENCES orders(id) ON DELETE CASCADE,     -- Ссылка на заказ
    product_id UUID REFERENCES products(id) ON DELETE CASCADE, -- Ссылка на товар
    quantity INT NOT NULL CHECK (quantity > 0),                -- Количество товара
    PRIMARY KEY (order_id, product_id)                         -- Составной первичный ключ
);

-- Индексы для ускорения поиска
-- CREATE INDEX idx_orders_customer ON orders(customer_id);
-- CREATE INDEX idx_orders_status ON orders(status);