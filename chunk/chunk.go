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
	var result = make([][2]int, n)
	var start int

	for i := range result {
		var end int = len / n
		if len%n > i {
			end++
		}

		result[i] = [2]int{start, start + end}

		start += end
	}

	return result
}
