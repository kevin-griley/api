-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "citext" WITH SCHEMA "public";
CREATE EXTENSION IF NOT EXISTS "pg_trgm" WITH SCHEMA "public";
CREATE EXTENSION IF NOT EXISTS "postgis" WITH SCHEMA "public";

CREATE TYPE "organization_type_enum" AS ENUM ('airline', 'carrier', 'warehouse');
CREATE TYPE "uld_status_enum" AS ENUM ('delivered', 'in_transit', 'in_warehouse');
CREATE TYPE "uld_type_enum" AS ENUM ('PMC', 'AKE', 'LD3', 'LD7');

-- Organization Table
CREATE TABLE IF NOT EXISTS "organizations" (
    "id" UUID PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "name" TEXT NOT NULL,
    "unique_url" CITEXT UNIQUE,
    "address" TEXT NOT NULL,
    "contact_info" TEXT NOT NULL,
    "organization_type" "organization_type_enum" NOT NULL
);


-- PMC Table
CREATE TABLE IF NOT EXISTS "uld_inventories" (
    "id" UUID PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "uld_number" VARCHAR(20) UNIQUE,
    "uld_type" "uld_type_enum" NOT NULL,
    "uld_status" "uld_status_enum" NOT NULL,
    "current_location_id" UUID NOT NULL, -- !FK warehouse, airline, carrier
    "current_location_type" "organization_type_enum" NOT NULL,
    "organization_id" UUID NOT NULL -- FK organization
);


CREATE TYPE "manifest_status_enum" AS ENUM ('draft', 'submitted', 'accepted', 'rejected');

-- Manifest Table
CREATE TABLE IF NOT EXISTS "delivery_manifests" (
    "id" UUID PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "manifest_date" TIMESTAMPTZ NOT NULL,
    "warehouse_id" UUID NOT NULL, -- FK warehouse
    "airline_id" UUID NOT NULL, -- FK airline
    "carrier_id" UUID NOT NULL, -- FK carrier
    "signature_info" TEXT NOT NULL,
    "manifest_status" "manifest_status_enum" NOT NULL,
    "created_by" UUID NOT NULL, -- FK user
    "organization_id" UUID NOT NULL -- FK organization
);

-- Manifest Items Table
CREATE TABLE IF NOT EXISTS "manifest_items" (
    "id" UUID PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "manifest_id" UUID NOT NULL, -- FK manifest
    "uld_inventory_id" UUID NOT NULL, -- FK uld_inventory
    "addition_info" TEXT NOT NULL,
    "organization_id" UUID NOT NULL -- FK organization
);

-- Warehouse Table
CREATE TABLE IF NOT EXISTS "warehouses" (
    "id" UUID PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "name" TEXT NOT NULL,
    "address" TEXT NOT NULL,
    "contact_info" TEXT NOT NULL,
    "organization_id" UUID NOT NULL -- FK organization
);

-- Airline Table
CREATE TABLE IF NOT EXISTS "airlines" (
    "id" UUID PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "name" TEXT NOT NULL,
    "address" TEXT NOT NULL,
    "contact_info" TEXT NOT NULL,
    "organization_id" UUID NOT NULL -- FK organization
);

-- Carrier Table
CREATE TABLE IF NOT EXISTS "carriers" (
    "id" UUID PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "name" TEXT NOT NULL,
    "address" TEXT NOT NULL,
    "contact_info" TEXT NOT NULL,
    "organization_id" UUID NOT NULL -- FK organization
);

-- User Table
CREATE TABLE IF NOT EXISTS "users" (
    "id" UUID PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "user_name" CITEXT UNIQUE,
    "email" CITEXT UNIQUE,
    "hashed_password" TEXT NOT NULL,
    "is_admin" BOOLEAN NOT NULL DEFAULT FALSE,
    "is_verified" BOOLEAN NOT NULL DEFAULT FALSE,
    "is_deleted" BOOLEAN NOT NULL DEFAULT FALSE,
    "last_request" TIMESTAMPTZ NOT NULL,
    "last_login" TIMESTAMPTZ NOT NULL,
    "failed_login_attempts" INT NOT NULL DEFAULT 0
);

CREATE TYPE "organization_status" AS ENUM ('pending', 'active', 'inactive');
CREATE TYPE "permissions_enum" AS ENUM ('uld.read', 'uld.write', 'manifest.read', 'manifest.write', 'user.read', 'user.write', 'organization.read', 'organization.write');

-- User Associations Table
CREATE TABLE IF NOT EXISTS "user_associations" (
    "id" UUID PRIMARY KEY,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL,
    "status" "organization_status" NOT NULL,
    "permissions" "permissions_enum"[],
    "user_id" UUID NOT NULL, -- FK user
    "organization_id" UUID NOT NULL, -- FK organization
    
    CONSTRAINT "unique_user_org" UNIQUE ("user_id", "organization_id")
);

ALTER TABLE "uld_inventories" ADD CONSTRAINT "fk_uld_organization" FOREIGN KEY ("organization_id") REFERENCES "organizations"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "delivery_manifests" ADD CONSTRAINT "fk_warehouse" FOREIGN KEY ("warehouse_id") REFERENCES "warehouses"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "delivery_manifests" ADD CONSTRAINT "fk_airline" FOREIGN KEY ("airline_id") REFERENCES "airlines"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "delivery_manifests" ADD CONSTRAINT "fk_carrier" FOREIGN KEY ("carrier_id") REFERENCES "carriers"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "delivery_manifests" ADD CONSTRAINT "fk_created_by" FOREIGN KEY ("created_by") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "delivery_manifests" ADD CONSTRAINT "fk_delivery_manifest_organization" FOREIGN KEY ("organization_id") REFERENCES "organizations"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "manifest_items" ADD CONSTRAINT "fk_manifest" FOREIGN KEY ("manifest_id") REFERENCES "delivery_manifests"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "manifest_items" ADD CONSTRAINT "fk_uld_inventory" FOREIGN KEY ("uld_inventory_id") REFERENCES "uld_inventories"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "manifest_items" ADD CONSTRAINT "fk_manifest_item_organization" FOREIGN KEY ("organization_id") REFERENCES "organizations"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "warehouses" ADD CONSTRAINT "fk_warehouse_organization" FOREIGN KEY ("organization_id") REFERENCES "organizations"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "airlines" ADD CONSTRAINT "fk_airline_organization" FOREIGN KEY ("organization_id") REFERENCES "organizations"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "carriers" ADD CONSTRAINT "fk_carrier_organization" FOREIGN KEY ("organization_id") REFERENCES "organizations"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "user_associations" ADD CONSTRAINT "fk_user" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "user_associations" ADD CONSTRAINT "fk_user_association_organization" FOREIGN KEY ("organization_id") REFERENCES "organizations"("id") ON DELETE CASCADE ON UPDATE CASCADE;

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
DROP TABLE IF EXISTS "uld_inventories";
DROP TABLE IF EXISTS "organizations";

DROP TYPE IF EXISTS "permissions_enum";
DROP TYPE IF EXISTS "organization_type_enum";
DROP TYPE IF EXISTS "uld_status_enum";
DROP TYPE IF EXISTS "uld_type_enum";
DROP TYPE IF EXISTS "manifest_status_enum";
DROP TYPE IF EXISTS "organization_status";
-- +goose StatementEnd
