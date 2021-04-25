package single

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var socketFile string

// Start start single
func Start(socketFile string) (started bool) {
	c, _ := net.Dial("unix", socketFile)
	if c != nil {
		return true
	}
	go func() {
		os.Remove(socketFile)
		l, err := net.Listen("unix", socketFile)
		if err != nil {
			panic(err)
		}
		defer l.Close()

		http.Serve(l, nil)
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Fatal("accept error:", err)
			} else {
				fmt.Println(conn)
				conn.Close()
			}
		}
	}()
	return false
}

// StartAuto auto start
func StartAuto() bool {
	h := md5.New()
	execName := os.Args[0]
	dir := filepath.Dir(execName)
	dirAbs, err := filepath.Abs(dir)
	if err == nil && dirAbs != dir {
		execName = filepath.Join(dirAbs, execName)
	}
	h.Write([]byte(execName))
	filename := hex.EncodeToString(h.Sum(nil))

	p, e := home()
	if e != nil {
		p = "./"
	}
	filepath := path.Join(p, "."+path.Base(execName)+"-"+filename+".sock")
	return Start(filepath)
}

func home() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support

	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}
