local tables = import 'tables.jsonnet';

local basicName = 'basicSchemaName';
local rareName = 'rareSchemaName';
local basicRareName = 'basicRareSchemaName';

{
  basic: {
    name: basicName,
    blueprint: [tables.basic],
  },

  rare: {
    name: rareName,
    blueprint: [tables.rare],
  },

  basic_rare: {
    name: basicRareName,
    blueprint: [
      tables.basic_rare,
      tables.basic,
      tables.rare,
    ],
  },
 
  zero:{},
}
