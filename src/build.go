package src
import (
  "fmt"
  "os"
  //"errors"
)

/*
Build functions for various env/module/style combinations
*/

func createModules() error {
  for _, name := range(modules) {
    path := tfDir+"/"+name
    err := os.MkdirAll(path, 0755)
    if err != nil {
      return err
    }
  }
  return nil
}

func buildMonolith() error {
  fmt.Println("Hello from monolith")
  err := createModules()
  if err != nil {
    return err
  }

  return nil
}

