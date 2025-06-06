/*
 * Bituin (Filipino for "star") - The MicroScript Package Manager
 * Copyright (c) 2025 Cyril John Magayaga
 *
 * It was originally written in Go programming language.
 */

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"
)

const (
	VERSION      = "v0.1.0"
	AUTHOR       = "Cyril John Magayaga"
	PREVIEW_FLAG = "--preview"
)

// Command constants
const (
	NEW         = "new"
	INIT        = "init"
	RUN         = "run"
	ADD         = "add"
	HELP        = "help"
	VERSION_CMD = "version"
	AUTHOR_CMD  = "author"
)

// Templates
const MAIN_MICROSCRIPT = `function main() {
    console.write("Hello, World!");
}

main();`

func getBituinToml(projectName string) string {
	return fmt.Sprintf(`[package]
name = "%s"
main_file = "src/main.microscript"`, projectName)
}

func printUsage() {
	fmt.Println("\033[32mUsage:\033[0m")
	fmt.Println("  \033[34mnew\033[0m [project_name]  - Create a new bituin package in a new directory")
	fmt.Println("  \033[34minit\033[0m [project_name] - Create a new bituin package in an existing directory")
	fmt.Println("  \033[34madd\033[0m [filename]      - Create a new MicroScript source file")
	fmt.Println("  \033[34mrun\033[0m [--preview] [filename] - Run the current project (optionally in preview mode)")
	fmt.Println("\n\033[32mOptions:\033[0m")
	fmt.Println("  \033[34mhelp\033[0m             - Show this help message")
	fmt.Println("  \033[34mversion\033[0m          - Show version information")
	fmt.Println("  \033[34mauthor\033[0m           - Show author information")
}

func createDirectoryStructure(projectPath string) error {
	directories := []string{
		projectPath,
		filepath.Join(projectPath, "src"),
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func createMainMicroscript(projectPath string) error {
	mainPath := filepath.Join(projectPath, "src", "main.microscript")
	return os.WriteFile(mainPath, []byte(MAIN_MICROSCRIPT), 0644)
}

func createBituinConfig(projectPath, projectName string) error {
	configPath := filepath.Join(projectPath, "bituin.toml")
	return os.WriteFile(configPath, []byte(getBituinToml(projectName)), 0644)
}

func addMicroscriptFile(filename string) {
	startTime := time.Now()

	// Ensure we're in a Bituin project directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	bituinTomlPath := filepath.Join(cwd, "bituin.toml")
	if _, err := os.Stat(bituinTomlPath); os.IsNotExist(err) {
		fmt.Println("Error: Not in a Bituin project directory")
		os.Exit(1)
	}

	// Create src directory if it doesn't exist
	srcDir := filepath.Join(cwd, "src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		fmt.Printf("Error creating src directory: %v\n", err)
		os.Exit(1)
	}

	// Create the new file with a basic template
	filePath := filepath.Join(srcDir, filename)
	template := `function main() {
    // Add your code here
}

main();`

	if err := os.WriteFile(filePath, []byte(template), 0644); err != nil {
		fmt.Printf("Error creating MicroScript file: %v\n", err)
		os.Exit(1)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("[%.3fs] Create file: %s\n", elapsed.Seconds(), filename)
}

func runProject(args []string) {
	startTime := time.Now()
	isPreview := false
	var targetFile string

	// Parse arguments
	for i := 1; i < len(args); i++ {
		if args[i] == PREVIEW_FLAG {
			isPreview = true
		} else {
			targetFile = args[i]
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	bituinTomlPath := filepath.Join(cwd, "bituin.toml")
	if _, err := os.Stat(bituinTomlPath); os.IsNotExist(err) {
		fmt.Println("Error: bituin.toml not found. Are you in a bituin project directory?")
		os.Exit(1)
	}

	configContent, err := os.ReadFile(bituinTomlPath)
	if err != nil {
		fmt.Printf("Error reading bituin.toml: %v\n", err)
		os.Exit(1)
	}

	var mainFile, mainFileName string

	// Determine main file
	if targetFile != "" {
		mainFileName = targetFile
		mainFile = filepath.Join(cwd, "src", mainFileName)

		// Update bituin.toml
		re := regexp.MustCompile(`main_file\s*=\s*"([^"]+)"`)
		updatedConfig := re.ReplaceAllString(string(configContent), fmt.Sprintf(`main_file = "src/%s"`, mainFileName))
		if err := os.WriteFile(bituinTomlPath, []byte(updatedConfig), 0644); err != nil {
			fmt.Printf("Error updating bituin.toml: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Use the main_file from bituin.toml
		re := regexp.MustCompile(`main_file\s*=\s*"([^"]+)"`)
		matches := re.FindStringSubmatch(string(configContent))

		if len(matches) == 0 {
			// Default to main.microscript if no main_file is specified
			mainFile = filepath.Join(cwd, "src", "main.microscript")
			mainFileName = "main.microscript"
		} else {
			mainFile = filepath.Join(cwd, matches[1])
			mainFileName = filepath.Base(mainFile)
		}
	}

	if _, err := os.Stat(mainFile); os.IsNotExist(err) {
		fmt.Printf("Error: Main file \"%s\" not found.\n", mainFile)
		os.Exit(1)
	}

	// Look for appropriate executable
	var executableName string
	if isPreview {
		executableName = "microscript-preview.exe"
		fmt.Printf("\033[90mRunning in preview mode: %s\033[0m\n", mainFileName)
	} else {
		executableName = "microscript.exe"
		fmt.Printf("\033[90mRunning: %s\033[0m\n", mainFileName)
	}

	microscriptExe := filepath.Join(cwd, "..", executableName)
	if _, err := os.Stat(microscriptExe); os.IsNotExist(err) {
		fmt.Printf("Error: %s not found in parent directory.\n", executableName)
		os.Exit(1)
	}

	// Execute the MicroScript file
	cmd := exec.Command(microscriptExe, "run", mainFile)
	output, err := cmd.CombinedOutput()

	elapsed := time.Since(startTime)

	if err != nil {
		fmt.Printf("\033[91m[%.3fs] Error: execution failed\033[0m\n", elapsed.Seconds())
		if len(output) > 0 {
			fmt.Print(string(output))
		}
		os.Exit(1)
	}

	mode := "preview"
	if !isPreview {
		mode = "regular"
	}
	fmt.Printf("\033[90m[%.3fs] Success: %s (%s mode)\033[0m\n", elapsed.Seconds(), mainFileName, mode)
	fmt.Println("\033[32m✓\033[0m Project executed successfully!")
	fmt.Println()
	if len(output) > 0 {
		fmt.Print(string(output))
	}
}

func createProject(projectName string, isNew bool) {
	var projectPath string

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if isNew {
		projectPath = filepath.Join(cwd, projectName)
	} else {
		projectPath = cwd
	}

	if isNew {
		if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
			fmt.Printf("Error: Directory \"%s\" already exists.\n", projectName)
			os.Exit(1)
		}
	}

	if err := createDirectoryStructure(projectPath); err != nil {
		fmt.Printf("Error creating directory structure: %v\n", err)
		os.Exit(1)
	}

	if err := createMainMicroscript(projectPath); err != nil {
		fmt.Printf("Error creating main.microscript: %v\n", err)
		os.Exit(1)
	}

	if err := createBituinConfig(projectPath, projectName); err != nil {
		fmt.Printf("Error creating bituin.toml: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Bituin project \"%s\" created successfully!\n", projectName)
	fmt.Println("\nTo get started:")

	if isNew {
		fmt.Printf("  cd %s\n", projectName)
	}

	fmt.Println("  bituin run")
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	command := args[0]

	switch command {
	case HELP:
		printUsage()
	case VERSION_CMD:
		fmt.Println(VERSION)
	case AUTHOR_CMD:
		fmt.Println(AUTHOR)
	case NEW:
		if len(args) < 2 {
			fmt.Println("Error: Project name required for new command")
			printUsage()
			os.Exit(1)
		}
		createProject(args[1], true)
	case INIT:
		if len(args) < 2 {
			fmt.Println("Error: Project name required for init command")
			printUsage()
			os.Exit(1)
		}
		createProject(args[1], false)
	case ADD:
		if len(args) < 2 {
			fmt.Println("Error: File name required for add command")
			printUsage()
			os.Exit(1)
		}
		addMicroscriptFile(args[1])
	case RUN:
		runProject(args)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}
