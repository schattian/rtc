package git

import (
	"context"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
)

// Index is the git-like representation of a group of a NON-ready-to-deliver changes
type Index struct {
	Id      int64     `json:"id,omitempty"`
	Changes []*Change `json:"changes,omitempty"`
}

// FetchChanges retrieves the changes from DB by its .ChangeIds and assigns them to .Changes field
func (idx *Index) FetchChanges(ctx context.Context, db *sqlx.DB) (err error) {
	rows, err := db.NamedQueryContext(ctx, `SELECT * FROM changes WHERE committed=FALSE AND index_id=:id`, idx)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		chg := Change{}
		err = rows.StructScan(chg)
		if err != nil {
			return
		}
		idx.Changes = append(idx.Changes, &chg)
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	return
}

// Add will attach the given change to the commit changes
// In case the change is invalid or is duplicated, it returns an error
func (idx *Index) Add(chg *Change) error {
	err := chg.Validate()
	if err != nil {
		return err
	}
	if idx.containsChange(chg) { // avoids duplication. It discards untracked changes in the comparison
		return nil
	}

	for i, otherChg := range idx.Changes {
		if Overrides(chg, otherChg) {
			idx.rmChangeByIndex(i)
		}
	}
	idx.Changes = append(idx.Changes, chg)
	return nil
}

// Rm deletes the given change from the commit
// This action is irreversible
func (idx *Index) Rm(chg *Change) error {
	for i, otherChg := range idx.Changes {
		if chg.Equals(otherChg) {
			idx.rmChangeByIndex(i)
			return nil
		}
	}
	return fmt.Errorf("change %v NOT FOUND", chg)
}

// rmChangeByIndex will delete without preserving order giving the desired index to delete
func (idx *Index) rmChangeByIndex(i int) {
	var lock sync.Mutex // Avoid overlapping itself with a parallel call
	lock.Lock()
	lastIndex := len(idx.Changes) - 1
	idx.Changes[i] = idx.Changes[lastIndex]
	idx.Changes[lastIndex] = nil // Notices the GC to rm the last elem to avoid mem-leak
	idx.Changes = idx.Changes[:lastIndex]
	lock.Unlock()
}

// containsChange verifies if the given change is already present, and triggering the **exactly** same action
func (idx *Index) containsChange(chg *Change) bool {
	for _, otherChg := range idx.Changes {
		if chg.Equals(otherChg) {
			return true
		}
	}
	return false
}
