CREATE TABLE flights(
    id                          uuid PRIMARY KEY,
    name                        varchar(100) not null,
    aircraft_id                 uuid not null,
    departure_airport_id        uuid not null,
    arrival_airport_id          uuid not null,
    departure_date              timestamptz not null,
    duration                    int not null,
    price_additional_baggage    int not null,
    price_seat_selection        int not null,
    is_international            bool not null,
    baggage_included            bool not null,
    pet_allowed                 bool not null,
    FOREIGN KEY (aircraft_id) REFERENCES aircrafts (id) ON DELETE CASCADE,
    FOREIGN KEY (departure_airport_id) REFERENCES airports (id) ON DELETE CASCADE,
    FOREIGN KEY (arrival_airport_id) REFERENCES airports (id) ON DELETE CASCADE
    );

CREATE INDEX idx_flights_airports_date ON flights(departure_airport_id, arrival_airport_id, departure_date DESC);