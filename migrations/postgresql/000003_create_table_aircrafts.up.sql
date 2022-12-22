CREATE TABLE IF NOT EXISTS aircrafts(
    id          uuid PRIMARY KEY,
    airline_id  uuid not null,
    name        varchar (100) not null,
    FOREIGN KEY (airline_id) REFERENCES airlines (id) ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_aircrafts_airline_id ON aircrafts(airline_id);