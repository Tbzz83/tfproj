package src
import (
  "fmt"
  "os"
  "errors"
  "github.com/MakeNowJust/heredoc"
  "strings"
)

/*
Build functions for various env/module/style combinations
*/

type Stack struct {
  Name string
  Description string
  f *Flags
}

type Layered struct {
  Name string 
  Description string
  f *Flags
}

func (p *Stack) Describe() {
  fmt.Println(p.Description)
}

func (p *Layered) Describe() {
  fmt.Println(p.Description)
}

// Create boilerplate .tf files found in module directories
// Typically main.tf, variables.tf, versions.tf
// Path is the directory of the module you want the boilerplate to be placed
// Function will check that existing .tf files already exist, and if they do 
// leave them alone. 'exclude' will skip files if the exclude string is provided
func moduleBoilerplate(path string, providers delimStringSlice) error {
  files := [...]string{"main.tf", "variables.tf", "versions.tf", "outputs.tf"}

  for _, name := range(files) {
    if name != "versions.tf" {
      filePath := path + "/" + name 

      err := touchFile(filePath)
      if err != nil {
        return err
      }

    } else {
      // For versions.tf we can put some basic boilerplate .tf code
      filePath := path+"/"+name

      err := versionsHeredoc(filePath, providers)
      if err != nil {
        return err
      }
    }
  }

  return nil
}

// Similar to moduleBoilerplate but for generic root files that will call modules
func rootBoilerplate(path, backend string) error {
  files := [...]string{"variables.tf", "outputs.tf"}

  for _, name := range(files) {
    filePath := path + "/" + name

    err := touchFile(filePath)
    if err != nil {
      return err
    }
  }

  err := backendHeredoc(path + "/" + "backend_config.tf", path, backend)
  if err != nil {
    return err
  }

  return err
}

// Makes sure style is formatted correctly, then call build() for the respective style requested if valid
func (f *Flags)buildStyle() error {
  var err error
  var project Project

  style := *f.Style
  styles := f.Styles
  describe := *f.Describe
  plan := *f.Plan

  switch strings.ToLower(style) {
  case "stack":
    project = &Stack{Name: style, Description: stackDescription(), f: f}
  case "layered":
    project = &Layered{Name: style, Description: layeredDescription(), f: f}
//  case "":
//    fmt.Print(warningString+" you have not provided a value for '--style'\n\n")
  default:
    errMsg := errorString+" '"+style+"' is not a valid option for '--style'\nOptions are: "

    // Adding options to string from global styles variable
    for _, s := range(styles) {
      if s == "" {continue}
      errMsg += fmt.Sprintf("'%s' ", s)
    }
    errMsg += "\n"

    err = errors.New(errMsg)
    return err
  }

  // If describe flag set print the description then return without building
  if describe {
    project.Describe()
    return nil
  }

  // If plan flag set, print the plan and then return without building
  if plan {
    project.Plan()
    return nil
  }
  
  // Final check to make sure project is populated
  if project == nil {
    return errors.New(errorString+" unknown error occurred with style '"+style+"'\n")
  }

  err = project.Build()

  return err
}

// For stack projects, a .tf file for each module should be created
// in one directory for every env. it should also put some basic boilerplate
// for sourcing the appropriate module. May also optionally create a 
// backend_config.tf file if the .tfstate is being store externally
func (s *Stack) Build() error {
  modules := *s.f.Modules
  tfDir := *s.f.TfDir
  envs := *s.f.Envs
  providers := *s.f.Providers
  backend := *s.f.Backend

  for _, moduleName := range(modules) {
    modulePath := tfDir+"/modules/"+moduleName

    err := createDir(modulePath)
    if err != nil {
      return err
    }

    // Create module boilerplate in each modules directory
    err = moduleBoilerplate(modulePath, providers)
    if err != nil {
      return err
    }

    if len(envs) > 0 {
      // create dirs for envs if they don't already exist
      // Create boilerplate files to source modules
      for _, envName := range(envs) {
        envPath := tfDir+"/envs/"+envName

        _, err := os.Stat(envPath)
        if os.IsNotExist(err) {
          err := createDir(envPath)
          if err != nil {
            return err
          }
        } 

        moduleRelPath := "../../modules/"+moduleName

        err = sourceModuleHeredoc(envPath+"/"+moduleName+".tf", moduleName, moduleRelPath)
        if err != nil {
          return err
        }

        err = rootBoilerplate(envPath, backend)
        if err != nil {
          return err
        }
      }
    } else {
      // User has not provided any envs and just wants files directly in root directory
      moduleRelPath := "modules/"+moduleName

      err = sourceModuleHeredoc(tfDir+"/"+moduleName+".tf", moduleName, moduleRelPath)
      if err != nil {
        return err
      }

      err = rootBoilerplate(tfDir, backend)
      if err != nil {
        return err
      }
    }
  }
  return nil
}

// A project where each module (like vm, vnet etc...) has an individual
// root directory dedicated to it, each with its own .tfstate files.
func (l *Layered) Build() error {
  modules := *l.f.Modules
  tfDir := *l.f.TfDir
  envs := *l.f.Envs
  providers := *l.f.Providers
  backend := *l.f.Backend

  for _, moduleName := range(modules) {

    modulePath := tfDir+"/modules/"+moduleName

    err := createDir(modulePath)
    if err != nil {
      return err
    }

    err = moduleBoilerplate(modulePath, providers)
    if err != nil {
      return err
    }

    if len(envs) > 0 {
      for _, envName := range(envs) {
        envPath := tfDir+"/envs/"+envName+"/"+moduleName

        _, err := os.Stat(envPath)
        if os.IsNotExist(err) {
          err := createDir(envPath)
          if err != nil {
            return err
          }
        } 

        moduleRelPath := "../../../modules/"+moduleName

        err = sourceModuleHeredoc(envPath+"/main.tf", moduleName, moduleRelPath)
        if err != nil {
          return err
        }

        err = rootBoilerplate(envPath, backend)
        if err != nil {
          return err
        }
      }
    } else {
      err := createDir(tfDir+"/"+moduleName)
      if err != nil {
        return err
      }

      // User has not provided any envs and just wants files directly in root directory
      moduleRelPath := "../modules/"+moduleName

      err = sourceModuleHeredoc(tfDir+"/"+moduleName+"/main.tf", moduleName, moduleRelPath)
      if err != nil {
        return err
      }

      err = rootBoilerplate(tfDir+"/"+moduleName, backend)
      if err != nil {
        return err
      }
    }
  }
  return nil
}

// helper function for touchFile. If you want to create an empty file, call touchFile as it
// checks for existing file first.
func createFile(path string) (*os.File, error) {
  _, err := os.Stat(path)

  // checks if file exists
  if err == nil { 
    return nil, os.ErrExist
  }
  return os.Create(path)
}

// Creates a file and checks if there is a file that already exists
func touchFile(path string) error {
  
  f, err := createFile(path)

  if os.IsExist(err) {
  } else if err != nil {
    return err
  }
  f.Close()
  return nil
}

// Creates the main.tf file in the root directory (directory that calls tf modules) and adds
// heredoc to source the appropriate tf module in tf code
func sourceModuleHeredoc(path string, moduleName string, modulePath string) error {
  doc := heredoc.Doc(`
  module "`+moduleName+`" {
    source = "`+modulePath+`"
  }`+"\n")

  f, err := createFile(path)

  if os.IsExist(err) {
    return nil
  }

  err = os.WriteFile(path, []byte(doc), 0755)
  
  f.Close()
  return err
}

// Function to validate whether the backend flag supplied by the user is valid.
// if it is valid, return the heredoc for backend_config.tf
func switchBackendHeredoc(path, tfDir, backend string) (string, error) {
  backendDoc := ""

  // just get the env/dev/etc.. part to use for the key path 
  key := path[len(tfDir)+1:]

  switch backend {
  case "azure", "azurerm":
    // azure heredoc
    backendDoc = heredoc.Doc(
    `terraform {
      backend "azurerm" {
        resource_group_name  = "YOUR_RESOURCE_GROUP_NAME"
        storage_account_name = "YOUR_SA_NAME"
        container_name       = "YOUR_CONTAINER_NAME"
        key                  = "`+key+`"
        use_azuread_auth     = true
      }
    }`)
  case "aws", "s3":
    // aws heredoc
    backendDoc = heredoc.Doc(
    `terraform {
      backend "s3" {
        bucket = "YOUR_BUCKET_NAME"
        key    = "`+key+`"
        region = "YOUR_REGION"
      }
    }`)
  default:
    return "", errors.New(errorString+" '"+backend+"' is not a valid backend source\n")
  }

  return backendDoc, nil
}

// Function that will create a backend_config.tf file and put the appropriate contents
// based on the global backend variable flag. (Eg. '--backend azurerm' will use azure storage as the presumed
// backend provider). If backend variable zeroed then return without touching backend_config.tf
func backendHeredoc(path, tfDir, backend string) error {
  if backend == "" {
    return nil
  }

  backendDoc, err := switchBackendHeredoc(path, tfDir, backend)
  if err != nil {
    return err
  }
  
  f, err := createFile(path)
  if os.IsExist(err) {
    return nil
  }
  err = os.WriteFile(path, []byte(backendDoc), 0755)

  f.Close()
  return err
}

// Function that will create a versions.tf file and put the appropriate contents
// based on the global providers variable flag. Allows you to provide multiple different
// providers simultaneously
func versionsHeredoc(path string, providers delimStringSlice) error {
  requiredProviders := ""

  for _, provider := range(providers) {
    s := new(equalDelimSlice)
    err := s.Set(provider)

    if err != nil {
      return err
    }

    if len(*s) > 2 {
      return errors.New(errorString+" too many values to unpack for provider '"+provider+"'\n")
    }

    var providerVersion string

    if len(*s) > 1 {
      // Can assume user has provided a version
      providerVersion = (*s)[1]
    } else {
      providerVersion = "..." //'...' shorthand for latest versions of any provider
    }

    switch (*s)[0] {
    case "aws":
      // aws provider doc
      requiredProviders += `
      aws = {
        source  = "hashicorp/aws"
        version = "`+providerVersion+`"
      }
      `
    case "azurerm", "azure":
      // azure provider doc
      requiredProviders += `
      azurerm = {
        source  = "hashicorp/azurerm"
        version = "`+providerVersion+`"
      }
      `
    default:
      return errors.New(errorString+" '"+provider+"' is not a valid provider\n")
    }
  }

  doc := heredoc.Doc(
  `terraform {
    required_providers {
      `+requiredProviders+`
    }
  }`+"\n")

  // Call createFile here instead of touchFile as you need the pointer to the file
  f, err := createFile(path)
  if os.IsExist(err) {
    // If the file exists
    return nil
  }

  err = os.WriteFile(path, []byte(doc), 0755)

  f.Close()
  return err
}

func createDir(path string) error {
  err := os.MkdirAll(path, 0755)
  if err != nil {
    return err
  }
  return nil
}
