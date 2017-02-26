package mirror

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func GetHomeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}

func GetMirrorDir() string {
	dir := os.Getenv("MIRROR_HOME")
	if dir == "" {
		dir = filepath.Join(GetHomeDir(), ".mirror.d")
	}
	return dir
}

func GetCADir() string {
	return filepath.Join(GetMirrorDir(), "ca")
}

func GetCertDir() string {
	return filepath.Join(GetMirrorDir(), "certs")
}

func GetUsername() string {
	u := "unknown"
	osUser := ""

	switch runtime.GOOS {
	case "darwin", "linux":
		osUser = os.Getenv("USER")
	case "windows":
		osUser = os.Getenv("USERNAME")
	}

	if osUser != "" {
		u = osUser
	}

	return u
}

// retryable will retry the given function over and over until a
// non-error is returned.
var retryableSleep = 2 * time.Second

func Retryable(f func() error, timeout time.Duration) error {
	startTimeout := time.After(timeout)
	for {
		var err error
		if err = f(); err == nil {
			return nil
		}

		// Create an error and log it
		err = fmt.Errorf("Retryable error: %s", err)
		log.Printf(err.Error())

		// Check if we timed out, otherwise we retry. It is safe to
		// retry since the only error case above is if the command
		// failed to START.
		select {
		case <-startTimeout:
			return err
		default:
			time.Sleep(retryableSleep)
		}
	}
}

func OutputFileContents(file string) (string, error) {
	f, err := ioutil.ReadFile(file)
	if err == nil {
		return string(f), nil
	}
	return "", err

}
