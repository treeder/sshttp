package system

import (
	"fmt"
	"os"
)

func Host() string {
	host, err := os.Hostname()
	if err != nil {
		return fmt.Sprint(err)
	}
	return host
}
