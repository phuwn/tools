package util

import "strings"

// FormatFuncName - short function name
func FormatFuncName(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}

// CleanPath - make a clean path from the normal path with format `folder/file.go`
func CleanPath(path string) string {
	if path == "" {
		return "."
	}
	// Strip trailing slashes.
	for len(path) > 0 && path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}
	// Find the last element
	var lastSlash, lastLeg int
	for i := range path {
		if path[i] == '/' {
			lastLeg = lastSlash
			lastSlash = i
		}
	}
	if lastLeg != 0 {
		return path[lastLeg+1:]
	}
	path = path[lastSlash+1:]
	// If empty, it had only slashes.
	if path == "" {
		return "/"
	}
	return path
}
