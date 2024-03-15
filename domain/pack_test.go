package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackSizeCanBeConfigured(t *testing.T) {
	c := NewConfig()
	c.AddPackSize(250)
	c.AddPackSize(500)
	c.AddPackSize(1000)

	assert.Equal(t, 3, c.Len(), "Len should match the number of pack sizes")
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

func TestWithTwoSizesAmountingToAnotherSize(t *testing.T) {
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
