CREATE TABLE IF NOT EXISTS users(
    id          uuid PRIMARY KEY,
    name        varchar (100) not null,
    email       varchar (300) unique not null,
    password    varchar (100) not null
    );