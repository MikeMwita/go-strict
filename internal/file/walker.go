package file

import (
	"io/fs"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

// Walker gets all files in the filterDir and directories below.
// You can:
//   - filter files like "*.go"
//   - give a max. directory depth
//     -1 = unrestricted deep
//     0 = only filterDir
//     1 = filterDir and 1 deeper
//     etc.
func Walker(filterDir, locFilter string, locDepth int) ([]string, error) {
	var err error
	fileList := []string{}
	filterDir = strings.Replace(filterDir, "\\", "/", -1) // for windows
	if !DirExists(filterDir) {
		return fileList, err
	}
	fsys := os.DirFS(filterDir)
	fs.WalkDir(fsys, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		_, fileName := path.Split(filePath)
		matched, err := filepath.Match(locFilter, fileName)
		if err != nil || !matched {
			return nil
		}
		filePath = strings.ReplaceAll(filePath, "\\", "/") // for windows
		if strings.Contains(filePath, "/vendor/") {
			return nil
		}
		if locDepth >= 0 {
			slashes := strings.Count(filePath, "/")
			if slashes > locDepth {
				return nil
			}
		}
		fileList = append(fileList, filterDir+"/"+filePath)
		return nil
	})
	return fileList, err
}

// TextWriter writes text content to a file.
func TextWriter(content, filePath string) error {
	err := ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

// RootDir gets the root directory.
func RootDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

// DirExists checks if a directory exists.
func DirExists(dirPath string) bool {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// Benchmark function
func BenchmarkWalker(b *testing.B) {
	rootDir, _ := RootDir()
	// create some files first
	testDir := rootDir + "/file/temp"

	defer os.RemoveAll(testDir)

	for i := 0; i < 10; i++ {
		_ = createTestFile(testDir)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ { // use b.N for looping
		_, _ = Walker(testDir, "*.txt", 0)
	}
}

// Helper function
func createTestFile(testDir string) error {
	rsuffix := randomString(5)
	testFile := testDir + "/testfile_" + rsuffix + ".txt"
	err := TextWriter("Marcello Holland", testFile)
	return err
}

// Other helper functions
func randomString(length int) string {
	bytes := make([]byte, length)

	for i := 0; i < length; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}

	return string(bytes)
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}
