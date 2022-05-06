// Package money includes classes for money manipulation,
// using big.Int and big.Rat.
// Nothing is ever interpreted or represented as float, even big.Float.
//
// Taken inspiration from http://github.com/Rhymond/go-money,
// but using big.Int instead of int.
//
package money

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/karelbilek/money/decimal"
)

type Currency struct {
	Name     string
	Decimals int
}

func (c Currency) String() string {
	if c.Name == "" {
		c.Name = "UNKNOWN_CURRENCY"
	}
	return fmt.Sprintf("%s with %d decimals", c.Name, c.Decimals)
}

func (c Currency) zero() *Money {
	return &Money{
		minorAmount: new(big.Int),
		currency:    c,
	}
}

// Money represents money value, paired with currency type.
// Uses BigInt inside, but does not export it outside.
// Always saves in minor amount (different scale across currencies)
type Money struct {
	minorAmount *big.Int
	currency    Currency
}

func (m Money) Currency() Currency {
	return m.currency
}

// getMinorAmount is a helper function, so there are no panics
// on default values
func (m Money) getMinorAmount() *big.Int {
	if m.minorAmount == nil {
		return new(big.Int)
	}

	return m.minorAmount
}

// ToMinor returns value in minor units ("cents") as a string.
// The result is always representing integer.
func (m Money) ToMinor() string {
	return m.getMinorAmount().String()
}

func (m Money) DebugString() string {
	code := m.currency.Name
	if code == "" {
		code = "UNKNOWN_CURRENCY"
	}

	s := m.FormatMajor(Formatter{",", ".", GroupSizeThree, 0})
	return s + " " + code
}

// ToMajor returns value in major units as a string.
func (m Money) ToMajor() string {
	return m.FormatMajor(Formatter{"", ".", GroupSizeNone, 0})
}

// MinorInt64 returns the minor amount as a int64, or error
func (m Money) MinorInt64() (int64, error) {
	bi := m.getMinorAmount()
	if bi.IsInt64() {
		return bi.Int64(), nil
	}
	return 0, fmt.Errorf("number %s cannot be represented as int64", m.DebugString())
}

type Parser struct {
	Currency
	GroupSep string
	DecSep   string
}

// FromMajor imports from string in major units and a currency code
func FromMajor(major string, p Parser) (*Money, error) {
	if p.GroupSep == p.DecSep {
		return nil, fmt.Errorf(
			"group and decimal separator cannot be the same, are %q and %q",
			p.GroupSep, p.DecSep)
	}

	major = strings.ReplaceAll(major, p.GroupSep, "")
	major = strings.ReplaceAll(major, p.DecSep, ".")

	d, err := decimal.FromString(major)
	if err != nil {
		return nil, err
	}

	return FromMajorDecimal(d, p.Currency)
}

// FromMinorInt imports from the given int, representing minor units
func FromMinorInt(i int64, currency Currency) *Money {
	f := big.NewInt(i)
	return &Money{minorAmount: f, currency: currency}
}

// FromMinor imports from string in minor units and a currency code
func FromMinor(minor string, currency Currency) (*Money, error) {
	d, err := decimal.FromString(minor)
	if err != nil {
		return nil, err
	}

	f, err := d.BigInt(0)
	if err != nil {
		return nil, err
	}

	return &Money{minorAmount: f, currency: currency}, nil
}

// FromMajorInt imports from the given int, representing major units
func FromMajorInt(i int64, c Currency) *Money {
	dec := c.Decimals
	exp := big.NewInt(int64(math.Pow10(dec)))
	d := big.NewInt(i)
	f := (new(big.Int)).Mul(d, exp)
	return &Money{minorAmount: f, currency: c}
}

// fromMajorDecimal imports from the given decimal, representing major units
func FromMajorDecimal(f decimal.Decimal, c Currency) (*Money, error) {
	dec := c.Decimals

	d, err := f.BigInt(dec)
	if err != nil {
		return nil, err
	}

	return &Money{
		minorAmount: d,
		currency:    c,
	}, nil
}
