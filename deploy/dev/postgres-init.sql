-- Bootstraps the local dev Postgres container with the roles and
-- database that env.dev.sh / `withoutings migrate` expect.
-- Runs once, the first time the container starts on a fresh volume.

CREATE ROLE wotsa LOGIN PASSWORD 'wotsa';
CREATE ROLE wotrw LOGIN PASSWORD 'wotrw';

CREATE DATABASE wot OWNER wotsa TEMPLATE template0 ENCODING 'utf8' LC_COLLATE 'C';
ALTER DATABASE wot SET search_path TO wot;
GRANT TEMPORARY ON DATABASE wot TO wotrw;

\c wot
CREATE SCHEMA IF NOT EXISTS wot AUTHORIZATION wotsa;
