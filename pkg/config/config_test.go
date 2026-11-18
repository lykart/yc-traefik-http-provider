package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func clearEnv(t *testing.T) {
	for _, k := range []string{"PORT", "LOG_LEVEL", "SA_KEY_PATH", "CERT_FETCH_INTERVAL_SECONDS", "CERTIFICATE_ID", "CERTIFICATE_IDS"} {
		t.Setenv(k, "")
	}
}

func writeConfig(t *testing.T, content string) string {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestIDs(t *testing.T) {
	cm := &CertificateManager{
		CertificateId:  "a",
		CertificateIds: []string{" b ", "a", "", "c", "b"},
	}
	if got, want := cm.IDs(), []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got := (*CertificateManager)(nil).IDs(); got != nil {
		t.Fatalf("got %v, want nil", got)
	}
}

func TestLoadDefaults(t *testing.T) {
	clearEnv(t)
	cfg, err := load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != 8080 || cfg.LogLevel != "info" || cfg.CertFetchIntervalSeconds != 300 {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
	if cfg.Features.CertificateManager.Enabled || len(cfg.Features.CertificateManager.IDs()) != 0 {
		t.Fatalf("certificate manager must be disabled: %+v", cfg.Features.CertificateManager)
	}
}

func TestLoadYAML(t *testing.T) {
	clearEnv(t)
	dir := writeConfig(t, `
port: 9090
saKeyPath: ./key.json
certFetchIntervalSeconds: 60
features:
  certificateManager:
    enabled: true
    certificateId: old
    certificateIds:
      - one
      - two
`)
	cfg, err := load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != 9090 || cfg.SAKeyPath != "./key.json" || cfg.CertFetchIntervalSeconds != 60 {
		t.Fatalf("unexpected config: %+v", cfg)
	}
	cm := cfg.Features.CertificateManager
	if got, want := cm.IDs(), []string{"old", "one", "two"}; !cm.Enabled || !reflect.DeepEqual(got, want) {
		t.Fatalf("enabled=%v ids=%v, want %v", cm.Enabled, got, want)
	}
}

func TestLoadFeaturesWithoutCertificateManager(t *testing.T) {
	clearEnv(t)
	t.Setenv("CERTIFICATE_ID", "x")
	cfg, err := load(writeConfig(t, "features: {}\n"))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := cfg.Features.CertificateManager.IDs(), []string{"x"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestLoadEnv(t *testing.T) {
	clearEnv(t)
	t.Setenv("PORT", "9000")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("SA_KEY_PATH", "/secrets/key.json")
	t.Setenv("CERT_FETCH_INTERVAL_SECONDS", "30")
	t.Setenv("CERTIFICATE_ID", "single")
	t.Setenv("CERTIFICATE_IDS", "a, b,,single")

	cfg, err := load(writeConfig(t, "port: 9090\nsaKeyPath: ./key.json\n"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != 9000 || cfg.LogLevel != "debug" || cfg.SAKeyPath != "/secrets/key.json" || cfg.CertFetchIntervalSeconds != 30 {
		t.Fatalf("unexpected config: %+v", cfg)
	}
	cm := cfg.Features.CertificateManager
	if got, want := cm.IDs(), []string{"single", "a", "b"}; !cm.Enabled || !reflect.DeepEqual(got, want) {
		t.Fatalf("enabled=%v ids=%v, want %v", cm.Enabled, got, want)
	}
}

func TestLoadInvalid(t *testing.T) {
	for name, env := range map[string][2]string{
		"port":              {"PORT", "abc"},
		"interval":          {"CERT_FETCH_INTERVAL_SECONDS", "5m"},
		"negative interval": {"CERT_FETCH_INTERVAL_SECONDS", "-1"},
	} {
		t.Run(name, func(t *testing.T) {
			clearEnv(t)
			t.Setenv(env[0], env[1])
			if _, err := load(t.TempDir()); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
