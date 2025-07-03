package src

import (
  "fmt"
  "flag"
  "errors"
  "strings"
  "os"
)

//TODO 
// Providers
// --provider aws (or azure etc...)
// DONE...

// Backend for tfstate
// --backend aws (or azure etc...) 

// Plan
// --plan like tf plan to show the tree structure that will be created

// Copy
// --copy <existing_module>,<new_module_1>,<new_module_2>
// copy should change key words in the new modules like the module source heredoc to match
// the new module name

// delimStringSlice allows reading a delimited string from the cli into
// a single flag according to allowed delimeters specified in delimSplit()
type delimStringSlice []string

// equalDelimString reads a string and splits it into a slice based on a
// '=' character. 
type equalDelimSlice []string 

type Project interface {
  Build() error
  Describe()
}

// Global variables 
var describe bool
var create bool
var modules delimStringSlice
var envs delimStringSlice
var providers delimStringSlice
var tfDir string
var style string

// Modify possible styles here. Need a blank option in case there is a 
// configuration chosen that doesn't require stack to be set
var styles = [...]string{"stack", ""} 

// Constants
const (
  yellow = "\033[33m"
  red = "\033[31m"
  reset = "\033[0m" 
  errorString = red+"\nError:"+reset
  warningString = yellow+"\nWarning:"+reset
)

// allow splits using these delimeters
func delimSplit(r rune) bool {
  return r == ':' || r == ',' || r == ';' || r == ' '
}

// split string based on = sign only
func equalSplit(r rune) bool {
  return r == '='
}

func (s *equalDelimSlice) Set(value string) error {
  *s = strings.FieldsFunc(value, equalSplit)
  if len(*s) == 0 {
    return errors.New(errorString+" invalid equal ('=') separated string in flag")
  }
  return nil
}

// Set function required by flag.Var. Instructs on how row input value
// from flag.Var should be handled and processesed for delimStringSlice types
func (s *delimStringSlice) Set(value string) error {
  *s = strings.FieldsFunc(value, delimSplit)
  if len(*s) == 0 {
    return errors.New(errorString+" invalid comma separated string in flag")
  }
  return nil
}

func (s *equalDelimSlice) String() string {
  return ""
}

func (s *delimStringSlice) String() string {
  return ""
}

// Initiliazing global flags
func describeInit() {
  flag.BoolVar(&describe, "describe", false, "Usage: --describe/-describe. Will describe the style specified by the '--style' flag")
}
func createInit() {
  flag.BoolVar(&create, "create", false, "Usage: --create/-create")
}
func moduleInit() {
  flag.Var(&modules, "modules", "Usage: --modules/-modules. Requires '--create' to be set")
}
func envsInit() {
  flag.Var(&envs, "envs", "Usage: --envs/-envs. Requires '--create' to be set")
}
func styleInit() {
  flag.StringVar(&style, "style", "", "Usage: --style/-style. Requires '--modules' to be set")
}
func providersInit() {
  flag.Var(&providers, "providers", "Usage: --providers/-providers. Requires '--create' to be set. Options are 'azure', 'aws'")
}
func tfDirInit() error {
  wd, err := os.Getwd()
  if err != nil {
    return err
  }
  flag.StringVar(&tfDir, "dir", wd, "Usage: --dir/-dir. determines the location of the terraform project")

  return nil
}

func stackDescription() string {
  return `
  ----Stack Project----
  A project type where modules are referred to by a single .tf file. 
  A stack based architecture with one environment called 'dev' and two 
  modules called 'vm' and 'vnet' might look like:
  stack/
  ├── modules/
  │   ├── vm/
  │   │   ├── main.tf
  │   │   ├── variables.tf
  │   │   ├── outputs.tf
  │   │   └── versions.tf
  │   └── vnet/
  │       ├── main.tf
  │       ├── variables.tf
  │       ├── outputs.tf
  │       └── versions.tf
  └── envs/
      └── dev/
          ├── vm.tf
          ├── vnet.tf
          ├── variables.tf
          └── outputs.tf
  `
}

func layeredDescription() string {
  return `
  ----Layered Project----
  A project where each module (like vm, vnet etc...) has an individual
  root directory dedicated to it, each with its own .tfstate files.
  A layered based architecture with one environment called 'dev' and two modules called 'vm' and 'vnet' might look like:
  layered/
  ├── modules/
  │   ├── vm/
  │   │   ├── main.tf
  │   │   ├── variables.tf
  │   │   ├── outputs.tf
  │   │   └── versions.tf
  │   └── vnet/
  │       ├── main.tf
  │       ├── variables.tf
  │       ├── outputs.tf
  │       └── versions.tf
  └── envs/
      └── dev/
          ├── vm/
          │   ├── main.tf
          │   ├── variables.tf
          │   └── outputs.tf
          └── vnet/
              ├── main.tf
              ├── variables.tf
              └── outputs.tf
  `
}


// Depends on specific flag checker
func dependsOnCreate() error {
  if !create {
    // throw error
    return errors.New(errorString+" '--create' flag not specified\n")
  }
  return nil
}

func dependsOnModules() error {
  if len(modules) == 0 {
    return errors.New(errorString+" '--modules' flag not specified\n")
  }
  return nil
}

// Calling flag initialization
func flagInit() {
  defer flag.Parse()
  createInit()
  moduleInit()
  envsInit()
  styleInit()
  tfDirInit()
  describeInit()
  providersInit()
}

// --main--
func Cli() {
  flagInit()

  if describe {
    if style == "" {
      fmt.Println(errorString + " no style specified. Please specify a style with the '--style' flag")
      fmt.Println()
      return 
    }
    err := buildStyle()
    if err != nil {
      fmt.Println(err)
    }
    return 
  }

  // Remove the last slash if it exists
  if string(tfDir[len(tfDir)-1]) == "/" || string(tfDir[len(tfDir)-1]) == "\\" {
    tfDir = tfDir[:len(tfDir)-1]
  }

  // Check that flags that depend on --create are being set
  if len(modules) > 0 || len(envs) > 0 {
    err := dependsOnCreate()
    if err != nil {
      fmt.Println(err)
      return
    }
  }

  // Check that flags that depend on --modules are being set
  if len(style) > 0 {
    err := dependsOnModules()
    if err != nil {
      fmt.Println(err)
      return
    }
  }
  
  if len(modules) > 0 {
    // Validate styles is correct then call build on the style
    err := buildStyle() 
    if err != nil {
      fmt.Println(err)
      return
    }
  }

  //testPrintFlags()
}


// Testing 
func testPrintFlags() {
  fmt.Println()
  fmt.Println("---test printing all flags---")
  fmt.Println("create bool:", create)
  fmt.Println("modules:", modules, "len modules:", len(modules))
  fmt.Println("envs:", envs, "len envs:", len(envs))
  fmt.Println("style:", style)
  fmt.Println("tfDir:", tfDir)
  fmt.Println("providers:", providers)
}

