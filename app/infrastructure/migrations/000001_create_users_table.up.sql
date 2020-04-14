BEGIN;
CREATE TABLE "user" (
    id serial PRIMARY KEY,
    first_name varchar(50) not null,
    last_name varchar(50) not null,
    active bool not null default true,
    created_at timestamptz not null default now(),
    updated_at timestamptz
);
COMMIT;