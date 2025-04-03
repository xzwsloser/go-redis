package config

import "testing"

func TestSetUpConfig(t *testing.T) {
	SetupConfig("../redis.conf")
	t.Log(*Properties)
}
