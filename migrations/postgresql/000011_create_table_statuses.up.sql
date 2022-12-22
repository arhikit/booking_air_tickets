CREATE TABLE IF NOT EXISTS statuses(
    id          SERIAL PRIMARY KEY,
    name        varchar (10) not null
    );

INSERT INTO statuses(name) VALUES ('Created'), ('Paid'), ('Canceled'), ('Refunded'), ('Registered'), ('Closed');