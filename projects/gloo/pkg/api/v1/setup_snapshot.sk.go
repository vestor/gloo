// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"fmt"

	"github.com/solo-io/solo-kit/pkg/utils/hashutils"
	"go.uber.org/zap"
)

type SetupSnapshot struct {
	Settings SettingsByNamespace
}

func (s SetupSnapshot) Clone() SetupSnapshot {
	return SetupSnapshot{
		Settings: s.Settings.Clone(),
	}
}

func (s SetupSnapshot) Hash() uint64 {
	return hashutils.HashAll(
		s.hashSettings(),
	)
}

func (s SetupSnapshot) hashSettings() uint64 {
	return hashutils.HashAll(s.Settings.List().AsInterfaces()...)
}

func (s SetupSnapshot) HashFields() []zap.Field {
	var fields []zap.Field
	fields = append(fields, zap.Uint64("settings", s.hashSettings()))

	return append(fields, zap.Uint64("snapshotHash", s.Hash()))
}

type SetupSnapshotStringer struct {
	Version  uint64
	Settings []string
}

func (ss SetupSnapshotStringer) String() string {
	s := fmt.Sprintf("SetupSnapshot %v\n", ss.Version)

	s += fmt.Sprintf("  Settings %v\n", len(ss.Settings))
	for _, name := range ss.Settings {
		s += fmt.Sprintf("    %v\n", name)
	}

	return s
}

func (s SetupSnapshot) Stringer() SetupSnapshotStringer {
	return SetupSnapshotStringer{
		Version:  s.Hash(),
		Settings: s.Settings.List().NamespacesDotNames(),
	}
}
