package system

import (
	"bytes"
	"log"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type Process struct {
	User string  `json:"user"`
	Name string  `json:"name"`
	Pid  int     `json:"pid"`
	Cpu  float64 `json:"cpu"`
}

type ProcessesStruct struct {
	UsedMem   float64   `json:"used_mem"`
	UsedCpu   float64   `json:"used_cpu"`
	Processes []Process `json:"processes"`
}

type ByCpu []Process

func (a ByCpu) Len() int           { return len(a) }
func (a ByCpu) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCpu) Less(i, j int) bool { return a[i].Cpu > a[j].Cpu }

func Processes() (plist ProcessesStruct) {
	cmd := exec.Command("ps", "aux")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	totalLoad := 0.0
	totalMem := 0.0
	for {
		line, err := out.ReadString('\n')
		if err != nil {
			break
		}
		tokens := strings.Split(line, " ")
		ft := make([]string, 0)
		for _, t := range tokens {
			if t != "" && t != "\t" {
				ft = append(ft, t)
			}
		}
		name := ft[10]
		user := ft[0]
		pid, err := strconv.Atoi(ft[1])
		if err != nil {
			continue
		}
		cpu, err := strconv.ParseFloat(ft[2], 64)

		totalLoad += cpu
		if err != nil {
			log.Fatal(err)
		}
		mem, err := strconv.ParseFloat(ft[3], 64)
		totalMem += mem
		if err != nil {
			log.Fatal(err)
		}
		plist.Processes = append(plist.Processes, Process{user, name, pid, cpu})
	}
	sort.Sort(ByCpu(plist.Processes))
	plist.UsedCpu = totalLoad
	plist.UsedMem = totalMem
	return plist
}
