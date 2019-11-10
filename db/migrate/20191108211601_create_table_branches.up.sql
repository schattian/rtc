CREATE TABLE IF NOT EXISTS branches (
   id integer PRIMARY KEY,
   name varchar(30) UNIQUE NOT NULL,
   index_id integer UNIQUE NOT NULL
   -- FOREIGN KEY index_id REFERENCES indeces(index_id)
);
