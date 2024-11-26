package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type APIResponseCEP struct {
	URL      string `json:"url"`
	Response string `json:"response"`
	Error    string `json:"error"`
}

type CEPBrasil struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service:"`
}

type viacep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {

	apis := []string{
		"https://viacep.com.br/ws/01001001/json/",
		"https://brasilapi.com.br/api/cep/v1/01153000",
	}

	chCep := make(chan APIResponseCEP, len(apis))

	for _, api := range apis {
		go getCEP(api, chCep)
	}

	select {
	case res := <-chCep:

		if res.Error == "" {
			// ler json response
			var enderecoDetails map[string]interface{}
			if err := json.Unmarshal([]byte(res.Response), &enderecoDetails); err != nil {
				fmt.Printf("Error ao ler JSON response: %v\n", err)
			} else {
				urlString := string(res.URL)
				fmt.Println("URL da requisicao: ", urlString)
				if strings.Contains(urlString, "brasilapi") {
					fmt.Printf("CEP: %s\n", enderecoDetails["cep"])
					fmt.Printf("UF: %s\n", enderecoDetails["state"])
					fmt.Printf("Cidade: %s\n", enderecoDetails["city"])
					fmt.Printf("Bairro: %s\n", enderecoDetails["neighborhood"])
					fmt.Printf("Logradouro: %s\n", enderecoDetails["street"])
					fmt.Printf("service: %s\n", enderecoDetails["service"])
				} else {
					fmt.Printf("Cep: %s\n", enderecoDetails["cep"])
					fmt.Printf("Logradouro: %s\n", enderecoDetails["logradouro"])
					fmt.Printf("Complemento: %s\n", enderecoDetails["complemento"])
					fmt.Printf("Unidade: %s\n", enderecoDetails["unidade"])
					fmt.Printf("Bairro: %s\n", enderecoDetails["bairro"])
					fmt.Printf("Localidade: %s\n", enderecoDetails["localidade"])
					fmt.Printf("UF: %s\n", enderecoDetails["uf"])
					fmt.Printf("Estado: %s\n", enderecoDetails["estado"])
					fmt.Printf("Regiao: %s\n", enderecoDetails["regiao"])
					fmt.Printf("Ibge: %s\n", enderecoDetails["ibge"])
					fmt.Printf("Gia: %s\n", enderecoDetails["gia"])
					fmt.Printf("DDD: %s\n", enderecoDetails["ddd"])
					fmt.Printf("Siafi: %s\n", enderecoDetails["siafi"])
				}

			}
		} else {
			fmt.Printf("URL: %s\nError: %s\n\n", res.URL, res.Error)
			return
		}

	case <-time.After(1 * time.Second):
		fmt.Println("Timeout ao processar APIs")
		return

	}

}

func getCEP(url string, ch chan<- APIResponseCEP) {
	req, err := http.Get(url)
	if err != nil {
		fmt.Println("Error ao processar request", err)
		return
	}

	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error ao ler response", err)

	}

	ch <- APIResponseCEP{URL: url, Response: string(res)}
}
