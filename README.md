# ⚡ Desafio Go Concorrência: Buscador de CEP Rápido

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![API](https://img.shields.io/badge/REST%20API-lightgrey?style=for-the-badge&logo=api)
![Multithreading](https://img.shields.io/badge/Concurrency-FF9900?style=for-the-badge&logo=go)

Este projeto soluciona o desafio de concorrência do curso Go Expert (Full Cycle), implementando um buscador de CEP que utiliza **Multithreading** para consultar duas APIs externas simultaneamente. O objetivo é acatar a resposta mais rápida e garantir um tempo limite de execução global.

---

### ✨ Destaques da Solução

A solução demonstra a proficiência em técnicas avançadas de concorrência em Go, focando em:

* **Busca Concorrente (Multithreading):** Utiliza *goroutines* para disparar requisições HTTP para as duas APIs de CEP de forma paralela.
* **Mecanismo de Descarte (Racer Pattern):** Emprega *channels* para comunicar o resultado da requisição mais rápida ao *main thread*, garantindo que a resposta mais lenta seja automaticamente descartada.
* **Contexto com Timeout:** O uso de `context.WithTimeout` estabelece um limite de tempo rígido para o processo inteiro (as duas buscas e o consumo do channel), garantindo que o programa não demore mais que o máximo estipulado.
* **Gestão de Erros de Tempo:** Em caso de timeout, o programa exibe uma mensagem clara de erro, atendendo ao requisito de limite de resposta.

---

### ⚙️ Tecnologias e APIs

* **Linguagem de Programação:** Go
* **Bibliotecas:**
    * `net/http`: Para a realização das requisições HTTP.
    * `context`: Para a limitação de tempo de resposta.
* **APIs Consumidas:**
    * **API 1 (BrasilAPI):** `https://brasilapi.com.br/api/cep/v1/{cep}`
    * **API 2 (ViaCEP):** `http://viacep.com.br/ws/{cep}/json/`

---

### 🎯 Requisitos do Desafio

| Requisito | Solução Implementada |
| :--- | :--- |
| **1. Busca Mais Rápida** | As duas requisições são lançadas em goroutines. A primeira resposta a chegar ao *channel* é processada, e o processo é finalizado. |
| **2. Limite de Tempo** | O `context.WithTimeout` limita todo o processo de busca e espera a **1 segundo**. |
| **3. Exibição Detalhada** | O resultado final é exibido no *command line*, mostrando o endereço completo e o nome da API (`BrasilAPI` ou `ViaCEP`) que forneceu a resposta. |
| **4. Tratamento de Timeout** | Se o tempo limite de 1 segundo for atingido, o `select` bloqueia a busca e o erro de timeout do contexto é exibido no console. |

---

### 🏁 Como Executar o Projeto

#### Pré-requisitos
* Go instalado (versão 1.18 ou superior).

1.  Clone o repositório:
    ```bash
    git clone [https://github.com/seu-usuario/seu-repositorio.git](https://github.com/seu-usuario/seu-repositorio.git)
    cd seu-repositorio
    ```

2.  Execute o arquivo principal, passando o CEP desejado como argumento        (Exemplo: 07263725):
    ```bash
    go run main.go | go run main.go --cep=07263725
    ```

## Curso Go Expert (Full Cycle)
<p>Feito com ♥ by Paulo H.G Pimentel</p>
