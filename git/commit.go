package git

import (
	"encoding/json"
	"fmt"
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
// In case the change is invalid or is already commited, it returns an error
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
	comm.add(chg)
	return nil
}

func (comm *Commit) ToJSON() (json.RawMessage, error) {
	mapComm := comm.ToKeyVal()
	bytes, err := json.Marshal(mapComm)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(bytes), nil
}

func (comm *Commit) ToKeyVal() map[string]interface{} {
	mapComm := make(map[string]interface{})
	for _, chg := range comm.Changes {
		chgMap := chg.ToKeyVal()
		for col, val := range chgMap {
			mapComm[col] = val
		}
	}
	return mapComm
}

// func (comm *Commit) ToJSON() (json.RawMessage, error) {
// 	for _, chg := range comm.Changes {
// 		jsChg := chg.ToJSON()

// 	}
// }

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

// Rm deletes the given change from the commit
// This action is irrevertible
func (comm *Commit) Rm(chg *Change) error {
	for idx, otherChg := range comm.Changes {
		if chg.Equals(otherChg) {
			comm.rmChangeByIndex(idx)
			return nil
		}
	}
	return fmt.Errorf("change %v NOT FOUND", chg)
}

// GroupBy splits the commit changes by the given comparator cryteria
// See that strategy MUST define an equivalence relation (reflexive, transitive, symmetric)
func (comm *Commit) GroupBy(strategy changesMatcher) (grpChanges [][]*Change) {
	var omitTrans []int // Omits the transitivity of the comparisons storing the <j> element
	// Notice that <i> will not be iterated another time, so it isn't useful
	for i, chg := range comm.Changes {
		if checkIntInSlice(omitTrans, i) { // iterate only if <i> wasnt checked (due to
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

// add a change to the commit (directly appending the change to the commit changes, wout validation)
func (comm *Commit) add(chg *Change) {
	comm.Changes = append(comm.Changes, chg)
}

// rmChangeByIndex will delete without preserving order giving the desired index to delete
func (comm *Commit) rmChangeByIndex(i int) {
	var lock sync.Mutex // Avoid overlapping itself with a paralell call
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
