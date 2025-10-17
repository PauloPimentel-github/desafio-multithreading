package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PauloPimentel-github/desafio-multithreading/dto"
)

func FetchBrasilAPI(ctx context.Context, cep string, channel chan<- dto.CEPResult) {
	log.Printf("Iniciando busca na BrasilAPI...")
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		channel <- dto.CEPResult{Source: "BrasilAPI", Error: fmt.Errorf("erro ao criar requisição: %w", err)}
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// Não envia erro para o channel se o contexto for cancelado (TIMEOUT ou outro resultado)
		if ctx.Err() != nil {
			return
		}
		channel <- dto.CEPResult{Source: "BrasilAPI", Error: fmt.Errorf("erro na requisição HTTP: %w", err)}
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		channel <- dto.CEPResult{Source: "BrasilAPI", Error: fmt.Errorf("status inesperado: %s", res.Status)}
		return
	}

	var data dto.BrasilAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		channel <- dto.CEPResult{Source: "BrasilAPI", Error: fmt.Errorf("erro ao decodificar JSON: %w", err)}
		return
	}

	// Envia o resultado para o channel
	channel <- dto.CEPResult{
		Source: "BrasilAPI",
		Street: data.Street,
		City:   data.City,
		State:  data.State,
	}
}

func FetchViaCEP(ctx context.Context, cep string, channel chan<- dto.CEPResult) {
	log.Printf("Iniciando busca na ViaCEP...")
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		channel <- dto.CEPResult{Source: "ViaCEP", Error: fmt.Errorf("erro ao criar requisição: %w", err)}
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return
		}
		channel <- dto.CEPResult{Source: "ViaCEP", Error: fmt.Errorf("erro na requisição HTTP: %w", err)}
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		channel <- dto.CEPResult{Source: "ViaCEP", Error: fmt.Errorf("status inesperado: %s", res.Status)}
		return
	}

	var data dto.ViaCEPResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		channel <- dto.CEPResult{Source: "ViaCEP", Error: fmt.Errorf("erro ao decodificar JSON: %w", err)}
		return
	}

	// Verifica se a ViaCEP retornou um erro interno
	if data.Erro {
		channel <- dto.CEPResult{Source: "ViaCEP", Error: fmt.Errorf("CEP não encontrado pela ViaCEP")}
		return
	}

	// Envia o resultado para o channel
	channel <- dto.CEPResult{
		Source: "ViaCEP",
		Street: data.Logradouro,
		City:   data.Localidade,
		State:  data.Uf,
	}
}
