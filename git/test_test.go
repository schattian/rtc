package git

import "github.com/sebach1/git-crud/internal/assist"

var (
	gChanges goldenChanges
)

func init() {
	assist.DecodeJsonnet("changes", &gChanges)
}

type goldenChanges struct {
	Regular variadicChanges `json:"regular,omitempty"`
	Rare    variadicChanges `json:"rare,omitempty"`

	Zero *Change `json:"zero,omitempty"`
}

type variadicChanges struct {
	None      *Change `json:"none,omitempty"`
	Table     *Change `json:"table,omitempty"`
	Column    *Change `json:"column,omitempty"`
	Value     *Change `json:"value,omitempty"`
	ID        *Change `json:"id,omitempty"`
	Entity    *Change `json:"entity,omitempty"`
	Untracked *Change `json:"untracked,omitempty"`
}
