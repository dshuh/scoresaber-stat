package webhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

// SmilehubHookConf smile hub 메신저 web hook 설정 구조체
type SmilehubHookConf struct {
	Use        bool   `json:"use" yaml:"use"`
	Url        string `json:"url" yaml:"url"`
	WorkerSize int    `json:"worker_size" yaml:"worker_size"`
}

type SmilehubHook struct {
	addr      *url.URL
	levels    []logrus.Level
	formatter logrus.Formatter

	smilehubHookCh chan MessageParam
}

type MessageParam struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Title     string `json:"title,omitempty"`
	TitleLink string `json:"title_link,omitempty"`
	Text      string `json:"text,omitempty"`
	ImageURL  string `json:"image_url,omitempty"`
	Color     string `json:"color,omitempty"`
}

func NewSmilehubHook(addr *url.URL, levels []logrus.Level, formatter logrus.Formatter, workerSize int) (*SmilehubHook, error) {
	// 채널 생성
	smilehubHookCh := make(chan MessageParam, workerSize)

	// 워커 생성
	for i := 0; i < workerSize; i++ {
		go func(idx int) {
			for param := range smilehubHookCh {
				resp, _ := SendSmilehub(addr, param)
				if resp != nil {
					_ = resp.Body.Close()
				}
			}
		}(i)
	}

	hook := &SmilehubHook{
		addr,
		levels,
		formatter,
		smilehubHookCh,
	}
	return hook, nil
}

func (hook *SmilehubHook) Levels() []logrus.Level {
	return hook.levels
}

// Fire smilehub web call 호출
// 비동기 처리로 진행 채널 버퍼 오버 플로우 시 메시지 드롭 (smilhub web hook 응답 지연 으로 이한 성능 저하)
func (hook *SmilehubHook) Fire(entry *logrus.Entry) error {
	d, err := entry.String()
	if err != nil {
		return err
	}

	select {
	case hook.smilehubHookCh <- MessageParam{Text: d}:
	default:
		return errors.New("smilehubHookCh is full")
	}
	return nil
}

// Send smile hub 메신져 에 메시지 를 전송 합니다
// 미리 등록된 웹훅 url 을 통해 text 메시지 를 전송
// 기본 text 기능만 구현(attachments 기능 필요시 추가 개발 필요)
func SendSmilehub(url *url.URL, param MessageParam) (*http.Response, error) {
	var b []byte
	b, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer(b)
	resp, err := http.Post(url.String(), "application/json", buff)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
