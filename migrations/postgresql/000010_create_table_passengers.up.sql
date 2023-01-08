CREATE TABLE passengers(
    id                      uuid PRIMARY KEY,
    user_id                 uuid not null,
    name_passenger          varchar (200) not null,
    identity_data_passenger varchar (500) not null,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );
