local schemas = import './schemas.jsonnet';

{
  local basic = self.basic,
  basic: {
    assigned_schema: schemas.basic.name,
    members: [
      { assigned_table: schemas.basic.blueprint[0].name },
    ],
  },
  rare: {
    assigned_schema: schemas.rare.name,
    members: [
      { assigned_table: schemas.rare.blueprint[0].name },
    ],
  },
  basic_rare: {
    assigned_schema: schemas.basic_rare.name,
    members: [
      { assigned_table: schemas.basic_rare.blueprint[0].name },
    ],
  },

  inconsistent:
    basic { members: [
      { assigned_table: schemas.basic_rare.blueprint[0].name },
    ] },

  zero_members:
    basic { members: [] },

  zero: {},
}
