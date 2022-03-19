CREATE OR REPLACE FUNCTION update_time_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.update_time = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE IF NOT EXISTS addresses (
        id serial NOT NULL PRIMARY KEY,
        city text NOT NULL,
        street text NOT NULL,
        house text NOT NULL
);

CREATE TABLE IF NOT EXISTS pharmacies (
        id serial NOT NULL PRIMARY KEY,
        address_id int UNIQUE NOT NULL REFERENCES addresses(id),
        name text UNIQUE NOT NULL,
        is_blocked bool NOT NULL DEFAULT false,
        creation_time timestamp NOT NULL DEFAULT NOW(),
        update_time timestamp NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_pharmacy_time BEFORE UPDATE
    ON pharmacies FOR EACH ROW EXECUTE PROCEDURE
    update_time_column();

CREATE TABLE IF NOT EXISTS complaints (
        id serial NOT NULL PRIMARY KEY,
        name text NOT NULL DEFAULT '',
        email text NOT NULL DEFAULT '',
        complaint text NOT NULL DEFAULT '',
        worker_name text NOT NULL DEFAULT '',
        pharmacy_id int REFERENCES pharmacies(id),
        creation_time timestamp NOT NULL DEFAULT NOW(),
        update_time timestamp NOT NULL DEFAULT NOW()
);

CREATE TRIGGER update_product_item_update_time BEFORE UPDATE
    ON complaints FOR EACH ROW EXECUTE PROCEDURE
    update_time_column();

