package config

import (
	"github.com/imdario/mergo"
)

// Config represents a simple configuration object the merger tool
type Config struct {
	Format                      string
	WithOverride                bool
	WithAppendSlice             bool
	WithOverwriteWithEmptyValue bool
	WithSliceDeepCopy           bool
}

// New creates a new configuration from the provided input
func New(input []interface{}) *Config {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	c := input[0].(map[string]interface{})

	return &Config{
		Format:                      c["format"].(string),
		WithAppendSlice:             c["with_append_slice"].(bool),
		WithOverride:                c["with_override"].(bool),
		WithOverwriteWithEmptyValue: c["with_overwrite_with_empty_value"].(bool),
		WithSliceDeepCopy:           c["with_slice_deep_copy"].(bool),
	}
}

// ToMergoConfig returns an array of mergo.Config functions based on the current value of the configuration.
func (c *Config) ToMergoConfig() []func(config *mergo.Config) {
	configs := make([]func(config *mergo.Config), 0)

	if c.WithAppendSlice {
		configs = append(configs, mergo.WithAppendSlice)
	}

	if c.WithOverride {
		configs = append(configs, mergo.WithOverride)
	}

	if c.WithOverwriteWithEmptyValue {
		configs = append(configs, mergo.WithOverwriteWithEmptyValue)
	}

	if c.WithSliceDeepCopy {
		configs = append(configs, mergo.WithSliceDeepCopy)
	}

	return configs
}
