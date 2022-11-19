-- execute this using database owner (i.e. wotsa, not postgres).

create schema if not exists wot;

grant usage on schema wot to wotrw;
grant temporary on database wot to wotrw;

alter default privileges
    in schema wot
    grant select, insert, update, delete on tables to wotrw;

alter default privileges
    in schema wot
    grant usage on sequences to wotrw;

alter default privileges
    in schema wot
    grant usage on types to wotrw;

alter default privileges
    in schema wot
    grant usage on types to wotrw;

alter default privileges
    in schema wot
    grant execute on routines to wotrw;

-- set search path for all _future_ connections
alter database wot set search_path to wot;

-- set search path in _this_ connection (needed for migrations to work properly).
set search_path to wot;