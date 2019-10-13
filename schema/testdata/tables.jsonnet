local rareName = 'rareTableName';
local regularName = 'tableName';
local columns = import 'columns.jsonnet';
local rareRegularName = 'rareRegularName';
local regularOptionKey = "regularOptionKey";
local rareOptionKey = "rareOptionKey";

{
  basic: {
    name: regularName,
    columns: [columns.basic],
    option_keys: [regularOptionKey],
  },

  rare: {
    name: rareName,
    columns: [columns.rare],
    option_keys: [rareOptionKey],
  },

  basic_rare: {
    name: rareRegularName,
    columns: [
      columns.basic,
      columns.rare,
    ],
    option_keys: [regularOptionKey, rareOptionKey],
  },

  zero: {},
}
