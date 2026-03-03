package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/yourusername/pokedexcli/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2"

var cache = pokecache.NewCache(5 * time.Minute)

type LocationAreaResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type ExploreResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func GetLocationAreas(pageURL *string) (LocationAreaResponse, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	if val, ok := cache.Get(url); ok {
		var locationResp LocationAreaResponse
		err := json.Unmarshal(val, &locationResp)
		return locationResp, err
	}

	res, err := http.Get(url)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	cache.Add(url, data)

	var locationResp LocationAreaResponse
	err = json.Unmarshal(data, &locationResp)
	return locationResp, err
}

func ExploreArea(areaName string) (ExploreResponse, error) {
	url := baseURL + "/location-area/" + areaName

	if val, ok := cache.Get(url); ok {
		var exploreResp ExploreResponse
		err := json.Unmarshal(val, &exploreResp)
		return exploreResp, err
	}

	res, err := http.Get(url)
	if err != nil {
		return ExploreResponse{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return ExploreResponse{}, err
	}

	cache.Add(url, data)

	var exploreResp ExploreResponse
	err = json.Unmarshal(data, &exploreResp)
	return exploreResp, err
}

func GetPokemon(name string) (Pokemon, error) {
	url := baseURL + "/pokemon/" + name

	if val, ok := cache.Get(url); ok {
		var pokemon Pokemon
		err := json.Unmarshal(val, &pokemon)
		return pokemon, err
	}

	res, err := http.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}

	cache.Add(url, data)

	var pokemon Pokemon
	err = json.Unmarshal(data, &pokemon)
	return pokemon, err
}
