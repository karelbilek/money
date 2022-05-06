package decimal

import (
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"

	"github.com/wadey/go-rounding"
)

// Decimal represents a decimal number with arbitrary precision,
// that is guaranteed to be terminating. (0.33 is ok, 0.33333333.... is not OK)
//
// Internally uses rationals from big.Rat for representing numbers; however,
// the big.Rat cannot be set directly (as it could be repeating, like 1/3)
//
// Decimal representation is guaranteed to be under 400 characters,
// 200 characters the decimal part and 200 characters fractional part
//
// All the From functions return value, all the methods accept value receivers.
//
// Using https://github.com/wadey/go-rounding
type Decimal struct {
	rat      *big.Rat
	original string
}

// FromInt creates a decimal from int
func FromInt(i int64) Decimal {
	// int64 is surely under 200 length
	return Decimal{rat: big.NewRat(i, 1), original: fmt.Sprintf("%d", i)}
}

const (
	maxDecimalLen = 200
	maxFracLen    = 200
)

type parsedString struct {
	sign    string
	dec     string
	frac    string
	expSign string
	exp     string
}

// String formatted parsed number
func (ps parsedString) String() string {
	s := ps.sign + ps.dec + "." + ps.frac
	if ps.exp != "" {
		s += "e" + ps.expSign + ps.exp
	}
	return s
}

// JSON number regex
// from https://stackoverflow.com/questions/13340717/json-numbers-regular-expression
// with additionally allowing leading 0s (as we do not interpret them as octals)
var intreg = regexp.MustCompile(`^(?P<oparen>[(])?(?P<sign>-)?(?P<dec>(?:[0-9]+))(?:\.(?P<frac>[0-9]+))?(?:[eE](?P<expSign>[+-]?)(?P<exp>[0-9]+))?(?P<cparen>[)])?$`)

func parseString(s string) *parsedString {
	// number needs to be parsed manually to do checks on exponent
	// because none of the standard libraries do it
	match := intreg.FindStringSubmatch(s)
	if match == nil {
		return nil
	}
	paren := 0
	p := parsedString{}
	for i, name := range intreg.SubexpNames() {
		switch name {
		case "sign":
			p.sign = match[i]
		case "dec":
			p.dec = match[i]
		case "frac":
			p.frac = match[i]
		case "expSign":
			p.expSign = match[i]
		case "exp":
			p.exp = match[i]
		case "oparen", "cparen":
			if len(match[i]) == 0 {
				continue
			}
			paren++
			if paren == 2 {
				p.sign = "-"
			}
		}
	}
	if paren == 1 { // unmatched parens
		return nil
	}
	return &p
}

// FromString returns decimal from string.
// Allows decimal representation like "0.1e8".
func FromString(s string) (d Decimal, err error) {
	defer func() {
		if err == nil {
			d.original = s
		}
	}()

	if s == "" {
		return Decimal{rat: big.NewRat(0, 1)}, nil
	}

	parsed := parseString(s)
	if parsed == nil {
		return Decimal{}, fmt.Errorf("%s is not a valid decimal amount", s)
	}

	exp := 0
	if parsed.exp != "" {
		p, err := strconv.Atoi(parsed.exp)
		if err != nil {
			// this should never happen (parsed.exp is guaranteed to be a number from the regex)
			return Decimal{}, fmt.Errorf("unexpected error %w", err)
		}
		exp = p
	}
	decLen := len(parsed.dec)
	if parsed.expSign != "-" {
		decLen += exp
	}
	if decLen > maxDecimalLen {
		return Decimal{}, fmt.Errorf("decimal length %d bigger than allowed %d", decLen, maxDecimalLen)
	}

	fracLen := len(parsed.frac)
	if parsed.expSign == "-" {
		fracLen += exp
	}
	if fracLen > maxFracLen {
		return Decimal{}, fmt.Errorf("fractional length %d bigger than allowed %d", fracLen, maxFracLen)
	}

	rat, ok := new(big.Rat).SetString(parsed.String())
	if !ok {
		return Decimal{}, fmt.Errorf("%s is not a valid decimal amount", s)
	}

	return Decimal{rat: rat}, nil
}

// BigRat exports number as big.Rat.
// Changing the returned big.Rat does not change the original decimal.
func (v Decimal) BigRat() *big.Rat {
	if v.rat == nil {
		return new(big.Rat)
	}
	return new(big.Rat).Set(v.rat)
}

// FracDecimals returns number of fractional decimals needed to write the number.
// Returns '0' for '0', '1' for '0.5', etc.
func (v Decimal) FracDecimals() (int, error) {
	if v.rat == nil {
		return 0, nil
	}

	if !rounding.Finite(v.rat) {
		// edgecase that should never happen
		// if decimal is created by From* methods
		// but better error than panic
		return 0, fmt.Errorf("rational %s is repeating", v.rat)
	}

	// it's surely finite, this is possible
	return rounding.FinitePrec(v.rat), nil
}

// BigInt returns this number, moved decimal point right by fracDecimals, as big.Int.
func (v Decimal) BigInt(fracDecimals int) (*big.Int, error) {
	if fracDecimals < 0 {
		return nil, fmt.Errorf("cannot have negative fracDecimals %d", fracDecimals)
	}
	hasFracDecimals, err := v.FracDecimals()
	if err != nil {
		return nil, err
	}
	if hasFracDecimals > fracDecimals {
		return nil, fmt.Errorf(
			"%s has %d decimals, only %d allowed",
			v.original,
			hasFracDecimals,
			fracDecimals,
		)
	}

	f := new(big.Rat).Mul(v.BigRat(), big.NewRat(int64(math.Pow10(fracDecimals)), 1))
	if !f.IsInt() {
		// should never happen
		return nil, fmt.Errorf("int expected")
	}
	// if f IsInt -> denominator is 1, numerator is the int
	return f.Num(), nil
}

// UnmarshalJSON unmarshals from "1.23" to amount.
func (v *Decimal) UnmarshalJSON(data []byte) error {
	n, err := FromString(string(data))
	if err != nil {
		return err
	}
	v.rat = n.rat
	return nil
}
