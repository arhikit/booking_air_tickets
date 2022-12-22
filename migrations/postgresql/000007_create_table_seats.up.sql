CREATE TABLE IF NOT EXISTS seats(
    id              uuid PRIMARY KEY,
    class_seats_id  uuid not null,
    number          varchar (20) not null,
    FOREIGN KEY (class_seats_id) REFERENCES classes_seats (id) ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_seats_class_seats_id ON seats(class_seats_id);