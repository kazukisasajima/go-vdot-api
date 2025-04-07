package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

// ANSIカラーコード
const (
	green  = "\033[32m"
	blue   = "\033[34m"
	yellow = "\033[33m"
	red    = "\033[31m"
	reset  = "\033[0m"
)

func Info(format string, a ...interface{}) {
	printWithColor(green, "INFO", format, a...)
}

func Warn(format string, a ...interface{}) {
	printWithColor(yellow, "WARN", format, a...)
}

func Error(format string, a ...interface{}) {
	printWithColor(red, "ERROR", format, a...)
}

func Debug(format string, a ...interface{}) {
	printWithColor(blue, "DEBUG", format, a...)
}

func printWithColor(color string, level string, format string, a ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, a...)
	fmt.Fprintf(os.Stdout, "%s[%s] [%s] %s%s\n", color, timestamp, level, message, reset)
}

func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c) // ここで実行（中でエラーが起きることもある）
			stop := time.Now()

			status := c.Response().Status
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					status = he.Code // 明示的に HTTPError のコードにする
				}
			}

			method := c.Request().Method
			path := c.Path()
			latency := stop.Sub(start)

			switch {
			case status >= 500:
				Error("%s %s %d (%s)", method, path, status, latency)
			case status >= 400:
				Warn("%s %s %d (%s)", method, path, status, latency)
			default:
				Info("%s %s %d (%s)", method, path, status, latency)
			}

			return err
		}
	}
}
