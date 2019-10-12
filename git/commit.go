package git

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/sebach1/git-crud/internal/integrity"
)

// Commit is the abstraction that takes the proposed changes to an entity
// Actually, it can just link one entity
type Commit struct {
	ID      int       `json:"id,omitempty"`
	Changes []*Change `json:"changes,omitempty"`

	Reviewer Collaborator `json:"reviewer,omitempty"`
}

// Add will attach the given change to the commit changes
// In case the change is invalid or is already committed, it returns an error
func (comm *Commit) Add(chg *Change) error {
	err := chg.Validate()
	if err != nil {
		return err
	}
	if comm.containsChange(chg) { // checks duplication. It discards untracked changes in the comparison
		return errDuplicatedChg
	}

	for idx, otherChg := range comm.Changes { // Then check for overrides
		if Overrides(chg, otherChg) {
			comm.rmChangeByIndex(idx)
		}
	}
	comm.Changes = append(comm.Changes, chg)
	return nil
}

// Rm deletes the given change from the commit
// This action is irreversible
func (comm *Commit) Rm(chg *Change) error {
	for idx, otherChg := range comm.Changes {
		if chg.Equals(otherChg) {
			comm.rmChangeByIndex(idx)
			return nil
		}
	}
	return fmt.Errorf("change %v NOT FOUND", chg)
}

// Unmarshal the commit onto a given data structure through a given format
// ? Use reflection during unmarshal to take tags format
func (comm *Commit) Unmarshal(data interface{}, format string) error {
	switch format {
	case "json":
		rawJSON, err := integrity.ToJSON(comm)
		if err != nil {
			return err
		}
		err = json.Unmarshal(rawJSON, data)
		if err != nil {
			return err
		}
	}
	return nil
}

// CommitFrom takes a io.ReadCloser as the guideline of a new commit
func CommitFrom(body io.ReadCloser) (comm *Commit, err error) {
	var bodyMap map[string]interface{}

	err = json.NewDecoder(body).Decode(&bodyMap)
	if err != nil {
		return nil, err
	}

	err = comm.FromMap(bodyMap)
	if err != nil {
		return nil, err
	}

	return
}

// FromMap decodes the commit from its map version
// Notice that FromMap() is reciprocal to ToMap(), so it doesn't assign a table
func (comm *Commit) FromMap(Map map[string]interface{}) error {
	maybeID := Map["id"]
	ID, ok := maybeID.(integrity.ID)
	if !ok && maybeID != nil {
		return errors.New("the ENTITY_ID is NOT an ID type")
	}
	if maybeID != nil {
		delete(Map, "id")
	}

	for col, val := range Map {
		chg := new(Change)
		chg.FromMap(map[string]interface{}{col: val})
		comm.Changes = append(comm.Changes, chg)
	}

	if !ID.IsNil() {
		for _, chg := range comm.Changes {
			chg.EntityID = ID
		}
	}
	return nil
}

// ToMap returns a map with the content of the commit, omitting unnecessary fields
// It takes every change.ToMap and merges onto the resultant map
func (comm *Commit) ToMap() map[string]interface{} {
	mapComm := make(map[string]interface{})
	for _, chg := range comm.Changes {
		chgMap := chg.ToMap()
		for col, val := range chgMap {
			mapComm[col] = val
		}
	}
	return mapComm
}

// ColumnNames retrieves all ColumnName for each change
func (comm *Commit) ColumnNames() (colNames []integrity.ColumnName) {
	for _, chg := range comm.Changes {
		colNames = append(colNames, chg.ColumnName)
	}
	return
}

// TableName checks the unification of the changes' TableNames, and returns an error if there are != 1 TableName
// Returns the representative tableName of the changes
func (comm *Commit) TableName() (tableName integrity.TableName, err error) {
	for _, chg := range comm.Changes {
		if tableName != "" {
			if chg.TableName != tableName {
				return "", errMixedTables
			}
			continue
		}
		tableName = chg.TableName
	}
	return
}

// Options checks the unification of the changes' Options, and returns an error if there are not shared
// Returns the representative Options of the changes
func (comm *Commit) Options() (opts Options, err error) {
	for _, chg := range comm.Changes {
		if opts != nil {
			if !reflect.DeepEqual(chg.Options, opts) {
				return nil, errMixedOpts
			}
			continue
		}
		opts = chg.Options
	}
	return
}

// Type checks the unification of the changes' Types, and returns an error if there are != 1 Type
// Returns the representative type of the changes
func (comm *Commit) Type() (commType integrity.CRUD, err error) {
	for _, chg := range comm.Changes {
		if commType != "" {
			if chg.Type != commType {
				return "", errMixedTypes
			}
			continue
		}
		commType = chg.Type
	}
	return
}

// GroupBy splits the commit changes by the given comparator criteria
// See that strategy MUST define an equivalence relation (reflexive, transitive, symmetric)
func (comm *Commit) GroupBy(strategy changesMatcher) (grpChanges [][]*Change) {
	var omitTrans []int // Omits the transitivity of the comparisons storing the <j> element
	// Notice that <i> will not be iterated another time, so it isn't useful
	for i, chg := range comm.Changes {
		if checkIntInSlice(omitTrans, i) { // iterate only if <i> wasn't checked (due to
			// equivalence relation property we can avoid them)
			continue
		}

		iChgs := []*Change{chg}

		for j, otherChg := range comm.Changes {

			if i < j { // Checks the groupability only for all inside
				//  the upper-strict triangular form the 1-d matrix
				if strategy(chg, otherChg) {
					iChgs = append(iChgs, otherChg)
					omitTrans = append(omitTrans, j)
				}
			}

		}

		grpChanges = append(grpChanges, iChgs)
	}
	return
}

// rmChangeByIndex will delete without preserving order giving the desired index to delete
func (comm *Commit) rmChangeByIndex(i int) {
	var lock sync.Mutex // Avoid overlapping itself with a parallel call
	lock.Lock()
	lastIndex := len(comm.Changes) - 1
	comm.Changes[i] = comm.Changes[lastIndex]
	comm.Changes[lastIndex] = nil // Notices the GC to rm the last elem to avoid mem-leak
	comm.Changes = comm.Changes[:lastIndex]
	lock.Unlock()
}

// containsChange verifies if the given change is already present, and triggering the **exactly** same action
func (comm *Commit) containsChange(chg *Change) bool {
	for _, otherChg := range comm.Changes {
		if chg.Equals(otherChg) {
			return true
		}
	}
	return false
}

func checkIntInSlice(slice []int, elem int) bool {
	for _, sliceElem := range slice {
		if sliceElem == elem {
			return true
		}
	}
	return false
}
