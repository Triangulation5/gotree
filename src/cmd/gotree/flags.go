package main

import "strings"

type stringListFlag struct {
	values []string
}

func (s *stringListFlag) String() string {
	return strings.Join(s.values, ",")
}

func (s *stringListFlag) Set(value string) error {
	parts := strings.Split(value, ",")
	for _, part := range parts {
		v := strings.TrimSpace(part)
		if v != "" {
			s.values = append(s.values, v)
		}
	}
	return nil
}

func (s *stringListFlag) Values() []string {
	out := make([]string, len(s.values))
	copy(out, s.values)
	return out
}
