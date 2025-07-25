package src

import (
  "fmt"
  "flag"
  "errors"
  "strings"
  "os"
)

//TODO
// Use options struct instead of flag globals
// left adjust usage so it looks better. see mongodump for an example

// delimStringSlice allows reading a delimited string from the cli into
// a single flag according to allowed delimeters specified in delimSplit()
type delimStringSlice []string

// equalDelimString reads a string and splits it into a slice based on a
// '=' character. 
type equalDelimSlice []string 

// The main interface for which various styles will call their respective methods
// If a new style is added it must implement at least all the methods of the Project interface
type Project interface {
  Build() error
  Describe()
  Plan()
}

// Global variables and flag initialization
var (
  describe bool
  create bool
  modules delimStringSlice
  envs delimStringSlice
  providers delimStringSlice
  versionBool bool
  backend string
  tfDir string
  style string
  plan bool
  // Modify possible styles here. Need a blank option in case there is a 
  // configuration chosen that doesn't require stack to be set
  styles = [...]string{"stack", "layered", ""} 
)


// Constants
const (
  blueDir = "\033[1;34m"
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

// necessary methods to support flag.Var but don't need to return anything
func (s *equalDelimSlice) String() string {
  return ""
}
func (s *delimStringSlice) String() string {
  return ""
}

// Initiliazing global flags
func describeFlag() {
  flag.BoolVar(&describe, "describe", false, "Usage: --describe/-describe\nWill describe the style specified by the '--style' flag")
}

func createFlag() {
  flag.BoolVar(&create, "create", false, "Usage: --create/-create\nCreates the specified project configuration")
}

func planFlag() {
  flag.BoolVar(&plan, "plan", false, "Usage: --plan/-plan\nWill illustrate a plan of the specified project configuration without creation")
}

func versionFlag() {
  flag.BoolVar(&versionBool, "version", false, "Usage: --version/-version\nPrint tfproj version")
}

func moduleFlag() {
  flag.Var(&modules, "modules", "Usage: --modules/-modules <module1,module2>\nDetermines the modules to be created. For example 'vm,vnet' will create two modules for each respectively. At least one module must be provided")
}

func envsFlag() {
  flag.Var(&envs, "envs", "Usage: --envs/-envs <env1,env2>\nDetermines the infrastructure environments to be created. Can be left blank if desired")
}

func styleFlag() {
  usageString := "Usage: --style/-style <styleName>\nDetermines the style of the project to be used.\nOptions are: "

  for _, s := range(styles) {
    if s == "" {continue}
    usageString += fmt.Sprintf("'%s' ", s)
  }

  flag.StringVar(&style, "style", "", usageString)
}

func providersFlag() {
  flag.Var(&providers, "providers", "Usage: --providers/-providers <provider_a=provider_a_version,provider_b=provider_b_version>\nPopulates versions.tf file sourcing providers at latest version using provided version after '='.\nIf no version is provided the latest version will be used by specifying the '...' version.\nOptions are: 'azure' (or 'azurerm') and 'aws'")
}

func backendFlag() {
  flag.StringVar(&backend, "backend", "", "Usage: --backend/-backend <azure|aws>\nCreates backend_config.tf files with boilerplate for your tfstate storage.\nBe sure to manually specify your storage locations by editing this file\nOptions are: 'azure' (or 'azurerm') or 'aws'")
}

func tfDirFlag() error {
  wd, err := os.Getwd()
  if err != nil {
    return err
  }

  flag.StringVar(&tfDir, "dir", wd, "Usage: --dir/-dir\ndetermines the location of the terraform project")

  return nil
}

func stackDescription() string {
  return `
  ----Stack Project----
  A project type where modules are referred to by a single .tf file. 
  A stack based architecture with one environment called 'dev' and two 
  modules called 'vm' and 'vnet' might look like:

  `+blueDir+`stack`+reset+`
  ├── `+blueDir+`modules`+reset+`
  │   ├── `+blueDir+`vm`+reset+`
  │   │   ├── main.tf
  │   │   ├── variables.tf
  │   │   ├── outputs.tf
  │   │   └── versions.tf
  │   └── `+blueDir+`vnet`+reset+`
  │       ├── main.tf
  │       ├── variables.tf
  │       ├── outputs.tf
  │       └── versions.tf
  └── `+blueDir+`envs`+reset+`
      └── `+blueDir+`dev`+reset+`
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

  `+blueDir+`layered`+reset+`
  ├── `+blueDir+`modules`+reset+`
  │   ├── `+blueDir+`vm`+reset+`
  │   │   ├── main.tf
  │   │   ├── variables.tf
  │   │   ├── outputs.tf
  │   │   └── versions.tf
  │   └── `+blueDir+`vnet`+reset+`
  │       ├── main.tf
  │       ├── variables.tf
  │       ├── outputs.tf
  │       └── versions.tf
  └── `+blueDir+`envs`+reset+`
      └── `+blueDir+`dev`+reset+`
          ├── `+blueDir+`vm`+reset+`
          │   ├── main.tf
          │   ├── variables.tf
          │   └── outputs.tf
          └── `+blueDir+`vnet`+reset+`
              ├── main.tf
              ├── variables.tf
              └── outputs.tf
  `
}


// Depends on specific flag checker
func dependsOnCreate() error {
  if !create {
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
  createFlag()
  moduleFlag()
  envsFlag()
  styleFlag()
  tfDirFlag()
  describeFlag()
  providersFlag()
  backendFlag()
  planFlag()
  versionFlag()
}

// ====main====
func Cli(version string) {
  flagInit()

  if versionBool {
    fmt.Printf("tfproj v%s\n", version)
  }

  if describe {

    err := buildStyle()
    if err != nil {
      fmt.Println(err)
    }

    return 
  }

  // Remove the last slash if it exists from tfDir global variable
  if string(tfDir[len(tfDir)-1]) == "/" || string(tfDir[len(tfDir)-1]) == "\\" {
    tfDir = tfDir[:len(tfDir)-1]
  }

  if plan {

    err := buildStyle()
    if err != nil {
      fmt.Println(err)
    }

    return
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
  fmt.Println("backend:", backend)
}

