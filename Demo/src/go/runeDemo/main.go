package main

func balancedStringSplit(s string) int {
	result := 0
	countR := 0
	countL := 0
	for i := 0; i < len(s); i++ {
		if s[i] == uint8('R') {
			countR++
		} else {
			countL++
		}
		if countR == countL {
			result++
			countR = 0
			countL = 0
		}

	}

	return result

}

func main() {
	println(balancedStringSplit("RLRRLLRLRL")) // 输出 4
	println(balancedStringSplit("RLLLLRRRLR")) // 输出 3
	println(balancedStringSplit("LLLLRRRR"))   // 输出 1
	println(balancedStringSplit("RLRRRLLRLL")) // 输出 2
}
