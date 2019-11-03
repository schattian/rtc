local schemas = import './schemas.jsonnet';
local changes = import './changes.jsonnet';
local teams = import './teams.jsonnet';
local CRUD = {
  create: changes.basic.crud.create,
  retrieve: changes.basic.crud.retrieve,
  update: changes.basic.crud.update,
  delete: changes.basic.crud.delete,
};
local chgToComm(x) = { changes: [x] };

{
  local basic = self.basic,

  basic: {
    team: teams.basic,
    commits: [
      chgToComm(changes.basic.none),
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
