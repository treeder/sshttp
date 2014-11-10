package system

import (
	"encoding/json"
)

type SystemStruct struct {
	Host string      `json:"host"`
	Disk interface{} `json:"disk"`
	Cpu  interface{} `json:"cpuinfo"`
	Load interface{} `json:"load"`
	Ram  interface{} `json:"ram"`
	Time string      `json:"time"`
}

func System(s string) []byte {
	m := SystemStruct{Host(), Disk(s), Cpuinfo(), Load(), Ram(), Now()}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return nil
	}
	return b
}
