local columns = import 'columns.jsonnet';
local tables = import 'tables.jsonnet';

//  Notice that regular & rare types are all from UPDATE operation type
//  due to perform exhaustive fields analysis, and being UPDATE which takes all the fields
local regularID = 1;
local rareID = 101;

local regularEntityID = '01EntityID';
local rareEntityID = '001EntityID';

local strType = 'string';
local intType = 'int';
local jsonType = 'json';
local bytesType = 'bytes';
local float32Type = 'float32';
local float64Type = 'float64';

local regularType = strType;
local rareType = intType;

local regularStringValue = 'regularValue';
local rareStringValue = 'rareValue';

local regularIntValue = 1001;
local regularFloat32Value = regularIntValue;
local regularFloat64Value = regularIntValue;
local rareIntValue = 9001;

local regularJSONValue = { embedded_value: { another_embedding: 'regularValue' } };
local rareJSONValue = { embedded_value: { another_embedding: 'rareValue' } };
local toCreate(x) = x { entity_id: '', type: 'create' };
local toRetrieve(x) = x { value_type: '', str_value: '', type: 'retrieve' };
local toUpdate(x) = x { type: 'update' };
local toDelete(x) = x { value_type: '', str_value: '', column_name: '', type: 'delete' };

local regularOptionKey = tables.regularOptionKey;
local rareOptionKey = tables.rareOptionKey;
local regularOptionValue = 'regularOptionValue';
local rareOptionValue = 'rareOptionValue';
local createCRUD(x) = {
  create: toCreate(x),
  retrieve: toRetrieve(x),
  update: toUpdate(x),
  delete: toDelete(x),
};

{
  local regular = self.regular,

  inconsistent: {
    crud: {
      create: regular.crud.create { column_name: '' },
      retrieve: regular.crud.retrieve { column_name: '' },
      update: regular.crud.update { column_name: '' },
      delete: regular.crud.delete { value_type: regular.none.value_type },
    },
    table_name: regular.none { table_name: '' },
    column_name: regular.none { column_name: '', entity_id: '' },  // Contains entity to avoid unclassifiable handlings
  },


  regular: {
    local base = self.none,

    none: {
      table_name: tables.basic.name,
      column_name: columns.basic.name,
      str_value: regularStringValue,
      value_type: regularType,
      id: regularID,
      entity_id: regularEntityID,
      options: {
        regularOptionKey: regularOptionValue,
      },
    },

    crud: createCRUD(base),

    table_name: base { table_name: tables.rare.name },

    column_name: base { column_name: columns.rare.name },

    str_value: base { str_value: rareStringValue, value_type: strType },

    int_value: base { int_value: regularIntValue, value_type: intType, str_value: '' },

    float32_value: base { float32_value: regularFloat32Value, value_type: float32Type },

    float64_value: base { float64_value: regularFloat64Value, value_type: float64Type },

    id: base { id: rareID },

    entity_id: base { entity_id: rareEntityID },

    json_value: base { json_value: regularJSONValue, value_type: jsonType, str_value: '' },

    clean_value: base { str_value: '', value_type: '' },

  },

  rare: {
    local base = self.none,

    none: {
      table_name: tables.rare.name,
      column_name: columns.rare.name,
      int_value: rareIntValue,
      value_type: rareType,
      id: rareID,
      entity_id: rareEntityID,
      options: {
        rareOptionKey: rareOptionValue,
      },
    },

    crud: createCRUD(base),

    table_name: base { table_name: tables.basic.name },

    column_name: base { column_name: columns.basic.name },

    int_value: base { int_value: regularIntValue },

    str_value: base { str_value: rareStringValue, value_type: strType, int_value: 0 },

    id: base { id: regularID },

    json_value: base { json_value: rareJSONValue, value_type: jsonType, int_value: 0 },

    entity_id: base { entity_id: regularEntityID },

    clean_value: base { int_value: 0, value_type: '' },
  },

  zero: {},
}
