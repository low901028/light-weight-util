package netutil

import (
	"fmt"
	"runtime"
)

func GetDefaultHost() (string, error){
	return "", fmt.Errorf("default host not supported on %s_%s", runtime.GOOS, runtime.GOARCH)
}

func GetDefaultInterfaces() (map[string]uint8, error) {
	return nil, fmt.Errorf("default host not supported on %s_%s", runtime.GOOS, runtime.GOARCH)
}
