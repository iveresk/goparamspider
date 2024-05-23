package main

import (
	"encoding/json"
	"math/rand/v2"
	"os"
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
		m.MessageType = "regular"
		m.Message = "There is no file " + filepath
		m.getLogger()
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

func getLiveParams(mode, url, jwt string, paramLevel int, delay time.Duration, verbose bool, payload Payloads) [][]LogMessage {
	var res [][]LogMessage
	var launchMethod Method

	for _, Mode := range payload.Mode {
		if mode == "day" {
			for _, day := range Mode.Day {
				for _, get := range day.Get {
					launchMethod = get
					res = append(res, intruder(url, jwt, "GET", paramLevel, delay, verbose, launchMethod))
				}
				for _, post := range day.Post {
					launchMethod = post
					res = append(res, intruder(url, jwt, "POST", paramLevel, delay, verbose, launchMethod))
				}
				for _, options := range day.Options {
					launchMethod = options
					res = append(res, intruder(url, jwt, "OPTIONS", paramLevel, delay, verbose, launchMethod))
				}
				for _, patch := range day.Patch {
					launchMethod = patch
					res = append(res, intruder(url, jwt, "PATCH", paramLevel, delay, verbose, launchMethod))
				}
			}
		} else {
			// Mode Night
			for _, day := range Mode.Night {
				for _, get := range day.Get {
					launchMethod = get
					res = append(res, intruder(url, jwt, "GET", paramLevel, delay, verbose, launchMethod))
				}
				for _, post := range day.Post {
					launchMethod = post
					res = append(res, intruder(url, jwt, "POST", paramLevel, delay, verbose, launchMethod))
				}
				for _, options := range day.Patch {
					launchMethod = options
					res = append(res, intruder(url, jwt, "PATCH", paramLevel, delay, verbose, launchMethod))
				}
				for _, patch := range day.Options {
					launchMethod = patch
					res = append(res, intruder(url, jwt, "OPTIONS", paramLevel, delay, verbose, launchMethod))
				}
			}
		}

	}
	return res
}

func intruder(url, jwt, method string, paramLevel int, delay time.Duration, verbose bool, payload Method) []LogMessage {
	var allLog LogMessage
	var params []string
	var fuzzParam []string
	var res []LogMessage

	protocol := "https://"
	url = protocol + url

	userAgent := getRandomUserAgent()

	for _, value := range payload.Routes {
		// forming a list of parameters according to a parameter level
		// highly not recommending something more than 2-3 levels, the HUGE # of requests
		params = nil
		for p := 0; p < paramLevel; p++ {
			for pi, param := range payload.Parameters {
				if paramLevel == 1 {
					params = append(params, "?"+param+"=FUZZ")
				} else {
					params[pi] = params[pi] + "&" + param + "=FUZZ"
				}

			}
		}
		time.Sleep(delay * time.Millisecond)
		// Checking default routes WITHOUT parameters
		allLog = dialHHTP(url+value, jwt, userAgent, method, verbose)
		if verbose {
			res = append(res, allLog)
		} else {
			if allLog.MessageType == "regular" {
				res = append(res, allLog)
			}
		}
		for p := 0; p < paramLevel; p++ {
			for i := 0; i < len(params); i++ {
				for _, payloads := range payload.Payloads {
					// Loading payloads into parameters
					fuzzParam = nil
					fuzzParam = strings.Split(params[i], "=")
					for j, _ := range fuzzParam {
						if fuzzParam[j] == "" || fuzzParam[j] == "FUZZ" || fuzzParam[j] == "=" {
							continue
						}
						params[i] = fuzzParam[j] + "=" + payloads
					}
					// Checking default routes WITH parameters
					time.Sleep(delay * time.Millisecond)
					allLog = dialHHTP(url+value+params[i], jwt, userAgent, method, verbose)
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
