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

func createDir(path string) error {
  err := os.MkdirAll(path, 0755)
  if err != nil {
    return err
  }
  return nil
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

// Doesn't need to be a pointer as it is (aside from the name) stateless
func buildMonolith() error {
  // create dirs for modules
  for _, name := range(modules) {
    path := tfDir+"/modules/"+name
    err := createDir(path)
    if err != nil {
      return err
    }
    err = moduleBoilerplate(path)
    if err != nil {
      return err
    }
  }

  // create dirs for envs
  for _, name := range(envs) {
    path := tfDir+"/envs/"+name
    err := createDir(path)
    if err != nil {
      return err
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

func versionsHeredoc(path string) error {
  bp := heredoc.Doc(`
  terraform {
    required_providers {}
    }
  }`+"\n")
  f, err := createFile(path)
  if os.IsExist(err) {
    return nil
  }

  err = os.WriteFile(path, []byte(bp), 0755)

  f.Close()
  return err
}




