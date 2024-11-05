# time-management

一個客製化的時間套件這是一個用於管理和操作時間的 Go 庫，提供了多種時間相關的功能，包括時間加速、模擬時間、格式化時間等。


## 功能

- 模擬時間至指定時刻，輕鬆測試不同時間點的行為
- 支援時間加速，方便模擬時間流逝的效果
- 可從伺服器獲取時間，避免使用本地時間
- 獲取當前時間，支持時間加速
- 獲取特定時區的時間
- 計算時間差
- 睡眠指定時間，支持時間加速
- 返回一個通道，指定時間後會發送一個時間，支持時間加速
- 解析和格式化時間字符串
- 將時間轉換為 UTC 或指定時區



## 安裝
使用 go get 安裝：

待修正
```bash
go get git-golang.yile808.com/web3/common/time-management.git
```

## 使用方法

獲取當前時間
```bash
package main

import (
    "fmt"
    "timeManagement"
)

func main() {
    provider := timeManagement.GetProvider()
    currentTime := provider.Now()
    fmt.Println("Current Time:", currentTime)
}


```
設置時間加速比例
```bash
provider.SetTimeScale(2.0) // 時間加速2倍 (預設 1.0)
```
設置模擬時間
```bash
mockTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
provider.SetMockTime(mockTime)
```
清除模擬時間
```bash
provider.ClearMockTime()
```
選擇是否使用伺服器時間 ( 伺服器需先引入此 library 以及開放 api 端口 )
```bash
timeManagement.SetUseServerTime(true, "http://example.com")
currentTime := timeManagement.Now()
fmt.Println("Server Time:", currentTime)
```
## API

TimeProvider 介面

- Now() time.Time：返回當前時間，支持時間加速
- NowInZone(location *time.Location) time.Time：返回特定時區的時間
- Since(t time.Time) time.Duration：當前時間 - 指定時間
- Until(t time.Time) time.Duration：指定時間 - 當前時間
- Sleep(d time.Duration)：睡眠指定時間，支持時間加速
- After(d time.Duration) <-chan time.Time：返回一個通道，指定時間後會發送一個時間，支持時間加速
- Parse(layout, value string) (time.Time, error)：解析時間字符串，返回 UTC 時間
- ParseInLocation(layout, value string, loc *time.Location) (time.Time, error)：解析指定時區的時間字符串，返回 UTC 時間
- Format(t time.Time, layout string) string：格式化時間為字符串
- UTC(t time.Time) time.Time：將任何時間轉換為 UTC
- In(t time.Time, location *time.Location) time.Time：將 UTC 時間轉換為指定時區
- Unix(t time.Time) int64：將時間轉換為 Unix 時間戳
- UnixMilli(t time.Time) int64：將時間轉換為 Unix 毫秒時間戳
- SetTimeScale(scale float64)：設置時間加速比例
- GetTimeScale() float64：獲取當前的時間加速比例
- ClearTimeScale()：清除時間加速比例，恢復為 1.0
- SetMockTime(t time.Time)：設置模擬時間
- ClearMockTime()：清除模擬時間

全局函數
- GetProvider() TimeProvider：獲取默認的時間提供者
- SetUseServerTime(use bool, url string)：設置是否使用伺服器時間
- Now() time.Time：返回當前時間，根據配置選擇使用本地時間或伺服器時間


## 系統架構圖

![系統架構圖]([https://git-golang.yile808.com/web3/common/time-management/-/raw/develop/docs/img/SystemArchitectureDiagram.png?ref_type=heads&inline=false](https://github.com/cheweic0055/timeManagement/blob/main/docs/img/SystemArchitectureDiagram.png))



