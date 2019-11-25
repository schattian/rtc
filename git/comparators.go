package git

import "reflect"

type changesMatcher func(*Change, *Change) bool

// AreCompatible checks if two changes could be joint to perform any action
func AreCompatible(chg, otherChg *Change) bool {
	if chg.TableName != otherChg.TableName {
		return false
	}
	if chg.EntityId != otherChg.EntityId {
		return false
	}
	if chg.Type == "create" { // In case of both of EntityIds are nil
		//  (see that the above comparison discards 2x checking)
		return false
	}
	if chg.Type != otherChg.Type {
		return false
	}
	if !reflect.DeepEqual(chg.Options, otherChg.Options) {
		return false
	}
	return true
}
