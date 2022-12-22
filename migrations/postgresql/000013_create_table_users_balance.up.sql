CREATE TABLE IF NOT EXISTS users_balance(
    id                      uuid PRIMARY KEY,
    user_id                 uuid not null,
    sum_purchases           int not null,
    sum_bonuses             int not null,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_users_balance_user_id ON users_balance(user_id);
