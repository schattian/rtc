local changes = import 'changes.jsonnet';
local schemas = import 'schemas.jsonnet';
local teams = import 'teams.jsonnet';

local chgToComm(x) = { changes: [x] };
local CRUD = {
  foo: {
    create: changes.foo.crud.create,
    retrieve: changes.foo.crud.retrieve,
    update: changes.foo.crud.update,
    delete: changes.foo.crud.delete,
  },
  bar: {
    create: changes.bar.crud.create,
    retrieve: changes.bar.crud.retrieve,
    update: changes.bar.crud.update,
    delete: changes.bar.crud.delete,
  },
};

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
        chgToComm(CRUD.foo.create),
        chgToComm(CRUD.foo.retrieve),
        chgToComm(CRUD.foo.update),
        chgToComm(CRUD.foo.delete),
      ],
  },

  crud: {
    create: foo {
      team: teams.foo,
      commits: [
        chgToComm(CRUD.foo.create),
      ],
    },

    retrieve: foo {
      team: teams.foo,
      commits: [
        chgToComm(CRUD.foo.retrieve),
      ],
    },

    update: foo {
      team: teams.foo,
      commits: [
        chgToComm(CRUD.foo.update),
      ],
    },

    delete: foo {
      team: teams.foo,
      commits: [
        chgToComm(CRUD.foo.delete),
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
