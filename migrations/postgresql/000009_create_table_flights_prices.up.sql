CREATE TABLE IF NOT EXISTS flights_prices(
    id              uuid PRIMARY KEY,
    flight_id       uuid not null,
    class_seats_id  uuid not null,
    price_ticket    int not null,
    FOREIGN KEY (flight_id) REFERENCES flights (id) ON DELETE CASCADE,
    FOREIGN KEY (class_seats_id) REFERENCES classes_seats (id) ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_flights_prices_flight_id ON flights_prices(flight_id);
CREATE INDEX IF NOT EXISTS idx_flights_prices_class_seats_id ON flights_prices(class_seats_id);
