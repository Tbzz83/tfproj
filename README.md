![GitHub release](https://img.shields.io/github/v/release/Tbzz83/tfproj)
![License](https://img.shields.io/github/license/Tbzz83/tfproj)

# tfproj

## What's tfproj?

Tfproj is a simple lightweight cli tool that sets up your terraform project structure based on a specific style.

Supports `--create` to generate the project structure and `--plan` to preview it without making changes.

## Installation
1. Go to [releases](https://github.com/Tbzz83/tfproj/releases)
2. Download the binary for your operating system and extract it
4. Make it executable (Linux/MacOS) with `chmod +x tfproj`
5. Move to path (optional but recommended) `mv tfproj ~/.local/bin`
6. Validate installation with `tfproj --version`

## Quickstart

```bash
tfproj --create --envs dev --modules vm,vnet --style layered --dir my-tfproj --providers azure=4.36.0 --backend azure
```

Creates a `layered` Terraform project with two modules (`vm`, `vnet`) in the `my-tfproj` directory, with Azure as backend and provider.

## Usage

Tfproj has three main options. 
1. `tfproj --plan` will show an example of the project structure to be created.
2. `tfproj --create` will build the project structure specified.
3. `tfproj --style` determines the manner in which the terraform project structure will be created. Current options are `stack` or `layered`. See below for details

### Style
#### Layered style
```bash
$ tfproj --plan --envs dev --modules vm,vnet,rg --dir tfDir --providers azure=4.36.0,aws --backend azure --style layered
tfDir
├── envs
│   └── dev
│       ├── vm
│       │   ├── main.tf
│       │   ├── variables.tf
│       │   ├── outputs.tf
│       │   └── backend_config.tf
│       ├── vnet
│       │   ├── main.tf
│       │   ├── variables.tf
│       │   ├── outputs.tf
│       │   └── backend_config.tf
│       └── rg
│           ├── main.tf
│           ├── variables.tf
│           ├── outputs.tf
│           └── backend_config.tf
└── modules
    ├── vm
    │   ├── main.tf
    │   ├── variables.tf
    │   ├── outputs.tf
    │   └── versions.tf
    ├── vnet
    │   ├── main.tf
    │   ├── variables.tf
    │   ├── outputs.tf
    │   └── versions.tf
    └── rg
        ├── main.tf
        ├── variables.tf
        ├── outputs.tf
        └── versions.tf
```
#### Stack style
```bash
$ tfproj --plan --envs dev --modules vm,vnet,rg --dir tfDir --providers azure=4.36.0,aws --backend azure --style stack
tfDir
├── envs
│   └── dev
│       ├── vm.tf
│       ├── vnet.tf
│       ├── rg.tf
│       ├── variables.tf
│       ├── outputs.tf
│       └── backend_config.tf
└── modules
    ├── vm
    │   ├── main.tf
    │   ├── variables.tf
    │   ├── outputs.tf
    │   └── versions.tf
    ├── vnet
    │   ├── main.tf
    │   ├── variables.tf
    │   ├── outputs.tf
    │   └── versions.tf
    └── rg
        ├── main.tf
        ├── variables.tf
        ├── outputs.tf
        └── versions.tf
```
### `--create` vs `--plan`
Using the `--create` flag instead of `--plan` will create the style you specify based on the options presented. No folders or files will be overwritten. If you have a pre-existing terraform project and you simply want to add more modules, existing files or directories will be skipped.

### `--backend`
when the `--backend` flag is specified, a file called `backend_config.tf` will be created and initially populated with some boilerplate code for using a remote terraform state location. For options on available backend providers see options below

### `--providers`
similarly to `--backend`, `--providers` populates the `versions.tf` file in the `modules` directory with some boilerplate code for any providers you specify. You can also specify a specific version of a provider and multiple providers with `--providers azure=4.36.0 aws`. If no provider version is specified it will default to `...` meaning the latest provider terraform can find, though this is *not* recommended.

## Flags
- `-backend` string  
      Usage: --backend/-backend <azure|aws>  
      Creates `backend_config.tf` files with boilerplate for your tfstate storage.  
      Be sure to manually specify your storage locations by editing this file.  
      Options are: `azure` (or `azurerm`) or `aws`  

- `-create`  
      Usage: --create/-create  
      Creates the specified project configuration  

- `-describe`  
      Usage: --describe/-describe  
      Will describe the style specified by the `--style` flag  

- `-dir` string  
      Usage: --dir/-dir  
      Determines the location of the Terraform project (default `/home/azeezoe/projects/tfproj`)  

- `-envs` value  
      Usage: --envs/-envs <env1,env2>  
      Determines the infrastructure environments to be created. Can be left blank if desired

- `-modules` value  
      Usage: --modules/-modules <module1,module2>  
      Determines the modules to be created. For example, `vm,vnet` will create two modules respectively.  
      At least one module must be provided

- `-plan`  
      Usage: --plan/-plan  
      Will illustrate a plan of the specified project configuration without creation  

- `-providers` value  
      Usage: --providers/-providers <provider_a=version_a,provider_b=version_b>  
      Populates `versions.tf` by sourcing providers using the specified versions.  
      If no version is provided, the latest version will be used by specifying the `'...'` version.  
      Options are: `azure` (or `azurerm`) and `aws`  

- `-style` string  
      Usage: --style/-style <styleName>  
      Determines the style of the project to be used.  
      Options are: `stack`, `layered`  

- `-version`  
      Usage: --version/-version  
      Print `tfproj` version


