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

//Response struct
type Response struct {
	StatusCode int
	Message    string
	Success    bool
}

//Parse payload
//TODO: maybe move this out, check first component
//for actual type check and then the rest would arguments
//in array form. So we could call this Exec from argv or
//whatever.
func (p *Payload) Parse() {
	arguments := p.arguments
	chunks := strings.Split(arguments, ":")
	p.action = chunks[0]
	p.containerID = chunks[1]
}

//Exec docker container health check
func Exec(arguments string) Response {

	payload := parsePayload(arguments)

	//wat
	cmd := exec.Command("docker", "inspect", "--format={{.State.Running}}", payload.containerID)
	out, err := cmd.CombinedOutput()

	if err != nil {
		message := fmt.Sprintf("UNKNOWN - The container \"%s\" does not exist.", payload.containerID)
		return Response{Success: false, StatusCode: 3, Message: message}
	}

	s := byteToString(out)
	b, err := strconv.ParseBool(s)

	if b == false {
		message := fmt.Sprintf("CRITICAL - The container \"%s\" is not running.", payload.containerID)
		return Response{Success: false, StatusCode: 2, Message: message}
	}

	if b == true {
		cmd := exec.Command("docker", "inspect", "--format={{.State.StartedAt}}", payload.containerID)
		out, err := cmd.CombinedOutput()

		if err != nil {
			errorMessage := fmt.Sprintf("Unknown error: %s", err)
			return Response{Success: false, StatusCode: 4, Message: errorMessage}
		}

		network := byteToString(out)

		cmd = exec.Command("docker", "inspect", "--format={{.NetworkSettings.IPAddress}}", payload.containerID)
		out, err = cmd.CombinedOutput()

		if err != nil {
			errorMessage := fmt.Sprintf("Unknown error: %s", err)
			return Response{Success: false, StatusCode: 4, Message: errorMessage}
		}

		started := byteToString(out)
		message := fmt.Sprintf("OK - The container \"%s\" is running. IP: %s, StartedAt: %s", payload.containerID, started, network)
		return Response{Success: true, StatusCode: 0, Message: message}
	}
	return Response{Success: false, StatusCode: 5, Message: "Unknown error"}
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
