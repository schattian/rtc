# Fabric

This subpkg generates structs used to decode the `JSON` or `map[string]interface{}` formats of a **git.Commit to a native go struct** (friendlier that managing commit changes and better than using safe bodies).

> Notice that _every generation needs a valid schema to be performed_.

An alternative to this could be using [gojson](https://github.com/ChimeraCoder/gojson), creating the templates for every endpoint manually.

## Getting started

The standard usage of this module would be:

1. Fabric the structs with `git-crud fabric <SCHEMA_PATH>`.
2. Convert the obtained git.Commit to `map[string]interface{}` using `.ToMap()`.
3. Decode the obtained map to the struct you generated in the first step.

Example:

Using the github literal, first we are going to create the native structs for this `git-crud fabric literals/github.yaml`.

Then, and after committing a few changes (see how to commit)[] we orchestrate the changes:

```go
import (
    "context"

    "github.com/sebach1/git-crud/git"
    "github.com/sebach1/git-crud/literals"
    "github.com/sebach1/git-crud/literals/github"
)

myOwner := &git.Owner{Project: literals.StdPlanisphere}

err := myOwner.Orchestrate(context.Background(), github.Community, "github", YOUR_COMMIT, git.AreCompatible())
// Avoid err checking <is an example>

myOwner.Summary

```
