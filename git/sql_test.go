package git

func isExcluded(exclusions []string, name string) bool {
	for _, exc := range exclusions {
		if exc == name {
			return true
		}
	}
	return false
}
