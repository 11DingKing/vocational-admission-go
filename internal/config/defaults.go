package config

func DefaultHeaders() map[string]string {
	return map[string]string{"Content-Type": "application/json", "Cache-Control": "no-store"}
}
func SupportedTimezones() []string { return []string{"Asia/Shanghai", "Asia/Urumqi", "UTC"} }
func IsProduction(c Config) bool   { return c.HTTPAddr != ":0" && c.DBPath != "file::memory:" }
func (c Config) WorkerInterval() int {
	if c.WorkerIntervalSeconds <= 0 {
		return 5
	}
	return c.WorkerIntervalSeconds
}
func (c Config) SessionTTL() int {
	if c.SessionTTLSeconds <= 0 {
		return 3600
	}
	return c.SessionTTLSeconds
}
