CREATE TABLE users_balance(
    id                      uuid PRIMARY KEY,
    user_id                 uuid not null,
    sum_purchases           int not null,
    sum_bonuses             int not null,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );