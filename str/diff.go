package str

// Diff returns the difference between the 2 slices
func Diff(sa, sb []string) []string {
	var diff []string

	for i := 0; i < 2; i++ {
		for _, s1 := range sa {
			found := false
			for _, s2 := range sb {
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
			sa, sb = sb, sa
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
