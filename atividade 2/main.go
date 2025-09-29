package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

func readFileRange(filePath string, startLine, endLine int, wg *sync.WaitGroup, wwg *sync.WaitGroup, errmu *sync.Mutex, warmu *sync.Mutex, infmu *sync.Mutex) {
	defer wg.Done()

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Erro abrindo arquivo: %v\n", err)
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
			fmt.Printf("Lines %d-%d: %s\n", startLine, endLine-1, linha)
			
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

	if err := sc.Err(); err != nil {
		fmt.Printf("Erro lendo arquivo: %v\n", err)
		return
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

func main() {
	numLines := 1000
	linesPerThread := numLines / 4

	var readWg sync.WaitGroup
	var writeWg sync.WaitGroup
	var errMu, warMu, infMu sync.Mutex

	for i := 0; i < 4; i++ {
		startLine := i * linesPerThread
		endLine := (i + 1) * linesPerThread

		// Adiciona o resto para a última thread
		if i == 3 {
			endLine += numLines % 4
		}

		fmt.Printf("Iniciando thread %d: linhas %d até %d\n", i+1, startLine, endLine-1)

		readWg.Add(1)
		go readFileRange("logs.txt", startLine, endLine, &readWg, &writeWg, &errMu, &warMu, &infMu)
	}

	// Aguarda todas as goroutines completarem
	readWg.Wait()

	// Aguarda as escritas terminarem
	writeWg.Wait()
	fmt.Printf("Programa concluído.")
}
