-- Item table
CREATE TABLE IF NOT EXISTS Item (
    id integer PRIMARY KEY,
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

/* DATE: 1999-01-08 04:05:06 America/Chicago */