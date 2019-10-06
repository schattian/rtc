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

{
  local regular = self.regular,
 
  inconsistent: {
    table: regular.none {table_name: ""},
    column: regular.none {column_name: ""},
  },


  regular: {
    local base = self.none,

    none: {
      table_name: regularTable,
      column_name: regularColumn,
      str_value: regularStringValue,
      type: regularType,
      id: regularID,
      entity_id: regularEntityID,
    },


    table: base { table_name: rareTable },

    column: base { column_name: rareColumn },

    str_value: base { str_value: rareStringValue, type: strType },

    int_value: base { int_value: regularIntValue, type: intType, str_value: '' },

    float32_value: base { float32_value: regularFloat32Value, type: float32Type },

    float64_value: base { float64_value: regularFloat64Value, type: float64Type },

    id: base { id: rareID },

    entity: base { entity_id: rareEntityID },

    json_value: base { json_value: regularJSONValue, type: jsonType, str_value: '' },

    clean_value: base { str_value: '', type: '' },

    untracked: base { entity_id: '' },
  },

  rare: {
    local base = self.none,

    none: {
      table_name: rareTable,
      column_name: rareColumn,
      int_value: rareIntValue,
      type: rareType,
      id: rareID,
      entity_id: rareEntityID,
    },

    table: base { table_name: regularTable },

    column: base { column_name: regularColumn },

    int_value: base { int_value: regularIntValue },

    str_value: base { str_value: rareStringValue, type: strType, int_value: 0 },

    id: base { id: regularID },

    json_value: base { json_value: rareJSONValue, type: jsonType, int_value: 0 },

    clean_value: base { int_value: 0, type: '' },

    entity: base { entity_id: regularEntityID },

    untracked: base { entity_id: '' },
  },

  zero: {},
}
