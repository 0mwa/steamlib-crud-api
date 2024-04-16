create table developers
(
    id       bigserial
        constraint developers_pk
            primary key,
    name     varchar(255),
    country  varchar(100),
    steam_id integer
        constraint devsteam_id_unique
            unique
);

alter table developers
    owner to postgres;

