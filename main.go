package main

import (
	"flag"
)

func main() {
	var (
		m       LogMessage
		payload Payloads
		headers map[string]string
	)
	// Checking all the params
	mode := flag.String("m", "day", "Mode of the tool usage defining the if it is day (PR) mode or night (Full) scan.")
	url := flag.String("u", "", "Target Domain.")
	log := flag.String("l", "./assets/", "Log filename path. The default value is ./assets/")
	jwt := flag.String("t", "", "The authorization token for the white-box testing.")
	header := flag.String("H", "", "Set of key=value pairs, to set up headers for a request")
	delay := flag.Duration("d", 1000, "The delay in Milliseconds between requests not to be blocked by WAF.")
	paramLevel := flag.Int("p", 1, "The count of params to be tested combined in line.")
	output := flag.Bool("f", false, "Flag to set output to the logging file $TARGET.txt")
	verbose := flag.Bool("v", false, "Flag to set verbose flag and record all debugging and rejected requests.")
	ssl := flag.Bool("s", false, "Flag to set http or https connection mode")
	// Parsing the flags
	flag.Parse()
	if *url == "" {
		m.Message = "No Valid Domain name was in parameters. Please check your line."
		m.MessageType = "helper"
		m.getLogger()
	}

	//Parsing Header parameters
	headers, err := parseHeaders(*header)
	if err != nil {
		m.Message = "Please check your Headers parameters line."
		m.MessageType = "helper"
		m.getLogger()
	}

	// Starting to load the payloads
	payload.readJSON("payloads")
	// Main algorithm runs here
	switch *mode {
	// Day mode is for quick check before PR, night mode is for full scale check
	case "day", "night":
		params := getLiveParams(*mode, *url, *jwt, *paramLevel, *delay, *verbose, *ssl, payload, headers)
		for _, request := range params {
			for _, method := range request {
				if *output {
					method.appendToFile(*url, *log)
				}
				method.getLogger()
			}
		}
	default:
		{
			m.Message = "No Valid Mode was specified in parameters. Please check your line."
			m.MessageType = "helper"
			m.getLogger()
		}
	}
}
