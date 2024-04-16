create table games
(
    id           bigserial
        constraint games_pk
            primary key,
    name         varchar(255),
    img          varchar(255),
    description  varchar(255),
    rating       integer,
    developer_id integer,
    publisher_id integer,
    steam_id     integer
        constraint steam_id_unique
            unique
);

alter table games
    owner to postgres;

