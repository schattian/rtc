local scales = {
  pull_requests: 100000,
  commits: 10000,
  changes: 1000,
  branches: 100,
  indices: 10,

  columns: 3000,
  tables: 300,
  schemas: 30,
};

local gen(scale, altn) = scales[scale] + altn;

{

  changeId(altn)::
    gen('changes', altn),

  branchId(altn)::
    gen('branches', altn),

  indexId(altn)::
    gen('indices', altn),

  commitId(altn)::
    gen('commits', altn),

  pullRequestId(altn)::
    gen('pull_requests', altn),

  columnId(altn)::
    gen('columns', altn),

  tableId(altn)::
    gen('columns', altn),

  schemaId(altn)::
    gen('schemas', altn),
}
