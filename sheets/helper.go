package sheets

func columnToLetter(column int) string {
	letter := ""
	for column > 0 {
		temp := (column - 1) % 26
		letter = string(rune(temp+65)) + letter
		column = (column - temp - 1) / 26
	}
	return letter
}

func letterToColumn(letter string) int {
	column := 0
	length := len(letter)
	for i := 0; i < length; i++ {
		column += (int(letter[i]) - 64) * pow(26, length-i-1)
	}
	return column
}

func pow(x, y int) int {
	result := 1
	for y > 0 {
		result *= x
		y--
	}
	return result
}
