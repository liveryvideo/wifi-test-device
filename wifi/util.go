package wifi

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Log struct {
	Value string
	Time  string
}

type Logs struct {
	Total     int
	Remaining int
	Logs      []Log
}

func init() {
	logFile, err := os.Create("out.log")
	if err != nil {
		log.Println("Could not create log file:", err)
	}
	log.SetPrefix(": ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)
	log.SetOutput(logFile)
}

func performCommand(command string, args ...string) ([]byte, error) {
	log.Println("Performing command:", command, args)
	cmd := exec.Command(command, args...)
	return cmd.CombinedOutput()
}

func printOutput(output []byte, err error) {
	raw := string(output)
	if len(strings.Trim(raw, " ")) == 0 {
		return
	}
	log.Println(raw)
	if err != nil {
		log.Println("Error: ", err)
	}
}

func FetchLogs(start int, end int) Logs {
	logFile, err := os.Open("out.log")

	data, err := ioutil.ReadAll(logFile)

	if err != nil {
		log.Println("Could not open log file:", err)
	}

	raw := string(data)
	if strings.LastIndex(raw, "\n") == len(raw)-1 {
		raw = raw[0 : len(raw)-1]
	}

	lines := strings.Split(raw, "\n")

	if start < 0 {
		start = 0
	}

	if end < 1 {
		end = len(lines)
	}

	if start > len(lines) {
		start = len(lines)
	}

	if end > len(lines) {
		end = len(lines)
	}

	total := len(lines)

	lines = lines[start:end]

	remaining := total - end

	logs := make([]Log, len(lines))

	for i, _ := range logs {
		args := strings.SplitN(lines[i], ": ", 2)

		var message string
		var time string
		if len(args) > 1 {
			time = args[0]
			message = args[1]
		} else {
			message = args[0]
		}

		if len(message) == 0 && len(time) == 0 {
			continue
		}

		logs[i] = Log{
			Value: message,
			Time:  time,
		}
	}

	logsContainer := Logs{
		Total:     total,
		Remaining: remaining,
		Logs:      logs,
	}

	return logsContainer
}
