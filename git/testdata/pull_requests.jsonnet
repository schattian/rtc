local schemas = import '../../schema/testdata/schemas.jsonnet';
local changes = import './changes.jsonnet';
local teams = import './teams.jsonnet';
local CRUD = {
  create: changes.regular.crud.create,
  retrieve: changes.regular.crud.retrieve,
  update: changes.regular.crud.update,
  delete: changes.regular.crud.delete,
};
local chgToComm(x) = { changes: [x] };

{
  local basic = self.basic,

  basic: {
    team: teams.basic,
    commits: [
      chgToComm(changes.regular.none),
    ],
  },

  full: basic {
    team: teams.basic,
    commits:
      [
        chgToComm(CRUD.create),
        chgToComm(CRUD.retrieve),
        chgToComm(CRUD.update),
        chgToComm(CRUD.delete),
      ],
  },

  crud: {
    create: basic {
      team: teams.basic,
      commits: [
        chgToComm(CRUD.create),
      ],
    },

    retrieve: basic {
      team: teams.basic,
      commits: [
        chgToComm(CRUD.retrieve),
      ],
    },

    update: basic {
      team: teams.basic,
      commits: [
        chgToComm(CRUD.update),
      ],
    },

    delete: basic {
      team: teams.basic,
      commits: [
        chgToComm(CRUD.delete),
      ],
    },
  },

  zero_commits: basic {
    commits: [],
  },

  zero_team: basic {
    team: {},
  },


  zero: {},
}
