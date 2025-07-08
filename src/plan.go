package src

import (
  "fmt" 
  //"errors"
)

func (*Stack) Plan() {
  homeDirs := [...]string{"envs", "modules"}

  if tfDir[len(tfDir)-1] == '/' {
    tfDir = tfDir[:len(tfDir)-1]
  }
  homeLevel := 0
  fmt.Println(blueDir+tfDir+reset)

  for i, dir := range(homeDirs) {
    // Root directories
    last := false
    if i == len(homeDirs) - 1 {
      last = true
    }
    printDir(blueDir+dir+reset, homeLevel, last, []int{})
    
    if dir == "envs" {
      envLevel := homeLevel + 1
      envFiles := []string{"variables.tf", "outputs.tf"}
      if len(backend) > 0 {
        envFiles = append(envFiles, "backend_config.tf")
      }
      for j, env := range(envs) {
        envFilesLines := []int{1,1}
        last := false
        if j == len(envs) - 1 {
          last = true
          envFilesLines = []int{1,0}
        } 
        printDir(blueDir+env+reset, envLevel, last, []int{1})

        for _, file := range(modules) {
          moduleFile := file+".tf"
          printDir(moduleFile, envLevel+1, false, envFilesLines)
        }

        for l, file := range(envFiles) {
          last := false
          if l == len(envFiles) - 1 {
            last = true
          }
          printDir(file, envLevel+1, last, envFilesLines)
        }
      }
    }

    // Module directories
    last = false
    if i == len(homeDirs) - 1 {
      last = true
    }
    
    if dir == "modules" {
      moduleLevel := homeLevel + 1
      moduleFiles := []string{"variables.tf", "outputs.tf", "versions.tf", "main.tf"}
      moduleLines := []int{}
      for j, module := range(modules) {
        moduleFilesLines := []int{0,1}
        last := false
        if j == len(modules) - 1 {
          last = true
          moduleFilesLines = []int{}
        } 
        printDir(blueDir+module+reset, moduleLevel, last, moduleLines)

        for k, file := range(moduleFiles) {
          last := false
          if j == len(envs) - 1 {
          }
          if k == len(moduleFiles) - 1 {
            last = true
          }
          printDir(file, moduleLevel+1, last, moduleFilesLines)
        }
      }
    }
  }
}

func printDir(name string, level int, last bool, lines []int) {
  offset := 4
  spaces := ""
  for range(offset*level) {
    spaces += " "
  }

  if len(lines) > 0 {
    spaces = ""
  }
  for _, v := range(lines) {
    if v == 1 {
      spaces += "│   "
    } else {
      spaces += "    "
    }
  }

  if last {
    fmt.Print(spaces+"└── "+name+"\n")
    return
  } 
  fmt.Print(spaces+"├── "+name+"\n")
}

func (*Layered) Plan() {

}
