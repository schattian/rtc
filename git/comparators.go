package git

type changesMatcher func(*Change, *Change) bool

// AreCompatible checks if two changes could be joint to perform any action
func AreCompatible(chg, otherChg *Change) bool {
	if chg.TableName != otherChg.TableName {
		return false
	}
	if chg.EntityID != otherChg.EntityID {
		return false
	}
	if chg.Type == "create" { // In case of both of EntityIDs are nil
		//  (see that the above comparison discards 2x checking)
		return false
	}
	if chg.Type != otherChg.Type {
		return false
	}
	return true
}
