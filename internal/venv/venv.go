package venv


import (
	"fmt"
	"bytes"
	"time"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"github.com/LeonSideln1kov/viper/internal/python"
)


const (
    venvDir    = ".venv"
    minPython  = "3.10"
)


func CreateVenv() {
	pythonPath, err := python.FindPython()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Printf("VIPER requires Python %s+ to create virtual environments\n", minPython)
		os.Exit(1)
	}

	if _, err_check := os.Stat(venvDir); os.IsNotExist(err_check) {
		cmd := exec.Command(pythonPath, "-m", "venv", venvDir)
	
		output, err_exec := cmd.CombinedOutput()
		if err_exec != nil {
			fmt.Printf("Failed to create virtual environment:\n")
            fmt.Printf("Command: %s\n", cmd.String())
            fmt.Printf("Output:\n%s\n", output)
            fmt.Printf("Error Details: %v\n", err)
            os.Exit(1)
		}
		fmt.Println("Virtual environment created successfully")
	} else {
		fmt.Printf("Directory %s already exists\n", venvDir)
	}
}


func PipPath() (string, error) {
	var path string

    if runtime.GOOS == "windows" {
        path = filepath.Join(venvDir, "Scripts", "pip.exe") 
    } else {
		path = filepath.Join(venvDir, "bin", "pip")
	}
    
	if _, err := os.Stat(path); os.IsNotExist(err) {
        return "", fmt.Errorf("pip not found in virtual environment")
    }

	return path, nil
}


func PythonPath() (string, error) {
	var path string

	switch runtime.GOOS {
	case "windows":
		path = filepath.Join(venvDir, "Scripts", "python.exe")
	default: // Unix-like systems (linux, darwin, etc.)
		path = filepath.Join(venvDir, "bin", "python")
	}

	// Verify the file actually exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("virtual environment not found at %s", path)
	}

	return path, nil
}


func InstallPackage(pkg string) error {
    pipPath, err := PipPath()
    if err != nil {
        return fmt.Errorf("venv pip missing: %w", err)
    }

    cmd := exec.Command(pipPath, "install", pkg)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	
	start := time.Now()
	err = cmd.Run()
	duration := time.Since(start).Round(time.Millisecond)
	if err != nil {
		fmt.Printf("🚨 Installation failed for %s (after %s)\n", pkg, duration)
		fmt.Printf("=== STDOUT ===\n%s\n", outBuf.String())
		fmt.Printf("=== STDERR ===\n%s\n", errBuf.String())
		return fmt.Errorf("pip install failed: %w", err)
	}

	fmt.Printf("✅ %s installed in %s\n", pkg , duration)
	return nil
}
