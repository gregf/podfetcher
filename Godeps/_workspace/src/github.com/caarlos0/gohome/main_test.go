package gohome

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	appName = "MyApp"
	appdata = "C:/Users/john-doe/AppData/Local/"
)

func TestCache(t *testing.T) {
	cache := Cache(appName)
	assert.NotEmpty(t, cache)
	assert.True(t, strings.HasSuffix(cache, appName))
}

func TestCacheDarwin(t *testing.T) {
	cache := cacheFor(appName, "darwin")
	assert.NotEmpty(t, cache)
	assert.Equal(t, home()+"/Library/Caches/"+appName, cache)
}

func TestCacheWindows(t *testing.T) {
	os.Setenv("APPDATA", appdata)
	defer os.Unsetenv("APPDATA")
	cache := cacheFor(appName, "windows")
	assert.NotEmpty(t, cache)
	assert.Equal(t, appdata+appName, cache)
}

func TestCacheXdgUnset(t *testing.T) {
	cache := cacheFor(appName, "linux")
	assert.NotEmpty(t, cache)
	assert.Equal(t, home()+"/.cache/"+appName, cache)
}

func TestCacheXdgSet(t *testing.T) {
	xdgCache := os.Getenv("HOME") + "/.xdgcache/"
	os.Setenv("XDG_CACHE_HOME", xdgCache)
	defer os.Unsetenv("XDG_CACHE_HOME")
	cache := cacheFor(appName, "linux")
	assert.NotEmpty(t, cache)
	assert.Equal(t, xdgCache+appName, cache)
}

func TestConfigDarwin(t *testing.T) {
	config := configFor(appName, "darwin")
	assert.NotEmpty(t, config)
	assert.Equal(t, home()+"/Library/Application Support/"+appName, config)
}

func TestConfigWindows(t *testing.T) {
	os.Setenv("APPDATA", appdata)
	defer os.Unsetenv("APPDATA")
	config := configFor(appName, "windows")
	assert.NotEmpty(t, config)
	assert.Equal(t, appdata+appName, config)
}

func TestConfigXdgUnset(t *testing.T) {
	config := configFor(appName, "linux")
	assert.NotEmpty(t, config)
	assert.Equal(t, home()+"/.config/"+appName, config)
}

func TestConfigXdgSet(t *testing.T) {
	xdgConfig := os.Getenv("HOME") + "/.xdgconfig/"
	os.Setenv("XDG_CONFIG_HOME", xdgConfig)
	defer os.Unsetenv("XDG_CONFIG_HOME")
	config := configFor(appName, "linux")
	assert.NotEmpty(t, config)
	assert.Equal(t, xdgConfig+appName, config)
}

func TestConfig(t *testing.T) {
	config := Config(appName)
	assert.NotEmpty(t, config)
	assert.True(t, strings.HasSuffix(config, appName))
}
