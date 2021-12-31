package telemetry_v2

type TelemetryConfig struct {
	Id      string `yaml:"id"`
	Source  string `yaml:"source"`
	Query   string `yaml:"query"`
	Summary string `yaml:"summary"`
	Data    []struct {
		Name   string `yaml:"metric_name"`
		Label  string `yaml:"label"`
		Value  string `yaml:"value"`
		Column string `yaml:"column"`
	}
}
