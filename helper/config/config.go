package config

import "path/filepath"

// Config is the basic configuration object
type Config struct {
	basedir    string
	Libdir     string
	Xtsversion string
}

// SetBasedir sets the root of the speedata xts source files
func (cfg *Config) SetBasedir(basedir string) {
	cfg.basedir = basedir
	cfg.Libdir = filepath.Join(basedir, "lib")

}

// Basedir returns the root of the speedata xts source files
func (cfg *Config) Basedir() string {
	return cfg.basedir
}

// NewConfig creates a new configuration struct
func NewConfig(basedir, version string) *Config {
	cfg := &Config{}
	cfg.SetBasedir(basedir)
	cfg.Xtsversion = version
	return cfg
}
