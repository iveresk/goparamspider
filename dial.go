package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func dialHHTP(url, jwt, useragent, method string, verbose bool, headers, dbody map[string]string) LogMessage {
	var (
		m              LogMessage
		httpconnection *http.Request
		err            error
	)
	marshaled, _ := json.Marshal(dbody)
	if dbody == nil {
		httpconnection, err = http.NewRequest(method, url, nil)
	} else {
		if method == "GET" {
			httpconnection, err = http.NewRequest(method, url, nil)
		} else {
			httpconnection, err = http.NewRequest(method, url, bytes.NewBuffer(marshaled))
		}
	}

	if err != nil {
		if verbose {
			m.Environment = "debugging"
		}
		m.MessageType = "error"
		m.Method = method
		m.Target = url
		m.Body = logBody(dbody)
		m.Message = "Invalid request " + url + " with UserAgent " + useragent
		return m
	}

	httpconnection.Header.Set("User-Agent", useragent)
	httpconnection.Header.Add("Accept", "application/json")
	httpconnection.Header.Add("Content-Type", "application/json")

	if jwt != "" {
		httpconnection.Header.Add("Authorization", "Bearer "+jwt)
	}
	// Adding headers to the request
	if headers != nil {
		for k, v := range headers {
			httpconnection.Header.Add(k, v)
		}
	}

	resp, err := http.DefaultClient.Do(httpconnection)
	if err != nil {
		if verbose {
			m.Environment = "debugging"
		}
		m.Method = method
		m.MessageType = "error"
		m.Target = url
		m.Body = logBody(dbody)
		m.Message = "Can not take a status code, maybe WAF is blocking the connect for the " + url +
			" with UserAgent " + useragent + " the Error is: " + err.Error()
		if resp != nil {
			if reqHeadersBytes, err := json.Marshal(resp.Header); err != nil {
				m.Header = "Can not take header fot the target " + url
			} else {
				m.Header = string(reqHeadersBytes)
			}
		}
		return m
	}
	m.MessageType = "regular"
	m.Target = url
	m.Method = method
	m.Body = logBody(dbody)
	m.Status = resp.StatusCode
	// Taking the Body of response
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		m.MessageType = "error"
		m.Message = "For the " + url + " HTTP Response Status: " +
			strconv.Itoa(resp.StatusCode) + " " +
			http.StatusText(resp.StatusCode) +
			". Can not read the response Body."
	} else {
		m.Message = "For the " + url + " HTTP Response Status: " +
			strconv.Itoa(resp.StatusCode) + " " +
			http.StatusText(resp.StatusCode) +
			". Response body is " + string(body)
	}
	// Taking the Header of response
	if reqHeadersBytes, err := json.Marshal(resp.Header); err != nil {
		m.Header = "Can not take header fot the target " + url
	} else {
		m.Header = string(reqHeadersBytes)
	}
	return m
}
