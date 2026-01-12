package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type yamlConfig struct {
	Developers []yamlDeveloper `yaml:"developers"`
}

type yamlDeveloper struct {
	Email string    `yaml:"email"`
	Name  string    `yaml:"name"`
	Keys  []yamlKey `yaml:"keys"`
}

type yamlKey struct {
	ID string `yaml:"id"`
}

// Policy note:
//   - For easy CI configuration, we treat config entries as "active now".
//   - added_at defaults to Unix epoch so historical policy can validate old commits
//     *as long as you don't rely on commit dates for security-critical gating*.
func SyncFromYAML(path string, reset bool) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var cfg yamlConfig
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return fmt.Errorf("parse yaml: %w", err)
	}

	if reset {
		if err := ResetAll(); err != nil {
			return err
		}
	}

	epoch := sql.NullTime{Time: time.Unix(0, 0).UTC(), Valid: true}
	none := sql.NullTime{Valid: false}

	for _, d := range cfg.Developers {
		if d.Email == "" {
			return fmt.Errorf("developer missing email")
		}
		if err := UpsertDeveloper(d.Email, d.Name, epoch, none); err != nil {
			return err
		}
		for _, k := range d.Keys {
			if k.ID == "" {
				return fmt.Errorf("developer %s has key with empty id", d.Email)
			}
			if err := InsertDeveloperKey(d.Email, k.ID, epoch, none); err != nil {
				return err
			}
		}
	}

	return nil
}
