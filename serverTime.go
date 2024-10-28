package timeManagement

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	useServerTime bool
	serverURL     string
	mu            sync.RWMutex
)

// TimeResponse 定義時間響應的結構
type TimeResponse struct {
	CurrentTime string `json:"currentTime"`
}

// SetUseServerTime 設置是否使用伺服器時間
func SetUseServerTime(use bool, url string) {
	mu.Lock()
	defer mu.Unlock()
	useServerTime = use
	serverURL = url
}

// Now 返回當前時間，根據配置選擇使用本地時間或伺服器時間
func Now() time.Time {
	mu.RLock()
	defer mu.RUnlock()

	if useServerTime {
		serverTime, err := getServerTime()
		if err == nil {
			return serverTime
		}
		fmt.Println("Error getting server time, falling back to local UTC time:", err)
	}

	return time.Now().UTC()
}

// getServerTime 從時間伺服器獲取當前時間
func getServerTime() (time.Time, error) {
	resp, err := http.Get(serverURL + "/time")
	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return time.Time{}, fmt.Errorf("failed to get time: %s", resp.Status)
	}

	var timeResponse TimeResponse
	if err := json.NewDecoder(resp.Body).Decode(&timeResponse); err != nil {
		return time.Time{}, err
	}

	serverTime, err := time.Parse(time.RFC3339Nano, timeResponse.CurrentTime)
	if err != nil {
		return time.Time{}, err
	}

	return serverTime, nil
}
