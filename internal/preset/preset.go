package preset

import (
	"fmt"
	"strings"
)

// Preset represents a prompt preset
type Preset struct {
	Name        string
	Description string
	Template    string
}

var defaultPresets = map[string]Preset{
	"eli5": {
		Name:        "ELI5",
		Description: "Explain Like I'm 5 - Simple, beginner-friendly explanations",
		Template:    "Explain this in simple terms that a 5-year-old could understand. Use everyday language and avoid jargon.",
	},
	"intermediate": {
		Name:        "Intermediate",
		Description: "Standard explanations for intermediate learners",
		Template:    "Provide a clear, structured explanation suitable for intermediate learners. Include examples and context.",
	},
	"advanced": {
		Name:        "Advanced",
		Description: "Detailed explanations for advanced learners",
		Template:    "Provide a comprehensive, detailed explanation with technical accuracy. Include nuances and advanced concepts.",
	},
}

// GetPreset returns a preset by name
func GetPreset(name string) (*Preset, error) {
	if name == "" {
		preset := defaultPresets["intermediate"]
		return &preset, nil
	}

	key := strings.ToLower(strings.TrimSpace(name))
	preset, ok := defaultPresets[key]
	if !ok {
		return nil, fmt.Errorf("preset '%s' not found", name)
	}

	return &preset, nil
}

// ListPresets returns all available presets
func ListPresets() []Preset {
	presets := make([]Preset, 0, len(defaultPresets))
	for _, preset := range defaultPresets {
		presets = append(presets, preset)
	}
	return presets
}

// FormatPrompt formats a prompt using the preset template
func FormatPrompt(preset *Preset, userPrompt string) string {
	if preset == nil {
		return userPrompt
	}
	return fmt.Sprintf("%s\n\nUser question: %s", preset.Template, userPrompt)
}
