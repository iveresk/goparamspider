package main

import (
	"errors"
	"strings"
	"testing"
)

func TestGetFuzz(t *testing.T) {
	parameters := []string{"v", "query"}
	expect := []string{"?v=FUZZ", "?query=FUZZ"}
	expect2 := []string{"?v=FUZZ", "?query=FUZZ", "?v=FUZZ&v=FUZZ", "?v=FUZZ&query=FUZZ", "?query=FUZZ&v=FUZZ", "?query=FUZZ&query=FUZZ"}
	var test []string
	test = getFUZZ(1, parameters)
	if len(test) != len(expect) {
		t.Errorf("expect %#v, got %#v", expect, test)
	}
	for i, v := range test {
		if v != expect[i] {
			t.Errorf("expect %#v, got %#v", expect, test)
		}
	}
	test = getFUZZ(2, parameters)
	if len(test) != len(expect2) {
		t.Errorf("expect %#v, got %#v", expect2, test)
	}
}

func TestParseHeaders(t *testing.T) {
	header := "key1=value1,key2=value2"
	expect := map[string]string{"key1": "value1", "key2": "value2"}
	test, _ := parseHeaders(header)
	for key, value := range test {
		if value != expect[key] {
			t.Errorf("expect %#v, got %#v", expect, test)
		}
	}
}

func TestReplaceFUZZ(t *testing.T) {
	parameters := []string{"v=FUZZ", "query=FUZZ"}
	payloads := []string{"%3D%3D", "%3D", "%27"}
	expect := []string{"v=%3D%3D", "v=%3D", "v=%27", "query=%3D%3D", "query=%3D", "query=%27"}
	test := replaceFUZZ(1, parameters, payloads)
	if len(test) != len(expect) {
		t.Errorf("expect %#v, got %#v", payloads, test)
	}
	for i, v := range test {
		if v != expect[i] {
			t.Errorf("expect %#v, got %#v", expect, test)
		}
	}
}

func TestParseHeadersNil(t *testing.T) {
	header := ""
	expect := map[string]string{"key1": "value1", "key2": "value2"}
	_, err := parseHeaders(header)
	if err == nil || err.Error() != "Invalid header format" {
		t.Errorf("expect %#v, got an error %#v", expect, err.Error())
	}
}

func TestGetRandomUserAgent(t *testing.T) {
	expect := openFile("./assets/UserAgent.txt")
	res := strings.Split(string(expect[:]), "\n")
	test := getRandomUserAgent()
	flag := false
	for i := 0; i < len(res); i++ {
		if test == res[i] {
			flag = true
		}
	}
	if !flag {
		t.Errorf("expected User Agent %v, not in the file", test)
	}
}

func TestRandRange(t *testing.T) {
	m1, m2 := 0, 10
	test := randRange(m1, m2)
	flag := false
	for i := m1; i < m2; i++ {
		if test == i {
			flag = true
		}
	}
	if !flag {
		t.Errorf("expect %d, not in range min = %d, max = %d", test, m1, m2)
	}
}

func TestOpenFileNil(t *testing.T) {
	expect := errors.New("open nil")
	test := openFile("somefile")
	if test != nil {
		t.Errorf("expect %#v, got %#v", expect, test)
	}
}

func TestOpenFile(t *testing.T) {
	test := openFile("./assets/UserAgent.txt")
	if test == nil {
		t.Errorf("expect %#v, got %#v", "nil", test)
	}
}
