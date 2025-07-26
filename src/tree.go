package src

import (
  //"fmt"
)

type TreeNode struct {
  // Last node determines if from the parent, this current
  // node is the last of the parents Children. This node itself 
  // may have Children
  LastNode bool
  Name string
  FileType string
  Children map[string]*TreeNode
}

// TODO
// This is certainly a bit of a mess to read. Open to feedback on how to clean up

// Build data structure that actually represents tree structure of
// files are directories logically, based on style
func (s *Stack)treeInit() {
  tfDir := *s.f.TfDir
  envs := *s.f.Envs
  modules := *s.f.Modules
  backend := *s.f.Backend
  rootBoilerplateFiles := s.f.RootBoilerplateFiles
  moduleBoilerplateFiles := s.f.ModuleBoilerplateFiles
  
  if backend != "" {rootBoilerplateFiles = append(rootBoilerplateFiles, "backend_config.tf")}

  tfDirFolders := []string{} 
  if len(envs) > 0 {tfDirFolders = append(tfDirFolders, "envs")}
  tfDirFolders = append(tfDirFolders, "modules") // Will always exist and looks best at end

  // Initialize head node
  var head = TreeNode{
  	Name: tfDir,
  	FileType: "dir",
  	LastNode: true,
  	Children: make(map[string]*TreeNode),
  }

  for i, dir := range(tfDirFolders) {
    var newNode = &TreeNode{
    	Name: dir,
    	FileType: "dir",
    	LastNode: false,
    	Children: make(map[string]*TreeNode),
    }

    if i == len(tfDirFolders) - 1 {newNode.LastNode = true}
    head.Children[dir] = newNode
  }

  if len(envs) == 0 {
    for _, mod := range(modules) {
      var newNode = &TreeNode{
      	Name: mod+".tf",
      	FileType: "file",
      	LastNode: false,
      }
      head.Children[mod] = newNode
    }

    for _, bp := range(rootBoilerplateFiles) {
      var newNode = &TreeNode{
      	Name: bp,
      	FileType: "file",
      	LastNode: false,
      }
      head.Children[bp] = newNode
    }

  } else {

    cur := head.Children["envs"]
    for i, env := range(envs) {
      var newNode = &TreeNode{
      	Name: env,
      	FileType: "dir",
      	LastNode: false,
      	Children: make(map[string]*TreeNode),
      }
      if i == len(envs) - 1 {newNode.LastNode = true}
      cur.Children[env] = newNode

      cur = cur.Children[env]

      for _, mod := range(modules) {
        var newNode = &TreeNode{
          Name: mod+".tf",
          FileType: "dir",
          LastNode: false,
          Children: make(map[string]*TreeNode),
        }
        cur.Children[mod] = newNode
      }

      for i, bp := range(rootBoilerplateFiles) {
        var newNode = &TreeNode{
          Name: bp,
          FileType: "file",
          LastNode: false,
        }
        if i == len(rootBoilerplateFiles) - 1 {newNode.LastNode = true}
        cur.Children[bp] = newNode
      }

      cur = head.Children["envs"]
    }
  }

  cur := head.Children["modules"]

  for i, mod := range(modules) {
    var newNode = &TreeNode{
      Name: mod,
      FileType: "dir",
      LastNode: false,
      Children: make(map[string]*TreeNode),
    }
    if i == len(moduleBoilerplateFiles) - 1 {newNode.LastNode = true}
    cur.Children[mod] = newNode

    cur = cur.Children[mod]

    for i, bp := range(moduleBoilerplateFiles) {
      var newNode = &TreeNode{
        Name: bp,
        FileType: "file",
        LastNode: false,
      }
      if i == len(moduleBoilerplateFiles) - 1 {newNode.LastNode = true}
      cur.Children[bp] = newNode
    }

    cur = head.Children["modules"]
  }

  s.TreeHead = &head
}

func (l *Layered)treeInit() {
  tfDir := *l.f.TfDir
  envs := *l.f.Envs
  modules := *l.f.Modules
  backend := *l.f.Backend
  rootBoilerplateFiles := l.f.RootBoilerplateFiles
  moduleBoilerplateFiles := l.f.ModuleBoilerplateFiles
  
  if backend != "" {rootBoilerplateFiles = append(rootBoilerplateFiles, "backend_config.tf")}

  tfDirFolders := []string{} 
  if len(envs) > 0 {tfDirFolders = append(tfDirFolders, "envs")}
  tfDirFolders = append(tfDirFolders, "modules") // Will always exist and looks best at end

  // Initialize head node
  var head = TreeNode{
  	Name: tfDir,
  	FileType: "dir",
  	LastNode: true,
  	Children: make(map[string]*TreeNode),
  }

  for i, dir := range(tfDirFolders) {
    var newNode = &TreeNode{
    	Name: dir,
    	FileType: "dir",
    	LastNode: false,
    	Children: make(map[string]*TreeNode),
    }

    if i == len(tfDirFolders) - 1 {newNode.LastNode = true}
    head.Children[dir] = newNode
  }

  if len(envs) == 0 {
    cur := &head
    for _, mod := range(modules) {
      var newNode = &TreeNode{
      	Name: mod,
      	FileType: "dir",
      	LastNode: false,
      	Children: make(map[string]*TreeNode),
      }
      cur.Children[mod] = newNode

      cur = cur.Children[mod]

      for _, bp := range(rootBoilerplateFiles) {
        var newNode = &TreeNode{
          Name: bp,
          FileType: "file",
          LastNode: false,
        }
        cur.Children[bp] = newNode
      }

      newNode = &TreeNode{
        Name: "main.tf",
        FileType: "file",
        LastNode: true,
      }

      cur.Children["main.tf"] = newNode
    }

  } else {
    cur := head.Children["envs"]

    for i, env := range(envs) {
      var newNode = &TreeNode{
      	Name: env,
      	FileType: "dir",
      	LastNode: false,
      	Children: make(map[string]*TreeNode),
      }
      if i == len(envs) - 1 {newNode.LastNode = true}
      cur.Children[env] = newNode

      cur = cur.Children[env]

      for _, mod := range(modules) {
        var newNode = &TreeNode{
          Name: mod,
          FileType: "dir",
          LastNode: false,
          Children: make(map[string]*TreeNode),
        }
        cur.Children[mod] = newNode

        cur = cur.Children[mod]

        for _, bp := range(rootBoilerplateFiles) {
          var newNode = &TreeNode{
            Name: bp,
            FileType: "file",
            LastNode: false,
          }
          cur.Children[bp] = newNode
        }

        newNode = &TreeNode{
          Name: "main.tf",
          FileType: "file",
          LastNode: true,
        }

        cur.Children["main.tf"] = newNode

        // I know :/
        cur = head.Children["envs"].Children[env]
      }

      cur = head.Children["envs"]
    }
  }

  cur := head.Children["modules"]

  for i, mod := range(modules) {
    var newNode = &TreeNode{
      Name: mod,
      FileType: "dir",
      LastNode: false,
      Children: make(map[string]*TreeNode),
    }
    if i == len(moduleBoilerplateFiles) - 1 {newNode.LastNode = true}
    cur.Children[mod] = newNode

    cur = cur.Children[mod]

    for i, bp := range(moduleBoilerplateFiles) {
      var newNode = &TreeNode{
        Name: bp,
        FileType: "file",
        LastNode: false,
        Children: make(map[string]*TreeNode),
      }
      if i == len(moduleBoilerplateFiles) - 1 {newNode.LastNode = true}
      cur.Children[bp] = newNode
    }

    cur = head.Children["modules"]
  }

  l.TreeHead = &head
}

func (s *Stack)printAll() {
  s.TreeHead.printAllRecurse(0)
}

func (l *Layered)printAll() {
  l.TreeHead.printAllRecurse(0)
}


