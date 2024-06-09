package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
)

type LogMessage struct {
	MessageType string `json:"message_type"`
	Message     string `json:"message"`
	Header      string `json:"header,omitempty"`
	Environment string `json:"environment"`
	Target      string `json:"target"`
	Status      int    `json:"status"`
	Method      string `json:"method"`
	Body        string `json:"body"`
}

func getHelper(message string) {
	usage := "Flags:\n"
	usage += " -m Mode of the tool usage defining the if it is day (PR) mode or night (Full) scan.\n"
	usage += " -u Target Domain.\n"
	usage += " -t The authorization token for the white-box testing.\n"
	usage += " -H Set of key=value pairs, to set up headers for a request. Example: -H key1=value1,key2=value2.\n"
	usage += " -B Set of key=value pairs, to set up body data for a request. Example: -B key1=value1,key2=value2.\n"
	usage += " -d - The delay between requests not to be blocked by WAF. Default value is 1000ms\n"
	usage += " -p - The count of params to be tested combined in line.\n"
	usage += " - f - Flag to set output to the logging file /var/log/syslog.\n"
	usage += " - v - Flag to set verbose flag and record all debugging and rejected requests.\n"
	usage += " - s - Flag to set http or https connection mode.\n"
	usage += "Example:\n"
	usage += " ./goparamspider -u domain.com -m day -t <token> -H key1=value1,key2=value2 -f -s\n"
	fmt.Println(usage)
	log.Fatal("The error is " + message)
}

func (m *LogMessage) getLogger() {
	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	switch m.MessageType {
	case "helper":
		getHelper(m.Message)
	case "fatal":
		slogger.Error(m.Message, "body", m.Body, "method", m.Method, "target", m.Target, "header", m.Header, "env", m.Environment)
	case "regular":
		slogger.Info(m.Message, "body", m.Body, "method", m.Method, "target", m.Target, "header", m.Header, "env", m.Environment)
	case "error":
		if m.Environment == "debugging" {
			slogger.Warn(m.Message, "body", m.Body, "method", m.Method, "target", m.Target, "header", m.Header, "env", m.Environment)
		}
	default:
		getHelper(m.Message)
	}

}

func (m *LogMessage) appendToFile(url, logFilePath string) {
	filename := logFilePath + "/" + url + ".log"
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	flogger := slog.New(slog.NewJSONHandler(f, nil))
	// if something is wrong with the file just print to the console
	if err != nil {
		m.getLogger()
	}
	defer f.Close()

	switch m.MessageType {
	case "helper":
		getHelper(m.Message)
	case "fatal":
		flogger.Error(m.Message, "body", m.Body, "method", m.Method, "target", m.Target, "header", m.Header, "env", m.Environment)
	case "regular":
		flogger.Info(m.Message, "body", m.Body, "method", m.Method, "target", m.Target, "header", m.Header, "env", m.Environment)
	case "error":
		if m.Environment == "debugging" {
			flogger.Warn(m.Message, "body", m.Body, "method", m.Method, "target", m.Target, "header", m.Header, "env", m.Environment)
		}
	default:
		getHelper(m.Message)
	}
}
