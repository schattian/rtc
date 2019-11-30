local altns = import 'altns.libsonnet';
local ids = import 'ids.libsonnet';

local columns = import 'columns.jsonnet';
local indices = import 'indices.jsonnet';
local tables = import 'tables.jsonnet';

local cleanVals(x) = x { str_value: '', int_value: 0, float_32_value: 0, float_64_value: 0, json_value: '', value_type: '' };
local assignVal(base, val, typ) = cleanVals(base) { [std.format('%s_value', typ)]: val, value_type: typ };

// CRUD
local toCreate(x) = x { entity_id: '', type: 'create' };
local toRetrieve(x) = x { value_type: '', str_value: '', type: 'retrieve' };
local toUpdate(x) = x { type: 'update' };
local toDelete(x) = cleanVals(x) { column_name: '', type: 'delete' };
local createCRUD(x) = {
  create: toCreate(x),
  retrieve: toRetrieve(x),
  update: toUpdate(x),
  delete: toDelete(x),
};

//  Notice that foo & bar types are all from UPDATE operation type
//  due to the need to perform exhaustive fields analysis, and being UPDATE which takes all the fields
{
  local foo = self.foo,
  local bar = self.bar,

  foo: {
    local base = self.none,
    local alt = $.bar.none,
    none: {
      table_name: tables.foo.name,
      column_name: columns.foo.name,
      str_value: 'fooChangeStringValue',
      value_type: columns.foo.type,
      id: ids.changeId(altns.foo),
      entity_id: 'fooChangeEntityId',
      options: {
        [tables.foo.option_keys[0]]: 'fooChangeOptionValue',
      },
      index_id: indices.foo.id,
    },
    id: base { id: alt.id },
    entity_id: base { entity_id: alt.entity_id },
    table_name: base { table_name: alt.table_name },
    column_name: base { column_name: alt.column_name },
    options: base { options: alt.options },
    index_id: base { index_id: alt.index_id },

    str_value: assignVal(base, bar.str_value.str_value, 'str'),
    int_value: assignVal(base, altns.foo * 2, 'int'),
    float_32_value: assignVal(base, altns.foo * 3.2, 'float32'),
    float_64_value: assignVal(base, altns.foo * 6.4, 'float64'),
    json_value: assignVal(base, 'fooChangeJsonValue', 'json'),
    clean_value: cleanVals(base),

    crud: createCRUD(base),
  },

  bar: {
    local base = self.none,
    local alt = $.foo.none,

    none: {
      table_name: tables.bar.name,
      column_name: columns.bar.name,
      int_value: altns.bar * 2,
      value_type: columns.bar.type,
      id: ids.changeId(altns.bar),
      entity_id: 'barChangeEntityId',
      options: {
        [tables.bar.option_keys[0]]: 'barChangeOptionValue',
      },
      index_id: indices.bar.id,
    },
    id: base { id: alt.id },
    entity_id: base { entity_id: alt.entity_id },
    table_name: base { table_name: alt.table_name },
    column_name: base { column_name: alt.column_name },
    options: base { options: alt.options },
    index_id: base { index_id: alt.index_id },

    str_value: assignVal(base, 'fooChangeStringValue', 'str'),
    int_value: assignVal(base, altns.bar * 2, 'int'),
    float_32_value: assignVal(base, altns.bar * 3.2, 'float32'),
    float_64_value: assignVal(base, altns.bar * 6.4, 'float64'),
    json_value: assignVal(base, 'barChangeJsonValue', 'json'),
    clean_value: cleanVals(base),

    crud: createCRUD(base),
  },

  inconsistent: {
    crud: {
      create: foo.crud.create { column_name: '' },
      retrieve: foo.crud.retrieve { column_name: '' },
      update: foo.crud.update { column_name: '' },
      delete: foo.crud.delete { value_type: foo.none.value_type },
    },
    table_name: foo.none { table_name: '' },
    column_name: foo.none { column_name: '', entity_id: '' },  // Contains entity to avoid unclassifiable handlings
  },

  zero: {},
}
