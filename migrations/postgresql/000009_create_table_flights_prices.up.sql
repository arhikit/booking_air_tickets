CREATE TABLE flights_prices(
    id              uuid PRIMARY KEY,
    flight_id       uuid not null,
    class_seats_id  uuid not null,
    price_ticket    int not null,
    FOREIGN KEY (flight_id) REFERENCES flights (id) ON DELETE CASCADE,
    FOREIGN KEY (class_seats_id) REFERENCES classes_seats (id) ON DELETE CASCADE
    );