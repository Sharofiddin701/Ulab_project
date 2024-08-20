-- Admin table migration
CREATE TABLE IF NOT EXISTS "admin" (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL,
    "phone_number" VARCHAR(100) NOT NULL,
    "username" varchar(25) not null,
    "password" VARCHAR(100) NOT NULL,
    "address" VARCHAR(255),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

-- Customer table migration
CREATE TABLE IF NOT EXISTS "customer" (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL,
    "phone_number" VARCHAR(100) NOT NULL,
    "address" VARCHAR(255),
    "email" VARCHAR(100) NOT NULL,
    "password" VARCHAR(100) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

-- Brand table migration
CREATE TABLE IF NOT EXISTS "brand" (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL,
    "brand_image" VARCHAR(255),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP
);

-- Category table migration
CREATE TABLE IF NOT EXISTS "category" (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL,
    "url" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP
);

-- Orders table migration
CREATE TABLE IF NOT EXISTS "orders" (
    "id" UUID PRIMARY KEY,
    "customer_id" UUID REFERENCES "customer"("id"),
    "shipping" VARCHAR(255) CHECK ("shipping" IN ('yetkazish', 'ozi olib ketishi')) NOT NULL,
    "payment" VARCHAR(255) CHECK ("payment" IN ('click', 'naxt')) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

-- Product table migration
CREATE TABLE IF NOT EXISTS "product" (
     "id" UUID PRIMARY KEY,
    "is_favourite" BOOLEAN NOT NULL,
    "image" VARCHAR(255),
    "name" VARCHAR(100) NOT NULL,
    "product_categoty" VARCHAR(100) NOT NULL,
    "price" DECIMAL(10, 2) NOT NULL,
    "price_with_discount" DECIMAL(10, 2) NOT NULL,
    "rating" INT NOT NULL,
    "order_count" INT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

-- Favorite table migration
CREATE TABLE IF NOT EXISTS "favorite" (
     "product_id" UUID REFERENCES "product"("id"),
     "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     "updated_at" TIMESTAMP
);


-- Order_products table migration
CREATE TABLE IF NOT EXISTS "order_products" (
    "id" UUID PRIMARY KEY,
    "order_id" UUID REFERENCES "orders"("id"),
    "product_id" UUID REFERENCES "product"("id"),
    "quantity" INT NOT NULL,
    "price" DECIMAL(10, 2) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP
);

-- Banner table migration
CREATE TABLE IF NOT EXISTS "banner" (
    "id" UUID PRIMARY KEY,
    "banner_image" VARCHAR(255),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP
);
