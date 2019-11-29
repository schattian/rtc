local schemas = import './schemas.jsonnet';
local tables = import './tables.jsonnet';

{
  foo: {
    assigned_schema: schemas.foo.name,
    members: [
      { assigned_table: tables.foo.name },
    ],
  },
  bar: {
    assigned_schema: schemas.bar.name,
    members: [
      { assigned_table: tables.bar.name },
    ],
  },
  foo_bar: {
    assigned_schema: schemas.foo_bar.name,
    members: [
      { assigned_table: tables.foo_bar.name },
    ],
  },

  inconsistent:
    $.foo { members: [{ assigned_table: schemas.foo_bar.blueprint[0].name }] },

  zero_members:
    $.foo { members: [] },

  zero: {},
}
