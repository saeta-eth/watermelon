package violations

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// LogPath returns the violations log path for a project
func LogPath(projectDir string) string {
	return filepath.Join(projectDir, ".watermelon", "violations.log")
}

// Read returns recent violations for a project
func Read(projectDir string) ([]string, error) {
	logPath := LogPath(projectDir)
	file, err := os.Open(logPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Clear removes the violations log
func Clear(projectDir string) error {
	logPath := LogPath(projectDir)
	if err := os.Remove(logPath); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	fmt.Println("Violations log cleared")
	return nil
}
