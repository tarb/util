package chunk

// Chunk returns a slice of index pairs that can be used to equally divide
// a group. For example a slice of len 11 being chunked into 3 chunks will
// return the indexs [[0 4],[4 8],[8 11]]
//
// nums := []int{0,1,2,3,4,5,6,7,8,9,10}
// for _, n := range Chunk(len(nums), 3) {
//     fmt.Println(nums[n[0]:n[1]])
// }
func Chunk(len, n int) [][2]int {
	result := make([][2]int, n)
	start := 0

	for i := range result {
		end := len / n
		if len%n > i {
			end++
		}

		result[i] = [2]int{start, start + end}

		start += end
	}

	return result
}

// NumChunk returns how many chunks you would need to fit n items into len
// For example,
//	a slice with len 1000 with n of 200 would return 5
//  a slice with len 801 with n of 200 would return 5
//  a slice with len 1001 with n of 200 would return 6
func NumChunk(len, n int) int {
	if len%n == 0 {
		return len / n
	}

	return len/n + 1
}
