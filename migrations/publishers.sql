create table publishers
(
    id       bigserial
        constraint publishers_pk
            primary key,
    name     varchar(255),
    country  varchar(100),
    steam_id integer
        constraint pubsteam_id_unique
            unique
);

alter table publishers
    owner to postgres;

