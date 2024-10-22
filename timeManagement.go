package timeManagement

import (
	"sync"
	"time"
)

// 常用的時間格式常量
const (
	DateFormat          = "2006-01-02"
	TimeFormat          = "15:04:05"
	DateTimeFormat      = "2006-01-02 15:04:05"
	DateTimeFormatTZ    = "2006-01-02T15:04:05Z07:00"
	DateTimeFormatMilli = "2006-01-02 15:04:05.000"
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

	// 返回UTC時間，支持時間加速
	Now() time.Time

	// 返回特定時區的時間
	NowInZone(location *time.Location) time.Time

	// 當前時間 - 指定時間
	Since(t time.Time) time.Duration

	// 指定時間 - 當前時間
	Until(t time.Time) time.Duration

	// 睡眠指定時間，支持時間加速
	Sleep(d time.Duration)

	// 返回一個通道，指定時間後會發送一個時間，支持時間加速
	After(d time.Duration) <-chan time.Time

	// 解析時間字符串，返回UTC時間
	Parse(layout, value string) (time.Time, error)

	// 解析指定時區的時間字符串，返回UTC時間
	ParseInLocation(layout, value string, loc *time.Location) (time.Time, error)

	// 格式化時間為字符串
	Format(t time.Time, layout string) string

	// 將任何時間轉換為UTC
	ToUTC(t time.Time) time.Time

	// 將UTC時間轉換為指定時區
	ToZone(t time.Time, location *time.Location) time.Time

	// 將時間轉換為Unix時間戳
	ToUnix(t time.Time) int64

	// 將時間轉換為Unix毫秒時間戳
	ToUnixMilli(t time.Time) int64
}

type realTimeProvider struct{}

func (r *realTimeProvider) Now() time.Time {
	mockTimeLock.RLock()
	defer mockTimeLock.RUnlock()

	if mockTime != nil {
		return mockTime.UTC()
	}

	if timeScale != 1.0 {
		realElapsed := time.Since(scaleStart)
		scaledElapsed := time.Duration(float64(realElapsed) * timeScale)
		return baseTime.Add(scaledElapsed).UTC()
	}

	return time.Now().UTC()
}

func (r *realTimeProvider) NowInZone(location *time.Location) time.Time {
	return r.Now().In(location)
}

func (r *realTimeProvider) Since(t time.Time) time.Duration {
	mockTimeLock.RLock()
	defer mockTimeLock.RUnlock()

	if mockTime != nil {
		return mockTime.Sub(t.UTC())
	}

	return r.Now().Sub(t.UTC())
}

func (r *realTimeProvider) Until(t time.Time) time.Duration {
	mockTimeLock.RLock()
	defer mockTimeLock.RUnlock()

	if mockTime != nil {
		return t.UTC().Sub(*mockTime)
	}

	return t.UTC().Sub(r.Now())
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
	mockTimeLock.RLock()
	defer mockTimeLock.RUnlock()

	if mockTime != nil {
		// 模擬定時器，立即返回一個已經過期的通道
		ch := make(chan time.Time, 1)
		ch <- *mockTime
		return ch
	}

	if timeScale != 1.0 {
		adjustedDuration := time.Duration(float64(d) / timeScale)
		return time.After(adjustedDuration)
	}
	return time.After(d)
}

func (r *realTimeProvider) Parse(layout, value string) (time.Time, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

func (r *realTimeProvider) ParseInLocation(layout, value string, loc *time.Location) (time.Time, error) {
	t, err := time.ParseInLocation(layout, value, loc)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

func (r *realTimeProvider) Format(t time.Time, layout string) string {
	return t.UTC().Format(layout)
}

func (r *realTimeProvider) ToUTC(t time.Time) time.Time {
	return t.UTC()
}

func (r *realTimeProvider) ToZone(t time.Time, location *time.Location) time.Time {
	return t.In(location)
}

func (r *realTimeProvider) ToUnix(t time.Time) int64 {
	return t.UTC().Unix()
}

func (r *realTimeProvider) ToUnixMilli(t time.Time) int64 {
	return t.UTC().UnixMilli()
}

var defaultProvider TimeProvider = &realTimeProvider{}

// GetProvider 返回當前的時間提供者
func GetProvider() TimeProvider {
	return defaultProvider
}

// SetTimeScale 設置時間加速比例
func SetTimeScale(scale float64) {
	if scale <= 0 {
		panic("Time scale must be positive")
	}
	// 先獲取當前時間
	currentTime := defaultProvider.Now()

	mockTimeLock.Lock()
	defer mockTimeLock.Unlock()

	baseTime = currentTime
	scaleStart = time.Now().UTC()
	timeScale = scale
}

func GetTimeScale() float64 {
	mockTimeLock.RLock()
	defer mockTimeLock.RUnlock()
	return timeScale
}

func ResetTimeScale() {
	SetTimeScale(1.0)
}

func SetMockTime(t time.Time) {
	mockTimeLock.Lock()
	defer mockTimeLock.Unlock()
	utcTime := t.UTC()
	mockTime = &utcTime
	timeScale = 1.0
}

func ClearMockTime() {
	mockTimeLock.Lock()
	defer mockTimeLock.Unlock()
	mockTime = nil
}
