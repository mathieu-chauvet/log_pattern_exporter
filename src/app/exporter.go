package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var logFilesToMonitor arrayFlags

func main() {
	flag.Var(&logFilesToMonitor, "lf", "log files to monitor")
	//filePath := "/home/mathieuchauvet/grok_exporter/webHook.log"
	pattern := flag.String("pattern", "ERROR", "pattern to search in files")
	outputFile := flag.String("output_file", "/var/tmp/log_file_metrics.prom", "destination folder for the result")
	flag.Parse()

	var arrayMetrics []string

	for _, logfile := range logFilesToMonitor {
		count, err := countOccurences(logfile, *pattern)

		if err != nil {
			log.Printf("failed to count in file %s: %s\n", logfile, err)
			promTxt := prometheusFormat(logfile, *pattern, -1)
			arrayMetrics = append(arrayMetrics, promTxt)
			continue
		}

		promTxt := prometheusFormat(logfile, *pattern, count)

		arrayMetrics = append(arrayMetrics, promTxt)
	}

	fmt.Println(arrayMetrics)
	writeToFile(arrayMetrics, *outputFile)

}

func writeToFile(metrics []string, outputFile string) {
	file, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, data := range metrics {
		_, _ = datawriter.WriteString(data + "\n")
	}

	datawriter.Flush()
	file.Close()
}

func countOccurences(logfile string, pattern string) (int, error) {
	file, err := os.Open(logfile)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer file.Close()
	count := 0
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), pattern) {
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return count, nil
}

func prometheusFormat(logfile string, pattern string, nbErrors int) string {
	res := fmt.Sprintf("pattern_in_log_count{logfile=\"%s\", pattern=\"%s\"} %d\n", logfile, pattern, nbErrors)
	return res
}
