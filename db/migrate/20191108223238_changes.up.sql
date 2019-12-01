CREATE EXTENSION hstore;
CREATE TABLE IF NOT EXISTS changes (
   id int PRIMARY KEY,
   type varchar(10),
   change_ids int[],
   table_name varchar(50) NOT NULL,
   column_name varchar(50),
   str_value text,
   int_value bigint,
   float_32_value decimal,
   float_64_value decimal,
   bytes_value text,
   entity_id varchar(100),
   value_type varchar(10),
   options hstore
);