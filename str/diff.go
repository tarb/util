package str

// Diff returns the difference between the 2 slices
func Diff(s1, s2 []string) []string {
	var diff []string

	for i := 0; i < 2; i++ {
		for _, s1 := range s1 {
			found := false
			for _, s2 := range s2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			if !found {
				diff = append(diff, s1)
			}
		}

		if i == 0 {
			s1, s2 = s2, s1
		}
	}

	return diff
}

//
func IndexSlice(s string, sl []string) int {
	for i := range sl {
		if sl[i] == s {
			return i
		}
	}

	return -1
}

//
func Filter(sl []string, f func(string) bool) []string {
	nl := make([]string, 0, len(sl))

	for _, s := range sl {
		if f(s) {
			nl = append(nl, s)
		}
	}

	return nl
}
