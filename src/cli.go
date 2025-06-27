package src

import (
  "fmt"
  "flag"
  "errors"
  "strings"
  "os"
)

// delimStringSlice allows reading a delimited string from the cli into
// a single flag according to allowed delimeters specified in delimSplit()
type delimStringSlice []string

// Global variables 
var create bool
var modules delimStringSlice
var envs delimStringSlice
var tfDir string
var style string

// Constants
const (
  red = "\033[31m"
  reset = "\033[0m" 
  redError = red+"\nError:"+reset
)

// allow splits using these delimeters
func delimSplit(r rune) bool {
  return r == ':' || r == ',' || r == ';' || r == ' '
}

// Set function required by flag.Var. Instructs on how row input value
// from flag.Var should be handled and processesed for delimStringSlice types
func (cs *delimStringSlice) Set(value string) error {
  *cs = strings.FieldsFunc(value, delimSplit)
  if len(*cs) == 0 {
    return errors.New("\tError, invalid comma separated string in flag")
  }
  return nil
}

func (cs *delimStringSlice) String() string {
  return "TEST"
}

// Initiliazing global flags
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
func tfDirInit() error {
  wd, err := os.Getwd()
  if err != nil {
    return err
  }
  flag.StringVar(&tfDir, "dir", wd, "Usage: --dir/-dir. determines the location of the terraform project. Default to current directory")

  return nil
}

// Makes sure style is formatted correctly
func validateStyles() error {
  var err error = nil
  switch style {
  case "monolith":
    err := buildMonolith()
    if err != nil {
      return err
    }
  case "":
  default:
    err = errors.New(redError+" '"+style+"' is not a valid option for 'style'.\nOptions are: 'monolithic'\n")
  }
  return err
}


// Depends on specific flag checker
func dependsOn(flagName string) error {
  if !create {
    // throw error
    return errors.New(redError+" '"+flagName+"' flag not specified\n")
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
}


// --main--
func Cli() {
  flagInit()

  // Remove the last slash if it exists
  if string(tfDir[len(tfDir)-1]) == "/" || string(tfDir[len(tfDir)-1]) == "\\" {
    tfDir = tfDir[:len(tfDir)-1]
  }
  

  // Check that flags that depend on --create are being set
  if len(modules) > 0 || len(envs) > 0 {
    err := dependsOn("--create")
    if err != nil {
      fmt.Println(err)
      return
    }
  }
  
  // Validate styles is correct then call build on the style
  err := validateStyles() 
  if err != nil {
    fmt.Println(err)
    return
  }

  testPrintFlags()
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
}

