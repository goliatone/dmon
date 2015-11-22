package health

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

//Payload struct
type Payload struct {
	arguments   string
	containerID string
	action      string
}

//Parse payload
func (p *Payload) Parse() {
	arguments := p.arguments
	chunks := strings.Split(arguments, ":")
	p.action = chunks[0]
	p.containerID = chunks[1]
}

//Exec docker container health check
func Exec(arguments string) (bool, int, string) {

	payload := parsePayload(arguments)

	//wat
	cmd := exec.Command("docker", "inspect", "--format={{.State.Running}}", payload.containerID)
	out, err := cmd.CombinedOutput()

	if err != nil {
		message := fmt.Sprintf("UNKNOWN - The container \"%s\" does not exist.", payload.containerID)
		return false, 3, message
	}

	s := byteToString(out)
	b, err := strconv.ParseBool(s)

	if b == false {
		message := fmt.Sprintf("CRITICAL - The container \"%s\" is not running.", payload.containerID)
		return false, 2, message
	}

	if b == true {
		cmd := exec.Command("docker", "inspect", "--format={{.State.StartedAt}}", payload.containerID)
		out, err := cmd.CombinedOutput()

		if err != nil {
			errorMessage := fmt.Sprintf("Unknown error: %s", err)
			return false, 4, errorMessage
		}

		network := byteToString(out)

		cmd = exec.Command("docker", "inspect", "--format={{.NetworkSettings.IPAddress}}", payload.containerID)
		out, err = cmd.CombinedOutput()

		if err != nil {
			errorMessage := fmt.Sprintf("Unknown error: %s", err)
			return false, 4, errorMessage
		}

		started := byteToString(out)
		message := fmt.Sprintf("OK - The container \"%s\" is running. IP: %s, StartedAt: %s", payload.containerID, started, network)

		return true, 0, message
	}

	return false, 5, "Unknown error"
}

func parsePayload(arguments string) Payload {
	arguments = strings.Trim(arguments, "\n")
	payload := Payload{arguments: arguments}
	payload.Parse()
	return payload
}

func byteToString(b []byte) string {
	s := string(b[:])
	s = strings.Trim(s, "\n")
	return s
}
