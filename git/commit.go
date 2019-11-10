package git

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"reflect"

	"github.com/jmoiron/sqlx"
	"github.com/sebach1/git-crud/integrity"
	"github.com/sebach1/git-crud/msh"
)

// Commit is the git-like representation of a group of a ready-to-deliver signed changes
type Commit struct {
	ID      int64     `json:"id,omitempty"`
	Changes []*Change `json:"changes,omitempty"`

	ChangeIDs []int        `json:"change_ids,omitempty"`
	Reviewer  Collaborator `json:"reviewer,omitempty"`
	Errored   bool
}

func (comm *Commit) SetID(id int64) {
	comm.ID = id
}

func (comm *Commit) Table() string {
	return "commits"
}

func (comm *Commit) Columns() []string {
	return []string{
		"id",
		"errored",
		"change_ids",
	}
}

func (comm *Commit) FetchChanges(ctx context.Context, db *sqlx.DB) (err error) {
	rows, err := db.QueryxContext(ctx, `SELECT * FROM changes WHERE id=ANY($1)`, comm.ChangeIDs)
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
	err = rows.Err()
	if err != nil {
		return err
	}

	return
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

// FromCloser takes a io.ReadCloser as the guideline of a new commit
func FromCloser(body io.ReadCloser) (comm *Commit, err error) {
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
		chg := &Change{}
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

func checkIntInSlice(slice []int, elem int) bool {
	for _, sliceElem := range slice {
		if sliceElem == elem {
			return true
		}
	}
	return false
}
