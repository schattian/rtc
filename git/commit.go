package git

import (
	"context"
	"encoding/json"
	"io"
	"reflect"

	"github.com/jmoiron/sqlx"
	"github.com/sebach1/rtc/integrity"
	"github.com/sebach1/rtc/msh"
)

// Commit is the git-like representation of a group of a ready-to-deliver signed changes
type Commit struct {
	Id      int64     `json:"id,omitempty"`
	Changes []*Change `json:"changes,omitempty"`

	Reviewer Collaborator `json:"reviewer,omitempty"`

	Errored bool `json:"errored,omitempty"`
}

func NewCommit(changes []*Change) *Commit {
	return &Commit{Changes: changes}
}

// FetchChanges retrieves the changes from DB by its .ChangeIds and assigns them to .Changes field
func (comm *Commit) FetchChanges(ctx context.Context, db *sqlx.DB) (err error) {
	rows, err := db.NamedQueryContext(ctx, `SELECT * FROM changes WHERE commit_id=:id`, comm)
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
		comm.Changes = append(comm.Changes, &chg)
	}
	return rows.Err()
}

// GroupBy splits the commit's changes by the given comparator criteria
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

// Unmarshal the commit onto a given data structure through a given format
// ? Use reflection during unmarshal to take tags format
func (comm *Commit) Unmarshal(data interface{}, format string) error {
	switch format {
	case "json":
		rawJSON, err := msh.ToJSON(comm)
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

// CommitFromCloser takes a io.ReadCloser as the guideline of a new commit
func CommitFromCloser(body io.ReadCloser) (comm *Commit, err error) {
	var bodyMap map[string]interface{}

	err = json.NewDecoder(body).Decode(&bodyMap)
	if err != nil {
		return nil, err
	}

	comm, err = CommitFromMap(bodyMap)
	if err != nil {
		return nil, err
	}

	return
}

// CommitFromMap decodes the commit from its map version
// Notice that Commit.FromMap() is reciprocal to ToMap(), so it doesn't assign a table
func CommitFromMap(Map map[string]interface{}) (comm *Commit, err error) {
	maybeId := Map["id"]
	Id, ok := maybeId.(integrity.Id)
	if !ok && maybeId != nil {
		return nil, errInvalidCommitId
	}
	if maybeId != nil {
		delete(Map, "id")
	}

	for col, val := range Map {
		chg := &Change{}
		err := chg.FromMap(map[string]interface{}{col: val})
		if err != nil {
			return nil, err
		}
		comm.Changes = append(comm.Changes, chg)
	}

	if !Id.IsNil() {
		for _, chg := range comm.Changes {
			chg.EntityId = Id
		}
	}
	return
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

func checkIntInSlice(slice []int, elem int) bool {
	for _, sliceElem := range slice {
		if sliceElem == elem {
			return true
		}
	}
	return false
}

func (comm *Commit) copy() *Commit {
	if comm == nil {
		return nil
	}
	newComm := new(Commit)
	*newComm = *comm
	var newChgs []*Change
	for _, chg := range comm.Changes {
		newChgs = append(newChgs, chg.copy())
	}
	newComm.Changes = newChgs
	return newComm
}
