package valueobject

import "slices"

type CPF string

func (cpf CPF) IsValid() bool {
	var sum int
	var rest int
	var weight int
	var firstDigit int
	var secondDigit int
	blockList := []string{
		"00000000000",
		"11111111111",
		"22222222222",
		"33333333333",
		"44444444444",
		"55555555555",
		"66666666666",
		"77777777777",
		"88888888888",
		"99999999999",
	}

	if slices.Contains(blockList, string(cpf)) {
		return false
	}

	for i := 1; i <= 9; i++ {
		weight = 11 - i
		sum += weight * int(cpf[i-1]-'0')
	}

	rest = (sum * 10) % 11

	if rest == 10 || rest == 11 {
		firstDigit = 0
	} else {
		firstDigit = rest
	}

	if firstDigit != int(cpf[9]-'0') {
		return false
	}

	sum = 0

	for i := 1; i <= 10; i++ {
		weight = 12 - i
		sum += weight * int(cpf[i-1]-'0')
	}

	rest = (sum * 10) % 11

	if rest == 10 || rest == 11 {
		secondDigit = 0
	} else {
		secondDigit = rest
	}

	if secondDigit != int(cpf[10]-'0') {
		return false
	}

	return true
}
