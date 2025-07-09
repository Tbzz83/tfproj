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
          printDir(moduleFile, homeIndentLevel, false, []bool{})
        }

        for _, file := range(envFiles) {
          printDir(file, homeIndentLevel, false, []bool{})
        }
      } else {
        printDir(blueDir+dir+reset, homeIndentLevel, lastDir, []bool{})
      }
      
      // skipped if len(envs) == 0
      for j, env := range(envs) {
        envFilesLines := []bool{true,true}
        lastEnv := false
        if j == len(envs) - 1 {
          lastEnv = true
          envFilesLines = []bool{true,false}
        } 
        printDir(blueDir+env+reset, envLevel, lastEnv, []bool{true})

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
      printDir(blueDir+dir+reset, homeIndentLevel, lastHomeDir, []bool{})
      moduleLevel := homeIndentLevel + 1
      moduleFiles := []string{"main.tf", "variables.tf", "outputs.tf", "versions.tf"}
      moduleLines := []bool{}

      for i, module := range(modules) {
        moduleFilesLines := []bool{false,true}
        lastModule := false

        if i == len(modules) - 1 {
          lastModule = true
          moduleFilesLines = []bool{}
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
func printDir(name string, indentLevel int, last bool, vertLines []bool) {
  offset := 4
  spaces := ""

  if len(vertLines) == 0 {
    for range(offset*indentLevel) {
      spaces += " "
    }
  } 
 
  // skipped if len(vertLines) == 0
  for _, vert := range(vertLines) {
    if vert {
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
          moduleIndentLevel := homeIndentLevel+1

          printDir(blueDir+module+reset, homeIndentLevel, false, []bool{})
          printDir("main.tf", moduleIndentLevel, false, []bool{true})

          for i, file := range(envFiles) {
            lastFile := false

            if i == len(envFiles) - 1 {
              lastFile = true
            }

            printDir(file, moduleIndentLevel, lastFile, []bool{true})
          }
        }
      } else {
        printDir(blueDir+dir+reset, homeIndentLevel, lastHome, []bool{})
      }

      for j, env := range(envs) {

        envFilesLines := []bool{true,true,true}
        lastEnv := false

        if j == len(envs) - 1 {
          lastEnv = true
          envFilesLines = []bool{true,false,true}
        } 

        printDir(blueDir+env+reset, envLevel, lastEnv, []bool{true})

        for k, file := range(modules) {
          moduleFile := blueDir+file+reset
          moduleLines := []bool{true,true}
          moduleLevel := envLevel+1
          lastModule := false

          if k == len(modules) - 1 {
            lastModule = true
          }

          if lastEnv {
            moduleLines[1] = false
          }

          printDir(moduleFile, moduleLevel, lastModule, moduleLines)

          if lastModule {
            envFilesLines[2] = false
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

      printDir(blueDir+dir+reset, homeIndentLevel, lastHome, []bool{})
      moduleLevel := homeIndentLevel + 1
      moduleFiles := []string{"main.tf", "variables.tf", "outputs.tf", "versions.tf"}
      moduleLines := []bool{}

      for j, module := range(modules) {
        moduleFilesLines := []bool{false,true}
        lastHome = false

        if j == len(modules) - 1 {
          lastHome = true
          moduleFilesLines = []bool{}
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
