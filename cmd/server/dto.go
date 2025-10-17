package dto

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
