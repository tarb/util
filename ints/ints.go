// Package ints implements what really should have been a std
// library package from the begining.
// How is there no asm optimized math.MaxInt?
package ints

// MinOf finds the minimum value of the arguments provided
func MinOf(vars ...int) int {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}

// MaxOf finds the maximum value of the arguments provided
func MaxOf(vars ...int) int {
	max := vars[0]

	for _, i := range vars {
		if max < i {
			max = i
		}
	}

	return max
}

// Min finds the minimum value of the arguments provided
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max finds the maximum value of the arguments provided
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
