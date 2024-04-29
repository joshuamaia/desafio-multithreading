package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Address struct {
	Address     interface {} 
	Origem      string
}



func searchAddress(apiUrl string, ch chan Address) {

	resp, err := http.Get(apiUrl)
	if err != nil {
		fmt.Printf("Erro ao fazer a requisição: %s\n", err)
		return
	}
	defer resp.Body.Close()

	var addressSearch interface{}
	if err := json.NewDecoder(resp.Body).Decode(&addressSearch); err != nil {
		fmt.Printf("Erro ao decodificar a resposta: %s\n", err)
		return
	}

	address := Address{addressSearch, apiUrl}

	ch <- address
}

func printPrettyJson(address interface {}) {

	json, err := json.MarshalIndent(address, "", "  ")
	if err != nil {	
		fmt.Printf("Erro ao codificar o JSON: %s\n", err)
		return
	}

	fmt.Println(string(json))
}

func main() {
	cep := "70150900" 
	ch := make(chan Address, 2)

	apiViaCep := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	apiBrasil := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%scep", cep)
	
	go searchAddress(apiViaCep, ch)
	go searchAddress(apiBrasil, ch)

	timeout := time.After(1 * time.Second)

	select {
	case address := <-ch:
		fmt.Println("Endereço:")
		printPrettyJson(address.Address)	
		fmt.Printf("Origem: %s\n", address.Origem)
	case <-timeout:
		fmt.Println("Erro: Timeout de 1 segundo excedido.")
	}
}
