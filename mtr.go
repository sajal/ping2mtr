//package ping2mtr attempts to reproduce output of what mtr -n --raw command would have produced
//using only the ping command. This is to be run on linux machines where it is not possible to
//install and run mtr.
package ping2mtr

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func sendpings(dest string, ttl int, result chan []string) {
	out := ping(dest, strconv.Itoa(10), "255")
	//log.Println(out)
	timings := make([]int, 0)
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "time=") {
			t := strings.Split(line, "time=")[1]
			ms := strings.Split(t, " ")[0]
			dur, _ := time.ParseDuration(ms + "ms")
			timings = append(timings, int(dur.Nanoseconds()/1000))
		}
	}
	output := make([]string, len(timings))
	for i, timing := range timings {
		output[i] = fmt.Sprintf("p %d %d", ttl, timing)
	}
	result <- output
}

func ping(dest, count, ttl string) string {
	//ping -t 3 -c 1 -W 1 8.8.8.8
	var cmd *exec.Cmd
	var out bytes.Buffer
	//var stderr bytes.Buffer
	//cmd.Stderr = &stderr
	cmd = exec.Command("ping", "-n", "-c", count, "-t", ttl, "-W", "1", dest)
	cmd.Stdout = &out
	cmd.Run()
	return out.String()
}

func detecthop(dest string, ttl int) string {
	result := ping(dest, "1", strconv.Itoa(ttl))
	//line := ""
	for _, line := range strings.Split(result, "\n") {
		if strings.Contains(line, "icmp_seq") {
			if strings.Contains(line, "From") {
				return strings.Split(line, " ")[1]
			} else {
				item := strings.Split(line, " ")[3]
				return strings.Split(item, ":")[0]
			}
		}
	}
	return ""
}

func discoverhops(dest string) []string {
	hops := make([]string, 0)
	unresponsive := 0
	pingchan := make(chan []string)
	for i := 1; i < 31; i++ {
		hop := detecthop(dest, i)
		log.Println(i, hop)
		if hop != "" {
			hops = append(hops, hop)
			go sendpings(hop, i, pingchan)
		}
		if hop == "" {
			unresponsive++
		} else if hop == dest {
			break
		} else {
			unresponsive = 0
		}
		if unresponsive == 3 {
			break
		}
	}
	for i, hop := range hops {
		hops[i] = fmt.Sprintf("h %d %s", i, hop)
		o := <-pingchan
		hops = append(hops, o...)
	}
	return hops
}

//dest must be an IP
//The result is what mtr -n --raw would produce
func Ping2MTR(dest string) string {
	hops := discoverhops(dest)
	return strings.Join(hops, "\n")
}
