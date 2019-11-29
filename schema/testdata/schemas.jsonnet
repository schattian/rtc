local tables = import 'tables.jsonnet';

{
  foo: {
    name: 'fooSchemaName',
    blueprint: [tables.foo],
  },

  bar: {
    name: 'barSchemaName',
    blueprint: [tables.bar],
  },

  foo_bar: {
    name: 'fooBarSchemaName',
    blueprint: [
      tables.foo_bar,
      tables.foo,
      tables.bar,
    ],
  },

  zero: {},
}
