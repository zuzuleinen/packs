package domain

import (
	"fmt"
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

	for i := len(c.config.packs) - 1; i >= 0; i-- {
		div := totalItems / c.config.packs[i]
		if div > 0 {
			factors[c.config.packs[i]] = div
			totalItems -= c.config.packs[i] * div
		}
	}
	if totalItems > 0 {
		factors[minItemsPerUnit]++
	}

	fmt.Println(factors)

	// regulate factors
	// if we have 250 X 2 and 250 X 2 exists as pack size use that one
	if factors[minItemsPerUnit] > 1 {
		sizeToCheck := factors[minItemsPerUnit] * minItemsPerUnit
		if _, ok := factors[sizeToCheck]; ok {
			factors[sizeToCheck]++
			factors[minItemsPerUnit] = 0
		}
	}

	var packs []Pack
	for items, qty := range factors {
		if qty > 0 {
			packs = append(packs, Pack{
				Qty:          qty,
				ItemsPerUnit: items,
			})
		}
	}

	return packs
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
