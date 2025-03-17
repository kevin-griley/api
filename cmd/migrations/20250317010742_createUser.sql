-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "citext" WITH SCHEMA "public";
CREATE EXTENSION IF NOT EXISTS "pg_trgm" WITH SCHEMA "public";
CREATE EXTENSION IF NOT EXISTS "postgis" WITH SCHEMA "public";

CREATE TABLE IF NOT EXISTS "User" (
    "id" TEXT NOT NULL DEFAULT gen_random_uuid(),
    "createdAt" TIMESTAMPTZ NOT NULL,
    "updatedAt" TIMESTAMPTZ NOT NULL,
    "userName" CITEXT,
    "email" CITEXT NOT NULL,
    "password" TEXT NOT NULL,
    "lastRequest" TIMESTAMPTZ NOT NULL,
    "isAdmin" BOOLEAN NOT NULL DEFAULT FALSE,
    "isVerified" BOOLEAN NOT NULL DEFAULT FALSE,
    "isDisabled" BOOLEAN NOT NULL DEFAULT FALSE,
    "isDeleted" BOOLEAN NOT NULL DEFAULT FALSE,

    CONSTRAINT "User_pkey" PRIMARY KEY ("id")

);

CREATE UNIQUE INDEX "User_email_key" ON "User"("email");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP EXTENSION IF EXISTS "citext" CASCADE;
DROP EXTENSION IF EXISTS "pg_trgm" CASCADE;
DROP EXTENSION IF EXISTS "postgis" CASCADE;

DROP TABLE IF EXISTS "User";
-- +goose StatementEnd
