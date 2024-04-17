create sequence sessions_id_seq
    as integer;

alter sequence sessions_id_seq owner to postgres;

create table sessions
(
    id         integer default nextval('sessions_id_seq'::regclass) not null
        constraint sessions_pk
            primary key,
    user_id    integer
        constraint sessions_users_id_fk
            references users,
    token      varchar(255)                                         not null
        constraint token_uk
            unique,
    expiration timestamp
);

alter table sessions
    owner to postgres;

