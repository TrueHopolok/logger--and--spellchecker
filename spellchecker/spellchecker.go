package spellchecker

func FindScore(s1, s2 string) int {
	if s1 == s2 {
		return 0
	}
	size1, size2 := len(s1)+1, len(s2)+1
	if size1 == 1 {
		return size2 - 1
	}
	if size2 == 1 {
		return size1 - 1
	}
	table := make([][]int, size1)
	for i := range size1 {
		table[i] = make([]int, size2)
		table[i][0] = i
	}
	for j := range size2 {
		table[0][j] = j
	}
	for i := 1; i < size1; i++ {
		for j := 1; j < size2; j++ {
			table[i][j] = min(table[i][j-1], table[i-1][j]) + 1
			if table[i][j] > table[i-1][j-1] {
				table[i][j] = table[i-1][j-1]
				if s1[i-1] != s2[j-1] {
					table[i][j]++
				}
			}
		}
	}
	return table[size1-1][size2-1]
}
