create table if not exists users (
    id bigserial primary key,
    email varchar(255) not null unique,
    password varchar(128) not null
);