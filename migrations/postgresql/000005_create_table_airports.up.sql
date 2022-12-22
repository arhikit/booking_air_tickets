CREATE TABLE IF NOT EXISTS airports(
    id          uuid PRIMARY KEY,
    city_id     uuid not null,
    name        varchar (300) not null,
    FOREIGN KEY (city_id) REFERENCES cities (id) ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_airports_city_id ON airports(city_id);