package slices

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

// ElementToFirst removes the element from the slice and places it as first element.
func ElementToFirst(slice []string, element string) []string {
	index := -1
	for i, v := range slice {
		if v == element {
			index = i
			break
		}
	}

	if index == -1 {
		return slice
	}

	newSlice := []string{element}

	newSlice = append(newSlice, remove(slice, index)...)

	return newSlice
}
