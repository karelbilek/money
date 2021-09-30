package money

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wadey/go-rounding"
)

func TestMoney_SameCurrency(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}
		om := Money{minorAmount: big.NewInt(-70000), currency: php}
		s := m.SameCurrency(om)
		require.True(t, s)
		s = om.SameCurrency(m)
		require.True(t, s)
	})

	t.Run("false", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}
		om := Money{minorAmount: big.NewInt(-70000), currency: vnd}
		s := m.SameCurrency(om)
		require.False(t, s)
		s = om.SameCurrency(m)
		require.False(t, s)
	})

	t.Run("false - default", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}
		om := Money{}
		s := m.SameCurrency(om)
		require.False(t, s)
		s = om.SameCurrency(m)
		require.False(t, s)
	})
}

func TestMoney_Equals(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}
		om := Money{minorAmount: big.NewInt(123), currency: php}
		s, err := m.Equals(om)
		require.NoError(t, err)
		require.True(t, s)
	})

	t.Run("false", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(-100), currency: php}
		om := Money{minorAmount: big.NewInt(123), currency: php}
		s, err := m.Equals(om)
		require.NoError(t, err)
		require.False(t, s)
		s, err = om.Equals(m)
		require.NoError(t, err)
		require.False(t, s)
	})

	t.Run("different currencies", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(100), currency: php}
		om := Money{minorAmount: big.NewInt(100), currency: vnd}
		s, err := m.Equals(om)
		require.EqualError(t, err, "currencies PHP with 2 decimals and VND with 0 decimals don't match")
		require.False(t, s)
		s, err = om.Equals(m)
		require.EqualError(t, err, "currencies VND with 0 decimals and PHP with 2 decimals don't match")
		require.False(t, s)
	})

	t.Run("false - default", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}
		om := Money{}
		s, err := m.Equals(om)
		require.EqualError(t, err, "currencies PHP with 2 decimals and UNKNOWN_CURRENCY with 0 decimals don't match")
		require.False(t, s)
		s, err = om.Equals(m)
		require.EqualError(t, err, "currencies UNKNOWN_CURRENCY with 0 decimals and PHP with 2 decimals don't match")
		require.False(t, s)
	})
}

func TestMoney_More(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(500), currency: php}
		om := Money{minorAmount: big.NewInt(499), currency: php}
		s, err := m.More(om)
		require.NoError(t, err)
		require.True(t, s)
	})

	t.Run("false - equal", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(499), currency: php}
		om := Money{minorAmount: big.NewInt(499), currency: php}
		s, err := m.More(om)
		require.NoError(t, err)
		require.False(t, s)
	})

	t.Run("false - lesser", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(499), currency: php}
		om := Money{minorAmount: big.NewInt(500), currency: php}
		s, err := m.More(om)
		require.NoError(t, err)
		require.False(t, s)
	})

	t.Run("different currencies", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(100), currency: php}
		om := Money{minorAmount: big.NewInt(123), currency: vnd}
		s, err := m.More(om)
		require.EqualError(t, err, "currencies PHP with 2 decimals and VND with 0 decimals don't match")
		require.False(t, s)
		s, err = om.More(m)
		require.EqualError(t, err, "currencies VND with 0 decimals and PHP with 2 decimals don't match")
		require.False(t, s)
	})

	t.Run("false - default", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}
		om := Money{}
		s, err := m.More(om)
		require.EqualError(t, err, "currencies PHP with 2 decimals and UNKNOWN_CURRENCY with 0 decimals don't match")
		require.False(t, s)
		s, err = om.More(m)
		require.EqualError(t, err, "currencies UNKNOWN_CURRENCY with 0 decimals and PHP with 2 decimals don't match")
		require.False(t, s)
	})
}

func TestMoney_MoreEqual(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(500), currency: php}
		om := Money{minorAmount: big.NewInt(499), currency: php}
		s, err := m.MoreEqual(om)
		require.NoError(t, err)
		require.True(t, s)
	})

	t.Run("true - equal", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(499), currency: php}
		om := Money{minorAmount: big.NewInt(499), currency: php}
		s, err := m.MoreEqual(om)
		require.NoError(t, err)
		require.True(t, s)
	})

	t.Run("false - lesser", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(499), currency: php}
		om := Money{minorAmount: big.NewInt(500), currency: php}
		s, err := m.MoreEqual(om)
		require.NoError(t, err)
		require.False(t, s)
	})

	t.Run("different currencies", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(100), currency: php}
		om := Money{minorAmount: big.NewInt(123), currency: vnd}
		s, err := m.MoreEqual(om)
		require.EqualError(t, err, "currencies PHP with 2 decimals and VND with 0 decimals don't match")
		require.False(t, s)
		s, err = om.MoreEqual(m)
		require.EqualError(t, err, "currencies VND with 0 decimals and PHP with 2 decimals don't match")
		require.False(t, s)
	})

	t.Run("false - default", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}
		om := Money{}
		s, err := m.MoreEqual(om)
		require.EqualError(t, err, "currencies PHP with 2 decimals and UNKNOWN_CURRENCY with 0 decimals don't match")
		require.False(t, s)
		s, err = om.MoreEqual(m)
		require.EqualError(t, err, "currencies UNKNOWN_CURRENCY with 0 decimals and PHP with 2 decimals don't match")
		require.False(t, s)
	})
}

func TestMoney_Less(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(499), currency: php}
		om := Money{minorAmount: big.NewInt(500), currency: php}
		s, err := m.Less(om)
		require.NoError(t, err)
		require.True(t, s)
	})

	t.Run("false - equal", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(499), currency: php}
		om := Money{minorAmount: big.NewInt(499), currency: php}
		s, err := m.Less(om)
		require.NoError(t, err)
		require.False(t, s)
	})

	t.Run("false - bigger", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(500), currency: php}
		om := Money{minorAmount: big.NewInt(499), currency: php}
		s, err := m.Less(om)
		require.NoError(t, err)
		require.False(t, s)
	})

	t.Run("different currencies", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(100), currency: php}
		om := Money{minorAmount: big.NewInt(123), currency: vnd}
		s, err := m.Less(om)
		require.EqualError(t, err, "currencies PHP with 2 decimals and VND with 0 decimals don't match")
		require.False(t, s)
		s, err = om.Less(m)
		require.EqualError(t, err, "currencies VND with 0 decimals and PHP with 2 decimals don't match")
		require.False(t, s)
	})

	t.Run("false - default", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}
		om := Money{}
		s, err := m.Less(om)
		require.EqualError(t, err, "currencies PHP with 2 decimals and UNKNOWN_CURRENCY with 0 decimals don't match")
		require.False(t, s)
		s, err = om.Less(m)
		require.EqualError(t, err, "currencies UNKNOWN_CURRENCY with 0 decimals and PHP with 2 decimals don't match")
		require.False(t, s)
	})
}

func TestMoney_LessMajor(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(499), currency: php}
		s := m.LessMajor(5)
		require.True(t, s)
	})

	t.Run("false - equal", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(500), currency: php}
		s := m.LessMajor(5)
		require.False(t, s)
	})

	t.Run("false - bigger", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(500), currency: php}
		s := m.LessMajor(2)
		require.False(t, s)
	})

	t.Run("false - default", func(t *testing.T) {
		m := &Money{}
		s := m.LessMajor(2)
		require.True(t, s)
	})
}

func TestMoney_MoreMajor(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(499), currency: php}
		s := m.MoreMajor(4)
		require.True(t, s)
	})

	t.Run("false - equal", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(500), currency: php}
		s := m.MoreMajor(5)
		require.False(t, s)
	})

	t.Run("false - smaller", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(500), currency: php}
		s := m.MoreMajor(7)
		require.False(t, s)
	})

	t.Run("false - default", func(t *testing.T) {
		m := &Money{}
		s := m.MoreMajor(2)
		require.False(t, s)
	})
}

func TestMoney_LessEqualMajor(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(499), currency: php}
		s := m.LessEqualMajor(5)
		require.True(t, s)
	})

	t.Run("true - equal", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(500), currency: php}
		s := m.LessEqualMajor(5)
		require.True(t, s)
	})

	t.Run("false - bigger", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(500), currency: php}
		s := m.LessEqualMajor(2)
		require.False(t, s)
	})

	t.Run("false - default", func(t *testing.T) {
		m := &Money{}
		s := m.LessEqualMajor(2)
		require.True(t, s)
	})
}

func TestMoney_MoreEqualMajor(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(499), currency: php}
		s := m.MoreEqualMajor(4)
		require.True(t, s)
	})

	t.Run("true - equal", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(500), currency: php}
		s := m.MoreEqualMajor(5)
		require.True(t, s)
	})

	t.Run("false - smaller", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(500), currency: php}
		s := m.MoreEqualMajor(7)
		require.False(t, s)
	})

	t.Run("false - default", func(t *testing.T) {
		m := &Money{}
		s := m.MoreEqualMajor(2)
		require.False(t, s)
	})
}

func TestMoney_BetweenMajor(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(499), currency: php}
		s, err := m.BetweenMajor(4, 6)
		require.NoError(t, err)
		require.True(t, s)
	})

	t.Run("true - equal min", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(500), currency: php}
		s, err := m.BetweenMajor(5, 6)
		require.NoError(t, err)
		require.True(t, s)
	})

	t.Run("true - equal max", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(600), currency: php}
		s, err := m.BetweenMajor(5, 6)
		require.NoError(t, err)
		require.True(t, s)
	})

	t.Run("false - below", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(500), currency: php}
		s, err := m.BetweenMajor(7, 8)
		require.NoError(t, err)
		require.False(t, s)
	})

	t.Run("false - above", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(900), currency: php}
		s, err := m.BetweenMajor(7, 8)
		require.NoError(t, err)
		require.False(t, s)
	})

	t.Run("false - default", func(t *testing.T) {
		m := &Money{}
		s, err := m.BetweenMajor(2, 3)
		require.NoError(t, err)
		require.False(t, s)
	})

	t.Run("error - min > max", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(500), currency: php}
		s, err := m.BetweenMajor(7, 6)
		require.EqualError(t, err, "minimal is bigger than maximal")
		require.False(t, s)
	})
}

func TestMoney_LessEqual(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(499), currency: php}
		om := Money{minorAmount: big.NewInt(500), currency: php}
		s, err := m.LessEqual(om)
		require.NoError(t, err)
		require.True(t, s)
	})

	t.Run("true - equal", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(499), currency: php}
		om := Money{minorAmount: big.NewInt(499), currency: php}
		s, err := m.LessEqual(om)
		require.NoError(t, err)
		require.True(t, s)
	})

	t.Run("false - bigger", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(500), currency: php}
		om := Money{minorAmount: big.NewInt(499), currency: php}
		s, err := m.LessEqual(om)
		require.NoError(t, err)
		require.False(t, s)
	})

	t.Run("different currencies", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(100), currency: php}
		om := Money{minorAmount: big.NewInt(123), currency: vnd}
		s, err := m.LessEqual(om)
		require.EqualError(t, err, "currencies PHP with 2 decimals and VND with 0 decimals don't match")
		require.False(t, s)
		s, err = om.LessEqual(m)
		require.EqualError(t, err, "currencies VND with 0 decimals and PHP with 2 decimals don't match")
		require.False(t, s)
	})

	t.Run("false - default", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}
		om := Money{}
		s, err := m.LessEqual(om)
		require.EqualError(t, err, "currencies PHP with 2 decimals and UNKNOWN_CURRENCY with 0 decimals don't match")
		require.False(t, s)
		s, err = om.LessEqual(m)
		require.EqualError(t, err, "currencies UNKNOWN_CURRENCY with 0 decimals and PHP with 2 decimals don't match")
		require.False(t, s)
	})
}

func TestMoney_Between(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(500), currency: php}
		min := Money{minorAmount: big.NewInt(499), currency: php}
		max := Money{minorAmount: big.NewInt(501), currency: php}

		s, err := m.Between(min, max)
		require.NoError(t, err)
		require.True(t, s)
	})

	t.Run("false - below", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(400), currency: php}
		min := Money{minorAmount: big.NewInt(499), currency: php}
		max := Money{minorAmount: big.NewInt(501), currency: php}

		s, err := m.Between(min, max)
		require.NoError(t, err)
		require.False(t, s)
	})

	t.Run("false - above", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(600), currency: php}
		min := Money{minorAmount: big.NewInt(499), currency: php}
		max := Money{minorAmount: big.NewInt(501), currency: php}

		s, err := m.Between(min, max)
		require.NoError(t, err)
		require.False(t, s)
	})

	t.Run("different currencies 1", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(500), currency: php}
		min := Money{minorAmount: big.NewInt(499), currency: vnd}
		max := Money{minorAmount: big.NewInt(501), currency: php}

		s, err := m.Between(min, max)
		require.EqualError(t, err, "currencies VND with 0 decimals and PHP with 2 decimals don't match")
		require.Equal(t, s, false)
	})

	t.Run("different currencies 2", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(500), currency: php}
		min := Money{minorAmount: big.NewInt(499), currency: php}
		max := Money{minorAmount: big.NewInt(501), currency: vnd}

		s, err := m.Between(min, max)
		require.EqualError(t, err, "currencies PHP with 2 decimals and VND with 0 decimals don't match")
		require.Equal(t, s, false)
	})

	t.Run("error - max > min", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(500), currency: php}
		min := Money{minorAmount: big.NewInt(502), currency: php}
		max := Money{minorAmount: big.NewInt(499), currency: php}

		s, err := m.Between(min, max)
		require.EqualError(t, err, "minimal is bigger than maximal")
		require.Equal(t, s, false)
	})
}

func TestMoney_IsZero(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(0), currency: php}
		s := m.IsZero()
		require.True(t, s)
	})

	t.Run("false - > 0", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(1), currency: php}
		s := m.IsZero()
		require.False(t, s)
	})

	t.Run("false - < 0", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(-1), currency: php}
		s := m.IsZero()
		require.False(t, s)
	})

	t.Run("default", func(t *testing.T) {
		m := &Money{}
		s := m.IsZero()
		require.True(t, s)
	})
}

func TestMoney_HasCents(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(12312), currency: php}
		s, err := m.HasCents()
		require.NoError(t, err)
		require.True(t, s)
	})

	t.Run("false", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123400), currency: php}
		s, err := m.HasCents()
		require.NoError(t, err)
		require.False(t, s)
	})

	t.Run("false, 0 decimal currency", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(1234), currency: vnd}
		s, err := m.HasCents()
		require.NoError(t, err)
		require.False(t, s)
	})

	t.Run("default", func(t *testing.T) {
		m := &Money{}
		s, err := m.HasCents()
		require.NoError(t, err)
		require.False(t, s)
	})
}

func TestMoney_IsPositive(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(1), currency: php}
		s := m.IsPositive()
		require.True(t, s)
	})

	t.Run("false - 0", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(0), currency: php}
		s := m.IsPositive()
		require.False(t, s)
	})

	t.Run("false - < 0", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(-1), currency: php}
		s := m.IsPositive()
		require.False(t, s)
	})

	t.Run("default", func(t *testing.T) {
		m := &Money{}
		s := m.IsPositive()
		require.False(t, s)
	})
}

func TestMoney_IsNegative(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(-1), currency: php}
		s := m.IsNegative()
		require.True(t, s)
	})

	t.Run("false - 0", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(0), currency: php}
		s := m.IsNegative()
		require.False(t, s)
	})

	t.Run("false - > 0", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(1), currency: php}
		s := m.IsNegative()
		require.False(t, s)
	})

	t.Run("default", func(t *testing.T) {
		m := &Money{}
		s := m.IsNegative()
		require.False(t, s)
	})
}

func TestMoney_Absolute(t *testing.T) {
	t.Run("normal - negative", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(-123), currency: php}
		abs := m.Absolute()
		require.Zero(t, abs.minorAmount.Cmp(big.NewInt(123)))
		require.EqualValues(t, abs.currency, php)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(-123)))
	})

	t.Run("normal - positive", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}
		abs := m.Absolute()
		require.Zero(t, abs.minorAmount.Cmp(big.NewInt(123)))
		require.EqualValues(t, abs.currency, php)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(123)))
	})

	t.Run("default", func(t *testing.T) {
		m := &Money{}
		abs := m.Absolute()
		require.Zero(t, abs.minorAmount.Cmp(big.NewInt(0)))
		require.EqualValues(t, abs.currency, Currency{})
		require.Nil(t, m.minorAmount)
	})
}

func TestMoney_Negative(t *testing.T) {
	t.Run("normal - negative", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(-123), currency: php}
		abs := m.Negative()
		require.Zero(t, abs.minorAmount.Cmp(big.NewInt(123)))
		require.EqualValues(t, abs.currency, php)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(-123)))
	})

	t.Run("normal - positive", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}
		abs := m.Negative()
		require.Zero(t, abs.minorAmount.Cmp(big.NewInt(-123)))
		require.EqualValues(t, abs.currency, php)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(123)))
	})

	t.Run("default", func(t *testing.T) {
		m := &Money{}
		abs := m.Negative()
		require.Zero(t, abs.minorAmount.Cmp(big.NewInt(0)))
		require.EqualValues(t, abs.currency, Currency{})

		require.Nil(t, m.minorAmount)
	})
}

func TestMoney_Add(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}

		om := Money{minorAmount: big.NewInt(101), currency: php}
		sum, err := m.Add(om)
		require.NoError(t, err)
		require.Zero(t, sum.minorAmount.Cmp(big.NewInt(224)))
		require.EqualValues(t, sum.currency, php)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(123)))
		require.Zero(t, om.minorAmount.Cmp(big.NewInt(101)))
	})

	t.Run("different currencies", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}

		om := Money{minorAmount: big.NewInt(101), currency: vnd}
		sum, err := m.Add(om)
		require.EqualError(t, err, "currencies PHP with 2 decimals and VND with 0 decimals don't match")
		require.Nil(t, sum)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(123)))
		require.Zero(t, om.minorAmount.Cmp(big.NewInt(101)))
	})

	t.Run("default", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}
		om := Money{}

		sum, err := m.Add(om)
		require.EqualError(t, err, "currencies PHP with 2 decimals and UNKNOWN_CURRENCY with 0 decimals don't match")
		require.Nil(t, sum)

		sum, err = om.Add(m)
		require.EqualError(t, err, "currencies UNKNOWN_CURRENCY with 0 decimals and PHP with 2 decimals don't match")
		require.Nil(t, sum)
	})
}

func TestMoney_Subtract(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}

		om := Money{minorAmount: big.NewInt(101), currency: php}
		sum, err := m.Subtract(om)
		require.NoError(t, err)
		require.Zero(t, sum.minorAmount.Cmp(big.NewInt(22)))
		require.EqualValues(t, sum.currency, php)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(123)))
		require.Zero(t, om.minorAmount.Cmp(big.NewInt(101)))
	})

	t.Run("different currencies", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}

		om := Money{minorAmount: big.NewInt(101), currency: vnd}
		sum, err := m.Subtract(om)
		require.EqualError(t, err, "currencies PHP with 2 decimals and VND with 0 decimals don't match")
		require.Nil(t, sum)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(123)))
		require.Zero(t, om.minorAmount.Cmp(big.NewInt(101)))
	})

	t.Run("default", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}
		om := Money{}

		sum, err := m.Subtract(om)
		require.EqualError(t, err, "currencies PHP with 2 decimals and UNKNOWN_CURRENCY with 0 decimals don't match")
		require.Nil(t, sum)

		sum, err = om.Subtract(m)
		require.EqualError(t, err, "currencies UNKNOWN_CURRENCY with 0 decimals and PHP with 2 decimals don't match")
		require.Nil(t, sum)
	})
}

func TestMoney_Multiply(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}

		p, err := m.Multiply(100)
		require.NoError(t, err)
		require.Zero(t, p.minorAmount.Cmp(big.NewInt(12300)))
		require.EqualValues(t, p.currency, php)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(123)))
	})

	t.Run("default", func(t *testing.T) {
		m := &Money{}
		p, err := m.Multiply(10)
		require.NoError(t, err)
		require.Zero(t, p.minorAmount.Cmp(big.NewInt(0)))
		require.EqualValues(t, p.currency, Currency{})
	})
}

func TestMoney_MultiplyRat(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}

		p, err := m.MultiplyRat("-1.9", rounding.HalfUp)
		require.NoError(t, err)
		require.Zero(t, p.minorAmount.Cmp(big.NewInt(-234)))
		require.EqualValues(t, p.currency, php)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(123)))
	})

	t.Run("error on weird", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}

		p, err := m.MultiplyRat("foo bar", rounding.HalfUp)
		require.EqualError(t, err, "foo bar is not a valid rational amount")
		require.Nil(t, p)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(123)))
	})

	t.Run("default", func(t *testing.T) {
		m := &Money{}
		p, err := m.MultiplyRat("10", rounding.HalfUp)
		require.NoError(t, err)
		require.Zero(t, p.minorAmount.Cmp(big.NewInt(0)))
		require.EqualValues(t, p.currency, Currency{})
	})
}

func TestMoney_MultiplyBigRat(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}

		p := m.MultiplyBigRat(big.NewRat(1, 3), rounding.HalfUp)
		require.Zero(t, p.minorAmount.Cmp(big.NewInt(41)))
		require.EqualValues(t, p.currency, php)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(123)))
	})

	t.Run("default", func(t *testing.T) {
		m := &Money{}
		p := m.MultiplyBigRat(big.NewRat(1, 3), rounding.HalfUp)
		require.Zero(t, p.minorAmount.Cmp(big.NewInt(0)))
		require.EqualValues(t, p.currency, Currency{})
	})
}

func TestMoney_Divide(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}

		p, err := m.Divide(5, rounding.HalfUp)
		require.NoError(t, err)
		require.Zero(t, p.minorAmount.Cmp(big.NewInt(25)))
		require.EqualValues(t, p.currency, php)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(123)))
	})

	t.Run("error - division by zero", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}
		p, err := m.Divide(0, rounding.HalfUp)
		require.EqualError(t, err, "division by zero")
		require.Nil(t, p)
	})

	t.Run("default", func(t *testing.T) {
		m := &Money{}
		p, err := m.Divide(10, rounding.HalfUp)
		require.NoError(t, err)
		require.Zero(t, p.minorAmount.Cmp(big.NewInt(0)))
		require.EqualValues(t, p.currency, Currency{})
	})
}

func TestMoney_DivideRat(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}

		p, err := m.DivideRat("-1.9", rounding.HalfUp)
		require.NoError(t, err)
		require.Zero(t, p.minorAmount.Cmp(big.NewInt(-65)))
		require.EqualValues(t, p.currency, php)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(123)))
	})

	t.Run("error on weird", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}

		p, err := m.DivideRat("foo bar", rounding.HalfUp)
		require.EqualError(t, err, "foo bar is not a valid rational amount")
		require.Nil(t, p)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(123)))
	})

	t.Run("error - division by zero", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}
		p, err := m.DivideRat("0", rounding.HalfUp)
		require.EqualError(t, err, "division by zero")
		require.Nil(t, p)
	})

	t.Run("default", func(t *testing.T) {
		m := &Money{}
		p, err := m.DivideRat("10", rounding.HalfUp)
		require.NoError(t, err)
		require.Zero(t, p.minorAmount.Cmp(big.NewInt(0)))
		require.EqualValues(t, p.currency, Currency{})
	})
}

func TestMoney_DivideBigRat(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}

		p, err := m.DivideBigRat(big.NewRat(3, -1), rounding.HalfUp)
		require.NoError(t, err)
		require.Zero(t, p.minorAmount.Cmp(big.NewInt(-41)))
		require.EqualValues(t, p.currency, php)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(123)))
	})

	t.Run("error - nil rat", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}
		p, err := m.DivideBigRat(nil, rounding.HalfUp)
		require.EqualError(t, err, "nil rational")
		require.Nil(t, p)
	})

	t.Run("error - 0 quotient", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}
		p, err := m.DivideBigRat(big.NewRat(0, 3), rounding.HalfUp)
		require.EqualError(t, err, "division by zero")
		require.Nil(t, p)
	})

	t.Run("default", func(t *testing.T) {
		m := &Money{}
		p, err := m.DivideBigRat(big.NewRat(1, 3), rounding.HalfUp)
		require.NoError(t, err)
		require.Zero(t, p.minorAmount.Cmp(big.NewInt(0)))
		require.EqualValues(t, p.currency, Currency{})
	})
}

func TestMoney_RoundToMajor(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: php}
		r := m.RoundToMajor(rounding.Down)

		require.Zero(t, r.minorAmount.Cmp(big.NewInt(100)))
		require.EqualValues(t, r.currency, php)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(123)))
	})

	t.Run("normal - 0 decimals", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(123), currency: vnd}
		r := m.RoundToMajor(rounding.Down)

		require.Zero(t, r.minorAmount.Cmp(big.NewInt(123)))
		require.EqualValues(t, r.currency, vnd)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(123)))
	})

	t.Run("default", func(t *testing.T) {
		m := &Money{}
		r := m.RoundToMajor(rounding.Down)

		// default is "zero"
		require.True(t, r.IsZero())
		require.EqualValues(t, r.currency, Currency{})
	})
}

func TestMoney_Split(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(100), currency: php}
		ss, err := m.Split(3)
		require.NoError(t, err)
		require.Len(t, ss, 3)

		require.Zero(t, ss[0].minorAmount.Cmp(big.NewInt(34)))
		require.Zero(t, ss[1].minorAmount.Cmp(big.NewInt(33)))
		require.Zero(t, ss[2].minorAmount.Cmp(big.NewInt(33)))
		for _, s := range ss {
			require.EqualValues(t, s.currency, php)
		}

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(100)))
	})

	t.Run("negative", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(-100), currency: php}
		ss, err := m.Split(3)
		require.NoError(t, err)
		require.Len(t, ss, 3)

		require.Zero(t, ss[0].minorAmount.Cmp(big.NewInt(-34)))
		require.Zero(t, ss[1].minorAmount.Cmp(big.NewInt(-33)))
		require.Zero(t, ss[2].minorAmount.Cmp(big.NewInt(-33)))
		for _, s := range ss {
			require.EqualValues(t, s.currency, php)
		}

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(-100)))
	})

	t.Run("<= 0", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(-100), currency: php}
		ss, err := m.Split(0)
		require.EqualError(t, err, "split must be higher than zero, is 0")
		require.Len(t, ss, 0)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(-100)))
	})

	t.Run("default", func(t *testing.T) {
		m := &Money{}
		ss, err := m.Split(1)

		require.NoError(t, err)
		require.Len(t, ss, 1)
		require.Zero(t, ss[0].minorAmount.Cmp(big.NewInt(0)))
		require.EqualValues(t, ss[0].currency, Currency{})
	})
}

func TestMoney_Allocate(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(100), currency: php}
		ss, err := m.Allocate(2, 1)
		require.NoError(t, err)
		require.Len(t, ss, 2)

		require.Zero(t, ss[0].minorAmount.Cmp(big.NewInt(67)))
		require.Zero(t, ss[1].minorAmount.Cmp(big.NewInt(33)))
		for _, s := range ss {
			require.EqualValues(t, s.currency, php)
		}

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(100)))
	})

	t.Run("negative", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(-100), currency: php}
		ss, err := m.Allocate(2, 1)
		require.NoError(t, err)
		require.Len(t, ss, 2)

		require.Zero(t, ss[0].minorAmount.Cmp(big.NewInt(-67)))
		require.Zero(t, ss[1].minorAmount.Cmp(big.NewInt(-33)))
		for _, s := range ss {
			require.EqualValues(t, s.currency, php)
		}

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(-100)))
	})

	t.Run("no ratios", func(t *testing.T) {
		m := Money{minorAmount: big.NewInt(-100), currency: php}
		ss, err := m.Allocate()
		require.EqualError(t, err, "no ratios specified")
		require.Len(t, ss, 0)

		require.Zero(t, m.minorAmount.Cmp(big.NewInt(-100)))
	})

	t.Run("default", func(t *testing.T) {
		m := &Money{}
		ss, err := m.Allocate(2)

		require.NoError(t, err)
		require.Len(t, ss, 1)
		require.Zero(t, ss[0].minorAmount.Cmp(big.NewInt(0)))
		require.EqualValues(t, ss[0].currency, Currency{})
	})
}
