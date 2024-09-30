CREATE TABLE IF NOT EXISTS "admin" (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL,
    "phone_number" VARCHAR(100) NOT NULL,
    "email" VARCHAR(100) NOT NULL,
    "password" VARCHAR(100) NOT NULL,
    "address" VARCHAR(255),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "customer" (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL,
    "surname" VARCHAR(255),
    "birthday" DATE, 
    "gender" VARCHAR(20),
    "phone_number" VARCHAR(100) NOT NULL,
    "password" VARCHAR(100) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "brand" (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL,
    "brand_image" VARCHAR(255),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "category" (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL,
    "url" VARCHAR(255) NOT NULL,
    "parent_id" UUID REFERENCES "category" ("id"),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP 
);



CREATE TYPE order_status AS ENUM ('yangi', 'tasdiqlandi', 'yetkazib berildi');
CREATE TYPE delivery_status AS ENUM ('kuryer', 'pochta');
CREATE TYPE payment_method AS ENUM ('naxt')
CREATE TYPE payment_status AS ENUM ('to`langan', 'kutilmoqda')

 
CREATE TABLE IF NOT EXISTS "orders" (
    "id" UUID PRIMARY KEY,
    "total_price" DECIMAL(10, 2) NOT NULL,  
    "status" order_status DEFAULT 'yangi',  
    "longtitude" DECIMAL(9, 6) NOT NULL,
    "latitude" DECIMAL(9, 6) NOT NULL,
    "address_name" VARCHAR(255) NOT NULL, 
    "delivery_status" VARCHAR(50) NOT NULL, 
    "delivery_cost" DECIMAL(10, 2) DEFAULT 0,  
    "payment_method" VARCHAR(50) NOT NULL,  
    "payment_status" VARCHAR(50) NOT NULL,  
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Buyurtma yaratilgan vaqt
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Buyurtma yangilangan vaqt
    "customer_id" UUID,  -- Foydalanuvchi jadvaliga bog'lanadi
    FOREIGN KEY ("customer_id") REFERENCES "customer"("id")  -- Foydalanuvchi jadvaliga bog'lanadi
);


CREATE TABLE IF NOT EXISTS "order_items" (
    "id" UUID PRIMARY KEY,
    "quantity" INT NOT NULL,
    "price" DECIMAL(10, 2) NOT NULL,  -- Mahsulotning bitta narxi
    "total" DECIMAL(10, 2),  -- Jami narx
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Mahsulot qo'shilgan vaqt
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Ma'lumot yangilangan vaqt
    "order_id" UUID REFERENCES "orders"("id"),  -- Buyurtma bilan bog'lanadi
    "product_id" UUID REFERENCES "product"("id"),  -- Mahsulot bilan bog'lanadi
    "color_id" UUID REFERENCES "color"("id")  -- Rang bilan bog'lanadi
);


CREATE TYPE product_status AS ENUM ('novinka', 'rasprodaja', 'vremennaya_skidka', '');
    "status" product_status,

CREATE TABLE IF NOT EXISTS "product" (
    "id" UUID PRIMARY KEY,
    "category_id" UUID REFERENCES "category"("id"),
    "brand_id" UUID REFERENCES "brand"("id"),
    "image" TEXT,
    "favorite" BOOLEAN,
    "name" VARCHAR(100) NOT NULL,s
    "price" DECIMAL(10, 2) NOT NULL,
    "with_discount" DECIMAL(10, 2),
    "rating" FLOAT NOT NULL,
    "status" product_status,
    "description" VARCHAR(1000),
    "order_count" INT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

ALTER TABLE "product" ADD COLUMN "brand_id" UUID REFERENCES "brand"("id");

CREATE TABLE IF NOT EXISTS "banner" (
    "id" UUID PRIMARY KEY,
    "banner_image" VARCHAR(255),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "color" (
    "id" UUID PRIMARY KEY,
    "product_id" UUID REFERENCES "product"("id"),
    "color_name" VARCHAR(100) NOT NULL,  
    "color_url" TEXT[], 
    "count" INT DEFAULT 0, 
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);


CREATE TABLE IF NOT EXISTS "location" (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "info" TEXT,
    "latitude" DECIMAL(9, 6) NOT NULL,
    "longitude" DECIMAL(9, 6) NOT NULL,
    "image" TEXT,
    "opens_at" VARCHAR(255),
    "closes_at"VARCHAR(255),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);


CREATE TYPE delivery_status AS ENUM ('kuryer', 'pochta');

CREATE TABLE IF NOT EXISTS "shipping_details" (
    "id" UUID PRIMARY KEY,
    "order_id" UUID REFERENCES "orders"("id"),  -- Buyurtma ID'si
    "delivery_status" VARCHAR(50) NOT NULL,  -- Yetkazib berish turi (masalan, 'kuryer', 'pochta')
    "delivery_cost" DECIMAL(10, 2) DEFAULT 0,  -- Yetkazib berish narxi
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);


-- CREATE TABLE IF NOT EXISTS payment_details (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     order_id UUID NOT NULL,  -- Buyurtma ID'si
--     payment_method VARCHAR(50) NOT NULL,  -- To'lov turi (masalan, 'kartochka', 'naqd')
--     payment_status VARCHAR(50) NOT NULL,  -- To'lov holati (masalan, 'to'langan', 'kutilmoqda')
--     payment_amount DECIMAL(10, 2) NOT NULL,  -- To'langan summa
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     FOREIGN KEY (order_id) REFERENCES orders(id)  -- Buyurtma bilan bog'lanadi
-- );

-- CREATE TABLE IF NOT EXISTS order_status_history (
--     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
--     order_id UUID NOT NULL,  -- Buyurtma ID'si
--     status VARCHAR(50) NOT NULL,  -- Buyurtma holati
--     changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Holat o'zgargan vaqt
--     FOREIGN KEY (order_id) REFERENCES orders(id)  -- Buyurtma bilan bog'lanadi
-- );




