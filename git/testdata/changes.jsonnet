local columns = import 'columns.jsonnet';
local tables = import 'tables.jsonnet';

//  Notice that foo & bar types are all from UPDATE operation type
//  due to perform exhaustive fields analysis, and being UPDATE which takes all the fields
local fooId = 1;
local barId = 101;

local fooEntityId = '01EntityId';
local barEntityId = '001EntityId';

local fooStringValue = 'fooValue';
local barStringValue = 'barValue';
local fooIntValue = 1001;
local fooFloat32Value = fooIntValue;
local fooFloat64Value = fooIntValue;
local barIntValue = 9001;
local fooJSONValue = { embedded_value: { another_embedding: 'fooValue' } };
local barJSONValue = { embedded_value: { another_embedding: 'barValue' } };

// oK
local fooOptionKey = tables.fooOptionKey;
local barOptionKey = tables.barOptionKey;
local fooOptionValue = 'fooOptionValue';
local barOptionValue = 'barOptionValue';

// CRUD
local toCreate(x) = x { entity_id: '', type: 'create' };
local toRetrieve(x) = x { value_type: '', str_value: '', type: 'retrieve' };
local toUpdate(x) = x { type: 'update' };
local toDelete(x) = x { value_type: '', str_value: '', column_name: '', type: 'delete' };
local createCRUD(x) = {
  create: toCreate(x),
  retrieve: toRetrieve(x),
  update: toUpdate(x),
  delete: toDelete(x),
};

{
  local foo = self.foo,
  local bar = self.bar,


  foo: {
    local base = self.none,

    none: {
      table_name: tables.foo.name,
      column_name: columns.foo.name,
      str_value: fooStringValue,
      value_type: columns.foo.type,
      id: fooId,
      entity_id: fooEntityId,
      options: {
        fooOptionKey: fooOptionValue,
      },
    },

    crud: createCRUD(base),

    table_name: base { table_name: tables.bar.name },

    column_name: base { column_name: columns.bar.name },

    str_value: base { str_value: barStringValue, value_type: "string" },

    int_value: base { int_value: fooIntValue, value_type: "int", str_value: '' },

    float_32_value: base { float_32_value: fooFloat32Value, value_type: "float32" },

    float_64_value: base { float_64_value: fooFloat64Value, value_type: "float64" },

    id: base { id: barId },

    entity_id: base { entity_id: barEntityId },

    json_value: base { json_value: fooJSONValue, value_type: "json", str_value: '' },

    clean_value: base { str_value: '', value_type: '' },
    
    options: base { options: bar.none.options },

  },

  bar: {
    local base = self.none,

    none: {
      table_name: tables.bar.name,
      column_name: columns.bar.name,
      int_value: barIntValue,
      value_type: columns.bar.type,
      id: barId,
      entity_id: barEntityId,
      options: {
        barOptionKey: barOptionValue,
      },
    },

    crud: createCRUD(base),

    table_name: base { table_name: tables.foo.name },

    column_name: base { column_name: columns.foo.name },

    int_value: base { int_value: fooIntValue },

    str_value: base { str_value: barStringValue, value_type: "string", int_value: 0 },

    id: base { id: fooId },

    json_value: base { json_value: barJSONValue, value_type: "json", int_value: 0 },

    entity_id: base { entity_id: fooEntityId },

    clean_value: base { int_value: 0, value_type: '' },
    
    options: base { options: foo.none.options },
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
