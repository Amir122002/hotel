-- консьерж, портье, горничная, техник, метрдотель, повар, официант, охранник, менеджер отеля.

create table workers
(
    id bigserial primary key,
    full_name text not null,
    login text unique not null ,
    password text not null ,
    job_title bigint not null references job(id),
    active boolean not null default true,
    create_ad timestamptz not null default current_timestamp,
    update_at timestamptz not null default current_timestamp,
    delete_at timestamptz
);

create table job
(
    id bigserial primary key,
    work text not null,
    active boolean not null default true,
    create_ad timestamptz not null default current_timestamp,
    update_at timestamptz not null default current_timestamp,
    delete_at timestamptz
);

create table working_hours
(
    id bigserial primary key,
    user_id bigint not null references workers(id),
    active boolean not null default true,
    start_work timestamptz not null default current_timestamp,
    finish_work timestamptz
);

create table workers_tokens
(
    id bigserial not null primary key,
    token text not null,
    user_id bigint not null references workers(id),
    active boolean not null default true,
    start_time timestamptz not null default current_timestamp,
    end_time timestamptz not null default current_timestamp + interval '1 hour'
);

create table clients
(
    id bigserial not null primary key ,
    fill_name text not null ,
    number bigint not null references hotel_rooms(id),
    active boolean not null default true,
    create_at timestamptz not null default current_timestamp,
    update_at timestamptz not null default current_timestamp,
    delete_at timestamptz
);


-- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --

create table menu
(
    id bigserial primary key,
    food text not null ,
    price bigint not null
);

create table orders
(
    id bigserial primary key ,
    number_order bigint not null ,
    food_id bigint not null references menu(id),
    table_id bigint not null references restaurant_tables(id),
    reservation bigint references clients(id)
);

create table restaurant_tables
(
    id bigserial primary key ,
    table_type bigint not null references restaurant_tables_type(id),
    number_table bigint not null ,
    active boolean not null default true

);

create table restaurant_tables_type
(
    id bigserial primary key ,
    name text not null

);

create table reservations
(
    id bigserial primary key ,
    number_room bigint not null references hotel_rooms(number_room),
    table_id bigint not null references restaurant_tables(id),
    time_of_reservation timestamptz not null
);

create table hotel_rooms
(
    id bigserial primary key ,
    number_room bigint not null ,
    hotel_room_types bigint references hotel_room_types(id),
    active boolean not null default true

);

create table hotel_room_types
(
    id bigserial primary key ,
    name text not null
);

create table requests
(
    id bigserial primary key ,
    client_id bigint not null references clients(id),
    service_sector bigint not null references service_sectors(id),
    date timestamptz not null default current_timestamp
);

create table service_sectors
(
    id bigserial primary key ,
    name_sector text

);

ALTER TABLE reservations
DROP COLUMN client_id;

ALTER TABLE orders
DROP COLUMN reservation;

ALTER TABLE clients
DROP COLUMN hotel_room_id;

ALTER TABLE clients
DROP COLUMN login;

ALTER TABLE clients
DROP COLUMN password;

ALTER TABLE requests
DROP COLUMN client_id;


ALTER TABLE hotel_rooms
DROP COLUMN client_id;
--
ALTER TABLE hotel_rooms
    ADD COLUMN client_id bigint REFERENCES clients(id);

ALTER TABLE requests
    ADD COLUMN hotel_room bigint REFERENCES hotel_rooms(id);



drop table clients_tokens;