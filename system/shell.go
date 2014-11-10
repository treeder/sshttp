package system

import (
	"log"
	"os/exec"
	"time"
)

type ShellStruct struct {
	Output string `json:"out"`
}

const timeout = 60

func Shell(s string) ShellStruct {
	log.Println("Executing command:", s)
	var sh ShellStruct
	c1 := make(chan string, 1)
	go func() {
		out, err := exec.Command("/bin/sh", "-c", s).CombinedOutput()
		if err != nil {
			c1 <- err.Error()
		} else {
			c1 <- string(out)
		}
	}()
	select {
	case res := <-c1:
		sh.Output = res
	case <-time.After(time.Second * timeout):
		sh.Output = "Command timed out after " + string(timeout) + " seconds"
	}
	return sh
}
