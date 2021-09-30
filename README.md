money
=====

Go package for representing money.

Example:

```go
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
```

Goals:
* explicit better than implicit
* simple API
* error on suspicious things
* well tested
* not really caring about performance. (not unbearably slow, but
  there might be more performant packages)

Basics:
* Money is, internally, big integer + currency.
* Currency is name + number of decimals.
* You cannot compare, add, subtract two Money values with different currencies.
* Both parsing and exporting functions always differentiate
  between "Major" and "Minor". Minor meaning "cents", major meaning "dollars".
* Division either needs "rounding"; or you need to call "split", which
  distributes extra cents round-robin.

I have found that it's better to not include
any code that ties the currency name to its ISO counterpart (for precision),
or its locale (for parsing and formatting).
As that's not explicit, programmers (me included) will get confused,
what exactly does the currency/locale affect.

Real example: when the ISO lists were used, Indonesian rupiah had 0 decimals,
but then it was needed to allow for 2 decimals for accounting. On the
other hand, I still don't want to allow infinite decimals, even for accounting.

So, everything needs to be always explicit, for the principle of
no surprises.

Uses github.com/wadey/go-rounding for big.Rat rounding.

Idea of split/allocate from github.com/rhymond/go-money

See `example_test.go` for simple usage.

(C) 2021 Karel Bilek, Chad Kunde

MIT license