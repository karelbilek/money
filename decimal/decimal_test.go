package decimal

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

// better test both pointer and value in the JSON
type WithAmounts struct {
	Value   Decimal  `json:"value"`
	Pointer *Decimal `json:"pointer,omitempty"`
}

func TestFromInt(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		d := FromInt(50)
		require.Zero(t, d.rat.Cmp(big.NewRat(50, 1)))
	})

	t.Run("zero", func(t *testing.T) {
		d := FromInt(0)
		require.Zero(t, d.rat.Cmp(big.NewRat(0, 1)))
	})
}

func TestFromString(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		d, err := FromString("10.1")
		require.NoError(t, err)
		require.Zero(t, d.rat.Cmp(big.NewRat(101, 10)))
	})

	t.Run("empty", func(t *testing.T) {
		d, err := FromString("")
		require.NoError(t, err)
		require.Zero(t, d.rat.Cmp(big.NewRat(0, 1)))
	})

	t.Run("no fractional", func(t *testing.T) {
		d, err := FromString("10/1")
		require.EqualError(t, err, "10/1 is not a valid decimal amount")
		require.Nil(t, d.rat)
	})

	t.Run("too big exponent", func(t *testing.T) {
		d, err := FromString("1e1000000")
		require.EqualError(t, err, "decimal length 1000001 bigger than allowed 200")
		require.Nil(t, d.rat)
	})

	t.Run("too small exponent", func(t *testing.T) {
		d, err := FromString("1e-1000000")
		require.EqualError(t, err, "fractional length 1000000 bigger than allowed 200")
		require.Nil(t, d.rat)
	})
}

func TestDecimal_FracDecimals(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		c := Decimal{rat: big.NewRat(-1125, 1000)}
		d, err := c.FracDecimals()
		require.NoError(t, err)
		require.Equal(t, 3, d)
	})

	t.Run("default", func(t *testing.T) {
		c := Decimal{}
		d, err := c.FracDecimals()
		require.NoError(t, err)
		require.Equal(t, 0, d)
	})

	t.Run("repeating", func(t *testing.T) {
		// note - this is testing edgecase that should not happen with From* constructors
		c := Decimal{rat: big.NewRat(1, 3)}
		d, err := c.FracDecimals()
		require.EqualError(t, err, "rational 1/3 is repeating")
		require.Equal(t, 0, d)
	})
}

func TestDecimal_BigRat(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		d := Decimal{rat: big.NewRat(-123, 500)}
		require.Zero(t, d.BigRat().Cmp(big.NewRat(-123, 500)))
	})

	t.Run("empty", func(t *testing.T) {
		d := Decimal{}
		require.Zero(t, d.BigRat().Cmp(big.NewRat(0, 1)))
	})
}

func TestDecimal_BigInt(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		d := Decimal{rat: big.NewRat(-123, 500)}
		b, err := d.BigInt(3)
		require.NoError(t, err)
		require.Zero(t, b.Cmp(big.NewInt(-246)))
	})

	t.Run("too few decimals to get int", func(t *testing.T) {
		d := Decimal{rat: big.NewRat(-123, 500), original: "-0.246"}
		b, err := d.BigInt(2)
		require.EqualError(t, err, "-0.246 has 3 decimals, only 2 allowed")
		require.Nil(t, b)
	})

	t.Run("negative fracDecimals", func(t *testing.T) {
		d := Decimal{rat: big.NewRat(-123, 500)}
		b, err := d.BigInt(-2)
		require.EqualError(t, err, "cannot have negative fracDecimals -2")
		require.Nil(t, b)
	})

	t.Run("repeating", func(t *testing.T) {
		// note - this is testing edgecase that should not happen with From* constructors
		d := Decimal{rat: big.NewRat(1, 3)}
		b, err := d.BigInt(2)
		require.EqualError(t, err, "rational 1/3 is repeating")
		require.Nil(t, b)
	})
}

func TestDecimal_UnmarshalJSON(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		w := WithAmounts{}
		err := json.Unmarshal([]byte(`{"value": -0.02,"pointer":12345}`), &w)
		require.NoError(t, err)
		require.NotNil(t, w.Value.rat)
		require.Zero(t, w.Value.rat.Cmp(big.NewRat(-2, 100)))
		require.NotNil(t, w.Pointer)
		require.NotNil(t, w.Pointer.rat)
		require.Zero(t, w.Pointer.rat.Cmp(big.NewRat(12345, 1)))
	})

	t.Run("exp format", func(t *testing.T) {
		w := WithAmounts{}
		err := json.Unmarshal([]byte(`{"value": -12.3e-2}`), &w)
		require.NoError(t, err)
		require.Nil(t, w.Pointer)
		require.NotNil(t, w.Value.rat)
		require.Zero(t, w.Value.rat.Cmp(big.NewRat(-123, 1000)))
	})

	t.Run("default", func(t *testing.T) {
		w := WithAmounts{}
		err := json.Unmarshal([]byte(`{}`), &w)
		require.NoError(t, err)
		require.Nil(t, w.Value.rat)
		require.Nil(t, w.Pointer)
	})

	t.Run("error - nonsense", func(t *testing.T) {
		w := WithAmounts{}
		err := json.Unmarshal([]byte(`{"value": ["some", "weird", "thing"]}`), &w)
		require.EqualError(t, err, `["some", "weird", "thing"] is not a valid decimal amount`)
		require.Nil(t, w.Value.rat)
		require.Nil(t, w.Pointer)
	})

	t.Run("error - string", func(t *testing.T) {
		w := WithAmounts{}
		err := json.Unmarshal([]byte(`{"value": "1.2"}`), &w)
		require.EqualError(t, err, "\"1.2\" is not a valid decimal amount")
		require.Nil(t, w.Value.rat)
		require.Nil(t, w.Pointer)
	})
}
