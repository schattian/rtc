CREATE TABLE IF NOT EXISTS commits (
   id integer PRIMARY KEY,
   errored bool,
   change_ids integer[]
);