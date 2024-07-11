package main

import (
    "bytes"
    "encoding/json"
    "fmt"
	"io"
	"net/http"
	"os"
    "strconv"
)

type RequestBody struct {
	FromDate string
	Days    int64
	Aid     int64
}

func main() {
	fromDate:=os.Getenv("FromDate")
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
    body := &RequestBody{
		FromDate: 	fromDate,
		Days: days,
		Aid: aid,
	}
	
	fmt.Println(fromDate,cookie,token,days,aid)
	
	client := &http.Client{}

    bodyBytes, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

    buffer := bytes.NewBuffer(bodyBytes)
	req, err := http.NewRequest("POST",
		"https://api.juejin.cn/event_api/v1/event/month_stat",
		buffer)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("cookie", cookie)
	req.Header.Set("x-secsdk-csrf-token", token)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(responseBody))
}
