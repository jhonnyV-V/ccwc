package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func parseArgs(flags map[string]bool) ([]string, bool) {
	args := os.Args[1:]
	filenames := []string{}
	fullReport := true
	for _, v := range args {
		if v[0] == '-' {
			if fullReport {
				fullReport = false
			}
			flags[v] = true
		} else {
			filenames = append(filenames, v)
		}
	}
	return filenames, fullReport
}

func readLine(in *bufio.Reader) (string, error) {
	var bytes []byte
	for {
		line, isPrefix, err := in.ReadLine()
		if err != nil {
			return "", err
		}
		bytes = append(bytes, line...)
		if !isPrefix {
			break
		}
	}
	return string(bytes), nil
}

func generateReport(
	nlines, nwords, nchar, nbytes int,
	flags map[string]bool,
	fullReport bool,
	filename string,
) string {
	var report string

	if fullReport {
		report = fmt.Sprintf("%d %d %d %s\n", nlines, nwords, nchar, filename)
	} else {
		if flags["-l"] {
			report += fmt.Sprintf(" %d", nlines)
		}
		if flags["-w"] {
			report += fmt.Sprintf(" %d", nwords)
		}
		if flags["-m"] {
			report += fmt.Sprintf(" %d", nchar)
		}
		if flags["-c"] {
			report += fmt.Sprintf(" %d", nbytes)
		}

		report = strings.TrimSpace(report)

		report += fmt.Sprintf(" %s\n", filename)
	}
	return report
}

func getNumOfWords(line string) int {
	words := strings.Split(strings.TrimSpace(line), " ")
	num := len(words)
	if num == 1 && words[0] == "" {
		return 0
	}
	return num
}

func read(file *bufio.Reader) (nlines, nwords, nchars, nbytes int) {
	nlines = 0
	nwords = 0
	nchars = 0
	nbytes = 0

	for {
		line, err := readLine(file)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		nlines++
		nwords += getNumOfWords(line)
		nchars += len([]rune(line))
		nbytes += len(line)
	}
	//EOF
	nchars++

	return nlines, nwords, nchars, nbytes
}

func main() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	var flags map[string]bool = make(map[string]bool, 4)

	filenames, fullReport := parseArgs(flags)

	var files []*bufio.Reader

	if len(filenames) == 0 {
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			files = append(files, bufio.NewReader(os.Stdin))
		} else {
			fmt.Println("No Valid Input")
			os.Exit(0)
		}
	} else {
		for _, v := range filenames {
			localFile, err := os.Open(v)
			if err != nil {
				panic(err)
			}
			files = append(files, bufio.NewReader(localFile))
		}
	}

	totalLines := 0
	totalWords := 0
	totalChars := 0
	totalBytes := 0

	for i, f := range files {
		var result string
		filename := ""
		if i < len(filenames) {
			filename = filenames[i]
		}
		nlines, nwords, nchars, nbytes := read(f)

		result = generateReport(
			nlines, nwords, nchars, nbytes,
			flags,
			fullReport,
			filename,
		)

		totalLines += nlines
		totalWords += nwords
		totalChars += nchars
		totalBytes += nbytes

		fmt.Print(result)
	}

	if len(files) > 1 {
		result := generateReport(
			totalLines,
			totalWords,
			totalChars,
			totalBytes,
			flags,
			fullReport,
			"total",
		)
		fmt.Print(result)
	}
}
