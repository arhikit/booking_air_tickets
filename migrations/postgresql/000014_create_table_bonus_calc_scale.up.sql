CREATE TABLE bonus_calc_scale(
    id                      SERIAL PRIMARY KEY,
    sum_purchases_to        int not null,
    sum_purchases_from      int not null,
    percent                 int not null
    );

INSERT INTO bonus_calc_scale(sum_purchases_to, sum_purchases_from, percent)
        VALUES (0, 9999, 0), (10000, 49999, 3), (50000, 0, 5);