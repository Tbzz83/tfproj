package src

import (
  "fmt" 
  //"errors"
)

func (*Stack) Plan() {
  homeDirs := [...]string{"envs/", "modules/"}

  if tfDir[len(tfDir)-1] != '/' {
    tfDir += "/"
  }
  homeLevel := 0
  fmt.Println(tfDir)

  for i, dir := range(homeDirs) {
    // Root directories
    last := false
    if i == len(homeDirs) - 1 {
      last = true
    }
    printDir(dir, homeLevel, last)
    
    if dir == "envs/" {
      envLevel := homeLevel + 1
      envFiles := []string{"variables.tf", "outputs.tf"}
      if len(backend) > 0 {
        envFiles = append(envFiles, "backend_config.tf")
      }
      for j, env := range(envs) {
        last := false
        if j == len(envs) - 1 {
          last = true
        }
        printLine(envLevel - 1)
        printDir(env+"/", envLevel, last)
        for k, file := range(envFiles) {
          last := false
          if k == len(envFiles) - 1 {
            last = true
          }
          printDir(file, envLevel+1, last)
        }
      }
    }
    // Module directories




  }
}

func (*Layered) Plan() {

}

func printLine(level int) {
  spaces := ""
  for range(3*level) {
    spaces += " "
  }
  fmt.Print("│")
}
func printDir(name string, level int, last bool) {
  spaces := ""
  for range(3*level) {
    spaces += " "
  }
  if last {
    fmt.Print(spaces+"└──"+name+"\n")
    return
  } 
  fmt.Print(spaces+"├──"+name+"\n")
}
