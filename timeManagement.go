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
	UTC(t time.Time) time.Time

	// 將UTC時間轉換為指定時區
	In(t time.Time, location *time.Location) time.Time

	// 將時間轉換為Unix時間戳
	Unix(t time.Time) int64

	// 將時間轉換為Unix毫秒時間戳
	UnixMilli(t time.Time) int64

	// 設置時間加速比例
	SetTimeScale(scale float64)

	// 獲取時間加速比例
	GetTimeScale() float64

	// 清除時間加速比例
	ClearTimeScale()

	// 設置模擬時間
	SetMockTime(t time.Time)

	// 清除模擬時間
	ClearMockTime()
}

type realTimeProvider struct {
	mockTime      *time.Time
	mockStartTime time.Time
	mockBaseTime  time.Time
	mockTimeLock  sync.RWMutex
	timeScale     float64
	baseTime      time.Time
	scaleStart    time.Time
}

var (
	instance *realTimeProvider
	once     sync.Once
)

// GetProvider 返回 TimeProvider 的單例實例
func GetProvider() TimeProvider {
	once.Do(func() {
		instance = &realTimeProvider{
			timeScale: 1.0,
		}
	})
	return instance
}

func (r *realTimeProvider) Now() time.Time {
	r.mockTimeLock.RLock()
	defer r.mockTimeLock.RUnlock()

	if r.mockTime != nil {
		// 計算從設置模擬時間開始經過的時間
		elapsed := time.Since(r.mockStartTime)
		if r.timeScale != 1.0 {
			elapsed = time.Duration(float64(elapsed) * r.timeScale)
		}
		return r.mockBaseTime.Add(elapsed).UTC()
	}

	if r.timeScale != 1.0 {
		realElapsed := time.Since(r.scaleStart)
		scaledElapsed := time.Duration(float64(realElapsed) * r.timeScale)
		return r.baseTime.Add(scaledElapsed).UTC()
	}

	return time.Now().UTC()
}

func (r *realTimeProvider) NowInZone(location *time.Location) time.Time {
	return r.Now().In(location)
}

func (r *realTimeProvider) Since(t time.Time) time.Duration {
	return r.Now().Sub(t.UTC())
}

func (r *realTimeProvider) Until(t time.Time) time.Duration {
	return t.UTC().Sub(r.Now())
}

func (r *realTimeProvider) Sleep(d time.Duration) {
	if r.timeScale != 1.0 {
		adjustedDuration := time.Duration(float64(d) / r.timeScale)
		time.Sleep(adjustedDuration)
		return
	}
	time.Sleep(d)
}

func (r *realTimeProvider) After(d time.Duration) <-chan time.Time {
	if r.timeScale != 1.0 {
		adjustedDuration := time.Duration(float64(d) / r.timeScale)
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

func (r *realTimeProvider) UTC(t time.Time) time.Time {
	return t.UTC()
}

func (r *realTimeProvider) In(t time.Time, location *time.Location) time.Time {
	return t.In(location)
}

func (r *realTimeProvider) Unix(t time.Time) int64 {
	return t.UTC().Unix()
}

func (r *realTimeProvider) UnixMilli(t time.Time) int64 {
	return t.UTC().UnixMilli()
}

func (r *realTimeProvider) SetTimeScale(scale float64) {
	if scale <= 0 {
		panic("Time scale must be positive")
	}
	// 先獲取當前時間
	currentTime := r.Now()

	r.mockTimeLock.Lock()
	defer r.mockTimeLock.Unlock()

	if r.mockTime != nil {
		// 更新模擬時間的基準時間和開始時間
		elapsed := time.Since(r.mockStartTime)
		r.mockBaseTime = r.mockBaseTime.Add(elapsed)
		r.mockStartTime = time.Now()
	}

	r.baseTime = currentTime
	r.scaleStart = time.Now().UTC()
	r.timeScale = scale
}

func (r *realTimeProvider) GetTimeScale() float64 {
	r.mockTimeLock.RLock()
	defer r.mockTimeLock.RUnlock()
	return r.timeScale
}

func (r *realTimeProvider) ClearTimeScale() {
	r.SetTimeScale(1.0)
}

func (r *realTimeProvider) SetMockTime(t time.Time) {
	r.mockTimeLock.Lock()
	defer r.mockTimeLock.Unlock()
	utcTime := t.UTC()
	r.mockBaseTime = utcTime
	r.mockStartTime = time.Now()
	r.mockTime = &utcTime
	r.timeScale = 1.0
}

func (r *realTimeProvider) ClearMockTime() {
	r.mockTimeLock.Lock()
	defer r.mockTimeLock.Unlock()
	r.mockTime = nil
	r.timeScale = 1.0
}
