local columns = import 'columns.jsonnet';
local schemas = import 'schemas.jsonnet';
local tables = import 'tables.jsonnet';

{
  foo: {
    name: schemas.foo.name,
    blueprint: [
      {
        name: tables.foo.name,
        columns: [{ name: columns.foo.name, type: "inconsistent" }],
        option_keys: tables.foo.option_keys,
      },
    ],
  },
}
