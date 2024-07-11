package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	body := os.Getenv("body")
	cookie := os.Getenv("COOKIE")
	token := os.Getenv("TOKEN")

	client := &http.Client{}

	buffer := bytes.NewBuffer([]byte(body))
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
