package namespace

// containsString checks if a string is in a slice of strings.
func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// mapsOfStringAreDifferent checks if two maps of strings are different.
func mapsOfStringAreDifferent(currentLabels, newLabels map[string]string) bool {

	if len(currentLabels) != len(newLabels) {
		return true
	}

	for key, value := range newLabels {
		if currentValue, exists := currentLabels[key]; !exists || currentValue != value {
			return true
		}
	}

	return false
}
