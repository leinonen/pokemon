package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type Pokemon struct {
	Weight int    `json:"weight"`
	Url    string `json:"url"`
	ID     int    `json:"id"`
}

type PokemonListItem struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

type PokemonList struct {
	Results []PokemonListItem `json:"results"`
}

func getPokemonWeight(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	body, _ := io.ReadAll(resp.Body)
	var poke Pokemon
	json.Unmarshal(body, &poke)

	return poke.Weight, nil

}

func main() {
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon?limit=500")
	if err != nil {
		fmt.Println("Error")
		return
	}
	body, _ := io.ReadAll(resp.Body)
	var pokemons PokemonList
	json.Unmarshal(body, &pokemons)

	var wg sync.WaitGroup
	weights := make(chan int, 10)
	for _, p := range pokemons.Results {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			weight, _ := getPokemonWeight(url)
			weights <- weight
		}(p.Url)
	}

	go func() {
		defer close(weights)
		wg.Wait()
	}()

	var count, totalWeight int
	for currentWeight := range weights {
		count++
		totalWeight += currentWeight

		fmt.Println("Avg", totalWeight/count, "Current", currentWeight, "Count", count)
	}
}
