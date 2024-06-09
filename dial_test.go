package main

import (
	"testing"
)

func TestDialNoConnection(t *testing.T) {
	var test LogMessage
	expect := LogMessage{
		MessageType: "error",
		Message:     "Can not take a status code, maybe WAF is blocking the connect for the https:// with UserAgent  the Error is: Get \"https:\": http: no Host in request URL",
		Environment: "debugging",
		Target:      "https://",
		Method:      "GET",
		Status:      0,
	}
	test = dialHHTP("https://", "", "", "GET", true, nil, nil)
	if test != expect {
		t.Errorf("expect %#v, got %#v", expect, test)
	}
}

func TestDialWAF(t *testing.T) {
	var test LogMessage
	expect := LogMessage{
		MessageType: "error",
		Message: "Can not take a status code, maybe WAF is blocking the connect for the https://gogle.com with " +
			"UserAgent  the Error is: Get \"https://gogle.com\": tls: failed to verify certificate: x509: " +
			"certificate is valid for www.google.com, not gogle.com",
		Environment: "debugging",
		Target:      "https://gogle.com",
		Method:      "GET",
		Status:      0,
	}
	test = dialHHTP("https://gogle.com", "", "", "GET", true, nil, nil)
	if test != expect {
		t.Errorf("expect %#v, got %#v", expect, test)
	}
}
