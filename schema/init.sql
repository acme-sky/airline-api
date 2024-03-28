CREATE TABLE airport (
    id SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    location VARCHAR(256) NOT NULL,
    latitude FLOAT,
    longitude FLOAT
);

CREATE TABLE flight (
    id SERIAL NOT NULL PRIMARY KEY,
    created_at timestamp NOT NULL,
    departaure_time timestamp NOT NULL,
    arrival_time timestamp NOT NULL,
    departaure_airport INT NOT NULL,
    arrival_airport INT NOT NULL,
    FOREIGN KEY(departaure_airport) REFERENCES airport(id),
    FOREIGN KEY(arrival_airport) REFERENCES airport(id)
);
