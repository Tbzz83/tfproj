package src

import (
  "fmt"
  "flag"
  "errors"
  "strings"
  "os"
)

//TODO
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

// Struct to hold flag variables
// Have to use pointers to variables to work with the Flags package
type Flags struct {
  Describe *bool
  Create *bool
  Modules *delimStringSlice
  Envs *delimStringSlice
  Providers *delimStringSlice
  VersionBool *bool
  Backend *string
  TfDir *string
  Style *string
  Plan *bool
  
  // Modify possible styles here. Need a blank option in case there is a 
  // configuration chosen that doesn't require stack to be set
  Styles []string
}

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
func (f *Flags)describeFlag() {
  flag.BoolVar(f.Describe, "describe", false, "Usage: --describe/-describe\nWill describe the style specified by the '--style' flag")
}

func (f *Flags)createFlag() {
  flag.BoolVar(f.Create, "create", false, "Usage: --create/-create\nCreates the specified project configuration")
}

func (f *Flags)planFlag() {
  flag.BoolVar(f.Plan, "plan", false, "Usage: --plan/-plan\nWill illustrate a plan of the specified project configuration without creation")
}

func (f *Flags)versionFlag() {
  flag.BoolVar(f.VersionBool, "version", false, "Usage: --version/-version\nPrint tfproj version")
}

func (f *Flags)moduleFlag() {
  flag.Var(f.Modules, "modules", "Usage: --modules/-modules <module1,module2>\nDetermines the modules to be created. For example 'vm,vnet' will create two modules for each respectively. At least one module must be provided")
}

func (f *Flags)envsFlag() {
  flag.Var(f.Envs, "envs", "Usage: --envs/-envs <env1,env2>\nDetermines the infrastructure environments to be created. Can be left blank if desired")
}

func (f *Flags)styleFlag() {
  usageString := "Usage: --style/-style <styleName>\nDetermines the style of the project to be used.\nOptions are: "

  for _, s := range(f.Styles) {
    if s == "" {continue}
    usageString += fmt.Sprintf("'%s' ", s)
  }

  flag.StringVar(f.Style, "style", "", usageString)
}

func (f *Flags)providersFlag() {
  flag.Var(f.Providers, "providers", "Usage: --providers/-providers <provider_a=provider_a_version,provider_b=provider_b_version>\nPopulates versions.tf file sourcing providers at latest version using provided version after '='.\nIf no version is provided the latest version will be used by specifying the '...' version.\nOptions are: 'azure' (or 'azurerm') and 'aws'")
}

func (f *Flags)backendFlag() {
  flag.StringVar(f.Backend, "backend", "", "Usage: --backend/-backend <azure|aws>\nCreates backend_config.tf files with boilerplate for your tfstate storage.\nBe sure to manually specify your storage locations by editing this file\nOptions are: 'azure' (or 'azurerm') or 'aws'")
}

func (f *Flags)tfDirFlag() error {
  wd, err := os.Getwd()
  if err != nil {
    return err
  }

  flag.StringVar(f.TfDir, "dir", wd, "Usage: --dir/-dir\ndetermines the location of the terraform project")

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
func (f *Flags)dependsOnCreate() error {
  if !(*f.Create) {
    return errors.New(errorString+" '--create' flag not specified\n")
  }
  return nil
}

func (f *Flags)dependsOnModules() error {
  if len((*f.Modules)) == 0 {
    return errors.New(errorString+" '--modules' flag not specified\n")
  }
  return nil
}

// Calling flag initialization
func flagInit() *Flags {
  defer flag.Parse()
  // Initialize vars
  var describe, create, versionBool, plan bool
  var modules, envs, providers delimStringSlice
  var backend, tfDir, style string

  var f = Flags{
    Describe: &describe,
    Create: &create, 
    Modules: &modules, 
    Envs: &envs, 
    Providers: &providers, 
    VersionBool: &versionBool, 
    Backend: &backend, 
    TfDir: &tfDir,
    Style: &style, 
    Plan: &plan,
    Styles: []string{"stack", "layered", ""},
  }

  f.createFlag()
  f.moduleFlag()
  f.envsFlag()
  f.styleFlag()
  f.tfDirFlag()
  f.describeFlag()
  f.providersFlag()
  f.backendFlag()
  f.planFlag()
  f.versionFlag()

  return &f
}

// ====main====
func Cli(version string) {
  f := flagInit()

  if *f.VersionBool {
    fmt.Printf("tfproj v%s\n", version)
  }

  if *f.Describe {

    err := f.buildStyle()
    if err != nil {
      fmt.Println(err)
    }

    return 
  }


  // Remove the last slash if it exists from tfDir global variable
  if string((*f.TfDir)[len((*f.TfDir))-1]) == "/" || string((*f.TfDir)[len((*f.TfDir))-1]) == "\\" {
    (*f.TfDir) = (*f.TfDir)[:len((*f.TfDir))-1]
  }

  if (*f.Plan) {

    err := f.buildStyle()
    if err != nil {
      fmt.Println(err)
    }

    return
  }

  // Check that flags that depend on --create are being set
  if len((*f.Modules)) > 0 || len((*f.Envs)) > 0 {
    err := f.dependsOnCreate()
    if err != nil {
      fmt.Println(err)
      return
    }
  }

  // Check that flags that depend on --(*f.Modules) are being set
  if len((*f.Style)) > 0 {
    err := f.dependsOnModules()
    if err != nil {
      fmt.Println(err)
      return
    }
  }
  
  if len((*f.Modules)) > 0 {
    // Validate (*f.Style)s is correct then call build on the (*f.Style)
    err := f.buildStyle() 
    if err != nil {
      fmt.Println(err)
      return
    }
  }

  //f.testPrintFlags()
}

// Testing 
func (f *Flags)testPrintFlags() {
  fmt.Println()
  fmt.Println("---test printing all flags---")
  fmt.Println("(*f.Create) bool:", (*f.Create))
  fmt.Println("(*f.Modules):", (*f.Modules), "len (*f.Modules):", len((*f.Modules)))
  fmt.Println("(*f.Envs):", (*f.Envs), "len (*f.Envs):", len((*f.Envs)))
  fmt.Println("(*f.Style):", (*f.Style))
  fmt.Println("(*f.TfDir):", (*f.TfDir))
  fmt.Println("(*f.Providers):", (*f.Providers))
  fmt.Println("(*f.Backend):", (*f.Backend))
}

