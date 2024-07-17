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

type LeetCodeRequestBody struct {
	OperationName string `json:"operationName"`
	Query         string `json:"query"`
}

type LeetCodeResponseData struct {
	Data struct {
		TodayRecord []struct {
			Date       string `json:"date"`
			UserStatus any    `json:"userStatus"`
			Question   struct {
				QuestionID         string  `json:"questionId"`
				FrontendQuestionID string  `json:"frontendQuestionId"`
				Difficulty         string  `json:"difficulty"`
				Title              string  `json:"title"`
				TitleCn            string  `json:"titleCn"`
				TitleSlug          string  `json:"titleSlug"`
				PaidOnly           bool    `json:"paidOnly"`
				FreqBar            any     `json:"freqBar"`
				IsFavor            bool    `json:"isFavor"`
				AcRate             float64 `json:"acRate"`
				Status             any     `json:"status"`
				SolutionNum        int     `json:"solutionNum"`
				HasVideoSolution   bool    `json:"hasVideoSolution"`
				TopicTags          []struct {
					Name           string `json:"name"`
					NameTranslated string `json:"nameTranslated"`
					ID             string `json:"id"`
				} `json:"topicTags"`
				Extra struct {
					TopCompanyTags []struct {
						ImgURL        string `json:"imgUrl"`
						Slug          string `json:"slug"`
						NumSubscribed int    `json:"numSubscribed"`
					} `json:"topCompanyTags"`
				} `json:"extra"`
			} `json:"question"`
			LastSubmission any `json:"lastSubmission"`
		} `json:"todayRecord"`
	} `json:"data"`
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

type DoubanResponseData struct {
	Subjects []struct {
		EpisodesInfo string `json:"episodes_info"`
		Rate         string `json:"rate"`
		CoverX       int    `json:"cover_x"`
		Title        string `json:"title"`
		URL          string `json:"url"`
		Playable     bool   `json:"playable"`
		Cover        string `json:"cover"`
		ID           string `json:"id"`
		CoverY       int    `json:"cover_y"`
		IsNew        bool   `json:"is_new"`
	} `json:"subjects"`
}

var Client *http.Client

func init() {
	Client = &http.Client{}
}

func main() {
	FeishuNotify(JueJinEvent())
	FeishuNotify(DoubanMoive())
	FeishuNotify(LeetCodeDaily())
}

func DoubanMovie2() *DoubanResponseData {
// https://movie.douban.com/cinema/nowplaying/shanghai/
	return nil
}

func DoubanMoive() *DoubanResponseData {
	var outData *DoubanResponseData
	conn := ClientConn{
		*Client,
		"GET",
		"https://movie.douban.com/j/search_subjects?type=movie&tag=%E7%83%AD%E9%97%A8&page_limit=10&page_start=0",
		map[string]string{
			"User-Agent":"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36",
			"Host":"movie.douban.com"},
		nil,
		&outData,
	}
	conn.Do()
	
	return outData
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

func LeetCodeDaily() *LeetCodeResponseData {
	body := &LeetCodeRequestBody{
		"questionOfToday",
		"\n    query questionOfToday {\n  todayRecord {\n    date\n    userStatus\n    question {\n      questionId\n      frontendQuestionId: questionFrontendId\n      difficulty\n      title\n      titleCn: translatedTitle\n      titleSlug\n      paidOnly: isPaidOnly\n      freqBar\n      isFavor\n      acRate\n      status\n      solutionNum\n      hasVideoSolution\n      topicTags {\n        name\n        nameTranslated: translatedName\n        id\n      }\n      extra {\n        topCompanyTags {\n          imgUrl\n          slug\n          numSubscribed\n        }\n      }\n    }\n    lastSubmission {\n      id\n    }\n  }\n}\n ",
	}
	
	var outData *LeetCodeResponseData
	conn := ClientConn{
		*Client,
		"POST",
		"https://leetcode.cn/graphql/",
		map[string]string{
			"Content-Type":	"application/json",
			"Host":"leetcode.cn",
			"User-Agent":"Apifox/1.0.0 (https://apifox.com)",
		},
		body,
		&outData,
	}
	conn.Do()
	
	return outData
}

func FeishuNotify(inData any) {
	var text string
	
    switch data := inData.(type) {
    case *JueJinEventResponseData:
		for _, d := range data.Datas {
			for _, event := range d.Events {
				text += fmt.Sprintf("%s %s %s\n", d.Date, event.Title, event.URL)
			}
		}
	case *DoubanResponseData:
        for _, movie := range data.Subjects {
            text += fmt.Sprintf("%s %s %s %s %v\n",movie.Title ,movie.Rate ,movie.EpisodesInfo,movie.URL,movie.IsNew)
        }
    case *LeetCodeResponseData:
        text += fmt.Sprintf("%s %s %s",
			data.Data.TodayRecord[0].Question.QuestionID,
			data.Data.TodayRecord[0].Question.TitleCn,
			data.Data.TodayRecord[0].Question.Difficulty,
			)
        
	}
	

	secret := os.Getenv("SECRET")
	timestamp := time.Now().Unix()
	sign, err := FeishuGenSign(secret, timestamp)
	if err != nil {
		panic(err)
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

	var req *http.Request
	var resp *http.Response
	var buffer *bytes.Buffer
	var bodyBytes []byte
	var err error
	
	if c.RequestBody != nil {
		bodyBytes, err = json.Marshal(c.RequestBody)
		if err != nil {
			panic(err)
		}
	}
	buffer = bytes.NewBuffer(bodyBytes)
	req, err = http.NewRequest(c.Method, c.URL, buffer)
	if err != nil {
		panic(err)
	}
	
	for k, v := range c.Header {
		req.Header.Set(k, v)
	}
	
	resp, err = c.C.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(responseBody, &c.ResponseData); err != nil {
		fmt.Println(string(responseBody))
		panic(err)
	}

	fmt.Printf("%#v\n", c.ResponseData)

	return &c.ResponseData
}
