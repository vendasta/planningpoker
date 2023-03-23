package main

import (
	"fmt"
	"math/rand"
)

var animals = []string{
	"badger",
	"bat",
	"bear",
	"dog",
	"moose",
	"crow",
	"goose",
	"mole",
	"rat",
	"dragon",
	"monkey",
	"snake",
	"tiger",
	"lion",
}

var adjectives = []string{
	"spicy",
	"rabid",
	"cheeky",
	"snarly",
	"sneaky",
	"lazy",
	"hungry",
	"sleepy",
	"grumpy",
	"sneezy",
	"bouncy",
	"fluffy",
	"silly",
}

func GenerateID() string {
	animal := animals[rand.Intn(len(animals))]
	adj := adjectives[rand.Intn(len(adjectives))]
	return fmt.Sprintf("%s %s", adj, animal)
}
