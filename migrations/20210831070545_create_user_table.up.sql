create table if not exists users (
    id serial primary key,
    email varchar(255) unique not null,
    username varchar(32) unique,
    password varchar(128) not null,
    first_name varchar(50),
    last_name varchar(50),
    created_at timestamp with time zone default current_timestamp not null,
    updated_at timestamp with time zone default current_timestamp not null,
    is_active boolean default true not null,
    is_archive boolean default false not null,
    is_superuser boolean default false not null
);