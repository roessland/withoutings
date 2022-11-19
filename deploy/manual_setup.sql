-- Execute this using superuser (e.g. postgres).

-- create "superadmin" user for migrations and admin work.
create role wotsa
    password '<mypass>'
    login;

-- create "readwrite" user for usage by webapp and worker.
create role wotrw
    password '<otherpass>'
    login;

-- to avoid "Wot" prefix ending up in generated SQLC repos.
alter role wotrw set search_path = 'wot';
alter role wotsa set search_path = 'wot';

-- PostgreSQL 15
-- create database wot
--     owner wotsa
--     template template0
--     encoding 'utf8'
--     locale 'en_US'
--     lc_collate = 'C'
--     icu_locale 'en_US_POSIX'
--     locale_provider icu;

-- PostgreSQL 14
create database wot
    owner wotsa
    template template0
    encoding 'utf8'
    lc_collate = 'C';

-- golang-migrate needs schema to exist when creating schema migrations table.
create schema if not exists wot;

-- set search path for all _future_ connections
alter database wot set search_path to wot;

grant temporary on database wot to wotrw;