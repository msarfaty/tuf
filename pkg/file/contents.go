package file

import (
	"io"
	"os"

	"github.com/google/go-cmp/cmp"
)

// determines if a file ends with a specific character
func FileEndsWith(name string, endsWith string) (bool, error) {
	file, err := os.Open(name)
	if err != nil {
		return false, err
	}

	fi, err := os.Stat(name)
	if err != nil {
		return false, err
	}
	ewb := []byte(endsWith)

	if fi.Size() < int64(len(ewb)) {
		return false, nil
	}
	file.Seek(int64(-1*len(ewb)), io.SeekEnd)
	cmp := make([]byte, len(ewb))
	_, err = file.Read(cmp)
	if err != nil {
		return false, err
	}
	return string(ewb) == string(cmp), nil
}

// Starting at or before a given position, deletes any occurrences of some bytes over N
func DeleteOverNOccurrences(content []byte, chars []byte, start int, occurrences int) []byte {
	length := len(chars)

	if start >= len(content) || len(chars) > len(content) {
		return content
	}

	// backtrack to first occurrence of chars from start pos
	found := false
	for range length {
		if start < 0 {
			return content
		}
		if start+length > len(content) {
			start -= 1
			continue
		}
		if cmp.Equal(content[start:start+length], chars) {
			found = true
			break
		}
		start -= 1
	}
	if !found {
		return content
	}

	// continuously backtrack until hit start of content or no more occurrences of chars
	for start >= 0 && cmp.Equal(content[start:start+length], chars) {
		start -= length
	}
	start += length

	// progress through n occurrences
	for range occurrences {
		if start < 0 {
			return content
		}
		if start+length >= len(content) {
			return content
		}
		if !cmp.Equal(content[start:start+length], chars) {
			return content
		}
		start += length
	}

	// track deletion window
	dt := start
	for {
		if dt > len(content) {
			return content
		}
		if dt+length > len(content) {
			break
		}
		if !cmp.Equal(content[dt:dt+length], chars) {
			break
		}
		dt += length
	}

	return append(content[:start], content[dt:]...)
}
