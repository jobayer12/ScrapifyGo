package utils

func RemoveDuplicates(items []string) []string {
	// Create a map to track unique items
	uniqueMap := make(map[string]bool)

	// Iterate over the items and add them to the map
	for _, item := range items {
		uniqueMap[item] = true
	}

	// Create a new slice to store the unique items
	var uniqueItems []string

	// Append the unique items to the new slice
	for item := range uniqueMap {
		uniqueItems = append(uniqueItems, item)
	}
	return uniqueItems

}
