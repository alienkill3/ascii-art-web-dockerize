package src

import "strings"

// isKnownFlag checks if an argument is a known flag (must start with --)
func isKnownFlag(arg string) bool {
	// Only recognize flags starting with --
	if !strings.HasPrefix(arg, "--") {
		return false
	}

	// Remove leading -- and check the flag name
	flagName := strings.TrimPrefix(arg, "--")

	// Handle flags with = separator
	if idx := strings.Index(flagName, "="); idx > 0 {
		flagName = flagName[:idx]
	}

	// List of known flag names
	knownFlags := map[string]bool{
		"color":  true,
		"align":  true,
		"s":      true,
		"substr": true,
		"output": true,
	}

	return knownFlags[flagName]
}

// isFlagName checks if a string is a known flag name (without dashes)
func isFlagName(name string) bool {
	knownFlags := map[string]bool{
		"color":  true,
		"align":  true,
		"s":      true,
		"substr": true,
		"output": true,
	}
	return knownFlags[name]
}

// needsFlagValue checks if a flag requires a value on the next argument
func needsFlagValue(arg string) bool {
	flags := map[string]bool{

		"--color": true,

		"--align": true,

		"--s": true,

		"--substr": true,

		"--output": true,
	}
	return flags[arg]
}
