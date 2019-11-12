local tables = import 'tables.jsonnet';

local fooName = 'fooSchemaName';
local barName = 'barSchemaName';
local fooBarName = 'fooBarSchemaName';

{
  foo: {
    name: fooName,
    blueprint: [tables.foo],
  },

  bar: {
    name: barName,
    blueprint: [tables.bar],
  },

  foo_bar: {
    name: fooBarName,
    blueprint: [
      tables.foo_bar,
      tables.foo,
      tables.bar,
    ],
  },
 
  zero:{},
}
