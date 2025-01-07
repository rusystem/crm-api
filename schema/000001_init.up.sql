-- Создаем последовательности
CREATE SEQUENCE warehouses_id_seq;
CREATE SEQUENCE planning_materials_id_seq;
CREATE SEQUENCE item_id_seq;
CREATE SEQUENCE purchased_materials_id_seq;
CREATE SEQUENCE planning_materials_archive_id_seq;
CREATE SEQUENCE purchased_materials_archive_id_seq;
CREATE SEQUENCE suppliers_id_seq;
CREATE SEQUENCE users_id_seq;
CREATE SEQUENCE companies_id_seq;
CREATE SEQUENCE refresh_tokens_id_seq;
CREATE SEQUENCE sections_id_seq;
CREATE SEQUENCE material_categories_id_seq;
CREATE
EXTENSION IF NOT EXISTS pgcrypto;

-- Создаем таблицы
CREATE TABLE "warehouses"
(
    "id"                 INT PRIMARY KEY DEFAULT nextval('warehouses_id_seq'),
    "name"               VARCHAR(255),
    "address"            VARCHAR(255),
    "responsible_person" VARCHAR(255),
    "phone"              VARCHAR(50),
    "email"              VARCHAR(255),
    "max_capacity"       INT,
    "current_occupancy"  INT,
    "other_fields"       JSONB,
    "country"            TEXT,
    "region"             TEXT,
    "comments"           TEXT,
    "created_at"         TIMESTAMP       DEFAULT (CURRENT_TIMESTAMP),
    "company_id"         INT
);

CREATE TABLE "planning_materials"
(
    "id"                       INT PRIMARY KEY DEFAULT nextval('planning_materials_id_seq'),
    "warehouse_id"             INT,
    "item_id"                  INT,
    "name"                     VARCHAR(255),
    "by_invoice"               VARCHAR(255),
    "article"                  VARCHAR(255),
    "product_category"         VARCHAR(255),
    "unit"                     VARCHAR(50),
    "total_quantity"           INT,
    "volume"                   INT,
    "price_without_vat"        DECIMAL,
    "total_without_vat"        DECIMAL,
    "supplier_id"              INT,
    "location"                 VARCHAR(255),
    "contract"                 DATE,
    "file"                     VARCHAR(255),
    "status"                   VARCHAR(255),
    "comments"                 TEXT,
    "reserve"                  VARCHAR(255),
    "received_date"            DATE,
    "last_updated"             DATE,
    "min_stock_level"          INT,
    "expiration_date"          DATE,
    "responsible_person"       VARCHAR(255),
    "storage_cost"             DECIMAL,
    "warehouse_section"        VARCHAR(255),
    "incoming_delivery_number" VARCHAR(255),
    "other_fields"             JSONB,
    "company_id"               INT
);

CREATE TABLE "purchased_materials"
(
    "id"                       INT PRIMARY KEY DEFAULT nextval('purchased_materials_id_seq'),
    "warehouse_id"             INT,
    "item_id"                  INT             DEFAULT nextval('item_id_seq'),
    "name"                     VARCHAR(255),
    "by_invoice"               VARCHAR(255),
    "article"                  VARCHAR(255),
    "product_category"         VARCHAR(255),
    "unit"                     VARCHAR(50),
    "total_quantity"           INT,
    "volume"                   INT,
    "price_without_vat"        DECIMAL,
    "total_without_vat"        DECIMAL,
    "supplier_id"              INT,
    "location"                 VARCHAR(255),
    "contract"                 DATE,
    "file"                     VARCHAR(255),
    "status"                   VARCHAR(255),
    "comments"                 TEXT,
    "reserve"                  VARCHAR(255),
    "received_date"            DATE,
    "last_updated"             DATE,
    "min_stock_level"          INT,
    "expiration_date"          DATE,
    "responsible_person"       VARCHAR(255),
    "storage_cost"             DECIMAL,
    "warehouse_section"        VARCHAR(255),
    "incoming_delivery_number" VARCHAR(255),
    "other_fields"             JSONB,
    "company_id"               INT
);

CREATE TABLE "planning_materials_archive"
(
    "id"                       INT PRIMARY KEY DEFAULT nextval('planning_materials_archive_id_seq'),
    "warehouse_id"             INT,
    "item_id"                  INT,
    "name"                     VARCHAR(255),
    "by_invoice"               VARCHAR(255),
    "article"                  VARCHAR(255),
    "product_category"         VARCHAR(255),
    "unit"                     VARCHAR(50),
    "total_quantity"           INT,
    "volume"                   INT,
    "price_without_vat"        DECIMAL,
    "total_without_vat"        DECIMAL,
    "supplier_id"              INT,
    "location"                 VARCHAR(255),
    "contract"                 DATE,
    "file"                     VARCHAR(255),
    "status"                   VARCHAR(255),
    "comments"                 TEXT,
    "reserve"                  VARCHAR(255),
    "received_date"            DATE,
    "last_updated"             DATE,
    "min_stock_level"          INT,
    "expiration_date"          DATE,
    "responsible_person"       VARCHAR(255),
    "storage_cost"             DECIMAL,
    "warehouse_section"        VARCHAR(255),
    "incoming_delivery_number" VARCHAR(255),
    "other_fields"             JSONB,
    "company_id"               INT
);

CREATE TABLE "purchased_materials_archive"
(
    "id"                       INT PRIMARY KEY DEFAULT nextval('purchased_materials_archive_id_seq'),
    "warehouse_id"             INT,
    "item_id"                  INT,
    "name"                     VARCHAR(255),
    "by_invoice"               VARCHAR(255),
    "article"                  VARCHAR(255),
    "product_category"         VARCHAR(255),
    "unit"                     VARCHAR(50),
    "total_quantity"           INT,
    "volume"                   INT,
    "price_without_vat"        DECIMAL,
    "total_without_vat"        DECIMAL,
    "supplier_id"              INT,
    "location"                 VARCHAR(255),
    "contract"                 DATE,
    "file"                     VARCHAR(255),
    "status"                   VARCHAR(255),
    "comments"                 TEXT,
    "reserve"                  VARCHAR(255),
    "received_date"            DATE,
    "last_updated"             DATE,
    "min_stock_level"          INT,
    "expiration_date"          DATE,
    "responsible_person"       VARCHAR(255),
    "storage_cost"             DECIMAL,
    "warehouse_section"        VARCHAR(255),
    "incoming_delivery_number" VARCHAR(255),
    "other_fields"             JSONB,
    "company_id"               INT
);

CREATE TABLE "suppliers"
(
    "id"                 INT PRIMARY KEY DEFAULT nextval('suppliers_id_seq'),
    "name"               VARCHAR(255),
    "legal_address"      VARCHAR(255),
    "actual_address"     VARCHAR(255),
    "warehouse_address"  VARCHAR(255),
    "contact_person"     VARCHAR(255),
    "phone"              VARCHAR(50),
    "email"              VARCHAR(255),
    "website"            VARCHAR(255),
    "contract_number"    VARCHAR(255),
    "product_categories" VARCHAR(255),
    "purchase_amount"    DECIMAL,
    "balance"            DECIMAL,
    "product_types"      INT,
    "comments"           TEXT,
    "files"              TEXT,
    "country"            VARCHAR(255),
    "region"             VARCHAR(255),
    "tax_id"             VARCHAR(255),
    "bank_details"       TEXT,
    "registration_date"  DATE,
    "payment_terms"      TEXT,
    "is_active"          BOOLEAN,
    "other_fields"       JSONB,
    "company_id"         INT
);

CREATE TABLE "users"
(
    "id"                          INT PRIMARY KEY DEFAULT nextval('users_id_seq'),
    "company_id"                  INT,
    "username"                    VARCHAR(50) UNIQUE  NOT NULL,
    "name"                        VARCHAR(50)         NOT NULL,
    "email"                       VARCHAR(100) UNIQUE NOT NULL,
    "phone"                       VARCHAR(255),
    "password_hash"               TEXT                NOT NULL,
    "created_at"                  TIMESTAMP       DEFAULT (CURRENT_TIMESTAMP),
    "updated_at"                  TIMESTAMP       DEFAULT (CURRENT_TIMESTAMP),
    "last_login"                  TIMESTAMP,
    "is_active"                   BOOLEAN         DEFAULT true,
    "role"                        VARCHAR(50)     DEFAULT 'user',
    "language"                    VARCHAR(50)     DEFAULT 'en',
    "country"                     VARCHAR(255),
    "is_approved"                 BOOLEAN,
    "is_send_system_notification" BOOLEAN,
    "sections"                    JSONB,
    "position"                    VARCHAR(255)
);

CREATE TABLE "companies"
(
    "id"          INT PRIMARY KEY DEFAULT nextval('companies_id_seq'),
    "name_ru"     VARCHAR(255),
    "name_en"     VARCHAR(255),
    "country"     VARCHAR(100),
    "address"     VARCHAR(255),
    "phone"       VARCHAR(50),
    "email"       VARCHAR(100),
    "website"     VARCHAR(100),
    "is_active"   BOOLEAN         DEFAULT true,
    "created_at"  TIMESTAMP       DEFAULT (CURRENT_TIMESTAMP),
    "updated_at"  TIMESTAMP       DEFAULT (CURRENT_TIMESTAMP),
    "is_approved" BOOLEAN,
    "timezone"    VARCHAR(255)
);

CREATE TABLE "refresh_tokens"
(
    "id"         INT PRIMARY KEY DEFAULT nextval('refresh_tokens_id_seq'),
    "user_id"    INT          NOT NULL,
    "company_id" INT,
    "roles"      VARCHAR(255) NOT NULL,
    "token"      VARCHAR(255) NOT NULL,
    "expires_at" TIMESTAMP    NOT NULL,
    "ip"         VARCHAR(255)
);

CREATE TABLE "sections"
(
    "id"   INT PRIMARY KEY DEFAULT nextval('sections_id_seq'),
    "name" VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE "material_categories"
(
    "id"          INT PRIMARY KEY DEFAULT nextval('material_categories_id_seq'),
    "name"        VARCHAR(255) NOT NULL,
    "company_id"  INT,
    "description" TEXT,
    "slug"        VARCHAR(255),
    "created_at"  TIMESTAMP       DEFAULT (CURRENT_TIMESTAMP),
    "updated_at"  TIMESTAMP       DEFAULT (CURRENT_TIMESTAMP),
    "is_active"   BOOLEAN         DEFAULT true,
    "img_url"     VARCHAR(255)
);

ALTER TABLE "planning_materials"
    ADD FOREIGN KEY ("warehouse_id") REFERENCES "warehouses" ("id");

ALTER TABLE "purchased_materials"
    ADD FOREIGN KEY ("warehouse_id") REFERENCES "warehouses" ("id");

ALTER TABLE "planning_materials_archive"
    ADD FOREIGN KEY ("warehouse_id") REFERENCES "warehouses" ("id");

ALTER TABLE "purchased_materials_archive"
    ADD FOREIGN KEY ("warehouse_id") REFERENCES "warehouses" ("id");

ALTER TABLE "planning_materials"
    ADD FOREIGN KEY ("supplier_id") REFERENCES "suppliers" ("id");

ALTER TABLE "purchased_materials"
    ADD FOREIGN KEY ("supplier_id") REFERENCES "suppliers" ("id");

ALTER TABLE "planning_materials_archive"
    ADD FOREIGN KEY ("supplier_id") REFERENCES "suppliers" ("id");

ALTER TABLE "purchased_materials_archive"
    ADD FOREIGN KEY ("supplier_id") REFERENCES "suppliers" ("id");

ALTER TABLE "users"
    ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "warehouses"
    ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "suppliers"
    ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "planning_materials"
    ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "purchased_materials"
    ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "planning_materials_archive"
    ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "purchased_materials_archive"
    ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "material_categories"
    ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

CREATE INDEX idx_suppliers_name ON "suppliers" ("name");

CREATE INDEX idx_materials_categories_name_lower ON material_categories (LOWER(name));
CREATE INDEX idx_materials_categories_company_id ON "material_categories" ("company_id");

CREATE INDEX idx_planning_materials_name ON planning_materials (LOWER(name));
CREATE INDEX idx_purchased_materials_name ON purchased_materials (LOWER(name));
CREATE INDEX idx_planning_materials_archive_name ON planning_materials_archive (LOWER(name));
CREATE INDEX idx_purchased_materials_archive_name ON purchased_materials_archive (LOWER(name));

CREATE INDEX idx_planning_materials_company_id ON planning_materials (company_id);
CREATE INDEX idx_purchased_materials_company_id ON purchased_materials (company_id);
CREATE INDEX idx_planning_materials_archive_company_id ON planning_materials_archive (company_id);
CREATE INDEX idx_purchased_materials_archive_company_id ON purchased_materials_archive (company_id);

DO
$$
DECLARE
company_id INT;
BEGIN
INSERT INTO "companies"
("name_ru",
 "name_en",
 "country",
 "address",
 "phone",
 "email",
 "website",
 "is_active",
 "is_approved",
 "timezone")
VALUES ('admin',
        'admin',
        'Kazakhstan',
        '',
        '',
        'admin@example.com',
        '',
        true,
        true,
        'Asia/Almaty') RETURNING id
INTO company_id;

INSERT INTO "users"
("company_id",
 "username",
 "name",
 "email",
 "phone",
 "password_hash",
 "role",
 "language",
 "country",
 "is_active",
 "is_approved",
 "is_send_system_notification",
 "sections",
 "position")
VALUES (company_id,
        'admin',
        'admin_name',
        'admin@example.com',
        '',
        crypt('pomogatoradmin', gen_salt('bf')),
        'admin',
        'ru',
        'Kazakhstan',
        true,
        true,
        true,
        '[
          "full_all_access"
        ]'::jsonb,
        'admin');
END $$;

INSERT INTO "sections" ("name")
VALUES ('full_all_access'),
       ('full_company_access'),
       ('full_access'),
       ('order_card_access'),
       ('production_data_access'),
       ('status_and_calculate_access'),
       ('purchase_planning_access') ON CONFLICT ("name") DO NOTHING;

ALTER TABLE "suppliers"
    ADD COLUMN "contract_date" TIMESTAMP;

ALTER TABLE "warehouses"
    ADD COLUMN "locality" TEXT;