package money_test

import (
	"fmt"
	"log"

	"github.com/karelbilek/money"

	"github.com/wadey/go-rounding"
)

var php = money.Currency{
	Name:     "PHP",
	Decimals: 2,
}

var phpParser = money.Parser{
	Currency: php,
	GroupSep: ",",
	DecSep:   ".",
}

var weirdFormatter = money.Formatter{
	GroupSep:    " ",
	DecSep:      ",",
	MinDecimals: 1,
	GroupSize:   money.GroupSizeIndian,
}

func ExampleMoney() {
	inStr := `-25.33`

	paid, err := money.FromMajor(inStr, phpParser)
	if err != nil {
		log.Fatal(err)
	}

	multiplied, err := paid.MultiplyRat("12345.4567", rounding.HalfUp)
	if err != nil {
		log.Fatal(err)
	}

	r := multiplied.RoundToMajor(rounding.Floor)

	p := r.FormatMajor(weirdFormatter)

	fmt.Println(p)

	// Output: -3 12 711,0
}
