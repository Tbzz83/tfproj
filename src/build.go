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
}

type Layered struct {
  Name string 
  Description string
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
func moduleBoilerplate(path string) error {
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
      err := versionsHeredoc(filePath)
      if err != nil {
        return err
      }
    }
  }
  return nil
}

// Similar to moduleBoilerplate but for generic root files that will call modules
func rootBoilerplate(path string) error {
  files := [...]string{"variables.tf", "outputs.tf"}
  for _, name := range(files) {
    filePath := path + "/" + name
    err := touchFile(filePath)
    if err != nil {
      return err
    }
  }
  return nil
}

// Makes sure style is formatted correctly, then call build() for the respective style requested if valid
func buildStyle() error {
  var err error
  var project Project
  switch strings.ToLower(style) {
  case "stack":
    project = &Stack{style, stackDescription()}
  case "layered":
    project = &Layered{style, layeredDescription()}
  case "":
    fmt.Print(warningString+" you have not provided a value for '--style'\n\n")
  default:
    errMsg := errorString+" '"+style+"' is not a valid option for '--style'\nOptions are: "
    for _, s := range(styles) {
      if s == "" {continue}
      errMsg += fmt.Sprintf("'%s' ", s)
    }
    errMsg += "\n"
    err = errors.New(errMsg)
  }

  // If describe flag set print the description then return without building
  if describe {
    project.Describe()
    return nil
  }
  
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
func (*Stack) Build() error {
  // create dirs for modules
  for _, moduleName := range(modules) {
    modulePath := tfDir+"/modules/"+moduleName
    err := createDir(modulePath)
    if err != nil {
      return err
    }
    // Create module boilerplate in each modules directory
    err = moduleBoilerplate(modulePath)
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

        err = rootBoilerplate(envPath)
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
      err = rootBoilerplate(tfDir)
      if err != nil {
        return err
      }
    }
  }
  return nil
}

// A project where each module (like vm, vnet etc...) has an individual
// root directory dedicated to it, each with its own .tfstate files.
func (* Layered) Build() error {
  for _, moduleName := range(modules) {
    modulePath := tfDir+"/modules/"+moduleName
    err := createDir(modulePath)
    if err != nil {
      return err
    }
    err = moduleBoilerplate(modulePath)
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

        err = rootBoilerplate(envPath)
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
      err = rootBoilerplate(tfDir+"/"+moduleName)
      if err != nil {
        return err
      }
    }
  }
  return nil
}

func createFile(path string) (*os.File, error) {
  _, err := os.Stat(path)
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

func versionsHeredoc(path string) error {
  requiredProviders := ""
  for _, prov := range(providers) {
    s := new(equalDelimSlice)
    err := s.Set(prov)
    if err != nil {
      return err
    }

    if len(*s) > 2 {
      return errors.New(errorString+" too many values to unpack for provider '"+prov+"'\n")
    }

    var provVersion string

    if len(*s) > 1 {
      // Can assume user has provided a version
      provVersion = (*s)[1]
    } else {
      provVersion = "..."//... will use the latest versions of any provider
    }
    switch (*s)[0] {
    case "aws":
      // aws provider doc
      requiredProviders += `
      aws = {
        source = "hashicorp/aws"
        version = "`+provVersion+`"
      }
      `
    case "azurerm":
      // azure provider doc
      requiredProviders += `
      azurerm = {
        source = "hashicorp/azurerm"
        version = "`+provVersion+`"
      }
      `
    default:
      return errors.New(errorString+" '"+prov+"' is not a valid provider\n")
    }
  }

  doc := heredoc.Doc(
  `terraform {
    required_providers {
      `+requiredProviders+`
    }
  }`+"\n")
  f, err := createFile(path)
  if os.IsExist(err) {
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
