BEGIN;
CREATE TABLE event (
    id serial PRIMARY KEY,
    name varchar(300) not null,
    "from" timestamptz not null,
    "to" timestamptz not null,
    user_id integer not null,
    active bool not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz,
    FOREIGN KEY (user_id) REFERENCES "user"(id)
);
COMMIT;