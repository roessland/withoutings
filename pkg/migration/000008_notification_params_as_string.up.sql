alter table notification
    rename column params to params_json;

alter table notification
    add column params text not null default '',
    alter column params_json drop not null;
