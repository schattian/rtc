local schemas = import '../../schema/testdata/schemas.jsonnet';
local changes = import './changes.jsonnet';
local teams = import './teams.jsonnet';
local CRUD = {
  create: changes.regular.crud.create,
  retrieve: changes.regular.crud.retrieve,
  update: changes.regular.crud.update,
  delete: changes.regular.crud.delete,
};

// local CRUDGen(exc) = {
//   team: teams.basic,
//   commits: [if exc.count(f) == 0 then CRUD.f for f in std.objectFields(CRUD)],
// };

{
  local basic = self.basic,

  basic: {
    team: teams.basic,
    commits: [
      changes.regular.none,
    ],
  },

  full: basic {
    commits: [
      CRUD
    ],
  },


  create: basic {
    commits: [
      changes.regular.crud.create,
    ],
  },

  retrieve: basic {
    commits: [
      CRUD.retrieve
    ],
  },

  update: basic {
    commits: [
      CRUD.update
    ],
  },

  delete: basic {
    commits: [
      CRUD.delete
    ],
  },


  zero_commits: basic {
    commits: [],
  },

  zero_team: basic {
    team: {},
  },


  zero: {},
}
