// This program takes the structured log output and makes it readable.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-json-experiment/json"
)

var service string

func init() {
	flag.StringVar(&service, "service", "", "filter which service to see")
}

func main() {
	flag.Parse()
	var b strings.Builder

	service := strings.ToLower(service)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()

		m := make(map[string]any)
		err := json.Unmarshal([]byte(s), &m)
		if err != nil {
			if service == "" {
				fmt.Println(s)
			}
			continue
		}
		//If a service filter was provided
		if service != "" && strings.ToLower(m["service"].(string)) != service {
			continue
		}
		traceID := "00000000-0000-0000-0000-000000000000"
		if v, ok := m["trace_id"]; ok && v != "" {
			traceID = fmt.Sprintf("%v", v)
		}
		//Build out the know portions of the log in the order
		b.Reset()
		b.WriteString(fmt.Sprintf("%s: %s: %s: %s: %s: %s: ",
			m["service"],
			m["time"],
			m["file"],
			m["level"],
			traceID,
			m["msg"],
		))

		//Add the rest of the keys ignoring the ones we already have in the logs
		for k, v := range m {
			switch k {
			case "service", "time", "file", "level", "trace_id", "msg":
				continue
			}
			b.WriteString(fmt.Sprintf("%s[%v]: ", k, v))
		}

		//Write the new log format, removing the last :
		out := b.String()
		fmt.Println(out[:len(out)-2])
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
