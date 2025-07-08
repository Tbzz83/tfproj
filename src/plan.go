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
    
    if dir == "envs" {
      envLevel := homeLevel + 1
      envFiles := []string{"variables.tf", "outputs.tf"}
      if len(backend) > 0 {
        envFiles = append(envFiles, "backend_config.tf")
      }
      if len(envs) == 0 {
        // no envs specified by user
        for _, file := range(modules) {
          moduleFile := file+".tf"
          printDir(moduleFile, homeLevel, false, []int{})
        }

        for _, file := range(envFiles) {
          printDir(file, homeLevel, false, []int{})
        }

      } else {
        printDir(blueDir+dir+reset, homeLevel, last, []int{})
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
      printDir(blueDir+dir+reset, homeLevel, last, []int{})
      moduleLevel := homeLevel + 1
      moduleFiles := []string{"main.tf", "variables.tf", "outputs.tf", "versions.tf"}
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

// Prints in a tree like structure. 'level' determines how many indents should
// be present. If you want backing lines to be printed as well, provide a populated
// 'lines' array. Eg. lines = []int{1,0,1} will tell the function to print a line `|`
// at the 0th and 2nd indent levels but not the 1st.
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
  homeDirs := [...]string{"envs", "modules"}

  if tfDir[len(tfDir)-1] == '/' {
    tfDir = tfDir[:len(tfDir)-1]
  }
  homeLevel := 0
  fmt.Println(blueDir+tfDir+reset)

  for i, dir := range(homeDirs) {
    // Root directories
    lastHome := false
    if i == len(homeDirs) - 1 {
      lastHome = true
    }
    
    if dir == "envs" {
      envLevel := homeLevel + 1
      envFiles := []string{"variables.tf", "outputs.tf"}
      if len(backend) > 0 {
        envFiles = append(envFiles, "backend_config.tf")
      }

      if len(envs) == 0 {
        // no user specified envs
        for _, module := range(modules) {
          printDir(blueDir+module+reset, homeLevel, false, []int{})

          printDir("main.tf", homeLevel+1, false, []int{1})

          for i, file := range(envFiles) {
            lastFile := false
            if i == len(envFiles) - 1 {
              lastFile = true
            }
            printDir(file, homeLevel+1, lastFile, []int{1})
          }
        }

      } else {
        printDir(blueDir+dir+reset, homeLevel, lastHome, []int{})
      }

      for j, env := range(envs) {
        envFilesLines := []int{1,1,1}
        lastEnv := false
        if j == len(envs) - 1 {
          lastEnv = true
          envFilesLines = []int{1,0,1}
        } 
        printDir(blueDir+env+reset, envLevel, lastEnv, []int{1})

        for k, file := range(modules) {
          moduleFile := blueDir+file+reset
          moduleLines := []int{1,1}
          lastModule := false
          if k == len(modules) - 1 {
            lastModule = true
          }

          if lastEnv {
            moduleLines[1] = 0
          }

          printDir(moduleFile, envLevel+1, lastModule, moduleLines)
          if lastModule {
            envFilesLines[2] = 0
          }
          printDir("main.tf", envLevel+2, false, envFilesLines)
          for l, file := range(envFiles) {
            lastFile := false
            if l == len(envFiles) - 1 {
              lastFile = true
            }
            printDir(file, envLevel+2, lastFile, envFilesLines)
          }
        }
      }
    }

    // Module directories
    last := false
    if i == len(homeDirs) - 1 {
      last = true
    }
    
    if dir == "modules" {
      printDir(blueDir+dir+reset, homeLevel, lastHome, []int{})
      moduleLevel := homeLevel + 1
      moduleFiles := []string{"main.tf", "variables.tf", "outputs.tf", "versions.tf"}
      moduleLines := []int{}
      for j, module := range(modules) {
        moduleFilesLines := []int{0,1}
        last = false
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
