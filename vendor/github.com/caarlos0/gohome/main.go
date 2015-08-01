package gohome

import (
	"os"
	"path/filepath"
	"runtime"
)

// Cache returns the correct folder to store your apps cache, according to the
// spec of each operating system.
func Cache(appName string) string {
	return cacheFor(appName, runtime.GOOS)
}

func cacheFor(appName, os string) string {
	switch os {
	case "darwin":
		return filepath.Join(home(), "Library", "Caches", appName)
	case "windows":
		return appData(appName)
	default:
		return filepath.Join(xdgCache(), appName)
	}
}

// Config returns the correct folder to store your apps configuration,
// according to the spec of each operating system.
func Config(appName string) string {
	return configFor(appName, runtime.GOOS)
}

func configFor(appName, os string) string {
	switch os {
	case "darwin":
		return filepath.Join(home(), "Library", "Application Support", appName)
	case "windows":
		return appData(appName)
	default:
		return filepath.Join(xdgConfig(), appName)
	}
}

func xdgConfig() string {
	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg == "" {
		return filepath.Join(home(), ".config")
	}
	return xdg
}

func xdgCache() string {
	xdg := os.Getenv("XDG_CACHE_HOME")
	if xdg == "" {
		return filepath.Join(home(), ".cache")
	}
	return xdg
}

func appData(appName string) string {
	return filepath.Join(os.Getenv("APPDATA"), appName)
}

func home() string {
	return os.Getenv("HOME")
}
