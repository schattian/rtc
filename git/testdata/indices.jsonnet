local altns = import 'altns.libsonnet';
local ids = import 'ids.libsonnet';

local changes = import 'changes.jsonnet';

{

  foo: {
    local change = changes.foo.none,
    id: ids.indexId(altns.foo),
    changes: [change],
  },

  bar: {
    local change = changes.foo.none,
    id: ids.indexId(altns.bar),
    changes: [change],
  },

  zero: {},

}
