package base

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"stove-gitlab.sginfra.net/backend/template/framework/util/log"
	"stove-gitlab.sginfra.net/backend/template/framework/web/config"
)

const StbRetCode = "APITraceRetCode"

var host string

func init() {
	host, _ = os.Hostname()
}

// New returns a new middleware handler with a default name and logger
func New() echo.MiddlewareFunc {
	return NewWithName("web")
}

// NewWithName returns a new middleware handler with the specified name
func NewWithName(name string) echo.MiddlewareFunc {
	return NewWithNameAndLogger("", name, logrus.StandardLogger())
}

// NewWithNameAndLogger returns a new middleware handler with the specified name
// and logger
func NewWithNameAndLogger(serviceName string, name string, l *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			entry := l.WithFields(logrus.Fields{
				"request": c.Request().RequestURI,
				"method":  c.Request().Method,
				"remote":  c.Request().RemoteAddr,
			})

			if serviceName != "" {
				entry = entry.WithField("servicename", serviceName)
			}

			if err := next(c); err != nil {
				c.Error(err)
			}

			if reqID := c.Request().Header.Get("X-Request-Id"); reqID != "" {
				entry = entry.WithField("request_id", reqID)
			}

			latency := time.Since(start)

			entry.WithFields(logrus.Fields{
				"status":                                c.Response().Status,
				"text_status":                           http.StatusText(c.Response().Status),
				"took":                                  latency,
				fmt.Sprintf("measure#%s.latency", name): latency.Nanoseconds(),
			}).Info()

			return nil
		}
	}
}

// APITraceLogger 안정화 지표용 에코 로거
func APITraceLogger(exculdeUrl []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			for _, url := range exculdeUrl {
				if strings.Contains(c.Request().URL.Path, url) {
					if err := next(c); err != nil {
						c.Error(err)
					}
					return nil
				}
			}

			start := time.Now()

			var reqBody []byte
			if c.Request().Body != nil {
				reqBody, _ = ioutil.ReadAll(c.Request().Body)
			}
			c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset

			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer

			if err := next(c); err != nil {
				c.Error(err)
			}

			latency := time.Since(start)
			logData := log.NewAPITraceData(c, log.DirectionIn, start, latency, reqBody, resBody)
			log.SetAPITraceLog(logData)

			return nil
		}
	}
}

// NewWithTimeFormat is new log with time format
func NewWithTimeFormat(timeFormat string) echo.MiddlewareFunc {
	return LogrusLogger(logrus.StandardLogger(), timeFormat)
}

// LogrusLogger is Another variant for better performance.
// With single log entry and time format.
func LogrusLogger(l *logrus.Logger, timeFormat string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			isError := false

			if err := next(c); err != nil {
				c.Error(err)
				isError = true
			}

			latency := time.Since(start)

			entry := l.WithFields(logrus.Fields{
				"server":  host,
				"path":    c.Request().RequestURI,
				"method":  c.Request().Method,
				"ip":      c.Request().RemoteAddr,
				"status":  c.Response().Status,
				"latency": latency,
				"time":    time.Now().Format(timeFormat),
			})

			if reqID := c.Request().Header.Get("X-Request-Id"); reqID != "" {
				entry = entry.WithField("request_id", reqID)
			}

			// Check middleware error
			if isError {
				entry.Error("error by handling request")
			} else {
				entry.Info("completed handling request")
			}

			return nil
		}
	}
}

// DetailInfoLogger 는 Request Body 및 Response Body 내용을 함께 출력하는 logger 입니다.
func DetailInfoLogger(serviceName string, name string, detail *config.DetailLog, l *logrus.Logger) echo.MiddlewareFunc {
	if detail.Request.Enable && detail.Request.MaxLen < 35 {
		detail.Request.MaxLen = 35
	}
	if detail.Response.Enable && detail.Response.MaxLen < 35 {
		detail.Response.MaxLen = 35
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			entry := l.WithFields(logrus.Fields{
				"request": c.Request().RequestURI,
				"method":  c.Request().Method,
				"remote":  c.Request().RemoteAddr,
			})

			if serviceName != "" {
				entry = entry.WithField("servicename", serviceName)
			}

			// Request (ref. body_dump.go of labstack/echo)
			reqBody := []byte{}
			if detail.Request.Enable {
				if c.Request().Body != nil { // Read
					reqBody, _ = ioutil.ReadAll(c.Request().Body)
				}
				c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset
			}

			// Response (ref. body_dump.go of labstack/echo)
			resBody := new(bytes.Buffer)
			if detail.Response.Enable {
				mw := io.MultiWriter(c.Response().Writer, resBody)
				writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
				c.Response().Writer = writer
			}

			if err := next(c); err != nil {
				c.Error(err)
			}

			if reqID := c.Request().Header.Get("X-Request-Id"); reqID != "" {
				entry = entry.WithField("request_id", reqID)
			}

			if resultCode := c.Get("resultCode"); resultCode != nil {
				entry = entry.WithField("result_code", resultCode)
			}

			if resultMessage := c.Get("resultMessage"); resultMessage != nil {
				entry = entry.WithField("result_message", resultMessage)
			}

			if customFieldArray := c.Get("customFieldArray"); customFieldArray != nil {
				if customField, ok := customFieldArray.([]string); ok {
					if len(customField) > 1 {
						// key, value 쌍 체크
						count := len(customField)
						if count%2 == 0 {
							for i := 0; i < count; i = i + 2 {
								entry = entry.WithField(customField[i], customField[i+1])
							}
						}
					}
				}
			}

			if detail.Request.Enable {
				if len(reqBody) > detail.Request.MaxLen {
					entry = entry.WithField("reqBody", string(reqBody[0:detail.Request.MaxLen])+"...[MaxLen]")
				} else {
					entry = entry.WithField("reqBody", string(reqBody[:]))
				}
			}
			if detail.Response.Enable {
				if len(resBody.String()) > detail.Response.MaxLen {
					entry = entry.WithField("resBody", resBody.String()[0:detail.Response.MaxLen]+"...[MaxLen]")
				} else {
					entry = entry.WithField("resBody", resBody.String())
				}
			}

			latency := time.Since(start)

			var latencyTime int64
			if detail.LatencyUnitMilliSec {
				latencyTime = latency.Nanoseconds() / 1000000
			} else {
				latencyTime = latency.Nanoseconds()
			}

			entry.WithFields(logrus.Fields{
				"status":                                c.Response().Status,
				"text_status":                           http.StatusText(c.Response().Status),
				"took":                                  latency,
				fmt.Sprintf("measure#%s.latency", name): latencyTime,
			}).Info()

			return nil
		}
	}
}

// ref. body_dump.go of labstack/echo
type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w *bodyDumpResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
