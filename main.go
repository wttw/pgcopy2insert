package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	copyRe := regexp.MustCompile(`^(?i)\s*copy\s+("[^"]+"|\S+)\s+\(([^)]+)\)\s+from\s+stdin`)
	scanner := bufio.NewScanner(os.Stdin)
	var copying = false
	var first = false
	for scanner.Scan() {
		line := scanner.Text()
		if copying {
			if strings.HasPrefix(line, `\.`) {
				copying = false
				fmt.Printf(";\n")
				continue
			}
			fields := strings.Split(line, "\t")
			for i, v := range fields {
				if v == `\N` {
					fields[i] = "NULL"
				} else {
					field := "'" + strings.ReplaceAll(v, "'", "''") + "'"
					if strings.Contains(field, `\`) {
						field = "E" + field
					}
					fields[i] = field
				}
			}
			if !first {
				fmt.Printf(",")
			}
			fmt.Printf("\n  (%s)", strings.Join(fields, ", "))
			first = false;
			continue
		}
		matches := copyRe.FindStringSubmatch(line)
		if matches != nil {
			copying = true
			first = true
			fmt.Printf("INSERT INTO %s (%s) VALUES", matches[1], matches[2])
			continue
		}
		fmt.Println(line)
	}
}
