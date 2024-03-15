package domain

import (
	"math"
	"sort"
)

// Calculator calculates the list of packs needed based on number of total items
type Calculator struct {
	config *Config
}

// NewCalculator creates a new Calculator based on Config
func NewCalculator(c *Config) *Calculator {
	return &Calculator{config: c}
}

// Pack represents a package to be added for order
type Pack struct {
	Qty          int
	ItemsPerUnit int
}

// Packs returns the number of Pack based on total items ordered
func (c *Calculator) Packs(totalItems int) []Pack {
	sort.Ints(c.config.sizes)

	minItemsPerUnit := c.config.sizes[0]

	factors := make(map[int]int)
	for _, v := range c.config.sizes {
		factors[v] = 0
	}

	for totalItems > 0 {
		if totalItems < minItemsPerUnit {
			factors[minItemsPerUnit]++
			totalItems -= minItemsPerUnit
			continue
		}
		closest := findClosest(totalItems, c.config.sizes)

		quotient := totalItems / closest

		if totalItems%closest == 0 {
			factors[closest] += quotient
			totalItems -= quotient * closest
		} else {
			// see if closest + minItemsPerUnit exists in factors
			if _, ok := factors[closest+minItemsPerUnit]; ok {
				factors[closest+minItemsPerUnit]++
				totalItems -= closest + minItemsPerUnit
			} else {
				factors[closest]++
				totalItems -= closest
			}
		}
	}

	var ps []Pack
	for k, v := range factors {
		if v > 0 {
			ps = append(ps, Pack{
				Qty:          v,
				ItemsPerUnit: k,
			})
		}
	}

	return ps
}

// findClosest finds the closest int to target from s
func findClosest(target int, s []int) int {
	closest := math.MaxInt
	for _, v := range s {
		if absDiff(v, target) < absDiff(target, closest) {
			closest = v
		}
	}
	return closest
}

// absDiff finds the absolute difference between a and b
func absDiff(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}
