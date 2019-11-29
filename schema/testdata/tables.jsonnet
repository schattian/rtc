local columns = import 'columns.jsonnet';

{
  local foo = self.foo,
  local bar = self.bar,

  foo: {
    name: 'fooTableName',
    columns: [columns.foo],
    option_keys: ['fooTableOptionKey'],
  },

  bar: {
    name: 'barTableName',
    columns: [columns.bar],
    option_keys: ['barTableOptionKey'],
  },

  foo_bar: {
    name: 'fooBarTableName',
    columns: [
      columns.foo,
      columns.bar,
    ],
    option_keys: [
      foo.option_keys[0],
      bar.option_keys[0],
    ],
  },

  zero: {},
}
