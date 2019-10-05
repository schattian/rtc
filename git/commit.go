package schematypes

// Commit is the abstraction that takes the proposed changes to an entity
// Actually, it can just link one entity
type Commit struct {
	ID      int       `json:"id,omitempty"`
	Changes []*Change `json:"changes,omitempty"`
}

// GroupBy splits the commit changes by the given comparator cryteria
// See that comparator MUST define an equivalence relation (reflexive, transitive, symmetric)
func (comm *Commit) GroupBy(comparator func(*Change, *Change) bool) (grpChanges [][]*Change) {
	var omitTrans []int // Omits the transitivity of the comparisons storing the <j> element
	// Notice that <i> will not be iterated another time, so it isn't useful
	for i, chg := range comm.Changes {
		if checkIntInSlice(omitTrans, i) { // iterate <-> <i> wasnt checked (due to
			// equivalence relation property we can avoid them)
			continue
		}

		iChgs := []*Change{chg}

		for j, otherChg := range comm.Changes {

			if i < j { // Checks the groupability only for all inside
				//  the upper-strict triangular form the 1-d matrix
				if comparator(chg, otherChg) {
					iChgs = append(iChgs, otherChg)
					omitTrans = append(omitTrans, j)
				}
			}

		}

		grpChanges = append(grpChanges, iChgs)
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
