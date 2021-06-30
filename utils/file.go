package utils

import "os"

func IsExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
