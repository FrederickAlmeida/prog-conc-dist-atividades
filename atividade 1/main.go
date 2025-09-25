package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

func readFileRange(filePath string, startLine, endLine int, errCh chan string, warnCh chan string, infoCh chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
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

			if strings.Contains(linha, "ERROR") {
				errCh <- linha
			}
			if strings.Contains(linha, "WARNING") {
				warnCh <- linha
			}
			if strings.Contains(linha, "INFO") {
				infoCh <- linha
			}
		}
		lineNum++
	}

	if err := sc.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
}

func main() {
	errCh := make(chan string, 10)
	warnCh := make(chan string, 10)
	infoCh := make(chan string, 10)

	numLines := 1000
	linesPerThread := numLines / 4

	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		startLine := i * linesPerThread
		endLine := (i + 1) * linesPerThread

		// Adiciona o resto para a última thread
		if i == 3 {
			endLine += numLines % 4
		}

		fmt.Printf("Iniciando thread %d: linhas %d até %d\n", i+1, startLine, endLine-1)

		wg.Add(1)
		go readFileRange("logs.txt", startLine, endLine, errCh, warnCh, infoCh, &wg)
	}

	// Aguarda todas as goroutines completarem
	wg.Wait()
}
