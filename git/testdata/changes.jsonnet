//  Notice that regular & rare types are all from UPDATE operation type
//  due to perform exhaustive fields analysis, and being UPDATE which takes all the fields
local regularID = 1;
local rareID = 101;

local regularEntityID = '01EntityID';
local rareEntityID = '001EntityID';

local regularTable = 'regularTable';
local rareTable = 'rareTable';

local regularColumn = 'regularColumn';
local rareColumn = 'rareColumn';


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
local toDelete(x) = x { value_type: '', str_value: '', column: '', type: 'delete' };

{
  local regular = self.regular,

  inconsistent: {
    table: regular.none { table_name: '' },
    column: regular.none { column_name: '', entity_id: ""},  // Contains entity to avoid unclassifiable handlings
  },


  regular: {
    local base = self.none,

    none: {
      table_name: regularTable,
      column_name: regularColumn,
      str_value: regularStringValue,
      value_type: regularType,
      id: regularID,
      entity_id: regularEntityID,
    },

    create: toCreate(base),
    retrieve: toRetrieve(base),
    update: toUpdate(base),
    delete: toDelete(base),

    table: base { table_name: rareTable },

    column: base { column_name: rareColumn },

    str_value: base { str_value: rareStringValue, value_type: strType },

    int_value: base { int_value: regularIntValue, value_type: intType, str_value: '' },

    float32_value: base { float32_value: regularFloat32Value, value_type: float32Type },

    float64_value: base { float64_value: regularFloat64Value, value_type: float64Type },

    id: base { id: rareID },

    entity: base { entity_id: rareEntityID },

    json_value: base { json_value: regularJSONValue, value_type: jsonType, str_value: '' },

    clean_value: base { str_value: '', value_type: '' },

    untracked: base { entity_id: '' },
  },

  rare: {
    local base = self.none,

    none: {
      table_name: rareTable,
      column_name: rareColumn,
      int_value: rareIntValue,
      value_type: rareType,
      id: rareID,
      entity_id: rareEntityID,
    },

    create: toCreate(base),
    retrieve: toRetrieve(base),
    update: toUpdate(base),
    delete: toDelete(base),


    table: base { table_name: regularTable },

    column: base { column_name: regularColumn },

    int_value: base { int_value: regularIntValue },

    str_value: base { str_value: rareStringValue, value_type: strType, int_value: 0 },

    id: base { id: regularID },

    json_value: base { json_value: rareJSONValue, value_type: jsonType, int_value: 0 },

    clean_value: base { int_value: 0, value_type: '' },

    entity: base { entity_id: regularEntityID },
    
    untracked: base { entity_id: '' },
  },

  zero: {},
}
