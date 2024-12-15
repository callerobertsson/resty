// Package utils editor function
package utils

import (
	"os"
	"os/exec"
)

// editFile opens editor on filePath for editing. File content text is returned.
func EditFile(filePath, editor string) (string, error) {
	// Look up editor path
	cmd, err := exec.LookPath(editor)
	if err != nil {
		return "", err
	}

	// Define process attributes
	var attr os.ProcAttr
	attr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	// Start the editor
	p, err := os.StartProcess(cmd, []string{cmd, filePath}, &attr)
	if err != nil {
		return "", err
	}
	_, _ = p.Wait()

	// Read file content
	text, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(text), nil
}
