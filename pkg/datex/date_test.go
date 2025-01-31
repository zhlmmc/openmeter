package datex_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/openmeterio/openmeter/openmeter/testutils"
	"github.com/openmeterio/openmeter/pkg/datex"
)

func TestISOOperations(t *testing.T) {
	t.Run("Parse", func(t *testing.T) {
		isoDuration := "P1Y2M3DT4H5M6S"

		period, err := datex.ISOString(isoDuration).Parse()
		require.NoError(t, err)

		now := testutils.GetRFC3339Time(t, "2020-01-01T00:00:00Z")
		expected := testutils.GetRFC3339Time(t, "2021-03-04T04:05:06Z")
		actual, precise := period.AddTo(now)
		assert.True(t, precise)
		assert.Equal(t, expected, actual)
	})

	t.Run("ParseError", func(t *testing.T) {
		isoDuration := "P1Y2M3DT4H5M6SX"
		_, err := datex.ISOString(isoDuration).Parse()
		assert.Error(t, err)
	})

	t.Run("Works with 0 duration", func(t *testing.T) {
		isoDuration := "P0D"

		period, err := datex.ISOString(isoDuration).Parse()
		require.NoError(t, err)

		now := testutils.GetRFC3339Time(t, "2020-01-01T00:00:00Z")
		expected := testutils.GetRFC3339Time(t, "2020-01-01T00:00:00Z")
		actual, precise := period.AddTo(now)
		assert.True(t, precise)
		assert.Equal(t, expected, actual)
	})

	t.Run("Adding periods", func(t *testing.T) {
		isoDuration1 := "PT5M"
		isoDuration2 := "PT1M1S"

		period1, err := datex.ISOString(isoDuration1).Parse()
		require.NoError(t, err)

		period2, err := datex.ISOString(isoDuration2).Parse()
		require.NoError(t, err)

		expectedS := "PT6M1S"
		expected, err := datex.ISOString(expectedS).Parse()
		require.NoError(t, err)

		actual, err := period1.Add(period2)
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("String representation", func(t *testing.T) {
		iso := datex.ISOString("PT5M")
		assert.Equal(t, "PT5M", iso.String())
	})
}

func TestNewPeriod(t *testing.T) {
	t.Run("all fields", func(t *testing.T) {
		period := datex.NewPeriod(1, 2, 0, 3, 4, 5, 6)
		assert.Equal(t, "P1Y2M3DT4H5M6S", period.String())
	})

	t.Run("zero values", func(t *testing.T) {
		period := datex.NewPeriod(0, 0, 0, 0, 0, 0, 0)
		assert.Equal(t, "P0D", period.String())
	})

	t.Run("negative values", func(t *testing.T) {
		period := datex.NewPeriod(-1, -2, 0, -3, -4, -5, -6)
		assert.Equal(t, "-P1Y2M3DT4H5M6S", period.String())
	})
}

func TestParsePtrOrNil(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		var iso *datex.ISOString
		period, err := iso.ParsePtrOrNil()
		require.NoError(t, err)
		assert.Nil(t, period)
	})

	t.Run("valid input", func(t *testing.T) {
		iso := datex.ISOString("P1Y")
		period, err := iso.ParsePtrOrNil()
		require.NoError(t, err)
		require.NotNil(t, period)
		assert.Equal(t, "P1Y", period.String())
	})

	t.Run("invalid input", func(t *testing.T) {
		iso := datex.ISOString("invalid")
		period, err := iso.ParsePtrOrNil()
		assert.Error(t, err)
		assert.Nil(t, period)
	})
}

func TestPeriodSubtract(t *testing.T) {
	t.Run("valid subtraction", func(t *testing.T) {
		p1, err := datex.ISOString("PT10M").Parse()
		require.NoError(t, err)
		p2, err := datex.ISOString("PT5M").Parse()
		require.NoError(t, err)

		result, err := p1.Subtract(p2)
		require.NoError(t, err)
		assert.Equal(t, "PT5M", result.String())
	})

	t.Run("negative result", func(t *testing.T) {
		p1, err := datex.ISOString("PT5M").Parse()
		require.NoError(t, err)
		p2, err := datex.ISOString("PT10M").Parse()
		require.NoError(t, err)

		result, err := p1.Subtract(p2)
		require.NoError(t, err)
		assert.Equal(t, "-PT5M", result.String())
	})
}

func TestPeriodsAlign(t *testing.T) {
	t.Run("aligned periods", func(t *testing.T) {
		larger, err := datex.ISOString("PT10M").Parse()
		require.NoError(t, err)
		smaller, err := datex.ISOString("PT5M").Parse()
		require.NoError(t, err)

		aligned, err := datex.PeriodsAlign(larger, smaller)
		require.NoError(t, err)
		assert.True(t, aligned)
	})

	t.Run("non-aligned periods", func(t *testing.T) {
		larger, err := datex.ISOString("PT11M").Parse()
		require.NoError(t, err)
		smaller, err := datex.ISOString("PT5M").Parse()
		require.NoError(t, err)

		aligned, err := datex.PeriodsAlign(larger, smaller)
		require.NoError(t, err)
		assert.False(t, aligned)
	})

	t.Run("smaller period is larger", func(t *testing.T) {
		larger, err := datex.ISOString("PT5M").Parse()
		require.NoError(t, err)
		smaller, err := datex.ISOString("PT10M").Parse()
		require.NoError(t, err)

		_, err = datex.PeriodsAlign(larger, smaller)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "smaller period is larger than larger period")
	})
}

func TestBetween(t *testing.T) {
	t.Run("5 minutes", func(t *testing.T) {
		start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2020, 1, 1, 0, 5, 0, 0, time.UTC)

		period := datex.Between(start, end)
		assert.Equal(t, "PT300S", period.String())
	})

	t.Run("zero duration", func(t *testing.T) {
		now := time.Now()
		period := datex.Between(now, now)
		assert.Equal(t, "P0D", period.String())
	})

	t.Run("negative duration", func(t *testing.T) {
		start := time.Date(2020, 1, 1, 0, 5, 0, 0, time.UTC)
		end := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

		period := datex.Between(start, end)
		assert.Equal(t, "-PT300S", period.String())
	})
}

func TestFromDuration(t *testing.T) {
	t.Run("5 minutes", func(t *testing.T) {
		duration := 5 * time.Minute
		period := datex.FromDuration(duration)
		assert.Equal(t, "PT5M", period.String())
	})

	t.Run("zero duration", func(t *testing.T) {
		period := datex.FromDuration(0)
		assert.Equal(t, "P0D", period.String())
	})

	t.Run("complex duration", func(t *testing.T) {
		duration := 1*time.Hour + 30*time.Minute + 45*time.Second
		period := datex.FromDuration(duration)
		assert.Equal(t, "PT5445S", period.String())
	})

	t.Run("negative duration", func(t *testing.T) {
		duration := -5 * time.Minute
		period := datex.FromDuration(duration)
		assert.Equal(t, "-PT5M", period.String())
	})
}

func TestPeriodEqual(t *testing.T) {
	t.Run("both nil", func(t *testing.T) {
		var p1, p2 *datex.Period
		assert.True(t, p1.Equal(p2))
	})

	t.Run("one nil", func(t *testing.T) {
		var p1 *datex.Period
		p2, _ := datex.ISOString("PT5M").Parse()
		assert.False(t, p1.Equal(&p2))
	})

	t.Run("equal periods", func(t *testing.T) {
		p1, _ := datex.ISOString("PT5M").Parse()
		p2, _ := datex.ISOString("PT5M").Parse()
		assert.True(t, (&p1).Equal(&p2))
	})

	t.Run("different periods", func(t *testing.T) {
		p1, _ := datex.ISOString("PT5M").Parse()
		p2, _ := datex.ISOString("PT10M").Parse()
		assert.False(t, (&p1).Equal(&p2))
	})
}

func TestISOStringPtrOrNil(t *testing.T) {
	t.Run("nil period", func(t *testing.T) {
		var p *datex.Period
		iso := p.ISOStringPtrOrNil()
		assert.Nil(t, iso)
	})

	t.Run("valid period", func(t *testing.T) {
		p, _ := datex.ISOString("PT5M").Parse()
		iso := (&p).ISOStringPtrOrNil()
		require.NotNil(t, iso)
		assert.Equal(t, "PT5M", iso.String())
	})
}

func TestMustParse(t *testing.T) {
	t.Run("valid period", func(t *testing.T) {
		period := datex.MustParse(t, "PT5M")
		assert.Equal(t, "PT5M", period.String())
	})
}
