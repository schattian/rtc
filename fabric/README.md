# Fabric

This subpkg generates structs used to decode the formats of a **git.Commit to a native go struct** (friendlier that managing commit changes and better than using unsafe maps).

> Notice that _every generation needs a valid schema to be performed_.

An alternative to this could be using [gojson](https://github.com/ChimeraCoder/gojson), creating the templates for every endpoint manually.

## Getting started

The standard usage of this module would be:

1. Fabric the structs with `git-crud fabric <SCHEMA_PATH>`.
2. Marshal the git.Commit using any utility provided in the msh pkg or use the `Mapable` provided abstraction to create your own marshaler.
3. Unmarshal the obtained map to the struct you generated in the first step.

Example:

Using the github literal, first we are going to create the native structs for this `git-crud fabric literals/github.yaml`.

Then, and after committing a few changes (see how to commit)[] we orchestrate the changes:

```go
// Notice it avoids err checking <is an example>

import (
    "encoding/json"
    "context"

    fabric "my_fabric_dir/github"
    "github.com/sebach1/git-crud/git"
    "github.com/sebach1/git-crud/literals"
    "github.com/sebach1/git-crud/literals/github"
    "github.com/sebach1/git-crud/msh"
)

myOwner := git.NewOwner(literals.StdPlanisphere)

myOwner.Orchestrate(context.Background(), github.Community, "github", YOUR_COMMIT, git.AreCompatible())
myOwner.Close()

for _, result := range myOwner.Summary {
    comm := CommitById(result.CommitId)// find the created commit with your DB implementation
    jsComm := msh.ToJSON(comm)
    var repo fabric.Repository{})
    json.Unmarshal(jsComm, &repo) // In case u are using json decoder
}

```
