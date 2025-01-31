package clock_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/openmeterio/openmeter/openmeter/testutils"
	"github.com/openmeterio/openmeter/pkg/clock"
)

func TestClock(t *testing.T) {
	clock.SetTime(testutils.GetRFC3339Time(t, "2024-06-30T15:39:00Z"))
	defer clock.ResetTime()

	now := clock.Now()
	diff := now.Sub(testutils.GetRFC3339Time(t, "2024-06-30T15:39:00Z"))
	if diff < 0 {
		diff = -diff
	}
	assert.True(t, diff < time.Second)
}

func TestFreezeTime(t *testing.T) {
	frozenTime := time.Date(2024, 6, 30, 15, 39, 0, 0, time.UTC)
	clock.FreezeTime(frozenTime)
	defer clock.UnFreeze()

	// Test multiple calls return same frozen time
	for i := 0; i < 5; i++ {
		now := clock.Now()
		assert.Equal(t, frozenTime, now)
		time.Sleep(100 * time.Millisecond)
	}
}

func TestUnFreeze(t *testing.T) {
	frozenTime := time.Date(2024, 6, 30, 15, 39, 0, 0, time.UTC)
	clock.FreezeTime(frozenTime)

	// Verify frozen time
	assert.Equal(t, frozenTime, clock.Now())

	// Unfreeze and verify time advances
	clock.UnFreeze()
	time.Sleep(100 * time.Millisecond)
	assert.True(t, clock.Now().After(frozenTime))
}

func TestSetAndResetTime(t *testing.T) {
	originalTime := clock.Now()
	newTime := time.Date(2024, 6, 30, 15, 39, 0, 0, time.UTC)

	// Set new time
	clock.SetTime(newTime)
	setTime := clock.Now()
	assert.True(t, setTime.Sub(newTime) < time.Second)

	// Reset time
	clock.ResetTime()
	resetTime := clock.Now()
	assert.True(t, resetTime.After(originalTime))
}

func TestConcurrentAccess(t *testing.T) {
	done := make(chan bool)
	go func() {
		clock.FreezeTime(time.Now())
		done <- true
	}()
	go func() {
		clock.Now()
		done <- true
	}()

	<-done
	<-done

	clock.UnFreeze()
}
