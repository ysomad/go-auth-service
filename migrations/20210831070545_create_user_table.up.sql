create table if not exists users (
    id serial primary key,
    email varchar(255) unique not null,
    password varchar(128) not null,
    first_name varchar(50),
    last_name varchar(50),
    created_at timestamp default current_timestamp not null,
    is_active boolean default true not null,
    is_superuser boolean default false not null
);