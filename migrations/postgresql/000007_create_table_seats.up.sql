CREATE TABLE seats(
    id              uuid PRIMARY KEY,
    class_seats_id  uuid not null,
    number          varchar (20) not null,
    FOREIGN KEY (class_seats_id) REFERENCES classes_seats (id) ON DELETE CASCADE
    );