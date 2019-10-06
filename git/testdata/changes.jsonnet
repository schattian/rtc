local regularID = 1;
local rareID = 101;

local regularEntityID = '01EntityID';
local rareEntityID = '001EntityID';

local regularTable = 'regularTable';
local rareTable = 'rareTable';

local regularColumn = 'regularColumn';
local rareColumn = 'rareColumn';

local regularStringValue = 'regularValue';
local rareStringValue = 'rareValue';
local regularIntValue = 1001;
local rareIntValue = 9001;

{
  regular: {
    local base = self.none,

    none: {
      table_name: regularTable,
      column_name: regularColumn,
      value: regularStringValue,
      id: regularID,
      entity_id: regularEntityID,
    },

    table: base { table_name: rareTable },

    column: base { column_name: rareColumn },

    value: base { value: rareStringValue },

    id: base { id: rareID },

    entity: base { entity_id: rareEntityID },

    untracked: base { entity_id: '' },
  },

  rare: {
    local base = self.none,

    none: {
      table_name: rareTable,
      column_name: rareColumn,
      value: rareStringValue,
      id: rareID,
      entity_id: rareEntityID,
    },

    table: base { table_name: regularTable },

    column: base { column_name: regularColumn },

    value: base { value: regularStringValue },

    id: base { id: regularID },

    entity: base { entity_id: regularEntityID },

    untracked: base { entity_id: '' },
  },

  zero: {},
}
