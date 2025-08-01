package files

import "os"

func IsFile(path string) (ok bool) {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func IsDir(path string) (ok bool) {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
