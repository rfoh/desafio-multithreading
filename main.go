package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Thread 1
func main() {
	// Neste desafio você terá que usar o que aprendemos com Multithreading e APIs para buscar o resultado mais rápido entre duas APIs distintas.

	// As duas requisições serão feitas simultaneamente para as seguintes APIs:

	// https://brasilapi.com.br/api/cep/v1/01153000 + cep

	// http://viacep.com.br/ws/" + cep + "/json/

	// Os requisitos para este desafio são:

	// - Acatar a API que entregar a resposta mais rápida e descartar a resposta mais lenta.

	// - O resultado da request deverá ser exibido no command line com os dados do endereço, bem como qual API a enviou.

	// - Limitar o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout deve ser exibido.

	c1 := make(chan string)
	c2 := make(chan string)

	if len(os.Args) < 2 {
		fmt.Println("Erro: CEP não fornecido. Use: go run main.go <CEP>")
		return
	}
	cep := os.Args[1]

	go func() {
		retorno, err := getCEP("https://brasilapi.com.br/api/cep/v1/" + cep)
		if err != nil {
			c1 <- err.Error()
		}
		c1 <- retorno
	}()

	go func() {
		retorno, err := getCEP("http://viacep.com.br/ws/" + cep + "/json/")
		if err != nil {
			c2 <- err.Error()
		}
		c2 <- retorno
	}()

	select {
	case resultado := <-c1:
		fmt.Println("Resultado BrasilAPI: ", resultado)
	case resultado := <-c2:
		fmt.Println("Resultado ViaCEP: ", resultado)
	case <-time.After(1 * time.Second):
		fmt.Println("Erro: Timeout de 1 segundo atingido.")
	}
}

func getCEP(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Erro: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Erro ao ler o corpo da resposta: %v", err)
	}

	retorno := string(body)
	return retorno, nil
}
