-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "citext" WITH SCHEMA "public";
CREATE EXTENSION IF NOT EXISTS "pg_trgm" WITH SCHEMA "public";
CREATE EXTENSION IF NOT EXISTS "postgis" WITH SCHEMA "public";

CREATE TYPE "organization_type_enum" AS ENUM ('airline', 'carrier', 'warehouse');

-- PMC Table
CREATE TABLE IF NOT EXISTS "pmc_inventories" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "uld_number" TEXT NOT NULL,
    "pmc_type" TEXT NOT NULL,
    "pmc_status" TEXT NOT NULL,
    "current_location_id" "organization_type_enum" NOT NULL,
    "current_location_type" TEXT NOT NULL,

    CONSTRAINT "pmc_pkey" PRIMARY KEY ("id")  
);

-- Manifest Table
CREATE TABLE IF NOT EXISTS "delivery_manifests" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "manifest_date" TIMESTAMPTZ NOT NULL,
    "warehouse_id" TEXT NOT NULL,
    "airline_id" TEXT NOT NULL,
    "carrier_id" TEXT NOT NULL,
    "signature_info" TEXT NOT NULL,
    "manifest_status" TEXT NOT NULL,
    "created_by" TEXT NOT NULL,

    CONSTRAINT "manifest_pkey" PRIMARY KEY ("id")

);

-- Manifest Items Table
CREATE TABLE IF NOT EXISTS "manifest_items" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "manifest_id" TEXT NOT NULL,
    "pmc_inventory_id" TEXT NOT NULL,
    "addition_info" TEXT NOT NULL,

    CONSTRAINT "manifest_item_pkey" PRIMARY KEY ("id")
);

-- Warehouse Table
CREATE TABLE IF NOT EXISTS "warehouses" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "name" TEXT NOT NULL,
    "address" TEXT NOT NULL,
    "contact_info" TEXT NOT NULL,

    CONSTRAINT "warehouse_pkey" PRIMARY KEY ("id")
);

-- Airline Table
CREATE TABLE IF NOT EXISTS "airlines" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "name" TEXT NOT NULL,
    "address" TEXT NOT NULL,
    "contact_info" TEXT NOT NULL,

    CONSTRAINT "airline_pkey" PRIMARY KEY ("id")
);

-- Carrier Table
CREATE TABLE IF NOT EXISTS "carriers" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "name" TEXT NOT NULL,
    "address" TEXT NOT NULL,
    "contact_info" TEXT NOT NULL,

    CONSTRAINT "carrier_pkey" PRIMARY KEY ("id")
);

-- User Table
CREATE TABLE IF NOT EXISTS "users" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "user_name" CITEXT NOT NULL,
    "email" CITEXT NOT NULL,
    "password" TEXT NOT NULL,
    "last_request" TIMESTAMPTZ NOT NULL,
    "is_admin" BOOLEAN NOT NULL DEFAULT FALSE,
    "is_verified" BOOLEAN NOT NULL DEFAULT FALSE,
    "is_deleted" BOOLEAN NOT NULL DEFAULT FALSE,

    CONSTRAINT "user_pkey" PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "user_user_name_key" ON "users"("user_name");
CREATE UNIQUE INDEX "user_email_key" ON "users"("email");

CREATE TYPE "permissions_enum" AS ENUM ('pmc.read', 'pmc.write', 'manifest.read', 'manifest.write', 'user.read', 'user.write', 'organization.read', 'organization.write');

-- User Associations Table
CREATE TABLE IF NOT EXISTS "user_associations" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "user_id" TEXT NOT NULL,
    "organization_id" TEXT NOT NULL,
    "organization_type" "organization_type_enum" NOT NULL,
    "permissions" "permissions_enum"[],

    CONSTRAINT "user_association_pkey" PRIMARY KEY ("id")

);

ALTER TABLE "delivery_manifests" ADD CONSTRAINT "fk_warehouse" FOREIGN KEY ("warehouse_id") REFERENCES "warehouses"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "delivery_manifests" ADD CONSTRAINT "fk_airline" FOREIGN KEY ("airline_id") REFERENCES "airlines"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "delivery_manifests" ADD CONSTRAINT "fk_carrier" FOREIGN KEY ("carrier_id") REFERENCES "carriers"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "delivery_manifests" ADD CONSTRAINT "fk_created_by" FOREIGN KEY ("created_by") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "manifest_items" ADD CONSTRAINT "fk_manifest" FOREIGN KEY ("manifest_id") REFERENCES "delivery_manifests"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "manifest_items" ADD CONSTRAINT "fk_pmc_inventory" FOREIGN KEY ("pmc_inventory_id") REFERENCES "pmc_inventories"("id") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "user_associations" ADD CONSTRAINT "fk_user" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "user_associations" ADD CONSTRAINT "unique_user_org" UNIQUE ("user_id", "organization_id", "organization_type");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP EXTENSION IF EXISTS "citext" CASCADE;
DROP EXTENSION IF EXISTS "pg_trgm" CASCADE;
DROP EXTENSION IF EXISTS "postgis" CASCADE;

DROP TABLE IF EXISTS "manifest_items";
DROP TABLE IF EXISTS "delivery_manifests";
DROP TABLE IF EXISTS "warehouses";
DROP TABLE IF EXISTS "airlines";
DROP TABLE IF EXISTS "carriers";
DROP TABLE IF EXISTS "user_associations";
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "pmc_inventories";

DROP TYPE IF EXISTS "organization_type_enum";
DROP TYPE IF EXISTS "permissions_enum";
DROP INDEX IF EXISTS "user_email_key";
DROP INDEX IF EXISTS "user_user_name_key";

-- +goose StatementEnd
