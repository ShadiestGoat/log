package log

// Doesn't log anything, but prints in the pretty debug color
// Doesn't include a prefix.
func PrintDebug(msg string, args ...any) {
	levelInfo[LL_DEBUG].Color.Printf(msg+"\n", args...)
}

// Doesn't log anything, but prints in the pretty success color
// Doesn't include a prefix.
func PrintSuccess(msg string, args ...any) {
	levelInfo[LL_SUCCESS].Color.Printf(msg+"\n", args...)
}

// Doesn't log anything, but prints in the pretty warning color
// Doesn't include a prefix.
func PrintWarn(msg string, args ...any) {
	levelInfo[LL_WARN].Color.Printf(msg+"\n", args...)
}
