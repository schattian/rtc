local rareName = 'rareTableName';
local regularName = 'tableName';
local columns = import 'columns.jsonnet';
local rareRegularName = 'rareRegularName';

{
  basic: {
    name: regularName,
    columns: [columns.basic],
  },

  rare: {
    name: rareName,
    columns: [columns.rare],
  },

  basic_rare: {
    name: rareRegularName,
    columns: [
      columns.basic,
      columns.rare,
    ],
  },

  zero: {},
}
