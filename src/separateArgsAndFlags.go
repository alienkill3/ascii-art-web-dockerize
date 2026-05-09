package src

import "strings"

// separateArgsAndFlags separates recognized flags from positional arguments.
// Only arguments with double-dash prefix (--) are treated as flags.
// Any argument starting with single dash (-) returns an error requiring double-dash format,
// unless the -- separator has been used to explicitly end flag parsing.
func separateArgsAndFlags(args []string) ([]string, []string, error) {
	var flagArgs []string
	var positionalArgs []string
	foundPositional := false
	foundExplicitSeparator := false // Track if we saw explicit -- separator

	// If only 1 argument and it looks like a flag, treat it as positional input
	if len(args) == 1 && strings.HasPrefix(args[0], "--") {
		return flagArgs, args, nil
	} else if len(args) == 1 && strings.HasPrefix(args[0], "-") {
		return nil, args, nil
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if foundPositional {
			// Once a positional arg is found, check if remaining args are single-dash style
			// But only if they match a known flag name
			if !foundExplicitSeparator && strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") {
				flagName := strings.TrimPrefix(arg, "-")
				// Handle flags with = separator
				if idx := strings.Index(flagName, "="); idx > 0 {
					flagName = flagName[:idx]
				}
				// If it matches a known flag name, it's an error (invalid syntax)
				if isFlagName(flagName) {
					return nil, nil, ErrUsage
				}
			}
			positionalArgs = append(positionalArgs, arg)
			continue
		}

		// Check for double dash which marks end of flags
		if arg == "--" && !foundExplicitSeparator {
			// Don't include the -- in flagArgs, just mark end of flags
			foundPositional = true
			foundExplicitSeparator = true // Mark that we explicitly saw --
			continue
		} else {
			// If we see -- again after we've already seen it, treat it as positional
			if arg == "--" && foundExplicitSeparator {
				positionalArgs = append(positionalArgs, arg)
				continue
			}
		}

		// If it starts with --, check if it's a known flag
		if strings.HasPrefix(arg, "--") {
			// Check if it looks like a known flag (color, align, s, substr, output)
			if isKnownFlag(arg) {
				flagArgs = append(flagArgs, arg)
				// If it's a flag with = (like --color=red), it's self-contained
				// If it's a flag without = but expects a value, get the next arg
				if !strings.Contains(arg, "=") && needsFlagValue(arg) {
					// i++ // Skip next iteration since we'll add it to flags
					// if i < len(args) {
					// 	flagArgs = append(flagArgs, args[i])
					// }
					return nil, nil, ErrUsage
				}
			} else {
				// Unknown flag-like argument, treat as positional
				foundPositional = true
				positionalArgs = append(positionalArgs, arg)
			}
		} else if strings.HasPrefix(arg, "-") {
			// Single dash argument: check if it matches a known flag name
			flagName := strings.TrimPrefix(arg, "-")
			// Handle flags with = separator
			if idx := strings.Index(flagName, "="); idx > 0 {
				flagName = flagName[:idx]
			}
			// If it matches a known flag name, it's an error (invalid syntax)
			if isFlagName(flagName) {
				return nil, nil, ErrUsage
			}
			// Unknown single-dash argument, treat as positional
			foundPositional = true
			positionalArgs = append(positionalArgs, arg)
		} else {
			// Doesn't start with -, it's positional
			foundPositional = true
			positionalArgs = append(positionalArgs, arg)
		}
	}

	return flagArgs, positionalArgs, nil
}
