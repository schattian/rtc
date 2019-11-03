local columns = import 'columns.jsonnet';

local basicName = 'basicTableName';
local rareName = 'rareTableName';
local basicRareName = 'basicRareTableName';

local basicOptionKey = 'basicOptionKey';
local rareOptionKey = 'rareOptionKey';

{
  basic: {
    name: basicName,
    columns: [columns.basic],
    option_keys: [basicOptionKey],
  },

  rare: {
    name: rareName,
    columns: [columns.rare],
    option_keys: [rareOptionKey],
  },

  basic_rare: {
    name: basicRareName,
    columns: [
      columns.basic,
      columns.rare,
    ],
    option_keys: [
      basicOptionKey,
      rareOptionKey,
    ],
  },

  zero: {},
}
