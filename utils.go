package main

func indexOf[E any](e *E, arr []*E) int {
	for i, a := range arr {
		if a == e {
			return i
		}
	}
	return -1
}

func contains[E any](e *E, arr []*E) bool {
	return indexOf(e, arr) > -1
}
