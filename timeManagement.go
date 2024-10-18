package timeManagement

import (
	"sync"
	"time"
)

var (
	mockTime     *time.Time
	mockTimeLock sync.RWMutex
	timeScale    float64 = 1.0
	baseTime     time.Time
	scaleStart   time.Time
)

// TimeProvider 提供所有時間相關的操作介面
type TimeProvider interface {
	Now() time.Time
	Since(t time.Time) time.Duration
	Until(t time.Time) time.Duration
	Sleep(d time.Duration)
	After(d time.Duration) <-chan time.Time
	Parse(layout, value string) (time.Time, error)
	Format(t time.Time, layout string) string
}

// realTimeProvider 實際時間的實現
type realTimeProvider struct{}

func (r *realTimeProvider) Now() time.Time {
	mockTimeLock.RLock()
	defer mockTimeLock.RUnlock()

	if mockTime != nil {
		return *mockTime
	}

	if timeScale != 1.0 {
		realElapsed := time.Since(scaleStart)
		scaledElapsed := time.Duration(float64(realElapsed) * timeScale)
		return baseTime.Add(scaledElapsed)
	}

	return time.Now()
}

func (r *realTimeProvider) Since(t time.Time) time.Duration {
	return r.Now().Sub(t)
}

func (r *realTimeProvider) Until(t time.Time) time.Duration {
	return t.Sub(r.Now())
}

func (r *realTimeProvider) Sleep(d time.Duration) {
	if timeScale != 1.0 {
		adjustedDuration := time.Duration(float64(d) / timeScale)
		time.Sleep(adjustedDuration)
		return
	}
	time.Sleep(d)
}

func (r *realTimeProvider) After(d time.Duration) <-chan time.Time {
	if timeScale != 1.0 {
		adjustedDuration := time.Duration(float64(d) / timeScale)
		return time.After(adjustedDuration)
	}
	return time.After(d)
}

func (r *realTimeProvider) Parse(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

func (r *realTimeProvider) Format(t time.Time, layout string) string {
	return t.Format(layout)
}

var defaultProvider TimeProvider = &realTimeProvider{}

// GetProvider 返回當前的時間提供者
func GetProvider() TimeProvider {
	return defaultProvider
}

// SetTimeScale 設置時間加速比例
// scale > 1.0 表示時間加速，例如 2.0 表示時間流逝速度是正常的兩倍
// scale < 1.0 表示時間減速，例如 0.5 表示時間流逝速度是正常的一半
// scale = 1.0 表示正常時間流逝
func SetTimeScale(scale float64) {
	if scale <= 0 {
		panic("Time scale must be positive")
	}

	mockTimeLock.Lock()
	defer mockTimeLock.Unlock()

	// 設置新的時間比例前，先保存當前時間作為基準點
	baseTime = defaultProvider.Now()
	scaleStart = time.Now()
	timeScale = scale
}

// GetTimeScale 獲取當前時間加速比例
func GetTimeScale() float64 {
	mockTimeLock.RLock()
	defer mockTimeLock.RUnlock()
	return timeScale
}

// ResetTimeScale 重置時間比例為正常速度
func ResetTimeScale() {
	SetTimeScale(1.0)
}

// SetMockTime 設置模擬時間（用於測試）
func SetMockTime(t time.Time) {
	mockTimeLock.Lock()
	defer mockTimeLock.Unlock()
	mockTime = &t
	timeScale = 1.0 // 重置時間比例
}

// ClearMockTime 清除模擬時間
func ClearMockTime() {
	mockTimeLock.Lock()
	defer mockTimeLock.Unlock()
	mockTime = nil
}

// 常用的時間格式常量
const (
	DateFormat          = "2006-01-02"
	TimeFormat          = "15:04:05"
	DateTimeFormat      = "2006-01-02 15:04:05"
	DateTimeFormatTZ    = "2006-01-02T15:04:05Z07:00"
	DateTimeFormatMilli = "2006-01-02 15:04:05.000"
)
