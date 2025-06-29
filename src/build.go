package src
import (
  //"fmt"
  "os"
  //"errors"
  "github.com/MakeNowJust/heredoc"
)

/*
Build functions for various env/module/style combinations
*/

type Monolith struct {
  name string
}

// Create boilerplate .tf files found in module directories
// Typically main.tf, variables.tf, versions.tf
// Path is the directory of the module you want the boilerplate to be placed
// Function will check that existing .tf files already exist, and if they do 
// leave them alone
func moduleBoilerplate(path string) error {
  files := [...]string{"main.tf", "variables.tf", "versions.tf"}
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

// For monolith projects, a .tf file for each module should be created
// in one directory for every env. it should also put some basic boilerplate
// for sourcing the appropriate module. May also optionally create a 
// backend_config.tf file if the .tfstate is being store externally
func (*Monolith) Build() error {
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
  doc := heredoc.Doc(`
  terraform {
    required_providers {}
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
