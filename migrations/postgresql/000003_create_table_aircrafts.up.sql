CREATE TABLE aircrafts(
    id          uuid PRIMARY KEY,
    airline_id  uuid not null,
    name        varchar (100) not null,
    FOREIGN KEY (airline_id) REFERENCES airlines (id) ON DELETE CASCADE
    );