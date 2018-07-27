package policy

import (
	"testing"

	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestGlobPattern_Matches(t *testing.T) {
	for _, tt := range []struct {
		name    string
		pattern string
		true    []string
		false   []string
	}{
		{
			name:    "all",
			pattern: "*",
			true:    []string{"", "1", "foo"},
			false:   nil,
		},
		{
			name:    "all prefixed",
			pattern: "glob:*",
			true:    []string{"", "1", "foo"},
			false:   nil,
		},
		{
			name:    "prefix",
			pattern: "master-*",
			true:    []string{"master-", "master-foo"},
			false:   []string{"", "foo-master"},
		},
	} {
		pattern := NewPattern(tt.pattern)
		assert.IsType(t, GlobPattern(""), pattern)
		t.Run(tt.name, func(t *testing.T) {
			for _, tag := range tt.true {
				assert.True(t, pattern.Matches(tag))
			}
			for _, tag := range tt.false {
				assert.False(t, pattern.Matches(tag))
			}
		})
	}
}

func TestSemverPattern_Matches(t *testing.T) {
	for _, tt := range []struct {
		name    string
		pattern string
		true    []string
		false   []string
	}{
		{
			name:    "all",
			pattern: "semver:*",
			true:    []string{"1", "1.0", "2.0.1-alpha.1"},
			false:   []string{"", "latest"},
		},
		{
			name:    "semver",
			pattern: "semver:~1",
			true:    []string{"v1", "1", "1.2", "1.2.3"},
			// FIXME(rndstr): should we always match latest?
			false:   []string{"", "latest", "2.0.0", ""},
		},
	} {
		pattern := NewPattern(tt.pattern)
		assert.IsType(t, SemverPattern{}, pattern)
		for _, tag := range tt.true {
			t.Run(fmt.Sprintf("%s[%q]", tt.name, tag), func(t *testing.T) {
				assert.True(t, pattern.Matches(tag))
			})
		}
		for _, tag := range tt.false {
			t.Run(fmt.Sprintf("%s[%q]", tt.name, tag), func(t *testing.T) {
				assert.False(t, pattern.Matches(tag))
			})
		}
	}
}
