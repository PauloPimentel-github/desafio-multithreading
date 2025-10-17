package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PauloPimentel-github/desafio-multithreading/internal/dto"
)

func main() {
	targetCEP := "07263725" // Você pode alterar este CEP para testar

	log.Printf("Buscando CEP %s nas duas APIs. Timeout de 1 segundo...", targetCEP)

	result := searchCEP(targetCEP)

	if result.Error != nil {
		fmt.Printf("\n!!! ERRO: %v\n", result.Error)
		return
	}

	fmt.Println("\n================================================")
	fmt.Printf("✔ RESPOSTA MAIS RÁPIDA ENCONTRADA\n")
	fmt.Printf("Fonte: %s\n", result.Source)
	fmt.Printf("Endereço: %s, %s - %s\n", result.Street, result.City, result.State)
	fmt.Println("================================================")
}

// ------------------------------------------------
// COORDENADOR DE CONCORRÊNCIA (MULTITHREADING)
// ------------------------------------------------

func searchCEP(cep string) dto.CEPResult {
	// 1. Cria um Contexto com Timeout de 1 segundo (requisito do desafio)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Garante que as goroutines parem se o resultado for encontrado ou o tempo expirar.

	// 2. Cria um Channel para receber o primeiro resultado
	// O buffer de 1 garante que a primeira goroutine que enviar o resultado não trave.
	resultChan := make(chan dto.CEPResult, 1)

	// 3. Lança as duas Goroutines
	go FetchBrasilAPI(ctx, cep, resultChan)
	go FetchViaCEP(ctx, cep, resultChan)

	// 4. Usa o 'select' para esperar pelo resultado ou pelo timeout
	select {
	case result := <-resultChan:
		// Retorna o resultado mais rápido
		return result
	case <-ctx.Done():
		// Retorna um erro se o timeout for atingido
		return dto.CEPResult{
			Error: fmt.Errorf("timeout excedido: A busca demorou mais de 1 segundo (%w)", ctx.Err()),
		}
	}
}
