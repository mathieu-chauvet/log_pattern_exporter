package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
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
	configFile := flag.String("parameter_file", "/etc/log_pattern_exporter.conf", "Optional file containing the list of log files to parse")
	substringLimit := flag.Int("substring_limit", -1, "search pattern only in the first X caracters of the line")
	flag.Parse()

	var arrayMetrics []string
	var filesToParse []string

	filesToParse = addFilesFromFlags(filesToParse)

	filesToParse = addFilesFromConfigFile(filesToParse, *configFile)

	arrayMetrics = searchPatternInFiles(filesToParse, pattern, arrayMetrics, *substringLimit)

	fmt.Println(arrayMetrics)
	writeToFile(arrayMetrics, *outputFile)

}

func searchPatternInFiles(filesToParse []string, pattern *string, arrayMetrics []string, substringLimit int) []string {
	// Add helpers
	arrayMetrics = append(arrayMetrics, prometheusHelpers(pattern))

	for _, logfile := range filesToParse {
		count, err := countOccurences(logfile, *pattern, substringLimit)

		if logfile == "" {
			continue
		}

		if err != nil {
			log.Printf("failed to count in file %s: %s\n", logfile, err)
			promTxt := prometheusFormat(logfile, *pattern, -1)
			arrayMetrics = append(arrayMetrics, promTxt)
			continue
		}

		promTxt := prometheusFormat(logfile, *pattern, count)

		arrayMetrics = append(arrayMetrics, promTxt)
	}
	return arrayMetrics
}

func prometheusHelpers(pattern *string) string {
	return "# HELP pattern_in_log_count The total number of occurences of " + *pattern + ".\n# TYPE pattern_in_log_count counter" + "\n"
}

func addFilesFromConfigFile(filesToParse []string, fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Println("no config file found in : " + fileName)
		return filesToParse
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), "#") {
			filesToParse = append(filesToParse, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return filesToParse
}

func addFilesFromFlags(filesToParse []string) []string {
	for _, logfile := range logFilesToMonitor {
		filesToParse = append(filesToParse, logfile)
	}
	return filesToParse
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

func countOccurences(logfile string, pattern string, substringLimit int) (int, error) {
	file, err := os.Open(logfile)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer file.Close()
	count := 0

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if substringLimit != -1 {
			if len(line) > substringLimit {
				line = line[:substringLimit]
			}
		}

		if strings.Contains(string(line), pattern) {
			count++
		}

	}

	return count, nil
}

func prometheusFormat(logfile string, pattern string, nbErrors int) string {
	normalizedPattern := normalizePattern(pattern)
	res := fmt.Sprintf("pattern_in_log_count_%s{logfile=\"%s\", pattern=\"%s\"} %d\n", normalizedPattern, logfile, pattern, nbErrors)
	return res
}

func normalizePattern(pattern string) string {

	pattern = strings.ReplaceAll(pattern, " ", "_")
	pattern = strings.ReplaceAll(pattern, "/", "_")
	pattern = strings.ReplaceAll(pattern, "\\", "_")
	return pattern
}
