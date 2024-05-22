package main

import (
	"flag"
)

func main() {
	var m LogMessage
	var payload Payloads
	// Checking all the params
	mode := flag.String("m", "day", "Mode of the tool usage defining the if it is day (PR) mode or night (Full) scan.")
	url := flag.String("u", "", "Target Domain.")
	jwt := flag.String("t", "", "The authorization token for the white-box testing.")
	delay := flag.Duration("d", 1000, "The delay in Milliseconds between requests not to be blocked by WAF.")
	paramLevel := flag.Int("l", 1, "The count of params to be tested combined in line.")
	output := flag.Bool("f", false, "Flag to set output to the logging file $TARGET.txt")
	verbose := flag.Bool("v", false, "Flag to set verbose flag and record all debugging and rejected requests.")
	flag.Parse()
	if *url == "" {
		m.Message = "No Valid Domain name was in parameters. Please check your line."
		m.MessageType = "helper"
		m.getLogger()
	}

	// Starting to load the payloads
	payload.readJSON("payloads")
	// Main algorithm runs here
	switch *mode {

	case "day", "night", "whitebox":
		params := getLiveParams(*mode, *url, *jwt, *paramLevel, *delay, *verbose, payload)
		for _, request := range params {
			for _, method := range request {
				if *output {
					method.appendToFile(*url)
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
