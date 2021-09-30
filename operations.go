package money

import (
	"fmt"
	"math"
	"math/big"

	"github.com/wadey/go-rounding"
)

// some operations

var bigOne = big.NewInt(1)

// temporary, before https://github.com/wadey/go-rounding/pull/6 is merged
func roundToInt(rat *big.Rat, method rounding.RoundingMode) *big.Int {
	c := new(big.Rat).Set(rat)

	rounding.Round(c, 0, method)

	if !c.IsInt() {
		// this should never happen
		panic(fmt.Errorf("unexpected unrounded rational"))
	}

	// as c is int, denominator is 1 and numerator is the int value
	return c.Num()
}

// SameCurrency check if given Money has the same currency
func (m Money) SameCurrency(om Money) bool {
	return m.currency == om.currency // note that comparing structs compares fields
}

func (m Money) assertSameCurrency(om Money) error {
	same := m.SameCurrency(om)

	if !same {
		return fmt.Errorf(
			"currencies %s and %s don't match",
			m.currency.String(),
			om.currency.String(),
		)
	}

	return nil
}

func (m Money) compare(om Money) int {
	return m.getMinorAmount().Cmp(om.getMinorAmount())
}

// Equals checks equality between two Money types.
// Returns error if different currencies.
func (m Money) Equals(om Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 0, nil
}

// More checks whether the value of Money is greater than the other.
// Returns error if different currencies.
func (m Money) More(om Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == 1, nil
}

// MoreEqual checks whether the value of Money is greater or equal than the other.
// Returns error if different currencies.
func (m Money) MoreEqual(om Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) >= 0, nil
}

// Less checks whether the value of Money is less than the other.
// Returns error if different currencies.
func (m Money) Less(om Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) == -1, nil
}

// Between checks whether the value of Money is between two values, both ends inclusive.
// Returns error if different currencies, or min > max
func (m Money) Between(min, max Money) (bool, error) {
	isInterval, err := min.LessEqual(max)
	if err != nil {
		return false, err
	}
	if !isInterval {
		return false, fmt.Errorf("minimal is bigger than maximal")
	}
	more, err := m.MoreEqual(min)
	if err != nil {
		return false, err
	}
	less, err := m.LessEqual(max)
	if err != nil {
		return false, err
	}
	return more && less, nil
}

// LessEqual checks whether the value of Money is less or equal than the other.
// Returns error if different currencies.
func (m Money) LessEqual(om Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.compare(om) <= 0, nil
}

// IsZero returns boolean of whether the value of Money is equals to zero.
func (m Money) IsZero() bool {
	return m.getMinorAmount().Cmp(new(big.Int)) == 0
}

// IsPositive returns boolean of whether the value of Money is positive.
func (m Money) IsPositive() bool {
	return m.getMinorAmount().Cmp(new(big.Int)) > 0
}

// IsNegative returns boolean of whether the value of Money is negative.
func (m Money) IsNegative() bool {
	return m.getMinorAmount().Cmp(new(big.Int)) < 0
}

// Absolute returns new Money struct from given Money using absolute monetary value.
func (m Money) Absolute() *Money {
	r := (new(big.Int)).Abs(m.getMinorAmount())
	return &Money{minorAmount: r, currency: m.currency}
}

// Negative returns new Money struct from given Money using negative monetary value.
func (m Money) Negative() *Money {
	r := (new(big.Int)).Neg(m.getMinorAmount())
	return &Money{minorAmount: r, currency: m.currency}
}

// Add returns new Money struct with value representing sum of Self and Other Money.
// Returns error if different currencies.
func (m Money) Add(om Money) (*Money, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return nil, err
	}

	r := (new(big.Int)).Add(m.getMinorAmount(), om.getMinorAmount())

	return &Money{minorAmount: r, currency: m.currency}, nil
}

// Subtract returns new Money struct with value representing difference of Self and Other Money.
func (m Money) Subtract(om Money) (*Money, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return nil, err
	}

	omNeg := (new(big.Int)).Neg(om.getMinorAmount())
	r := (new(big.Int)).Add(m.getMinorAmount(), omNeg)

	return &Money{minorAmount: r, currency: m.currency}, nil
}

// Multiply returns new Money struct with value representing Self multiplied value by multiplier.
func (m Money) Multiply(mul int64) (*Money, error) {
	r := (new(big.Int)).Mul(m.getMinorAmount(), big.NewInt(mul))
	return &Money{minorAmount: r, currency: m.currency}, nil
}

// Divide divides. NOTE: You might want to use Split instead.
func (m Money) Divide(d int64, round rounding.RoundingMode) (*Money, error) {
	if d == 0 {
		return nil, fmt.Errorf("division by zero")
	}

	rat := (new(big.Rat)).SetFrac(bigOne, big.NewInt(d))

	return m.MultiplyBigRat(rat, round), nil
}

// MultiplyRat returns new Money struct with value representing Self multiplied value by multiplier,
// that is a rational string.
// All inputs that are allowed in big.Rat are allowed - including "1/3"
func (m Money) MultiplyRat(ratString string, round rounding.RoundingMode) (*Money, error) {
	rat, ok := (new(big.Rat)).SetString(ratString)
	if !ok {
		return nil, fmt.Errorf("%s is not a valid rational amount", ratString)
	}

	return m.MultiplyBigRat(rat, round), nil
}

// MultiplyBigRat returns new Money struct with value representing Self multiplied value by
// big.Rat value
func (m Money) MultiplyBigRat(rat *big.Rat, round rounding.RoundingMode) *Money {
	if rat == nil {
		return m.currency.zero()
	}

	a := (new(big.Rat)).SetInt(m.getMinorAmount())

	c := a.Mul(a, rat)
	minor := roundToInt(c, round)

	return &Money{minorAmount: minor, currency: m.currency}
}

// DivideRat returns new Money struct with value representing Self divided value by divider,
// that is a rational string..
// All inputs that are allowed in big.Rat are allowed -
// including "1/3" - which would multiply by 3.
func (m Money) DivideRat(ratString string, round rounding.RoundingMode) (*Money, error) {
	rat, ok := (new(big.Rat)).SetString(ratString)
	if !ok {
		return nil, fmt.Errorf("%s is not a valid rational amount", ratString)
	}
	if rat.Sign() == 0 {
		return nil, fmt.Errorf("division by zero")
	}

	return m.DivideBigRat(rat, round)
}

// DivideBigRat returns new Money struct with value representing Self divided value by big.Rat.
func (m Money) DivideBigRat(rat *big.Rat, round rounding.RoundingMode) (*Money, error) {
	if rat == nil {
		return nil, fmt.Errorf("nil rational")
	}
	if rat.Sign() == 0 {
		return nil, fmt.Errorf("division by zero")
	}

	a := (new(big.Rat)).SetInt(m.getMinorAmount())
	c := a.Quo(a, rat)
	minor := roundToInt(c, round)
	return &Money{minorAmount: minor, currency: m.currency}, nil
}

// HasCents returns true if the value is not rounded to major units
func (m Money) HasCents() (bool, error) {
	dec := m.currency.Decimals

	if dec == 0 {
		return false, nil
	}

	rounded := m.RoundToMajor(rounding.Down)

	k, err := m.Subtract(*rounded)
	if err != nil {
		return false, err
	}
	isZero := k.IsZero()
	return !isZero, nil
}

// LessMajor checks whether the value of Money is less than an int, in major value.
func (m Money) LessMajor(i int64) bool {
	other := FromMajorInt(i, m.currency)
	return m.compare(*other) == -1
}

// MoreMajor checks whether the value of Money is more than an int, in major value.
func (m Money) MoreMajor(i int64) bool {
	other := FromMajorInt(i, m.currency)

	return m.compare(*other) == 1
}

// LessEqualMajor checks whether the value of Money is less or equal than an int, in major value.
func (m Money) LessEqualMajor(i int64) bool {
	other := FromMajorInt(i, m.currency)

	return m.compare(*other) <= 0
}

// MoreEqualMajor checks whether the value of Money is more or equal than an int, in major value.
func (m Money) MoreEqualMajor(i int64) bool {
	other := FromMajorInt(i, m.currency)

	return m.compare(*other) >= 0
}

// BetweenMajor checks whether the value of Money is between two ints, in major value.
// Returns error if min > max
func (m Money) BetweenMajor(min, max int64) (bool, error) {
	minM := FromMajorInt(min, m.currency)
	maxM := FromMajorInt(max, m.currency)

	is, err := m.Between(*minM, *maxM)
	if err != nil {
		return false, err
	}
	return is, nil
}

// LessMinor checks whether the value of Money is less than an int, in minor value.
func (m Money) LessMinor(i int64) bool {
	other := FromMinorInt(i, m.currency)

	return m.compare(*other) == -1
}

// MoreMinor checks whether the value of Money is more than an int, in minor value.
func (m Money) MoreMinor(i int64) bool {
	other := FromMinorInt(i, m.currency)

	return m.compare(*other) == 1
}

// LessEqualMinor checks whether the value of Money is less or equal than an int, in minor value.
func (m Money) LessEqualMinor(i int64) bool {
	other := FromMinorInt(i, m.currency)

	return m.compare(*other) <= 0
}

// MoreEqualMinor checks whether the value of Money is more or equal than an int, in minor value.
func (m Money) MoreEqualMinor(i int64) bool {
	other := FromMinorInt(i, m.currency)

	return m.compare(*other) >= 0
}

// BetweenMinor checks whether the value of Money is between two ints, in minor value.
// Returns error if min > max
func (m Money) BetweenMinor(min, max int64) (bool, error) {
	minM := FromMinorInt(min, m.currency)
	maxM := FromMinorInt(max, m.currency)

	is, err := m.Between(*minM, *maxM)
	if err != nil {
		return false, err
	}
	return is, nil
}

// RoundToMajor returns new Money struct with value rounded to major unit and with a given
// strategy
func (m Money) RoundToMajor(round rounding.RoundingMode) *Money {
	dec := m.currency.Decimals

	if dec == 0 {
		return &m
	}

	exp := big.NewInt(int64(math.Pow10(dec)))
	rat := (new(big.Rat)).SetFrac(m.getMinorAmount(), exp)

	i := roundToInt(rat, round)
	i = i.Mul(i, exp)

	return &Money{minorAmount: i, currency: m.currency}
}

// Split returns slice of Money structs with split Self value in given number.
// After division leftover pennies will be distributed round-robin amongst the parties.
// This means that parties listed first can receive more pennies than ones that are listed later.
func (m Money) Split(n int) ([]*Money, error) {
	if n <= 0 {
		return nil, fmt.Errorf("split must be higher than zero, is %d", n)
	}

	fl := m.getMinorAmount()
	neg := m.IsNegative()

	if neg {
		fl = (new(big.Int)).Neg(fl)
	}

	quo, rem := (new(big.Int)).QuoRem(fl, big.NewInt(int64(n)), new(big.Int))
	ms := make([]*Money, 0, n)

	for i := 0; i < n; i++ {
		ms = append(ms, &Money{minorAmount: (new(big.Int)).Set(quo), currency: m.currency})
	}

	l := int(rem.Int64())
	// Add leftovers to the first parties.
	for p := 0; l != 0; p++ {
		ms[p].minorAmount = ms[p].minorAmount.Add(ms[p].minorAmount, bigOne)
		l--
	}

	if neg {
		for i := 0; i < n; i++ {
			ms[i].minorAmount = ms[i].minorAmount.Neg(ms[i].minorAmount)
		}
	}

	return ms, nil
}

// Allocate returns slice of Money structs with split Self value in given ratios.
// It lets split money by given ratios without losing pennies and as Split operations distributes
// leftover pennies amongst the parties with round-robin principle.
func (m Money) Allocate(rs ...int) ([]*Money, error) {
	if len(rs) == 0 {
		return nil, fmt.Errorf("no ratios specified")
	}

	fl := m.getMinorAmount()
	neg := m.IsNegative()

	if neg {
		fl = (new(big.Int)).Neg(fl)
	}

	// Calculate sum of ratios.
	var sum int
	for _, r := range rs {
		sum += r
	}

	total := new(big.Int)
	ms := make([]*Money, 0, len(rs))

	for _, r := range rs {
		mul := (new(big.Int)).Mul(fl, big.NewInt(int64(r)))
		quo := (new(big.Int)).Quo(mul, big.NewInt(int64(sum)))

		party := &Money{
			minorAmount: quo,
			currency:    m.currency,
		}

		ms = append(ms, party)
		total = total.Add(total, quo)
	}

	// Calculate leftover value and divide to first parties.
	lo := (new(big.Int)).Sub(fl, total)

	l := lo.Int64()
	// Add leftovers to the first parties.
	for p := 0; l != 0; p++ {
		ms[p].minorAmount = ms[p].minorAmount.Add(ms[p].minorAmount, bigOne)
		l--
	}

	if neg {
		for i := 0; i < len(rs); i++ {
			ms[i].minorAmount = ms[i].minorAmount.Neg(ms[i].minorAmount)
		}
	}

	return ms, nil
}
