package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Estrutura genérica de resposta para unificar os resultados.
type CEPResult struct {
	Source string
	Street string
	City   string
	State  string
	Error  error
}

// ------------------------------------------------
// ESTRUTURAS DE DADOS PARA AS APIs
// ------------------------------------------------

// BrasilAPI - Estrutura de resposta (simplificada)
type BrasilAPIResponse struct {
	Street string `json:"street"`
	City   string `json:"city"`
	State  string `json:"state"`
}

// ViaCEP - Estrutura de resposta (simplificada)
type ViaCEPResponse struct {
	Logradouro string `json:"logradouro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
	Erro       bool   `json:"erro,omitempty"` // Campo opcional que ViaCEP usa para erro
}

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

func searchCEP(cep string) CEPResult {
	// 1. Cria um Contexto com Timeout de 1 segundo (requisito do desafio)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Garante que as goroutines parem se o resultado for encontrado ou o tempo expirar.

	// 2. Cria um Channel para receber o primeiro resultado
	// O buffer de 1 garante que a primeira goroutine que enviar o resultado não trave.
	resultChan := make(chan CEPResult, 1)

	// 3. Lança as duas Goroutines
	go fetchBrasilAPI(ctx, cep, resultChan)
	go fetchViaCEP(ctx, cep, resultChan)

	// 4. Usa o 'select' para esperar pelo resultado ou pelo timeout
	select {
	case result := <-resultChan:
		// Retorna o resultado mais rápido
		return result
	case <-ctx.Done():
		// Retorna um erro se o timeout for atingido
		return CEPResult{
			Error: fmt.Errorf("timeout excedido: A busca demorou mais de 1 segundo (%w)", ctx.Err()),
		}
	}
}

// ------------------------------------------------
// FUNÇÕES DE BUSCA (GOROUTINES)
// ------------------------------------------------

func fetchBrasilAPI(ctx context.Context, cep string, channel chan<- CEPResult) {
	time.Sleep(time.Second * 2)
	log.Printf("Iniciando busca na BrasilAPI...")
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		channel <- CEPResult{Source: "BrasilAPI", Error: fmt.Errorf("erro ao criar requisição: %w", err)}
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// Não envia erro para o channel se o contexto for cancelado (TIMEOUT ou outro resultado)
		if ctx.Err() != nil {
			return
		}
		channel <- CEPResult{Source: "BrasilAPI", Error: fmt.Errorf("erro na requisição HTTP: %w", err)}
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		channel <- CEPResult{Source: "BrasilAPI", Error: fmt.Errorf("status inesperado: %s", res.Status)}
		return
	}

	var data BrasilAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		channel <- CEPResult{Source: "BrasilAPI", Error: fmt.Errorf("erro ao decodificar JSON: %w", err)}
		return
	}

	// Envia o resultado para o channel
	channel <- CEPResult{
		Source: "BrasilAPI",
		Street: data.Street,
		City:   data.City,
		State:  data.State,
	}
}

func fetchViaCEP(ctx context.Context, cep string, channel chan<- CEPResult) {
	log.Printf("Iniciando busca na ViaCEP...")
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		channel <- CEPResult{Source: "ViaCEP", Error: fmt.Errorf("erro ao criar requisição: %w", err)}
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		channel <- CEPResult{Source: "ViaCEP", Error: fmt.Errorf("erro na requisição HTTP: %w", err)}
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		channel <- CEPResult{Source: "ViaCEP", Error: fmt.Errorf("status inesperado: %s", res.Status)}
		return
	}

	var data ViaCEPResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		channel <- CEPResult{Source: "ViaCEP", Error: fmt.Errorf("erro ao decodificar JSON: %w", err)}
		return
	}

	// Verifica se a ViaCEP retornou um erro interno
	if data.Erro {
		channel <- CEPResult{Source: "ViaCEP", Error: fmt.Errorf("CEP não encontrado pela ViaCEP")}
		return
	}

	// Envia o resultado para o channel
	channel <- CEPResult{
		Source: "ViaCEP",
		Street: data.Logradouro,
		City:   data.Localidade,
		State:  data.Uf,
	}
}
