package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"sync"
	"runtime"
)

func readFileRange(filePath string, startLine, endLine int, errCh chan string, warnCh chan string, infoCh chan string, wg *sync.WaitGroup) {
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
		return
	}
}

func writeToFile(fileName string, ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Create(fileName)
	if err != nil {
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
	// testando diferentes gomaxprocs
	procs := []int{1, 4, 6}
	for _, i := range procs {
		runtime.GOMAXPROCS(i)
		fmt.Printf("Iniciando teste com GOMAXPROCS = %d\n", i)
        var executionTime time.Duration = 0
		for range 30 {
			startTime := time.Now()

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
		
			endTime := time.Now()
			executionTime += endTime.Sub(startTime)
			
			// remover os arquivos criados
			os.Remove("error.txt")
			os.Remove("warning.txt")
			os.Remove("info.txt")
		}
		fmt.Printf("Tempo de execução médio: %v\n", executionTime / 30)
	}
}
