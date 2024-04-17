create sequence users_id_seq
    as integer;

alter sequence users_id_seq owner to postgres;

create table users
(
    id     integer default nextval('users_id_seq'::regclass) not null
        constraint users_pk
            primary key,
    login  varchar(255)                                      not null
        constraint users_login_uk
            unique,
    passwd varchar(255)                                      not null
);

alter table users
    owner to postgres;

