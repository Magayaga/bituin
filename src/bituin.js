/*
 * Bituin (Filipino for "star") - The MicroScript Package Manager
 * Copyright (c) 2025 Cyril John Magayaga
 * 
 * It was originally written in JavaScript programming language.
 * 
*/
import fs from "fs";
import path from "path";
import { exec } from "child_process";

const COMMANDS = {
  NEW: "new",
  INIT: "init",
  RUN: "run",
  ADD: "add",
  HELP: "help",
  VERSION: "version",
  AUTHOR: "author"
};

const VERSION = "v0.1.0";
const AUTHOR = "Cyril John Magayaga";

const TEMPLATES = {
  MAIN_MICROSCRIPT: `function main() {
    console.write("Hello, World!");
}

main();`,
  BITUIN_TOML: (projectName) => `[package]
name = "${projectName}"
main_file = "src/main.microscript"`
};

function printUsage() {
  console.log("\x1b[32mUsage:\x1b[0m");
  console.log(`  \x1b[34mnew\x1b[0m [project_name]  - Create a new bituin package in a new directory`);
  console.log(`  \x1b[34minit\x1b[0m [project_name] - Create a new bituin package in an existing directory`);
  console.log(`  \x1b[34madd\x1b[0m [filename]      - Create a new MicroScript source file`);
  console.log(`  \x1b[34mrun\x1b[0m                 - Run the current project`);
  console.log("\n\x1b[32mOptions:\x1b[0m");
  console.log(`  \x1b[34mhelp\x1b[0m             - Show this help message`);
  console.log(`  \x1b[34mversion\x1b[0m          - Show version information`);
  console.log(`  \x1b[34mauthor\x1b[0m           - Show author information`);
}

function createDirectoryStructure(projectPath) {
  const directories = [
    projectPath,
    path.join(projectPath, "src"),
  ];

  for (const dir of directories) {
    if (!fs.existsSync(dir)) {
      fs.mkdirSync(dir, { recursive: true });
    }
  }
}

function createMainMicroscript(projectPath) {
  fs.writeFileSync(
    path.join(projectPath, "src", "main.microscript"),
    TEMPLATES.MAIN_MICROSCRIPT
  );
}

function createBituinConfig(projectPath, projectName) {
  fs.writeFileSync(
    path.join(projectPath, "bituin.toml"),
    TEMPLATES.BITUIN_TOML(projectName)
  );
}

function addMicroscriptFile(filename) {
  const startTime = process.hrtime();
  
  try {
    // Ensure we're in a Bituin project directory
    if (!fs.existsSync(path.join(process.cwd(), "bituin.toml"))) {
      console.error("Error: Not in a Bituin project directory");
      process.exit(1);
    }

    // Create src directory if it doesn't exist
    const srcDir = path.join(process.cwd(), "src");
    if (!fs.existsSync(srcDir)) {
      fs.mkdirSync(srcDir, { recursive: true });
    }

    // Create the new file with a basic template
    const filePath = path.join(srcDir, filename);
    const template = `function main() {
    // Add your code here
}

main();`;

    fs.writeFileSync(filePath, template);

    const timeDiff = process.hrtime(startTime);
    const timeInSeconds = (timeDiff[0] + timeDiff[1] / 1e9).toFixed(3);
    
    console.log(`[${timeInSeconds}s] Create file: ${filename}`);
  } catch (error) {
    console.error("Error creating MicroScript file:", error);
    process.exit(1);
  }
}

function runProject() {
  const startTime = process.hrtime();
  const bituinTomlPath = path.join(process.cwd(), "bituin.toml");
  const args = process.argv.slice(2);
  
  if (!fs.existsSync(bituinTomlPath)) {
    console.error("Error: bituin.toml not found. Are you in a bituin project directory?");
    process.exit(1);
  }

  try {
    const configContent = fs.readFileSync(bituinTomlPath, "utf-8");
    let mainFile;
    let mainFileName;

    // Check if a specific file was provided as an argument
    if (args.length > 1) {
      // Use the provided filename
      mainFileName = args[1];
      mainFile = path.join(process.cwd(), "src", mainFileName);
      
      // Update bituin.toml with the new main_file
      const updatedConfig = configContent.replace(
        /main_file\s*=\s*"([^"]+)"/,
        `main_file = "src/${mainFileName}"`
      );
      fs.writeFileSync(bituinTomlPath, updatedConfig, "utf-8");
    } else {
      // Use the main_file from bituin.toml
      const mainFileMatch = configContent.match(/main_file\s*=\s*"([^"]+)"/);
      
      if (!mainFileMatch) {
        // Default to main.microscript if no main_file is specified
        mainFile = path.join(process.cwd(), "src", "main.microscript");
        mainFileName = "main.microscript";
      } else {
        mainFile = path.join(process.cwd(), mainFileMatch[1]);
        mainFileName = path.basename(mainFile);
      }
    }
    
    if (!fs.existsSync(mainFile)) {
      console.error(`Error: Main file "${mainFile}" not found.`);
      process.exit(1);
    }

    // Look for microscript.exe in current and parent directory
    let microscriptExe = path.join(process.cwd(), "microscript.exe");
    if (!fs.existsSync(microscriptExe)) {
      microscriptExe = path.join(process.cwd(), "..", "microscript.exe");
    }
    
    if (!fs.existsSync(microscriptExe)) {
      console.error("Error: microscript.exe not found in current or parent directory.");
      process.exit(1);
    }

    // Show checking message
    process.stdout.write(`\x1b[90mcheck: ${mainFileName}\x1b[0m\r`);

    // Execute the MicroScript file
    exec(`${microscriptExe} run "${mainFile}"`, (error, stdout, stderr) => {
      const timeDiff = process.hrtime(startTime);
      const timeInSeconds = (timeDiff[0] + timeDiff[1] / 1e9).toFixed(3);

      // Clear the checking message
      process.stdout.write("\x1b[2K");
      
      if (error) {
        console.error(`Execution error: ${error}`);
        return;
      }
      if (stderr) {
        console.error(stderr);
      }

      // Print execution time and success message
      console.log(`\x1b[90m[${timeInSeconds}s] check: ${mainFileName}\x1b[0m`);
      console.log("\x1b[32m✓\x1b[0m Project executed successfully!\n");
      console.log(stdout);
    });
  } catch (error) {
    console.error("Error running project:", error);
    process.exit(1);
  }
}

function createProject(projectName, isNew = true) {
  const projectPath = isNew 
    ? path.join(process.cwd(), projectName) 
    : process.cwd();

  if (isNew && fs.existsSync(projectPath)) {
    console.error(`Error: Directory "${projectName}" already exists.`);
    process.exit(1);
  }

  try {
    createDirectoryStructure(projectPath);
    createMainMicroscript(projectPath);
    createBituinConfig(projectPath, projectName);
    
    // Copy microscript.exe to the project directory
    const microscriptSource = path.join(process.cwd(), "microscript.exe");
    if (fs.existsSync(microscriptSource)) {
      fs.copyFileSync(
        microscriptSource,
        path.join(projectPath, "microscript.exe")
      );
    }
    
    console.log(`Bituin project "${projectName}" created successfully!`);
    console.log(`\nTo get started:`);
    
    if (isNew) {
      console.log(`  cd ${projectName}`);
    }
    
    console.log("  bituin run");
  } catch (error) {
    console.error("Error creating project:", error);
    process.exit(1);
  }
}

function main() {
  const args = process.argv.slice(2);
  
  if (args.length === 0) {
    printUsage();
    process.exit(1);
  }

  const command = args[0];

  switch (command) {
    case COMMANDS.HELP:
      printUsage();
      break;
    case COMMANDS.VERSION:
      console.log(VERSION);
      break;
    case COMMANDS.AUTHOR:
      console.log(AUTHOR);
      break;
    case COMMANDS.NEW:
      if (args.length < 2) {
        console.error("Error: Project name required for new command");
        printUsage();
        process.exit(1);
      }
      createProject(args[1], true);
      break;
    case COMMANDS.INIT:
      if (args.length < 2) {
        console.error("Error: Project name required for init command");
        printUsage();
        process.exit(1);
      }
      createProject(args[1], false);
      break;
    case COMMANDS.ADD:
      if (args.length < 2) {
        console.error("Error: File name required for add command");
        printUsage();
        process.exit(1);
      }
      addMicroscriptFile(args[1]);
      break;
    case COMMANDS.RUN:
      runProject();
      break;
    default:
      console.error(`Unknown command: ${command}`);
      printUsage();
      process.exit(1);
  }
}

main();
