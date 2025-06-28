package src

import (
  "fmt"
  "flag"
  "errors"
  "strings"
  "os"
)

//TODO 
//Providers

// delimStringSlice allows reading a delimited string from the cli into
// a single flag according to allowed delimeters specified in delimSplit()
type delimStringSlice []string

// Global variables 
var create bool
var modules delimStringSlice
var envs delimStringSlice
var tfDir string
var style string
// Modify possible styles here
var styles = [...]string{"monolith", ""} 

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

// Set function required by flag.Var. Instructs on how row input value
// from flag.Var should be handled and processesed for delimStringSlice types
func (cs *delimStringSlice) Set(value string) error {
  *cs = strings.FieldsFunc(value, delimSplit)
  if len(*cs) == 0 {
    return errors.New(errorString+" invalid comma separated string in flag")
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
  flag.StringVar(&tfDir, "dir", wd, "Usage: --dir/-dir. determines the location of the terraform project")

  return nil
}

// Makes sure style is formatted correctly, then call build() for the respective style requested if valid
func buildStyle() error {
  var err error
  switch style {
  case "monolith":
    err = buildMonolith()
  case "":
    fmt.Println(warningString+" you have not provided a value for '--style'")
  default:
    errMsg := errorString+" '"+style+"' is not a valid option for '--style'\nOptions are: "
    for _, s := range(styles) {
      if s == "" {continue}
      errMsg += fmt.Sprintf("%q ", s)
    }
    err = errors.New(errMsg)
    return err
  }
  return nil
}


// Depends on specific flag checker
func dependsOn(flagName string) error {
  if !create {
    // throw error
    return errors.New(errorString+" '"+flagName+"' flag not specified\n")
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

  // Check that flags that depend on --modules are being set
  if len(style) > 0 {
    err := dependsOn("--modules")
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

  // 
  //err := buildMonolith()

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

