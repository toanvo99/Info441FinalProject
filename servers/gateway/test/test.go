package main

import (
	"fmt"

	"github.com/mtslzr/pokeapi-go"
)

func main() {
	list, err := pokeapi.Pokemon("bulbasaur")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(list)
}
