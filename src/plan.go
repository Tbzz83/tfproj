package src

import (
  "fmt" 
  //"errors"
)

func (*Stack) Plan() {
  homeDirs := [...]string{"envs", "modules"}
  homeIndentLevel := 0

  fmt.Println(blueDir+tfDir+reset)

  for i, dir := range(homeDirs) {
    // Root directories
    lastDir := false
    if i == len(homeDirs) - 1 {
      lastDir = true
    }
    
    if dir == "envs" {
      envLevel := homeIndentLevel + 1
      envFiles := []string{"variables.tf", "outputs.tf"}

      if len(backend) > 0 {
        envFiles = append(envFiles, "backend_config.tf")
      }

      if len(envs) == 0 {
        // no envs specified by user
        for _, file := range(modules) {
          moduleFile := file+".tf"
          printDir(moduleFile, homeIndentLevel, false, []int{})
        }

        for _, file := range(envFiles) {
          printDir(file, homeIndentLevel, false, []int{})
        }
      } else {
        printDir(blueDir+dir+reset, homeIndentLevel, lastDir, []int{})
      }
      
      // skipped if len(envs) == 0
      for j, env := range(envs) {
        envFilesLines := []int{1,1}
        lastEnv := false
        if j == len(envs) - 1 {
          lastEnv = true
          envFilesLines = []int{1,0}
        } 
        printDir(blueDir+env+reset, envLevel, lastEnv, []int{1})

        for _, file := range(modules) {
          moduleFile := file+".tf"
          printDir(moduleFile, envLevel+1, false, envFilesLines)
        }

        for k, file := range(envFiles) {
          lastFile := false
          fileLevel := envLevel + 1

          if k == len(envFiles) - 1 {
            lastFile = true
          }

          printDir(file, fileLevel, lastFile, envFilesLines)
        }
      }
    }

    // Module directories
    lastHomeDir := false
    if i == len(homeDirs) - 1 {
      lastHomeDir = true
    }
    
    if dir == "modules" {
      printDir(blueDir+dir+reset, homeIndentLevel, lastHomeDir, []int{})
      moduleLevel := homeIndentLevel + 1
      moduleFiles := []string{"main.tf", "variables.tf", "outputs.tf", "versions.tf"}
      moduleLines := []int{}

      for i, module := range(modules) {
        moduleFilesLines := []int{0,1}
        lastModule := false

        if i == len(modules) - 1 {
          lastModule = true
          moduleFilesLines = []int{}
        } 

        printDir(blueDir+module+reset, moduleLevel, lastModule, moduleLines)

        for j, file := range(moduleFiles) {
          lastFile := false
          if i == len(envs) - 1 {
          }
          if j == len(moduleFiles) - 1 {
            lastFile = true
          }
          printDir(file, moduleLevel+1, lastFile, moduleFilesLines)
        }
      }
    }
  }
}

// Prints in a tree like structure. 'intendLevel' determines how many indents should
// be present. If you want backing vertical lines to be printed as well, provide a populated
// 'lines' array. Eg. vertLines = []int{1,0,1} will tell the function to print a line `|`
// at the 0th and 2nd indent levels but not the 1st.
func printDir(name string, indentLevel int, last bool, vertLines []int) {
  offset := 4
  spaces := ""

  if len(vertLines) == 0 {
    for range(offset*indentLevel) {
      spaces += " "
    }
  } 
 
  // skipped if len(vertLines) == 0
  for _, v := range(vertLines) {
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
  homeIndentLevel := 0

  fmt.Println(blueDir+tfDir+reset)

  for i, dir := range(homeDirs) {
    // Root directories
    lastHome := false
    if i == len(homeDirs) - 1 {
      lastHome = true
    }
    
    if dir == "envs" {
      envLevel := homeIndentLevel + 1
      envFiles := []string{"variables.tf", "outputs.tf"}

      if len(backend) > 0 {
        envFiles = append(envFiles, "backend_config.tf")
      }

      if len(envs) == 0 {
        // no user specified envs
        for _, module := range(modules) {

          printDir(blueDir+module+reset, homeIndentLevel, false, []int{})
          printDir("main.tf", homeIndentLevel+1, false, []int{1})

          for i, file := range(envFiles) {
            lastFile := false

            if i == len(envFiles) - 1 {
              lastFile = true
            }

            printDir(file, homeIndentLevel+1, lastFile, []int{1})
          }
        }
      } else {
        printDir(blueDir+dir+reset, homeIndentLevel, lastHome, []int{})
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
          moduleLevel := envLevel+1
          lastModule := false

          if k == len(modules) - 1 {
            lastModule = true
          }

          if lastEnv {
            moduleLines[1] = 0
          }

          printDir(moduleFile, moduleLevel, lastModule, moduleLines)

          if lastModule {
            envFilesLines[2] = 0
          }

          printDir("main.tf", envLevel+2, false, envFilesLines)

          for l, file := range(envFiles) {
            lastFile := false
            fileLevel := moduleLevel+1

            if l == len(envFiles) - 1 {
              lastFile = true
            }

            printDir(file, fileLevel, lastFile, envFilesLines)
          }
        }
      }
    }

    // Module directories
    lastHome = false
    if i == len(homeDirs) - 1 {
      lastHome = true
    }
    
    if dir == "modules" {

      printDir(blueDir+dir+reset, homeIndentLevel, lastHome, []int{})
      moduleLevel := homeIndentLevel + 1
      moduleFiles := []string{"main.tf", "variables.tf", "outputs.tf", "versions.tf"}
      moduleLines := []int{}

      for j, module := range(modules) {
        moduleFilesLines := []int{0,1}
        lastHome = false

        if j == len(modules) - 1 {
          lastHome = true
          moduleFilesLines = []int{}
        } 

        printDir(blueDir+module+reset, moduleLevel, lastHome, moduleLines)

        for k, file := range(moduleFiles) {
          lastFile := false
          fileLevel := moduleLevel+1

          if j == len(envs) - 1 {
          }

          if k == len(moduleFiles) - 1 {
            lastFile = true
          }

          printDir(file, fileLevel, lastFile, moduleFilesLines)
        }
      }
    }
  }
}
