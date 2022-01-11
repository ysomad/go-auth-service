create table if not exists accounts(
    id uuid primary key default gen_random_uuid(),
    email varchar(255) unique not null,
    username varchar(16) unique not null,
    password varchar(128) not null,
    created_at timestamp with time zone default current_timestamp not null,
    updated_at timestamp with time zone default current_timestamp not null,
    is_archive boolean default false not null
);
