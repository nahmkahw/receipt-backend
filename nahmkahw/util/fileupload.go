package util

import (
    "os"
)

func CreateUploadsDir(path string) error {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return os.MkdirAll(path, os.ModePerm)
    }
    return nil
}
