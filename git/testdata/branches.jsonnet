local altns = import 'altns.libsonnet';
local ids = import 'ids.libsonnet';

local indices = import 'indices.jsonnet';

{
  foo: {
    id: ids.branchId(altns.foo),
    name: 'fooBranchName',
    index_id: indices.foo.id,
    index: indices.foo,
  },
  bar: {
    id: ids.branchId(altns.bar),
    name: 'barBranchName',
    index_id: indices.bar.id,
    index: indices.bar,
  },

  zero: {},
}
