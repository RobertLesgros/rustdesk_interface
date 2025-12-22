package config

type Logger struct {
	Path         string
	Level        string
	ReportCaller bool `mapstructure:"report-caller"`
}

// Audit configuration for structured security audit logging
type Audit struct {
	Enabled  bool   `mapstructure:"enabled"`
	FilePath string `mapstructure:"file-path"`
}
