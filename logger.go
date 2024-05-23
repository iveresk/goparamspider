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
	Environment string `json:"environment"`
	Target      string `json:"target"`
	Status      int    `json:"status"`
}

func getHelper(message string) {
	usage := "Flags:\n"
	usage += " -m Mode of the tool usage defining the if it is day (PR) mode or night (Full) scan.\n"
	usage += " -u Target Domain.\n"
	usage += " -t The authorization token for the white-box testing.\n"
	usage += " -d - The delay between requests not to be blocked by WAF. Default value is 1000ms\n"
	usage += " -p - The count of params to be tested combined in line.\n"
	usage += " - f - Flag to set output to the logging file /var/log/syslog.\n"
	usage += " - v - Flag to set verbose flag and record all debugging and rejected requests.\n"
	usage += "Example:\n"
	usage += " ./goparamspider -m day -u domain.com \n"
	fmt.Println(usage)
	log.Fatal("The error is " + message)
}

func (m *LogMessage) getLogger() {
	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	switch m.MessageType {
	case "helper":
		getHelper(m.Message)
	case "fatal":
		slogger.Error(m.Target, m.Environment, m.Message)
	case "regular":
		slogger.Info(m.Target, m.Environment, m.Message)
	case "error":
		if m.Environment == "debugging" {
			slogger.Warn(m.Target, m.Environment, m.Message)
		}
	default:
		getHelper(m.Message)
	}

}

func (m *LogMessage) appendToFile(url, logFilePath string) {
	filename := logFilePath + "/" + url + ".log"
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	// if something is wrong with the file just print to the console
	if err != nil {
		m.getLogger()
	}

	defer f.Close()

	if _, err = f.WriteString(m.Message + "\n"); err != nil {
		m.getLogger()
	}
}
