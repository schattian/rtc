package git

import (
	"math/rand"

	"github.com/sebach1/git-crud/internal/test/assist"
)

var (
	gChanges goldenChanges
)

func init() {
	assist.DecodeJsonnet("changes", &gChanges)
}

type goldenChanges struct {
	Regular      variadicChanges `json:"regular,omitempty"`
	Rare         variadicChanges `json:"rare,omitempty"`
	Inconsistent variadicChanges `json:"inconsistent,omitempty"`

	Zero *Change `json:"zero,omitempty"`
}

type variadicChanges struct {
	None   *Change `json:"none,omitempty"`
	Table  *Change `json:"table,omitempty"`
	Column *Change `json:"column,omitempty"`

	StrValue     *Change `json:"str_value,omitempty"`
	IntValue     *Change `json:"int_value,omitempty"`
	Float32Value *Change `json:"float32_value,omitempty"`
	Float64Value *Change `json:"float64_value,omitempty"`
	JSONValue    *Change `json:"json_value,omitempty"`
	CleanValue   *Change `json:"clean_value,omitempty"`

	ID        *Change `json:"id,omitempty"`
	Entity    *Change `json:"entity,omitempty"`
	Untracked *Change `json:"untracked,omitempty"`
}

func randChg(chgs ...*Change) *Change {
	return chgs[rand.Intn(len(chgs)-1)]
}
