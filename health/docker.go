package health

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const (
	inspectRunning = "inspect --format={{.State.Running}}"
	inspectStarted = "inspect --format={{.State.StartedAt}}"
	inspectNetwork = "inspect --format={{.NetworkSettings.IPAddress}}"
)

//Payload struct
type Payload struct {
	Arguments string
	Target    string
	Action    string
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
func (p *Payload) Parse() error {
	arguments := p.Arguments
	chunks := strings.Split(arguments, ":")

	if len(chunks) != 2 {
		return errors.New("Command formatted improperly. Expected command:arg1(:argN)")
	}

	p.Action = chunks[0]
	p.Target = chunks[1]

	return nil
}

//Exec docker container health check
func Exec(arguments string) (Response, error) {

	payload, err := parsePayload(arguments)

	if err != nil {
		return Response{}, err
	}

	//wat
	cmd := exec.Command("docker", inspectRunning, payload.Target)
	out, err := cmd.CombinedOutput()

	if err != nil {
		message := fmt.Sprintf("UNKNOWN - The container \"%s\" does not exist.", payload.Target)
		return Response{Success: false, StatusCode: 3, Message: message}, nil
	}

	s := byteToString(out)
	b, err := strconv.ParseBool(s)

	if b == false {
		message := fmt.Sprintf("CRITICAL - The container \"%s\" is not running.", payload.Target)
		return Response{Success: false, StatusCode: 2, Message: message}, nil
	}

	if b == true {
		cmd := exec.Command("docker", inspectStarted, payload.Target)
		out, err := cmd.CombinedOutput()

		if err != nil {
			errorMessage := fmt.Sprintf("Unknown error: %s", err)
			return Response{Success: false, StatusCode: 4, Message: errorMessage}, nil
		}

		network := byteToString(out)

		cmd = exec.Command("docker", inspectNetwork, payload.Target)
		out, err = cmd.CombinedOutput()

		if err != nil {
			errorMessage := fmt.Sprintf("Unknown error: %s", err)
			return Response{Success: false, StatusCode: 4, Message: errorMessage}, nil
		}

		started := byteToString(out)
		message := fmt.Sprintf("OK - The container \"%s\" is running. IP: %s, StartedAt: %s", payload.Target, started, network)
		return Response{Success: true, StatusCode: 0, Message: message}, nil
	}
	return Response{}, errors.New("Health check unhandled.")
}

func parsePayload(arguments string) (Payload, error) {
	arguments = strings.Trim(arguments, "\n")
	payload := Payload{Arguments: arguments}
	err := payload.Parse()
	return payload, err
}

func byteToString(b []byte) string {
	s := string(b[:])
	s = strings.Trim(s, "\n")
	return s
}
