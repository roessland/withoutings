package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// curl --header "Authorization: Bearer YOUR_ACCESS_TOKEN" --data "action=getsummary&startdateymd=2020-07-01&enddateymd=2020-07-02" 'https://wbsapi.withings.net/v2/sleep '

func SleepGetSummary() {
	form := url.Values{}
	form.Add("action", "getsummary")
	form.Add("startdateymd", "2021-09-01")
	form.Add("enddateymd", "2021-10-01")
	form.Add("data_fields", "total_sleep_time")

	req, err := http.NewRequest(
		"POST",
		"https://wbsapi.withings.net/v2/sleep",
		nil)
	if err != nil {
		log.Println(err)
	}
	req.URL.RawQuery = form.Encode()
	req.Header.Add("Authorization", "Bearer " + token.AccessToken)
	fmt.Println(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(respBody))
	fmt.Println(resp.Header)
}

