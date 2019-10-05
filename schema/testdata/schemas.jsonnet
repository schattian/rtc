local tables = import 'tables.jsonnet';
local regularName = 'basicSchemaName';
local rareName = 'rareSchemaName';
local basicRareName = 'basicRareSchemaName';

{
  basic: {
    name: regularName,
    blueprint: [tables.basic],
  },

  rare: {
    name: rareName,
    blueprint: [tables.rare],
  },

  basic_rare: {
    name: basicRareName,
    blueprint: [
      tables.basic,
      tables.rare,
    ],
  },
 
  zero:{},
}