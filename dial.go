package main

import (
	"net/http"
	"strconv"
)

func dial(url, useragent, method string, verbose bool) LogMessage {
	var m LogMessage
	protocol := "https://"

	hconn, err := http.NewRequest(method, protocol+url, nil)
	if err != nil {
		if verbose {
			m.Environment = "debugging"
		}
		m.MessageType = "error"
		m.Target = url
		m.Message = "Connection refused by the source " + url + "with UserAgent " + useragent
		return m
	}
	hconn.Header.Set("User-Agent", useragent)

	resp, err := http.DefaultClient.Do(hconn)
	if err != nil {
		if verbose {
			m.Environment = "debugging"
		}
		m.MessageType = "error"
		m.Target = url
		m.Message = "Can not take a status code, maybe WAF is blocking the connect for the " + url +
			"with UserAgent " + useragent
		return m
	} else {
		m.MessageType = "regular"
		m.Target = url
		m.Status = resp.StatusCode
		m.Message = "HTTP Response Status: " + strconv.Itoa(resp.StatusCode) + " " + http.StatusText(resp.StatusCode)
		return m
	}
}
