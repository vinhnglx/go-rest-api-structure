create table if not exists products (
    id serial primary key,
    name varchar(100) not null,
    description text,
    price numeric(10, 2) not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);