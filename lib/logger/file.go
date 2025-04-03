package logger

import (
	"fmt"
	"os"
)

func checkNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

func checkPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

func mkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func isNotExistMkDir(src string) error {
	if notExist := checkNotExist(src); notExist {
		if err := mkDir(src); err != nil {
			return err
		}
	}

	return nil
}

func mustOpen(filename, dir string) (*os.File, error) {
	perm := checkPermission(dir)
	if perm {
		return nil, fmt.Errorf("permission denied dir: %v", dir)
	}

	err := isNotExistMkDir(dir)
	if err != nil {
		return nil, fmt.Errorf("erro during make dir: %v", dir)
	}

	f, err := os.OpenFile(dir+string(os.PathSeparator)+filename,
		os.O_APPEND|os.O_RDWR|os.O_CREATE,
		0644)
	if err != nil {
		return nil, fmt.Errorf("fail to open file , err: %v", err)
	}

	return f, nil
}
