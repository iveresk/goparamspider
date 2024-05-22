package main

import (
	"net/http"
	"strconv"
)

func dial(url, jwt, useragent, method string, verbose bool) LogMessage {
	var m LogMessage

	hconn, err := http.NewRequest(method, url, nil)
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
	if jwt != "" {
		hconn.Header.Set("Authorization", "Bearer "+jwt)
	}

	resp, err := http.DefaultClient.Do(hconn)
	if err != nil {
		if verbose {
			m.Environment = "debugging"
		}
		m.MessageType = "error"
		m.Target = url
		m.Message = "Can not take a status code, maybe WAF is blocking the connect for the " + url +
			"with UserAgent " + useragent + "the Error is: " + err.Error()
		return m
	} else {
		m.MessageType = "regular"
		m.Target = url
		m.Status = resp.StatusCode
		m.Message = "HTTP Response Status: " + strconv.Itoa(resp.StatusCode) + " " + http.StatusText(resp.StatusCode)
		return m
	}
}
