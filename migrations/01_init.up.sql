CREATE TABLE url (
    id uuid primary key,
    short_url varchar(8) unique not null,
    long_url varchar not null
)