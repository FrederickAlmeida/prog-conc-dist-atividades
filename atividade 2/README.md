Implemente um sistema de análise de logs onde múltiplos threads leem e classificam registros de um log por categoria (erro, info, warning). Os resultados são agrupados em estruturas separadas.

SETUP:
- coloque o arquivo de logs no mesmo diretório que o arquivo main.go
- coloque numLines igual ao número de linhas do arquivo
- estando no diretório do arquivo main.go e logs.txt, rode:
go run main.go