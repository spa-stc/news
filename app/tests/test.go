package tests

import (
	"net"
	"strconv"
	"testing"

	"github.com/spa-stc/newsletter/server/profile"
)

func GetTestingProfile(t *testing.T) *profile.Profile {
	dir := t.TempDir()
	mode := "development"
	port := getUnusedPort()
	driver := "sqlite"

	return &profile.Profile{
		Dir:    dir,
		Env:    mode,
		Port:   strconv.Itoa(port),
		Driver: driver,
		DSN:    "newsletter.db",
	}
}

func getUnusedPort() int {
	// Get a random unused port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	// Get the port number
	port := listener.Addr().(*net.TCPAddr).Port
	return port
}
