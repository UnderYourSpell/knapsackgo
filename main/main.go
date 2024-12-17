package main

import (
	"fmt"
	"math/rand"
	"sort"
)

type Genomer interface {
	getGene()
	setWeight(newWeight int)
	setGene(newGene []int)
	getWeight() int
	calcFitness()
}

type Gene struct {
	weight int
	gene   []int
}

func (x Gene) getWeight() int {
	return x.weight
}

func (x *Gene) setWeight(newWeight int) {
	x.weight = newWeight
}

func (x Gene) getGene() {
	fmt.Println("Gene:", x.gene)
}

func (x *Gene) setGene(newGene []int) {
	x.gene = newGene
}

// calculates fitness
func (x *Gene) calcFitness(items []int, numItems int, target int) {
	totalWeight := 0
	for i := 0; i < numItems; i++ {
		totalWeight += x.gene[i] * items[i]
	}
	if totalWeight > target {
		totalWeight = 0
	}
	x.weight = totalWeight
}

// Stochastic Universal Selection
// Based on code from C++ Implenentation of TSP
func SUSSelction(genes *[]Gene, parents *[]Gene, popSize int) {
	N := popSize / 2
	F := 0
	keepIndices := []int{}
	for i := 0; i < popSize; i++ {
		F += (*genes)[i].weight
	}
	P := F / N
	start := rand.Intn(P)
	for i := 0; i < N; i++ {
		keepIndices = append(keepIndices, start+i*P)
	}
	for i := 0; i < N; i++ {
		value := keepIndices[i]
		j := 0
		fitnessSum := 0
		for j < len(*genes) {
			fitnessSum += (*genes)[j].getWeight()
			if fitnessSum > value {
				break
			}
			j++
		}
		*parents = append(*parents, (*genes)[j])
	}
}

// Single Point Crossover
func SPX(g1 Gene, g2 Gene, parents *[]Gene, numItems int) {
	X := rand.Intn(numItems)
	newG1 := Gene{
		weight: 0,
		gene:   []int{},
	}
	newG2 := Gene{
		weight: 0,
		gene:   []int{},
	}
	gene1Front := g1.gene[0:X]
	gene1Back := g1.gene[X:]
	gene2Front := g2.gene[0:X]
	gene2Back := g2.gene[X:]

	newG1.gene = append(append([]int{}, gene1Front...), gene2Back...)
	newG2.gene = append(append([]int{}, gene2Front...), gene1Back...)
	*parents = append(*parents, newG1, newG2)
}

func swapMutate(g *Gene, numItems int) {
	threshold := 1
	for i := 0; i < numItems; i++ {
		randMutate := rand.Float32()
		if randMutate <= float32(threshold) {
			if g.gene[i] == 1 {
				g.gene[i] = 0
			} else {
				g.gene[i] = 1
			}
		}
	}
}

func main() {
	const numItems = 100
	const maxWeight = 100
	const popSize = 24
	const maxGenerations = 100
	const target = 5000

	itemList := [numItems]int{}
	var genePool []Gene

	//initializing weights for items
	for i := 0; i < numItems; i++ {
		itemList[i] = rand.Intn(maxWeight)
	}

	//Initalize gene pool
	for i := 0; i < popSize; i++ {
		var geneArray []int
		for j := 0; j < numItems; j++ {
			if rand.Float32() >= 0.5 {
				geneArray = append(geneArray, 1)
			} else {
				geneArray = append(geneArray, 0)
			}
		}

		newGene := Gene{
			weight: 0,
			gene:   geneArray,
		}
		newGene.calcFitness(itemList[:], numItems, target)
		genePool = append(genePool, newGene)
	}

	for p := 0; p < maxGenerations; p++ {
		//sort genePool
		sort.Slice(genePool, func(i, j int) bool {
			return genePool[i].weight > genePool[j].weight
		})

		var parents []Gene

		//Selection
		SUSSelction(&genePool, &parents, popSize)
		//Breeding
		for i := 0; i < popSize/2; i += 2 {
			SPX(parents[i], parents[i+1], &parents, numItems)
		}

		for i := 0; i < len(parents); i++ {
			swapMutate(&parents[i], numItems)
			parents[i].calcFitness(itemList[:], numItems, target)
		}

		genePool = append(genePool, parents...)
		sort.Slice(genePool, func(i, j int) bool {
			return genePool[i].weight > genePool[j].weight
		})
		genePool = genePool[:popSize]
	}
	//end timing here

	sort.Slice(genePool, func(i, j int) bool {
		return genePool[i].weight > genePool[j].weight
	})
	fmt.Println("Final Best Gene:", genePool[0].weight)
}
