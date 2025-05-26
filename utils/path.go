package utils

import "os"

func EnsureFolderExists(path string) error {
	if err := PathExists(path); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func PathExists(path string) error {
	_, err := os.Stat(path)
	return err
}
