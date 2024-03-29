package utils

func RemoveDuplicateValues(data []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range data {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
