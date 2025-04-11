CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    role VARCHAR(255) NOT NULL
);

CREATE TABLE pvz (
    id UUID PRIMARY KEY,
    -- registrationDate TIMESTAMP NOT NULL DEFAULT now()
    registrationDate TIMESTAMPTZ NOT NULL,
    city TEXT NOT NULL CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань'))
);

CREATE TABLE receptions (
    id UUID PRIMARY KEY,
    date_time TIMESTAMPTZ NOT NULL,
    pvz_id UUID NOT NULL REFERENCES pvz(id),
    status TEXT NOT NULL CHECK (status IN ('in_progress', 'close'))
);

CREATE TABLE products (
    id UUID PRIMARY KEY,
    date_time TIMESTAMPTZ NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь')),
    reception_id UUID NOT NULL REFERENCES receptions(id) ON DELETE CASCADE
);

