package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type JueJinEventRequestBody struct {
	FromDate string `json:"from_date"`
	Days     int64  `json:"days"`
	Aid      int64  `json:"aid"`
}

type ClientConn struct {
	C            http.Client
	Method       string
	URL          string
	Header       map[string]string
	RequestBody  any
	ResponseData any
}

type JueJinEventResponseData struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	Datas  []struct {
		Date   string `json:"date"`
		Count  int    `json:"count"`
		Events []struct {
			Title   string `json:"title"`
			URL     string `json:"url"`
			EventID string `json:"event_id"`
			IDType  int    `json:"id_type"`
		} `json:"events"`
	} `json:"data"`
}

type FeishuRequestBody struct {
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
	MsgType   string `json:"msg_type"`
	Content   struct {
		Text string `json:"text"`
	} `json:"content"`
}

type FeishuResponseData struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

var Client *http.Client

func init() {
	Client = &http.Client{}
}

func main() {
	FeishuNotify(JueJinEvent())
}

func JueJinEvent() *JueJinEventResponseData {
	fromDate := time.Now().Format("2006-01-02")
	cookie := os.Getenv("COOKIE")
	token := os.Getenv("TOKEN")
	days, err := strconv.ParseInt(os.Getenv("Days"), 10, 64)
	if err != nil {
		panic(err)
	}
	aid, err := strconv.ParseInt(os.Getenv("Aid"), 10, 64)
	if err != nil {
		panic(err)
	}
	body := &JueJinEventRequestBody{
		FromDate: fromDate,
		Days:     days,
		Aid:      aid,
	}

	var outData *JueJinEventResponseData
	conn := ClientConn{
		*Client,
		"POST",
		"https://api.juejin.cn/event_api/v1/event/month_stat",
		map[string]string{"Content-Type": "application/json", "cookie": cookie, "x-secsdk-csrf-token": token},
		body,
		&outData,
	}
	conn.Do()

	return outData
}

func FeishuNotify(inData any) {
	data := inData.(*JueJinEventResponseData)

	secret := os.Getenv("SECRET")
	timestamp := time.Now().Unix()
	sign, err := FeishuGenSign(secret, timestamp)
	if err != nil {
		panic(err)
	}

	var text string
	for _, d := range data.Datas {
		for _, event := range d.Events {
			text += fmt.Sprintf("%s %s %s\n", d.Date, event.Title, event.URL)
		}
	}

	body := &FeishuRequestBody{
		Timestamp: strconv.FormatInt(timestamp, 10),
		Sign:      sign,
		MsgType:   "text",
		Content: struct {
			Text string `json:"text"`
		}{Text: text},
	}

	var outData FeishuResponseData

	conn := ClientConn{
		*Client,
		"POST",
		"https://open.feishu.cn/open-apis/bot/v2/hook/0304a27d-8440-4253-b774-7da6f4eef447",
		map[string]string{"Content-Type": "application/json"},
		body,
		outData,
	}
	conn.Do()
}

func FeishuGenSign(secret string, timestamp int64) (string, error) {
	//timestamp + key 做sha256, 再进行base64 encode
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret

	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}

func (c *ClientConn) Do() any {
	fmt.Printf("%#v\n", c.RequestBody)

	bodyBytes, err := json.Marshal(c.RequestBody)
	if err != nil {
		panic(err)
	}
	buffer := bytes.NewBuffer(bodyBytes)
	req, err := http.NewRequest(c.Method, c.URL, buffer)
	if err != nil {
		panic(err)
	}

	for k, v := range c.Header {
		req.Header.Set(k, v)
	}

	resp, err := c.C.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(responseBody, &c.ResponseData); err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", c.ResponseData)

	return &c.ResponseData
}
