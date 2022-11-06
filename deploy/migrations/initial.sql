-- Execute this using superuser (e.g. postgres).

create role wotsa
    password '<mypass>'
    login;

create role wotrw
    password '<otherpass>'
    login;

create database wot
    owner wotsa
    template template0
    encoding 'utf8'
    locale 'en_US'
    lc_collate = 'C'
    icu_locale 'en_US_POSIX'
    locale_provider icu;
