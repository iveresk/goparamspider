package main

import (
	"io"
	"net/http"
	"strconv"
)

func dialHHTP(url, jwt, useragent, method string, verbose bool) LogMessage {
	var m LogMessage

	httpconnection, err := http.NewRequest(method, url, nil)
	if err != nil {
		if verbose {
			m.Environment = "debugging"
		}
		m.MessageType = "error"
		m.Target = url
		m.Message = "Connection refused by the source " + url + "with UserAgent " + useragent
		return m
	}
	httpconnection.Header.Set("User-Agent", useragent)
	if jwt != "" {
		httpconnection.Header.Add("Authorization", "Bearer "+jwt)
	}

	resp, err := http.DefaultClient.Do(httpconnection)
	if err != nil {
		if verbose {
			m.Environment = "debugging"
		}
		m.MessageType = "error"
		m.Target = url
		m.Message = "Can not take a status code, maybe WAF is blocking the connect for the " + url +
			"with UserAgent " + useragent + "the Error is: " + err.Error()
		return m
	}
	m.MessageType = "regular"
	m.Target = url
	m.Status = resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		m.MessageType = "error"
		m.Message = "For the " + url + "HTTP Response Status: " +
			strconv.Itoa(resp.StatusCode) + " " +
			http.StatusText(resp.StatusCode) +
			". Can not read the response Body."
	} else {
		m.Message = "For the " + url + "HTTP Response Status: " +
			strconv.Itoa(resp.StatusCode) + " " +
			http.StatusText(resp.StatusCode) +
			". Response body is " + string(body)
	}
	return m
}
