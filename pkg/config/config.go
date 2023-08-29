package config

type Config struct {
	Middlewares Middlewares `yaml:"middlewares"`
}

type Middlewares struct {
	BodyFilter MiddlewareConfig[BodyFilterConfig] `validate:"omitempty" yaml:"bodyFilter,omitempty"` //nolint:tagliatelle,lll // valid tag
}

type MiddlewareConfig[T any] struct {
	Enabled bool `yaml:"enabled"`
	Config  []T  `validate:"required_if=Enabled true,dive,required" yaml:"config"`
}

type BodyFilterConfig struct {
	Paths   []BodyFilterConfigPaths `validate:"required,gt=0,dive,gt=0"           yaml:"paths"`
	Methods []string                `validate:"required,gt=0,dive,gt=0,uppercase" yaml:"methods"`
	Filter  string                  `validate:"required"                          yaml:"filter"`
}

type BodyFilterConfigPaths struct {
	Path string `yaml:"path"`
	Type string `validate:"oneof=glob" yaml:"type"`
}
