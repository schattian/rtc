local columns = import 'columns.jsonnet';

local fooName = 'fooTableName';
local barName = 'barTableName';
local fooBarName = 'fooBarTableName';

local fooOptionKey = 'fooOptionKey';
local barOptionKey = 'barOptionKey';

{
  foo: {
    name: fooName,
    columns: [columns.foo],
    option_keys: [fooOptionKey],
  },

  bar: {
    name: barName,
    columns: [columns.bar],
    option_keys: [barOptionKey],
  },

  foo_bar: {
    name: fooBarName,
    columns: [
      columns.foo,
      columns.bar,
    ],
    option_keys: [
      fooOptionKey,
      barOptionKey,
    ],
  },

  zero: {},
}
