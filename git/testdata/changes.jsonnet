local columns = import 'columns.jsonnet';
local tables = import 'tables.jsonnet';

//  Notice that basic & rare types are all from UPDATE operation type
//  due to perform exhaustive fields analysis, and being UPDATE which takes all the fields
local basicID = 1;
local rareID = 101;

local basicEntityID = '01EntityID';
local rareEntityID = '001EntityID';

local basicStringValue = 'basicValue';
local rareStringValue = 'rareValue';
local basicIntValue = 1001;
local basicFloat32Value = basicIntValue;
local basicFloat64Value = basicIntValue;
local rareIntValue = 9001;
local basicJSONValue = { embedded_value: { another_embedding: 'basicValue' } };
local rareJSONValue = { embedded_value: { another_embedding: 'rareValue' } };

// oK
local basicOptionKey = tables.basicOptionKey;
local rareOptionKey = tables.rareOptionKey;
local basicOptionValue = 'basicOptionValue';
local rareOptionValue = 'rareOptionValue';

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
  local basic = self.basic,
  local rare = self.rare,


  basic: {
    local base = self.none,

    none: {
      table_name: tables.basic.name,
      column_name: columns.basic.name,
      str_value: basicStringValue,
      value_type: columns.basic.type,
      id: basicID,
      entity_id: basicEntityID,
      options: {
        basicOptionKey: basicOptionValue,
      },
    },

    crud: createCRUD(base),

    table_name: base { table_name: tables.rare.name },

    column_name: base { column_name: columns.rare.name },

    str_value: base { str_value: rareStringValue, value_type: "string" },

    int_value: base { int_value: basicIntValue, value_type: "int", str_value: '' },

    float32_value: base { float32_value: basicFloat32Value, value_type: "float32" },

    float64_value: base { float64_value: basicFloat64Value, value_type: "float64" },

    id: base { id: rareID },

    entity_id: base { entity_id: rareEntityID },

    json_value: base { json_value: basicJSONValue, value_type: "json", str_value: '' },

    clean_value: base { str_value: '', value_type: '' },
    
    options: base { options: rare.none.options },

  },

  rare: {
    local base = self.none,

    none: {
      table_name: tables.rare.name,
      column_name: columns.rare.name,
      int_value: rareIntValue,
      value_type: columns.rare.type,
      id: rareID,
      entity_id: rareEntityID,
      options: {
        rareOptionKey: rareOptionValue,
      },
    },

    crud: createCRUD(base),

    table_name: base { table_name: tables.basic.name },

    column_name: base { column_name: columns.basic.name },

    int_value: base { int_value: basicIntValue },

    str_value: base { str_value: rareStringValue, value_type: "string", int_value: 0 },

    id: base { id: basicID },

    json_value: base { json_value: rareJSONValue, value_type: "json", int_value: 0 },

    entity_id: base { entity_id: basicEntityID },

    clean_value: base { int_value: 0, value_type: '' },
    
    options: base { options: basic.none.options },
  },

  inconsistent: {
    crud: {
      create: basic.crud.create { column_name: '' },
      retrieve: basic.crud.retrieve { column_name: '' },
      update: basic.crud.update { column_name: '' },
      delete: basic.crud.delete { value_type: basic.none.value_type },
    },
    table_name: basic.none { table_name: '' },
    column_name: basic.none { column_name: '', entity_id: '' },  // Contains entity to avoid unclassifiable handlings
  },

  zero: {},
}
