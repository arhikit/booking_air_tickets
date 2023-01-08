CREATE TABLE classes_seats(
    id           uuid PRIMARY KEY,
    aircraft_id  uuid not null,
    name         varchar (300) not null,
    count_seats  int not null,
    width        int not null,
    pitch        int not null,
    count_in_row int not null,
    FOREIGN KEY (aircraft_id) REFERENCES aircrafts (id) ON DELETE CASCADE
    );