package git

import (
	"context"
	"sync"

	"github.com/sebach1/rtc/internal/store"

	"github.com/jmoiron/sqlx"
)

// Index is the git-like representation of a group of a NON-ready-to-deliver changes
type Index struct {
	Id      int64     `json:"id,omitempty"`
	Changes []*Change `json:"changes,omitempty"`
}

// add will attach the given change to the index changes
// In case the change is invalid or is duplicated, it returns an error
// Its reciprocal to idx.Rm() -excepting for the generated id, obviously-
func (idx *Index) Add(ctx context.Context, db *sqlx.DB, chg *Change) error {
	err := idx.FetchUncommittedChanges(ctx, db)
	if err != nil {
		return err
	}
	err = idx.add(chg)
	if err != nil {
		return err
	}
	return store.InsertIntoDB(ctx, db, chg)
}

func (idx *Index) add(chg *Change) error {
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

// Rm deletes the given change
// This action is irreversible
// Its reciprocal to idx.Add() -excepting for the generated Id, obviously-
func (idx *Index) Rm(ctx context.Context, db *sqlx.DB, chg *Change) error {
	err := store.DeleteFromDB(ctx, db, chg)
	if err != nil {
		return err
	}
	idx.rm(chg)
	return nil
}

func (idx *Index) rm(chg *Change) {
	for i, otherChg := range idx.Changes {
		if chg.Equals(otherChg) {
			idx.rmChangeByIndex(i)
			return
		}
	}
}

// FetchUncommittedChanges retrieves the changes from DB by its .ChangeIds and assigns them to .Changes field
// It filters committed changes in query
func (idx *Index) FetchUncommittedChanges(ctx context.Context, db *sqlx.DB) (err error) {
	rows, err := db.NamedQueryContext(ctx, `SELECT * FROM changes WHERE commit_id=0 AND index_id=:id`, idx)
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
		return
	}
	return
}

// FetchChanges retrieves the changes from DB by its .ChangeIds and assigns them to .Changes field
func (idx *Index) FetchChanges(ctx context.Context, db *sqlx.DB) (err error) {
	rows, err := db.NamedQueryContext(ctx, `SELECT * FROM changes WHERE index_id=:id`, idx)
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
		return
	}
	return
}

// Commit returns a persisted commit with the index's uncommitted changes.
func (idx *Index) Commit(ctx context.Context, db *sqlx.DB) ([]*Commit, error) {
	err := idx.FetchUncommittedChanges(ctx, db)
	if err != nil {
		return nil, err
	}
	comms, err := idx.commit(ctx, db)
	if err != nil {
		return nil, err
	}
	return comms, nil
}

func (idx *Index) commit(ctx context.Context, db *sqlx.DB) ([]*Commit, error) {
	var comms []*Commit
	comm := NewCommit(idx.Changes)

	var batch []store.Storable
	for _, changes := range comm.GroupBy(AreCompatible) {
		comm := &Commit{Changes: changes}
		err := store.InsertIntoDB(ctx, db, comm)
		if err != nil {
			return nil, err
		}

		for _, chg := range changes {
			chg.CommitId = comm.Id
			batch = append(batch, chg)
		}
		comms = append(comms, comm)
	}
	err := store.UpdateIntoDB(ctx, db, batch...)
	if err != nil {
		return nil, err
	}
	return comms, nil
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
