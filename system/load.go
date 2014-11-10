package system

import (
	"io/ioutil"
)

type LoadStruct struct {
	Avg1 float64 `json:"avg1"`
	Avg2 float64 `json:"avg2"`
	Avg3 float64 `json:"avg3"`
}

func Load() interface{} {
	b, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return map[string]string{"msg": "Unsupported"}
	}
	return LoadStruct{toFloat(b[0:4]), toFloat(b[5:9]), toFloat(b[10:14])}
}
