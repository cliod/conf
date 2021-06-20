package conf

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// RootPath return project root path
func RootPath() (path string) {
	// default linux/mac os
	var (
		sp = "/"
		ss []string
	)
	if runtime.GOOS == "windows" {
		sp = "\\"
	}

	stdout, _ := exec.Command("go", "env", "GOMOD").Output()
	path = string(bytes.TrimSpace(stdout))
	if path != "" {
		ss = strings.Split(path, sp)
		ss = ss[:len(ss)-1]
		path = strings.Join(ss, sp) + sp
		return
	}

	// GOPATH
	fileDir, _ := os.Getwd()
	path = os.Getenv("GOPATH") // < go 1.17 use
	ss = strings.Split(fileDir, path)
	if path != "" {
		ss2 := strings.Split(ss[1], sp)
		path += sp
		for i := 1; i < len(ss2); i++ {
			path += ss2[i] + sp
			if Exists(path) {
				return path
			}
		}
	}
	return
}

// Exists file/dir exit
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
