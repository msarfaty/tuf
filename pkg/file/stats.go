package file

import "os"

// whether or not the file specified is empty
func FileIsEmpty(name string) (bool, error) {
	stats, err := os.Stat(name)
	if err != nil {
		return false, err
	}

	return stats.Size() == 0, nil
}
