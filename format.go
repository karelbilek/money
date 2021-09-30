package money

import (
	"math/big"
	"strings"
)

type Formatter struct {
	GroupSep    string
	DecSep      string
	GroupSize   GroupSize
	MinDecimals int
}

// GroupSize gives thousand chunk size for i-th thousands,
// starting with 0. If returns <= 0 -> stops separating
type GroupSize func(i int) int

// GroupSizeNone does no separation
func GroupSizeNone(_ int) int {
	return 0
}

// GroupSizeThree always splits on three
func GroupSizeThree(_ int) int {
	return 3
}

// GroupSizeIndian is indian number system - first sep 3, then all 2
func GroupSizeIndian(i int) int {
	if i == 0 {
		return 3
	}
	return 2
}

// FormatMajor takes the decimal and nicely formats it with
// thousandsSep as a thousand separator.
// Note that no fractional decimals are printed if decimal has none
func (m Money) FormatMajor(f Formatter) string {
	if m.minorAmount == nil {
		return m.currency.zero().FormatMajor(f)
	}

	a := m.minorAmount
	isNeg := false

	if a.Sign() == -1 {
		isNeg = true
		a = new(big.Int).Neg(m.minorAmount)
	}

	allStr := a.String()
	var intStr, decStr string

	decimals := m.currency.Decimals
	if decimals == 0 {
		intStr = allStr
		decStr = ""
	} else {
		if len(allStr) < decimals {
			intStr = "0"
			decStr = strings.Repeat("0", decimals-len(allStr)) + allStr
		} else {
			intStr = allStr[0 : len(allStr)-decimals]
			decStr = allStr[len(allStr)-decimals:]
			if intStr == "" {
				intStr = "0"
			}
		}
		decStr = strings.TrimRight(decStr, "0")
	}

	if len(decStr) < f.MinDecimals {
		decStr += strings.Repeat("0", f.MinDecimals-len(decStr))
	}
	if decStr != "" {
		decStr = f.DecSep + decStr
	}

	sign := ""

	i := 0
	b := len(intStr)
THOUSANDS:
	for {
		currSize := f.GroupSize(i)
		if currSize <= 0 {
			break THOUSANDS
		}
		i++
		b -= currSize
		if b > 0 {
			intStr = intStr[:b] + f.GroupSep + intStr[b:]
		} else {
			break THOUSANDS
		}
	}

	if isNeg {
		sign = "-"
	}

	return sign + intStr + decStr
}
