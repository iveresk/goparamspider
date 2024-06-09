package main

import (
	"encoding/json"
	"errors"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

type Payloads struct {
	Mode []struct {
		Day []struct {
			Get []struct {
				Routes     []string `json:"routes"`
				Parameters []string `json:"parameters"`
				Payloads   []string `json:"payloads"`
			} `json:"GET"`
			Post []struct {
				Routes     []string `json:"routes"`
				Parameters []string `json:"parameters"`
				Payloads   []string `json:"payloads"`
			} `json:"POST"`
			Options []struct {
				Routes     []string `json:"routes"`
				Parameters []string `json:"parameters"`
				Payloads   []string `json:"payloads"`
			} `json:"OPTIONS"`
			Patch []struct {
				Routes     []string `json:"routes"`
				Parameters []string `json:"parameters"`
				Payloads   []string `json:"payloads"`
			} `json:"PATCH"`
		} `json:"day"`
		Night []struct {
			Get []struct {
				Routes     []string `json:"routes"`
				Parameters []string `json:"parameters"`
				Payloads   []string `json:"payloads"`
			} `json:"GET"`
			Post []struct {
				Routes     []string `json:"routes"`
				Parameters []string `json:"parameters"`
				Payloads   []string `json:"payloads"`
			} `json:"POST"`
			Options []struct {
				Routes     []string `json:"routes"`
				Parameters []string `json:"parameters"`
				Payloads   []string `json:"payloads"`
			} `json:"OPTIONS"`
			Patch []struct {
				Routes     []string `json:"routes"`
				Parameters []string `json:"parameters"`
				Payloads   []string `json:"payloads"`
			} `json:"PATCH"`
		} `json:"night"`
	} `json:"mode"`
}

type Method struct {
	Routes     []string `json:"routes"`
	Parameters []string `json:"parameters"`
	Payloads   []string `json:"payloads"`
}

func openFile(filepath string) []byte {
	var m LogMessage
	// reading the device-keywords.json file
	content, err := os.ReadFile(filepath)
	if err != nil {
		m.MessageType = "fatal"
		m.Message = "There is no file " + filepath
		m.getLogger()
		return nil
	}
	return content
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func (Payloads *Payloads) readJSON(assets string) {
	var m LogMessage
	filepath := "./assets/" + assets + ".json"
	content := openFile(filepath)
	err := json.Unmarshal(content, &Payloads)
	if err != nil || Payloads.Mode == nil {
		m.MessageType = "regular"
		m.Message = "Check if the file " + filepath + " is in json format"
		m.getLogger()
	}
}

func getRandomUserAgent() string {
	content := openFile("./assets/UserAgent.txt")
	res := strings.Split(string(content[:]), "\n")
	randomid := randRange(1, len(res))
	return res[randomid-1]
}

func getFUZZ(paramlevel int, payload []string) []string {
	var res []string
	for p := 0; p < paramlevel; p++ {
		for pi, param := range payload {
			if p < 1 {
				res = append(res, "?"+param+"=FUZZ")
			} else {
				for pj, _ := range payload {
					splitted := strings.Split(res[pj], "?")
					if len(splitted) > 1 {
						for s := 1; s < len(splitted); s += 2 {
							res = append(res, res[pi]+"&"+splitted[s])
						}
					} else {
						res = append(res, res[pi]+"&"+res[pj])
					}
				}
			}
		}
	}
	return res
}

func replaceFUZZ(paramlevel int, params, payloads []string) []string {
	var res []string
	var fuzzParam []string
	for p := 0; p < paramlevel; p++ {
		for _, param := range params {
			for _, payload := range payloads {
				fuzzParam = nil
				fuzzParam = strings.Split(param, "=")
				for j, _ := range fuzzParam {
					if fuzzParam[j] == "" || fuzzParam[j] == "FUZZ" || fuzzParam[j] == "=" {
						continue
					}
					res = append(res, fuzzParam[j]+"="+payload)
				}
			}
		}
	}
	return res
}

func checkBodyFuzz(body map[string]string) bool {
	for _, v := range body {
		if strings.Contains(v, "FUZZ") {
			return true
		}
	}
	return false
}

func bodyFuzz(body map[string]string, payloads []string) map[string]string {
	var res = make(map[string]string)
	for k, v := range body {
		if strings.Contains(v, "FUZZ") {
			for i, payload := range payloads {
				res[k+strconv.Itoa(i)] = strings.Replace(v, "FUZZ", payload, 1)
			}
		}
	}
	return res
}

func logBody(body map[string]string) string {
	res := ""
	for k, v := range body {
		res = res + k + ":" + v + ","
	}
	return res
}

func parseHeaders(headers string) (map[string]string, error) {
	var (
		res        = make(map[string]string)
		splitcomma []string
		splitequal []string
	)
	splitcomma = strings.Split(headers, ",")
	if len(splitcomma) < 1 {
		return nil, errors.New("Invalid header format")
	}
	for _, v := range splitcomma {
		splitequal = strings.Split(v, "=")
		if len(splitequal) != 2 {
			return nil, errors.New("Invalid header format")
		}
		res[splitequal[0]] = splitequal[1]
	}
	return res, nil
}

func makeAttack(mode, url, jwt string, paramLevel int, delay time.Duration, verbose, ssl bool, payload Payloads, headers, body map[string]string) [][]LogMessage {
	var (
		res          [][]LogMessage
		launchMethod Method
	)

	for _, Mode := range payload.Mode {
		if mode == "day" {
			for _, day := range Mode.Day {
				for _, get := range day.Get {
					launchMethod = get
					res = append(res, intruder(url, jwt, "GET", paramLevel, delay, verbose, ssl, launchMethod, headers, body))
				}
				for _, post := range day.Post {
					launchMethod = post
					res = append(res, intruder(url, jwt, "POST", paramLevel, delay, verbose, ssl, launchMethod, headers, body))
				}
				for _, options := range day.Options {
					launchMethod = options
					res = append(res, intruder(url, jwt, "OPTIONS", paramLevel, delay, verbose, ssl, launchMethod, headers, body))
				}
				for _, patch := range day.Patch {
					launchMethod = patch
					res = append(res, intruder(url, jwt, "PATCH", paramLevel, delay, verbose, ssl, launchMethod, headers, body))
				}
				return res
			}
		} else {
			// Mode Night
			for _, day := range Mode.Night {
				for _, get := range day.Get {
					launchMethod = get
					res = append(res, intruder(url, jwt, "GET", paramLevel, delay, verbose, ssl, launchMethod, headers, body))
				}
				for _, post := range day.Post {
					launchMethod = post
					res = append(res, intruder(url, jwt, "POST", paramLevel, delay, verbose, ssl, launchMethod, headers, body))
				}
				for _, options := range day.Patch {
					launchMethod = options
					res = append(res, intruder(url, jwt, "PATCH", paramLevel, delay, verbose, ssl, launchMethod, headers, body))
				}
				for _, patch := range day.Options {
					launchMethod = patch
					res = append(res, intruder(url, jwt, "OPTIONS", paramLevel, delay, verbose, ssl, launchMethod, headers, body))
				}
				return res
			}
		}

	}
	return res
}

func intruder(url, jwt, method string, paramLevel int, delay time.Duration, verbose, ssl bool, payload Method, headers, body map[string]string) []LogMessage {
	var (
		allLog          LogMessage
		params, fuzzeds []string
		res             []LogMessage
		fuzzbody        map[string]string
	)

	if ssl {
		// changing the default connection protocol if needed
		url = "https://" + url
	} else {
		url = "http://" + url
	}
	// Taking random User-Agent
	userAgent := getRandomUserAgent()
	// forming a list of parameters according to a parameter level
	// highly not recommending something more than 2 levels, the HUGE # of requests
	params = getFUZZ(paramLevel, payload.Parameters)
	fuzzeds = replaceFUZZ(paramLevel, params, payload.Payloads)
	// Time to FUZZ body if it is in param
	if body != nil && checkBodyFuzz(body) {
		fuzzbody = bodyFuzz(body, payload.Payloads)
	}

	for _, route := range payload.Routes {
		time.Sleep(delay * time.Millisecond)
		// Checking default routes WITHOUT parameters
		if fuzzbody == nil {
			allLog = dialHHTP(url+route, jwt, userAgent, method, verbose, headers, body)
			if verbose {
				res = append(res, allLog)
			} else {
				if allLog.MessageType == "regular" {
					res = append(res, allLog)
				}
			}
		} else {
			for k, v := range fuzzbody {
				for k2, _ := range body {
					if strings.Contains(k, k2) {
						body[k2] = v
					}
				}
				allLog = dialHHTP(url+route, jwt, userAgent, method, verbose, headers, body)
				if verbose {
					res = append(res, allLog)
				} else {
					if allLog.MessageType == "regular" {
						res = append(res, allLog)
					}
				}
			}
		}

		// Testing requests with parameters
		for _, fuzzed := range fuzzeds {
			// Checking default routes WITH parameters
			time.Sleep(delay * time.Millisecond)
			if fuzzbody == nil {
				allLog = dialHHTP(url+route+fuzzed, jwt, userAgent, method, verbose, headers, body)
				if verbose {
					res = append(res, allLog)
				} else {
					if allLog.MessageType == "regular" {
						res = append(res, allLog)
					}
				}
			} else {
				for k, v := range fuzzbody {
					for k2, _ := range body {
						if strings.Contains(k, k2) {
							body[k2] = v
						}
					}
					allLog = dialHHTP(url+route+fuzzed, jwt, userAgent, method, verbose, headers, body)
					if verbose {
						res = append(res, allLog)
					} else {
						if allLog.MessageType == "regular" {
							res = append(res, allLog)
						}
					}
				}
			}
		}
	}
	return res
}
