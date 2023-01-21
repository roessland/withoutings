-- execute this using database owner (i.e. wotsa, not postgres).

create schema if not exists wot;

grant usage on schema wot to wotrw;

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