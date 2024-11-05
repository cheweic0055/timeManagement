package timeManagement

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetUseServerTime(t *testing.T) {
	tests := []struct {
		useServerTime bool
		serverURL     string
	}{
		{true, "http://example.com"},
		{false, ""},
	}

	for _, tt := range tests {
		SetUseServerTime(tt.useServerTime, tt.serverURL)

		mu.RLock()
		assert.Equal(t, tt.useServerTime, useServerTime, "useServerTime should match")
		assert.Equal(t, tt.serverURL, serverURL, "serverURL should match")
		mu.RUnlock()
	}
}

func TestServerNow(t *testing.T) {
	mockTime := time.Now().UTC().Add(1 * time.Hour).Format(time.RFC3339Nano)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"currentTime":"` + mockTime + `"}`))
	}))
	defer server.Close()

	tests := []struct {
		useServerTime bool
		serverURL     string
		expectError   bool
	}{
		{true, server.URL, false},
		{false, "", false},
		{true, "http://invalid-url", true},
	}

	for _, tt := range tests {
		SetUseServerTime(tt.useServerTime, tt.serverURL)

		currentTime := Now()

		if tt.useServerTime && !tt.expectError {
			expectedTime, err := time.Parse(time.RFC3339Nano, mockTime)
			require.NoError(t, err, "parsing mockTime should not produce an error")
			assert.Equal(t, expectedTime, currentTime, "expected server time should match current time")
		} else if !tt.useServerTime {
			assert.WithinDuration(t, time.Now().UTC(), currentTime, time.Millisecond, "expected local UTC time should be within 1ms of current time")
		}
	}
}
