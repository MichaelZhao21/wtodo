-- Item table
CREATE TABLE IF NOT EXISTS Item (
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL,
    due timestamp with time zone,
    start timestamp with time zone,
    length smallint,
    priority smallint,
    finished boolean
);

-- Tag table
CREATE TABLE IF NOT EXISTS Tag (
    item_id integer PRIMARY KEY,
    name varchar(50)
);

-- Test insert into item
INSERT INTO Item
VALUES (DEFAULT, 'Test Item', '2022-04-20 23:59:00 America/Chicago', NULL, 0, 2, 'False');

-- Insert one item
INSERT INTO Item VALUES (DEFAULT, $1, $2, $3, $4, $5, $6);

-- Select all items
SELECT * FROM Item

/* DATE: 1999-01-08 04:05:06 America/Chicago */