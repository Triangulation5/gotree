package main

import (
	"os"
	"testing"
)

func TestResolveSortMode(t *testing.T) {
	tests := []struct {
		name      string
		sortBy    string
		oldSorted bool
		oldExt    bool
		want      string
		wantErr   bool
	}{
		{name: "default", want: "name"},
		{name: "legacy sorted", oldSorted: true, want: "name"},
		{name: "legacy ext", oldExt: true, want: "ext"},
		{name: "explicit mtime", sortBy: "mtime", want: "mtime"},
		{name: "invalid", sortBy: "wat", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := resolveSortMode(tc.sortBy, tc.oldSorted, tc.oldExt)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("expected %s, got %s", tc.want, got)
			}
		})
	}
}

func TestResolveColorEnabled(t *testing.T) {
	original, hadOriginal := os.LookupEnv("NO_COLOR")
	t.Cleanup(func() {
		if hadOriginal {
			_ = os.Setenv("NO_COLOR", original)
			return
		}
		_ = os.Unsetenv("NO_COLOR")
	})

	if err := os.Unsetenv("NO_COLOR"); err != nil {
		t.Fatalf("unset NO_COLOR: %v", err)
	}

	enabled, err := resolveColorEnabled("auto")
	if err != nil || !enabled {
		t.Fatalf("expected auto to enable color without NO_COLOR, got enabled=%v err=%v", enabled, err)
	}

	if err := os.Setenv("NO_COLOR", "1"); err != nil {
		t.Fatalf("set NO_COLOR: %v", err)
	}
	enabled, err = resolveColorEnabled("auto")
	if err != nil || enabled {
		t.Fatalf("expected auto to disable color with NO_COLOR, got enabled=%v err=%v", enabled, err)
	}

	enabled, err = resolveColorEnabled("mono")
	if err != nil || enabled {
		t.Fatalf("expected mono to disable color, got enabled=%v err=%v", enabled, err)
	}

	enabled, err = resolveColorEnabled("color")
	if err != nil || !enabled {
		t.Fatalf("expected color to force enable, got enabled=%v err=%v", enabled, err)
	}
}
