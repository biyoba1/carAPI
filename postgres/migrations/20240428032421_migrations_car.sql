-- +goose Up
CREATE TABLE IF NOT EXISTS car (
    id SERIAL PRIMARY KEY,
    reg_num VARCHAR(20) NOT NULL,
    mark VARCHAR(50) NOT NULL,
    model VARCHAR(50) NOT NULL,
    year VARCHAR(50) NOT NULL,
    owner_name VARCHAR(50) NOT NULL,
    owner_surname VARCHAR(50) NOT NULL,
    owner_patronymic VARCHAR(50)
    );

-- +goose Down
drop table car;
