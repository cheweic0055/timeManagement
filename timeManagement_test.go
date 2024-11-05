package timeManagement

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNow(t *testing.T) {
	provider := GetProvider()
	now := provider.Now()
	assert.False(t, now.IsZero(), "Expected non-zero time")
}

func TestNowInZone(t *testing.T) {
	provider := GetProvider()
	location, err := time.LoadLocation("America/New_York")
	require.NoError(t, err, "Failed to load location")
	now := provider.NowInZone(location)
	assert.Equal(t, location, now.Location(), "Expected location to match")
}

func TestSince(t *testing.T) {
	provider := GetProvider()
	start := provider.Now()
	time.Sleep(10 * time.Millisecond)
	duration := provider.Since(start)
	assert.GreaterOrEqual(t, duration, 10*time.Millisecond, "Expected duration >= 10ms")
}

func TestUntil(t *testing.T) {
	provider := GetProvider()
	future := provider.Now().Add(10 * time.Millisecond)
	duration := provider.Until(future)
	assert.GreaterOrEqual(t, duration, 9*time.Millisecond, "Expected duration >= 9ms")
}

func TestSleep(t *testing.T) {
	provider := GetProvider()

	tests := []struct {
		sleepDuration time.Duration
		minDuration   time.Duration
	}{
		{10 * time.Millisecond, 10 * time.Millisecond},
		{50 * time.Millisecond, 50 * time.Millisecond},
		{100 * time.Millisecond, 100 * time.Millisecond},
		{0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.sleepDuration.String(), func(t *testing.T) {
			start := provider.Now()
			provider.Sleep(tt.sleepDuration)
			duration := provider.Since(start)
			assert.GreaterOrEqual(t, duration, tt.minDuration, "Expected sleep duration >= %v", tt.minDuration)
		})
	}
}

func TestAfter(t *testing.T) {
	provider := GetProvider()

	tests := []struct {
		delayDuration time.Duration
		timeout       time.Duration
	}{
		{10 * time.Millisecond, 20 * time.Millisecond},
		{50 * time.Millisecond, 100 * time.Millisecond},
		{100 * time.Millisecond, 200 * time.Millisecond},
		{0, 10 * time.Millisecond},
	}

	for _, tt := range tests {
		t.Run(tt.delayDuration.String(), func(t *testing.T) {
			ch := provider.After(tt.delayDuration)
			select {
			case <-ch:
				// success
				assert.True(t, true, "Expected to receive from channel")
			case <-time.After(tt.timeout):
				assert.Fail(t, "Expected to receive from channel within %v", tt.timeout)
			}
		})
	}
}

func TestParse(t *testing.T) {
	provider := GetProvider()
	layout := DateTimeFormat
	value := "2023-01-01 12:00:00"
	parsedTime, err := provider.Parse(layout, value)
	require.NoError(t, err, "Failed to parse time")
	expectedTime, _ := time.Parse(layout, value)
	assert.True(t, parsedTime.Equal(expectedTime.UTC()), "Expected parsed time to match")
}

func TestParseInLocation(t *testing.T) {
	provider := GetProvider()
	layout := DateTimeFormat
	value := "2023-01-01 12:00:00"
	location, err := time.LoadLocation("America/New_York")
	require.NoError(t, err, "Failed to load location")
	parsedTime, err := provider.ParseInLocation(layout, value, location)
	require.NoError(t, err, "Failed to parse time")
	expectedTime, _ := time.ParseInLocation(layout, value, location)
	assert.True(t, parsedTime.Equal(expectedTime.UTC()), "Expected parsed time to match")
}

func TestFormat(t *testing.T) {
	provider := GetProvider()
	now := provider.Now()
	formatNow := provider.Format(now, DateTimeFormat)
	parseNow, _ := time.Parse(DateTimeFormat, formatNow)
	formatted := provider.Format(now, DateTimeFormat)
	parsedTime, _ := time.Parse(DateTimeFormat, formatted)
	assert.True(t, parsedTime.Equal(parseNow.UTC()), "Expected formatted time to match")
}

func TestUTC(t *testing.T) {
	provider := GetProvider()
	now := provider.Now()
	utcTime := provider.UTC(now)
	assert.True(t, utcTime.Equal(now.UTC()), "Expected UTC time to match")
}

func TestIn(t *testing.T) {
	provider := GetProvider()
	location, err := time.LoadLocation("America/New_York")
	require.NoError(t, err, "Failed to load location")
	now := provider.Now()
	inTime := provider.In(now, location)
	assert.Equal(t, location, inTime.Location(), "Expected location to match")
}

func TestUnix(t *testing.T) {
	provider := GetProvider()
	now := provider.Now()
	unix := provider.Unix(now)
	assert.Equal(t, now.UTC().Unix(), unix, "Expected Unix timestamp to match")
}

func TestUnixMilli(t *testing.T) {
	provider := GetProvider()
	now := provider.Now()
	unixMilli := provider.UnixMilli(now)
	assert.Equal(t, now.UTC().UnixMilli(), unixMilli, "Expected Unix milli timestamp to match")
}

func TestSetTimeScale(t *testing.T) {
	provider := GetProvider()
	provider.SetTimeScale(2.0)
	assert.Equal(t, 2.0, provider.GetTimeScale(), "Expected time scale to be 2.0")
	provider.ClearTimeScale()
	assert.Equal(t, 1.0, provider.GetTimeScale(), "Expected time scale to be reset to 1.0")
}

func TestSetMockTime(t *testing.T) {
	provider := GetProvider()
	mockTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	provider.SetMockTime(mockTime)
	now := provider.Now()
	formatNow := provider.Format(now, DateTimeFormat)
	parseNow, _ := time.Parse(DateTimeFormat, formatNow)
	assert.True(t, parseNow.Equal(mockTime), "Expected mock time to match")
	provider.ClearMockTime()
}

func TestSingleton(t *testing.T) {
	provider1 := GetProvider()
	provider2 := GetProvider()

	assert.Equal(t, provider1, provider2, "Expected both providers to be the same instance")
}
