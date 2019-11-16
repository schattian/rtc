package git

import (
	"math/rand"

	"github.com/sebach1/git-crud/internal/test/assist"
	"github.com/sebach1/git-crud/schema"
)

var (
	gChanges      goldenChanges
	gTeams        goldenTeams
	gPullRequests goldenPullRequests

	gSchemas schema.GoldenSchemas
	gTables  schema.GoldenTables
	gColumns schema.GoldenColumns
)

func init() {
	assist.DecodeJsonnet("changes", &gChanges)
	assist.DecodeJsonnet("pull_requests", &gPullRequests)
	assist.DecodeJsonnet("teams", &gTeams)

	assist.DecodeJsonnet("schemas", &gSchemas)
	assist.DecodeJsonnet("tables", &gTables)
	assist.DecodeJsonnet("columns", &gColumns)
}

type goldenChanges struct {
	Foo          variadicChanges `json:"foo,omitempty"`
	Bar          variadicChanges `json:"bar,omitempty"`
	Inconsistent variadicChanges `json:"inconsistent,omitempty"`

	Zero *Change `json:"zero,omitempty"`
}
type variadicChanges struct {
	None       *Change `json:"none,omitempty"`
	Id         *Change `json:"id,omitempty"`
	EntityId   *Change `json:"entity_id,omitempty"`
	TableName  *Change `json:"table_name,omitempty"`
	ColumnName *Change `json:"column_name,omitempty"`

	StrValue     *Change `json:"str_value,omitempty"`
	IntValue     *Change `json:"int_value,omitempty"`
	Float32Value *Change `json:"float32_value,omitempty"`
	Float64Value *Change `json:"float64_value,omitempty"`
	JSONValue    *Change `json:"json_value,omitempty"`
	CleanValue   *Change `json:"clean_value,omitempty"`

	Options *Change `json:"options,omitempty"`

	ChgCRUD `json:"crud,omitempty"`
}

func randChg(chgs ...*Change) *Change {
	return chgs[rand.Intn(len(chgs)-1)]
}

type goldenPullRequests struct {
	Foo *PullRequest `json:"foo,omitempty"`

	Full *PullRequest `json:"full,omitempty"`

	PrCRUD `json:"crud,omitempty"`

	ZeroCommits *PullRequest `json:"zero_commits,omitempty"`
	ZeroTeam    *PullRequest `json:"zero_team,omitempty"`

	Zero *PullRequest `json:"zero,omitempty"`
}

type ChgCRUD struct {
	Create   *Change `json:"create,omitempty"`
	Retrieve *Change `json:"retrieve,omitempty"`
	Update   *Change `json:"update,omitempty"`
	Delete   *Change `json:"delete,omitempty"`
}

type PrCRUD struct {
	Create   *PullRequest `json:"create,omitempty"`
	Retrieve *PullRequest `json:"retrieve,omitempty"`
	Update   *PullRequest `json:"update,omitempty"`
	Delete   *PullRequest `json:"delete,omitempty"`
}

type goldenTeams struct {
	Foo    *Team `json:"foo,omitempty"`
	Bar    *Team `json:"bar,omitempty"`
	FooBar *Team `json:"foo_bar,omitempty"`

	Inconsistent *Team `json:"inconsistent,omitempty"`

	ZeroMembers *Team `json:"zero_members,omitempty"`
	Zero        *Team `json:"zero,omitempty"`
}
