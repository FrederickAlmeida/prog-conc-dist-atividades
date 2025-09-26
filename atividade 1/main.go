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
				errCh <- linha
			case strings.HasPrefix(formattedString, "[WARNING]"):
				warnCh <- linha
			case strings.HasPrefix(formattedString, "[INFO]"):
				infoCh <- linha
			}
		}
		lineNum++
	}

	if err := sc.Err(); err != nil {
		fmt.Printf("Erro lendo arquivo: %v\n", err)
		return
	}
}

func writeToFile(fileName string, ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Erro criando arquivo %s: %v\n", fileName, err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for msg := range ch {
		writer.WriteString(msg + "\n")
	}
	writer.Flush()
}

func main() {
	errCh := make(chan string, 10)
	warnCh := make(chan string, 10)
	infoCh := make(chan string, 10)

	numLines := 1000
	linesPerThread := numLines / 4

	var readWg sync.WaitGroup
	var writeWg sync.WaitGroup

	writeWg.Add(3)
	go writeToFile("error.txt", errCh, &writeWg)
	go writeToFile("warning.txt", warnCh, &writeWg)
	go writeToFile("info.txt", infoCh, &writeWg)

	for i := 0; i < 4; i++ {
		startLine := i * linesPerThread
		endLine := (i + 1) * linesPerThread

		// Adiciona o resto para a última thread
		if i == 3 {
			endLine += numLines % 4
		}

		fmt.Printf("Iniciando thread %d: linhas %d até %d\n", i+1, startLine, endLine-1)

		readWg.Add(1)
		go readFileRange("logs.txt", startLine, endLine, errCh, warnCh, infoCh, &readWg)
	}

	// Aguarda todas as goroutines completarem
	readWg.Wait()

	// Fecha os canais após a leitura
	close(errCh)
	close(warnCh)
	close(infoCh)

	// Aguarda as escritas terminarem
	writeWg.Wait()
	fmt.Printf("Programa concluído.")
}
