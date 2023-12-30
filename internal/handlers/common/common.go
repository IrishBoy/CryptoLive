package common

func AddString(set map[string]bool, s string) {
	set[s] = true
}

// Function to get unique values from the set
func GetUniqueValues(set map[string]bool) []string {
	uniqueValues := make([]string, 0, len(set))
	for key := range set {
		uniqueValues = append(uniqueValues, key)
	}
	return uniqueValues
}
