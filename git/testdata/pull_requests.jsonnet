local schemas = import './schemas.jsonnet';
local changes = import './changes.jsonnet';
local teams = import './teams.jsonnet';
local CRUD = {
  create: changes.foo.crud.create,
  retrieve: changes.foo.crud.retrieve,
  update: changes.foo.crud.update,
  delete: changes.foo.crud.delete,
};
local chgToComm(x) = { changes: [x] };

{
  local foo = self.foo,

  foo: {
    team: teams.foo,
    commits: [
      chgToComm(changes.foo.none),
    ],
  },

  full: foo {
    team: teams.foo,
    commits:
      [
        chgToComm(CRUD.create),
        chgToComm(CRUD.retrieve),
        chgToComm(CRUD.update),
        chgToComm(CRUD.delete),
      ],
  },

  crud: {
    create: foo {
      team: teams.foo,
      commits: [
        chgToComm(CRUD.create),
      ],
    },

    retrieve: foo {
      team: teams.foo,
      commits: [
        chgToComm(CRUD.retrieve),
      ],
    },

    update: foo {
      team: teams.foo,
      commits: [
        chgToComm(CRUD.update),
      ],
    },

    delete: foo {
      team: teams.foo,
      commits: [
        chgToComm(CRUD.delete),
      ],
    },
  },

  zero_commits: foo {
    commits: [],
  },

  zero_team: foo {
    team: {},
  },


  zero: {},
}
