package domain

// Config holds the list of all configured sizes
type Config struct {
	sizes []int
}

// NewConfig initializes a new
func NewConfig() *Config {
	var c Config
	c.sizes = make([]int, 0)
	return &c
}

// AddPackSize adds a new size to configuration
func (c *Config) AddPackSize(size int) {
	c.sizes = append(c.sizes, size)
}

// Len returns the number of configured sizes
func (c *Config) Len() int {
	return len(c.sizes)
}
