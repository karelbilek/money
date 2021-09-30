package money

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

var php = Currency{
	Name:     "PHP",
	Decimals: 2,
}

var vnd = Currency{
	Name:     "VND",
	Decimals: 0,
}

var bhd = Currency{
	Name:     "BHD",
	Decimals: 3,
}

func TestMoney_ToMinor(t *testing.T) {
	t.Run("normal - 2 decimals", func(t *testing.T) {
		j := &Money{minorAmount: big.NewInt(12345), currency: php}
		r := j.ToMinor()
		require.Equal(t, "12345", r)
	})

	t.Run("normal - 0 decimals", func(t *testing.T) {
		j := &Money{minorAmount: big.NewInt(12345), currency: vnd}
		r := j.ToMinor()
		require.Equal(t, "12345", r)
	})

	t.Run("default", func(t *testing.T) {
		j := &Money{}
		r := j.ToMinor()
		require.Equal(t, "0", r)
	})
}

func TestMoney_ToMajor(t *testing.T) {
	t.Run("normal - 2 decimals", func(t *testing.T) {
		j := &Money{minorAmount: big.NewInt(12345), currency: php}
		r := j.ToMajor()
		require.Equal(t, "123.45", r)
	})

	t.Run("normal - 0 decimals", func(t *testing.T) {
		j := &Money{minorAmount: big.NewInt(12345), currency: vnd}
		r := j.ToMajor()
		require.Equal(t, "12345", r)
	})

	t.Run("default not panics", func(t *testing.T) {
		j := &Money{}
		r := j.ToMajor()
		require.Equal(t, "0", r)
	})
}

var defFormatter = Formatter{
	GroupSep:    ",",
	DecSep:      ".",
	GroupSize:   GroupSizeIndian,
	MinDecimals: 2,
}

func TestMoney_FormatMajor(t *testing.T) {
	t.Run("default not panics", func(t *testing.T) {
		j := &Money{}
		r := j.FormatMajor(defFormatter)
		require.Equal(t, "0.00", r)
	})

	t.Run("normal - 2 decimals, thousand sep", func(t *testing.T) {
		j := &Money{minorAmount: big.NewInt(-123456789012), currency: php}
		r := j.FormatMajor(Formatter{
			GroupSep:    "  ",
			DecSep:      ",",
			GroupSize:   GroupSizeThree,
			MinDecimals: 2,
		})
		require.Equal(t, "-1  234  567  890,12", r)
	})

	t.Run("normal - 2 decimals, thousand sep, short", func(t *testing.T) {
		j := &Money{minorAmount: big.NewInt(-12), currency: php}
		r := j.FormatMajor(Formatter{
			GroupSep:    "_",
			DecSep:      ".",
			GroupSize:   GroupSizeThree,
			MinDecimals: 2,
		})
		require.Equal(t, "-0.12", r)
	})

	t.Run("normal - 2 decimals, indian sep", func(t *testing.T) {
		j := &Money{minorAmount: big.NewInt(-123456789000), currency: php}
		r := j.FormatMajor(Formatter{
			GroupSep:    ",",
			DecSep:      "..",
			GroupSize:   GroupSizeIndian,
			MinDecimals: 2,
		})
		require.Equal(t, "-1,23,45,67,890..00", r)
	})
}

func TestMoney_DebugString(t *testing.T) {
	t.Run("normal - 2 decimals", func(t *testing.T) {
		j := &Money{minorAmount: big.NewInt(12345), currency: php}
		r := j.DebugString()
		require.Equal(t, "123.45 PHP", r)
	})

	t.Run("normal - 0 decimals", func(t *testing.T) {
		j := &Money{minorAmount: big.NewInt(12345), currency: vnd}
		r := j.DebugString()
		require.Equal(t, "12,345 VND", r)
	})

	t.Run("default not panics", func(t *testing.T) {
		j := &Money{}
		r := j.DebugString()
		require.Equal(t, "0 UNKNOWN_CURRENCY", r)
	})
}

func TestFromMinor(t *testing.T) {
	t.Run("normal - 2 decimals", func(t *testing.T) {
		m, err := FromMinor("12345", php)
		require.NoError(t, err)
		require.NotNil(t, m)
		require.EqualValues(t, m.currency, php)
		require.NotNil(t, m.minorAmount)
		require.Zero(t, m.minorAmount.Cmp(big.NewInt(12345)))
	})

	t.Run("normal - 0 decimals, negative", func(t *testing.T) {
		m, err := FromMinor("-12345", vnd)
		require.NoError(t, err)
		require.NotNil(t, m)
		require.EqualValues(t, m.currency, vnd)
		require.NotNil(t, m.minorAmount)
		require.Zero(t, m.minorAmount.Cmp(big.NewInt(-12345)))
	})

	t.Run("nonsense", func(t *testing.T) {
		m, err := FromMinor("foo 123", vnd)
		require.EqualError(t, err, "foo 123 is not a valid decimal amount")
		require.Nil(t, m)
	})

	t.Run("frac format", func(t *testing.T) {
		m, err := FromMinor("2/5", vnd)
		require.EqualError(t, err, "2/5 is not a valid decimal amount")
		require.Nil(t, m)
	})

	t.Run("1e5 format", func(t *testing.T) {
		m, err := FromMinor("1e5", vnd)
		require.NoError(t, err)
		require.NotNil(t, m)
		require.EqualValues(t, m.currency, vnd)
		require.NotNil(t, m.minorAmount)
		require.Zero(t, m.minorAmount.Cmp(big.NewInt(100000)))
	})

	t.Run("decimals", func(t *testing.T) {
		m, err := FromMinor("1.234", vnd)
		require.EqualError(t, err, "1.234 has 3 decimals, only 0 allowed")
		require.Nil(t, m)
	})

	t.Run("leading zeroes", func(t *testing.T) {
		m, err := FromMinor("0012345", php)
		require.NoError(t, err)
		require.NotNil(t, m)
		require.EqualValues(t, m.currency, php)
		require.NotNil(t, m.minorAmount)
		require.Zero(t, m.minorAmount.Cmp(big.NewInt(12345)))
	})

	t.Run("leading zeroes - 0", func(t *testing.T) {
		m, err := FromMinor("000", php)
		require.NoError(t, err)
		require.NotNil(t, m)
		require.EqualValues(t, m.currency, php)
		require.NotNil(t, m.minorAmount)
		require.Zero(t, m.minorAmount.Cmp(big.NewInt(0)))
	})
}

func TestFromMajor(t *testing.T) {
	t.Run("normal - 2 decimals", func(t *testing.T) {
		m, err := FromMajor("123.4", Parser{Currency: php, GroupSep: " ", DecSep: "."})
		require.NoError(t, err)
		require.NotNil(t, m)
		require.EqualValues(t, m.currency, php)
		require.NotNil(t, m.minorAmount)
		require.Zero(t, m.minorAmount.Cmp(big.NewInt(12340)))
	})

	t.Run("normal - 3 decimals", func(t *testing.T) {
		m, err := FromMajor("123.456", Parser{Currency: bhd, GroupSep: " ", DecSep: "."})
		require.NoError(t, err)
		require.NotNil(t, m)
		require.EqualValues(t, m.currency, bhd)
		require.NotNil(t, m.minorAmount)
		require.Zero(t, m.minorAmount.Cmp(big.NewInt(123456)))
	})

	t.Run("normal - 0 decimals", func(t *testing.T) {
		m, err := FromMajor("12 345", Parser{Currency: vnd, GroupSep: " ", DecSep: "."})
		require.NoError(t, err)
		require.NotNil(t, m)
		require.EqualValues(t, m.currency, vnd)
		require.NotNil(t, m.minorAmount)
		require.Zero(t, m.minorAmount.Cmp(big.NewInt(12345)))
	})

	t.Run("normal - negative", func(t *testing.T) {
		m, err := FromMajor("-123.45", Parser{Currency: php, GroupSep: " ", DecSep: "."})
		require.NoError(t, err)
		require.NotNil(t, m)
		require.EqualValues(t, m.currency, php)
		require.NotNil(t, m.minorAmount)
		require.Zero(t, m.minorAmount.Cmp(big.NewInt(-12345)))
	})

	t.Run("1e5 format", func(t *testing.T) {
		m, err := FromMajor("1e5", Parser{Currency: php, GroupSep: " ", DecSep: "."})
		require.NoError(t, err)
		require.NotNil(t, m)
		require.EqualValues(t, m.currency, php)
		require.NotNil(t, m.minorAmount)
		require.Zero(t, m.minorAmount.Cmp(big.NewInt(10000000)))
	})

	t.Run("1e-20", func(t *testing.T) {
		m, err := FromMajor("1e-20", Parser{Currency: vnd, GroupSep: " ", DecSep: "."})
		require.EqualError(t, err, "1e-20 has 20 decimals, only 0 allowed")
		require.Nil(t, m)
	})

	t.Run("nonsense", func(t *testing.T) {
		m, err := FromMajor("foo 123", Parser{Currency: vnd, GroupSep: " ", DecSep: "."})
		require.EqualError(t, err, "foo123 is not a valid decimal amount")
		require.Nil(t, m)
	})

	t.Run("too many decimals - 2 decimal currency", func(t *testing.T) {
		m, err := FromMajor("1.234", Parser{Currency: php, GroupSep: " ", DecSep: "."})
		require.EqualError(t, err, "1.234 has 3 decimals, only 2 allowed")
		require.Nil(t, m)
	})

	t.Run("too many decimals - 0 decimal currency", func(t *testing.T) {
		m, err := FromMajor("1.2", Parser{Currency: vnd, GroupSep: " ", DecSep: "."})
		require.EqualError(t, err, "1.2 has 1 decimals, only 0 allowed")
		require.Nil(t, m)
	})

	t.Run("normal - leading zeroes", func(t *testing.T) {
		m, err := FromMajor("0010.20", Parser{Currency: php, GroupSep: " ", DecSep: "."})
		require.NoError(t, err)
		require.NotNil(t, m)
		require.EqualValues(t, m.currency, php)
		require.NotNil(t, m.minorAmount)
		require.Zero(t, m.minorAmount.Cmp(big.NewInt(1020)))
	})

	t.Run("normal - switched commas", func(t *testing.T) {
		m, err := FromMajor("-123.4567,89", Parser{Currency: php, GroupSep: ".", DecSep: ","})
		require.NoError(t, err)
		require.NotNil(t, m)
		require.EqualValues(t, m.currency, php)
		require.NotNil(t, m.minorAmount)
		require.Zero(t, m.minorAmount.Cmp(big.NewInt(-123456789)))
	})

	t.Run("equal separators", func(t *testing.T) {
		m, err := FromMajor("1.2", Parser{Currency: vnd, GroupSep: ".", DecSep: "."})
		require.EqualError(t, err, "group and decimal separator cannot be the same, are \".\" and \".\"")
		require.Nil(t, m)
	})

	t.Run("normal - extra zero decimals", func(t *testing.T) {
		m, err := FromMajor("123.400000", Parser{Currency: php, GroupSep: " ", DecSep: "."})
		require.NoError(t, err)
		require.NotNil(t, m)
		require.EqualValues(t, m.currency, php)
		require.NotNil(t, m.minorAmount)
		require.Zero(t, m.minorAmount.Cmp(big.NewInt(12340)))
	})
}

func TestFromMajorInt(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		m := FromMajorInt(123, php)
		require.NotNil(t, m)
		require.EqualValues(t, m.currency, php)
		require.NotNil(t, m.minorAmount)
		require.Zero(t, m.minorAmount.Cmp(big.NewInt(12300)))
	})

	t.Run("normal - 3 decimals", func(t *testing.T) {
		m := FromMajorInt(123, bhd)
		require.NotNil(t, m)
		require.EqualValues(t, m.currency, bhd)
		require.NotNil(t, m.minorAmount)
		require.Zero(t, m.minorAmount.Cmp(big.NewInt(123000)))
	})

	t.Run("normal - 0 decimals", func(t *testing.T) {
		m := FromMajorInt(-123, vnd)
		require.NotNil(t, m)
		require.EqualValues(t, m.currency, vnd)
		require.NotNil(t, m.minorAmount)
		require.Zero(t, m.minorAmount.Cmp(big.NewInt(-123)))
	})
}

func TestNegative(t *testing.T) {
	t.Run("negative signed", func(t *testing.T) {
		m, err := FromMinor("-10000", php)
		require.NoError(t, err)
		require.NotNil(t, m)
		require.EqualValues(t, m.currency, php)
		require.NotNil(t, m.minorAmount)
		require.Zero(t, m.minorAmount.Cmp(big.NewInt(-10000)))
	})

	t.Run("negative paren", func(t *testing.T) {
		m, err := FromMinor("(10000)", php)
		require.NoError(t, err)
		require.NotNil(t, m)
		require.EqualValues(t, m.currency, php)
		require.NotNil(t, m.minorAmount)
		require.Zero(t, m.minorAmount.Cmp(big.NewInt(-10000)))
	})
}
