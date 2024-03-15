package domain

import (
	"math"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Config struct {
	packs []int
}

func NewConfig() *Config {
	var c Config
	c.packs = make([]int, 0)
	return &c
}

func (c *Config) AddPackSize(p int) {
	c.packs = append(c.packs, p)
}

func (c *Config) Len() int {
	return len(c.packs)
}

func TestPackSizeCanBeConfigured(t *testing.T) {
	c := NewConfig()
	c.AddPackSize(250)
	c.AddPackSize(500)
	c.AddPackSize(1000)

	assert.Equal(t, 3, c.Len(), "Len should match the number of pack sizes")
}

type Calculator struct {
	config *Config
}

func NewCalculator(c *Config) *Calculator {
	return &Calculator{config: c}
}

type Pack struct {
	Qty          int
	ItemsPerUnit int
}

func (c *Calculator) Packs(totalItems int) []Pack {
	sort.Ints(c.config.packs)

	minItemsPerUnit := c.config.packs[0]

	factors := make(map[int]int)
	for _, v := range c.config.packs {
		factors[v] = 0
	}

	for totalItems > 0 {
		if totalItems < minItemsPerUnit {
			factors[minItemsPerUnit]++
			totalItems -= minItemsPerUnit
			continue
		}
		clst := findClosest(totalItems, c.config.packs)

		div := totalItems / clst
		mod := totalItems % clst

		if mod != 0 {
			// see if clst + minItemsPerUnit exists in factors
			if _, ok := factors[clst+minItemsPerUnit]; ok {
				factors[clst+minItemsPerUnit]++
				totalItems -= clst + minItemsPerUnit
			} else {
				factors[clst]++
				totalItems -= clst
			}
		} else {
			factors[clst] += div
			totalItems = totalItems - div*clst
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

func findClosest(target int, s []int) int {
	closest := math.MaxInt
	for _, v := range s {
		if absDiff(v, target) < absDiff(target, closest) {
			closest = v
		}
	}
	return closest
}

func absDiff(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

func TestCalculator(t *testing.T) {
	c := NewConfig()
	c.AddPackSize(250)
	c.AddPackSize(500)
	c.AddPackSize(1000)
	c.AddPackSize(2000)
	c.AddPackSize(5000)

	calc := NewCalculator(c)

	packs := calc.Packs(1)
	assert.Len(t, packs, 1)
	assert.True(t, packExists(1, 250, packs))

	packs = calc.Packs(250)
	assert.Len(t, packs, 1)
	assert.True(t, packExists(1, 250, packs))

	packs = calc.Packs(251)
	assert.Len(t, packs, 1)
	assert.True(t, packExists(1, 500, packs))

	packs = calc.Packs(500)
	assert.Len(t, packs, 1)
	assert.True(t, packExists(1, 500, packs))

	packs = calc.Packs(501)
	assert.Len(t, packs, 2)
	assert.True(t, packExists(1, 500, packs))
	assert.True(t, packExists(1, 250, packs))

	packs = calc.Packs(1000)
	assert.Len(t, packs, 1)
	assert.True(t, packExists(1, 1000, packs))

	packs = calc.Packs(1500)
	assert.Len(t, packs, 2)
	assert.True(t, packExists(1, 1000, packs))
	assert.True(t, packExists(1, 500, packs))

	packs = calc.Packs(12_001)
	assert.Len(t, packs, 3)
	assert.True(t, packExists(2, 5000, packs))
	assert.True(t, packExists(1, 2000, packs))
	assert.True(t, packExists(1, 250, packs))
}

func TestAnother(t *testing.T) {
	c := NewConfig()
	c.AddPackSize(250)
	c.AddPackSize(500)
	c.AddPackSize(750)
	c.AddPackSize(1000)
	c.AddPackSize(2000)
	c.AddPackSize(5000)

	calc := NewCalculator(c)

	packs := calc.Packs(750)
	assert.Len(t, packs, 1)
	assert.True(t, packExists(1, 750, packs))

	packs = calc.Packs(501)
	assert.Len(t, packs, 1)
	assert.True(t, packExists(1, 750, packs))
}

func packExists(qty int, itemsPerUnit int, bs []Pack) bool {
	for _, v := range bs {
		if v.Qty == qty && v.ItemsPerUnit == itemsPerUnit {
			return true
		}
	}
	return false
}
