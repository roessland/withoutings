-- execute this using database owner (i.e. wotsa, not postgres).

create schema wot;

-- drop schema if exists public; -- optional

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

alter database wot set search_path to wot;