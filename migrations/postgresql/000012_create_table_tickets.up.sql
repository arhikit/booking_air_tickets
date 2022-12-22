CREATE TABLE IF NOT EXISTS tickets(
    id                          uuid PRIMARY KEY,
    status_id                   int not null,
    status_timestamp            timestamptz not null,
    flight_id                   uuid not null,
    user_id                     uuid not null,
    passenger_id                uuid not null,
    class_seats_id              uuid not null,
    seat_id                     uuid,
    count_additional_baggage    int not null,
    price                       int not null,
    paid_with_bonuses           int not null,
    accrued_bonuses             int not null,
    FOREIGN KEY (status_id) REFERENCES statuses (id) ON DELETE CASCADE,
    FOREIGN KEY (flight_id) REFERENCES flights (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (passenger_id) REFERENCES passengers (id) ON DELETE CASCADE,
    FOREIGN KEY (class_seats_id) REFERENCES classes_seats (id) ON DELETE CASCADE,
    FOREIGN KEY (seat_id) REFERENCES seats (id) ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_tickets_status_id ON tickets(status_id);
CREATE INDEX IF NOT EXISTS idx_tickets_flight_id ON tickets(flight_id);
CREATE INDEX IF NOT EXISTS idx_tickets_user_id ON tickets(user_id);
CREATE INDEX IF NOT EXISTS idx_tickets_passenger_id ON tickets(passenger_id);
CREATE INDEX IF NOT EXISTS idx_tickets_class_seats_id ON tickets(class_seats_id);
