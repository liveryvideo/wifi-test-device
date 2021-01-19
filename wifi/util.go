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

func FetchLogs() []Log {
	logFile, err := os.Open("out.log")
	data, err := ioutil.ReadAll(logFile)
	if err != nil {
		log.Println("Could not open log file:", err)
	}

	lines := strings.Split(string(data), "\n")
	logs := make([]Log, len(lines)-1)

	for i := range logs {
		args := strings.SplitN(lines[i], ": ", 2)

		var message string
		var time string
		if len(args) > 1 {
			time = args[0]
			message = args[1]
		} else {
			message = args[0]
		}

		logs[i] = Log{
			Value: message,
			Time:  time,
		}
	}

	return logs
}
