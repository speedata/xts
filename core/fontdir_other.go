//go:build !windows && !darwin

package core

// FontFolder
func FontFolder() (string, error) {
	return "", nil
}
