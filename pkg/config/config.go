package config

type Config struct {
	Middlewares Middlewares `yaml:"middlewares"`
}

type Middlewares struct {
	BodyFilter MiddlewareConfig[BodyFilterConfig] `yaml:"bodyFilter,omitempty"` //nolint:tagliatelle // valid tag
}

type MiddlewareConfig[T any] struct {
	Enabled bool `yaml:"enabled"`
	Config  []T  `yaml:"config"`
}

type BodyFilterConfig struct {
	Paths   []BodyFilterConfigPaths `yaml:"paths"`
	Methods []string                `yaml:"methods"`
	Filter  string                  `yaml:"filter"`
}

type BodyFilterConfigPaths struct {
	Path string `yaml:"path"`
	Type string `yaml:"type"`
}
