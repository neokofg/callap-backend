package utils

func GetOrderedIds(id1, id2 string) (string, string) {
	if id1 < id2 {
		return id1, id2
	}
	return id2, id1
}
