package log

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
)

// APITrace 안정화 지표 설정 구조체, 스토브 시스템에서 사용하는(정의된) 안정화 지표용 로그를 저장하기 위한 설정
//
// ServiceName: 서비스명, 로그에 servicename 필드에 출력되는 정보
// Use: 안정화 지표 로그 활성화
// Path: 로그 파일 경로
// KeepDay: 로그 파일 보과 기간 (기본값: 7 일)
// RotateTime: 파일 로테이션 시간 (기본값: 24 시간)
// Console: 로그 콘솔 출력 활성화
// ExcludeUrl: 로그에서 제외할 URL 경로 ["/healthcheck/_check"]
// IncludeBodyFor: 로그에서 body 내용 출력 허용할 content-type (request, response)
type APITrace struct {
	ServiceName    string   `json:"service_name" yaml:"service_name"`
	Use            bool     `json:"use" yaml:"use"`
	Path           string   `json:"path" yaml:"path"`
	KeepDay        uint     `json:"keep_day" yaml:"keep_day"`
	RotateTime     uint     `json:"rotate_time_h" yaml:"rotate_time_h"`
	Console        bool     `json:"console" yaml:"console"`
	ExcludeUrl     []string `json:"exclude_url" yaml:"exclude_url"`
	IncludeBodyFor []string `json:"include_body_for" yaml:"include_body_for"`
}

type APITraceData struct {
	ServiceName string      `json:"servicename"`
	Timestamp   time.Time   `json:"timestamp"` // UTC timezone
	Direction   string      `json:"direction"` // in, out
	Type        string      `json:"type"`      // http
	Event       Event       `json:"event"`
	Destination Destination `json:"destination"`
	Source      Source      `json:"source"`
	Url         Url         `json:"url"`
	Http        Http        `json:"http"`
}

type Event struct {
	Start    time.Time     `json:"start"`    // UTC timezone
	End      time.Time     `json:"end"`      // UTC timezone
	Duration time.Duration `json:"duration"` // nanoseconds
}

type Destination struct {
	Domain string `json:"domain,omitempty"` // Request.Host
	IP     string `json:"ip"`               // Request.Host
	Port   int    `json:"port"`
}

type Source struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type Url struct {
	Domain string `json:"domain"` // Request.Host
	Full   string `json:"full"`
	Path   string `json:"path"`
	Query  string `json:"query,omitempty"`
}

type Http struct {
	Version  string   `json:"version"`
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

type Request struct {
	Method  string              `json:"method"`
	Headers map[string][]string `json:"headers,omitempty"`
	Body    *Body               `json:"body,omitempty"`
}

type Response struct {
	StatusCode   int                 `json:"status_code"`
	StatusPhrase string              `json:"status_phrase,omitempty"` // 처리 안함
	Headers      map[string][]string `json:"headers,omitempty"`
	Body         *Body               `json:"body,omitempty"`
}

type Body struct {
	Bytes   int    `json:"bytes,omitempty"`
	Content string `json:"content,omitempty"`
}

const (
	DirectionIn  = "in"  // 서버로 들어오는 요청
	DirectionOut = "out" // 외부로 호출하는 요청
)

func (s APITraceData) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}

// 안정화지표 로그 데이터 생성
func NewAPITraceData(c echo.Context, direction string, start time.Time, latency time.Duration, reqBody []byte, resBody *bytes.Buffer) *APITraceData {

	// Event
	end := start.Add(latency)

	// Destination
	var serverPortNum int
	serverHost, serverPort, err := net.SplitHostPort(c.Request().Host)
	if err != nil {
		serverHost = c.Request().Host
		hostAddr := strings.Split(c.Echo().Server.Addr, ":")
		if len(hostAddr) < 2 {
			serverPortNum = 0
		} else {
			serverPortNum, _ = strconv.Atoi(hostAddr[1])
		}
	} else {
		serverPortNum, _ = strconv.Atoi(serverPort)
	}

	// Source
	var remotePortNum int
	remoteHost, remotePort, remoteErr := net.SplitHostPort(c.Request().RemoteAddr)
	if remoteErr != nil {
		remoteHost = c.Request().RemoteAddr
		remotePortNum = 0
	} else {
		remotePortNum, _ = strconv.Atoi(remotePort)
	}

	// Url
	full := fmt.Sprintf("%s://%s%s", c.Scheme(), c.Request().Host, c.Request().RequestURI)

	// Http
	version := fmt.Sprintf("%v.%v", c.Request().ProtoMajor, c.Request().ProtoMinor)

	// Http.Request
	request := Request{}
	request.Method = c.Request().Method

	if len(c.Request().Header) > 0 {
		request.Headers = make(map[string][]string)
		for key, value := range c.Request().Header {
			request.Headers[strings.ToLower(key)] = value
		}
	}

	reqLen := len(reqBody)
	if reqLen > 0 {
		reqData := Body{}
		reqData.Bytes = reqLen
		contentType := c.Request().Header.Get(echo.HeaderContentType)
		for _, include := range gLogger.conf.APITrace.IncludeBodyFor {
			if strings.Contains(contentType, include) {
				reqData.Content = fmt.Sprintf("%v", string(reqBody))
				break
			}
		}
		request.Body = &reqData
	}

	// Http.Response
	response := Response{}
	response.StatusCode = c.Response().Status

	if len(c.Response().Header()) > 0 {
		response.Headers = make(map[string][]string)
		for key, value := range c.Response().Header() {
			response.Headers[strings.ToLower(key)] = value
		}
	}

	resLen := resBody.Len()
	if resLen > 0 {
		resData := Body{}
		resData.Bytes = resLen
		contentType := c.Response().Header().Get(echo.HeaderContentType)
		for _, include := range gLogger.conf.APITrace.IncludeBodyFor {
			if strings.Contains(contentType, include) {
				switch c.Response().Header().Get(echo.HeaderContentEncoding) {
				case "gzip":
					var buf bytes.Buffer
					err = decompressGzip(resBody.Bytes(), &buf)
					if err == nil {
						resData.Content = buf.String()
					}
				default:
					resData.Content = resBody.String()
				}
				break
			}
		}
		response.Body = &resData
	}

	return &APITraceData{
		ServiceName: gLogger.conf.APITrace.ServiceName,
		Timestamp:   start,
		Type:        "http",
		Direction:   direction,
		Event: Event{
			Start:    start,
			End:      end,
			Duration: latency,
		},
		Destination: Destination{
			Domain: serverHost,
			IP:     serverHost,
			Port:   serverPortNum,
		},
		Source: Source{
			IP:   remoteHost,
			Port: remotePortNum,
		},
		Url: Url{
			Domain: serverHost,
			Full:   full,
			Path:   c.Request().URL.Path,
			Query:  c.QueryString(),
		},
		Http: Http{
			Version:  version,
			Request:  request,
			Response: response,
		},
	}
}

// APITraceLogger 로직 안정화 지표 로그
func SetAPITraceLog(data *APITraceData) {
	if gLogger == nil {
		fmt.Println(data.String())
		return
	}

	gLogger.apiLogger.WithFields(logrus.Fields{
		"servicename": data.ServiceName,
		"timestamp":   data.Timestamp,
		"type":        data.Type,
		"direction":   data.Direction,
		"event":       data.Event,
		"destination": data.Destination,
		"source":      data.Source,
		"url":         data.Url,
		"http":        data.Http,
	}).Print()
}

func decompressGzip(data []byte, w io.Writer) error {
	r, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer r.Close()

	_, err = io.Copy(w, r)
	return err
}
