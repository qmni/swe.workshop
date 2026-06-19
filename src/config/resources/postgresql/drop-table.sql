-- Aufruf:
-- docker compose exec db bash
-- psql --dbname=player --username=player --file=/sql/drop-table.sql

SET search_path TO public;

DROP TABLE IF EXISTS player CASCADE;
DROP TABLE IF EXISTS guild CASCADE;

DROP TYPE IF EXISTS "PlayerStatus" CASCADE;
DROP TYPE IF EXISTS "PlayerClass" CASCADE;
