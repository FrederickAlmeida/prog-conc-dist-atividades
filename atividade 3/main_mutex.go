package main

import (
	"bufio"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
	"fmt"
)

func readFileRange(filePath string, startLine, endLine int, wg *sync.WaitGroup, wwg *sync.WaitGroup, errmu *sync.Mutex, warmu *sync.Mutex, infmu *sync.Mutex) {
	defer wg.Done()

	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	lineNum := 0

	for sc.Scan() {
		if lineNum >= endLine {
			break
		}

		if lineNum >= startLine {
			linha := sc.Text()
			formattedString := strings.ToUpper(strings.TrimSpace(linha))

			switch {
			case strings.HasPrefix(formattedString, "[ERROR]"):
				wwg.Add(1)
				go writeToErrFile(linha, errmu, wwg)
			case strings.HasPrefix(formattedString, "[WARNING]"):
				wwg.Add(1)
				go writeToWarningFile(linha, warmu, wwg)
			case strings.HasPrefix(formattedString, "[INFO]"):
				wwg.Add(1)
				go writeToInfoFile(linha, infmu, wwg)
			}
		}
		lineNum++
	}
}

func writeToErrFile(content string, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	fileName := "error.txt"

	mu.Lock()
	defer mu.Unlock()

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Erro abrindo/criando arquivo %s: %v\n", fileName, err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content + "\n")
	if err != nil {
		fmt.Printf("Erro escrevendo no arquivo: %v\n", err)
	}
}

func writeToInfoFile(content string, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	fileName := "info.txt"

	mu.Lock()
	defer mu.Unlock()

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Erro abrindo/criando arquivo %s: %v\n", fileName, err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content + "\n")
	if err != nil {
		fmt.Printf("Erro escrevendo no arquivo: %v\n", err)
	}
}

func writeToWarningFile(content string, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	fileName := "warning.txt"

	mu.Lock()
	defer mu.Unlock()

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Erro abrindo/criando arquivo %s: %v\n", fileName, err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content + "\n")
	if err != nil {
		fmt.Printf("Erro escrevendo no arquivo: %v\n", err)
	}
}
func runOnce() time.Duration {
	start := time.Now()

	_, _ = os.Create("error.txt")
	_, _ = os.Create("warning.txt")
	_, _ = os.Create("info.txt")

	numLines := 1000
	linesPerThread := numLines / 4

	var readWg sync.WaitGroup
	var writeWg sync.WaitGroup
	var errMu, warMu, infMu sync.Mutex

	for i := 0; i < 4; i++ {
		startLine := i * linesPerThread
		endLine := (i + 1) * linesPerThread
		if i == 3 {
			endLine += numLines % 4
		}

		readWg.Add(1)
		go readFileRange("logs.txt", startLine, endLine, &readWg, &writeWg, &errMu, &warMu, &infMu)
	}

	readWg.Wait()
	writeWg.Wait()

	os.Remove("error.txt")
	os.Remove("warning.txt")
	os.Remove("info.txt")

	return time.Since(start)
}

func main() {
	// Testar com GOMAXPROCS = 1, 2, 6
	for _, procs := range []int{1, 2, 4, 6} {
		runtime.GOMAXPROCS(procs)
		var total time.Duration

		for i := 1; i <= 30; i++ {
			duration := runOnce()
			total += duration
		}

		println("Tempo médio com GOMAXPROCS =", procs, ":", (total.Milliseconds() / 30), "ms")
	}
}
