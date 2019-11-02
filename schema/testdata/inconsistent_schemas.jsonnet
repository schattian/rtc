local columns = import 'columns.jsonnet';
local schemas = import 'schemas.jsonnet';
local tables = import 'tables.jsonnet';

{
  basic: {
    name: schemas.basic.name,
    blueprint: [
      {
        name: tables.basic.name,
        columns: [{ name: columns.basic.name, type: "inconsistent" }],
        option_keys: tables.basic.option_keys,
      },
    ],
  },
}
