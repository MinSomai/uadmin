package main

import (
	"fmt"
	"github.com/uadmin/uadmin/preloaded"
	"github.com/uadmin/uadmin/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/otiai10/copy"
	"golang.org/x/mod/modfile"
)

// Help is the command line help for the cli tool
const Help = `Usage: uadmin COMMAND [--src]
This tools helps you prepare a folder for a new project or update static files and templates

Commands:
  prepare         Generates folders and prepares static and templates
  version         Shows the version of uAdmin

Arguments:
  --src           If you want to copy static files and templates from src folder

Get full documentation online:
https://uadmin-docs.readthedocs.io/en/latest/
`

func main() {
	args := os.Args

	// Check if there are no args
	if len(args) < 2 {
		fmt.Println(Help)
		return
	}

	// Check if the first arg is not a command
	if strings.HasSuffix(args[1], "-") {
		fmt.Println("ERROR: Invalid sytax. Please provide a command")
		fmt.Println(Help)
		return
	}
	command := args[1]

	// flag.Parse()

	// prepapre command
	if command == "prepare" {
		var dst string
		var src string
		// First ge the path
		ex, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		// Generate folders
		folderList := []string{"models", "api", "views", "media"}
		for _, v := range folderList {
			dst = filepath.Join(ex, v)
			if _, err = os.Stat(dst); os.IsNotExist(err) {
				err = os.MkdirAll(dst, os.FileMode(0744))
				if err != nil {
					utils.Trail(utils.WARNING, "Unable to create \"%s\" folder: %s", v, err)
				} else {
					utils.Trail(utils.OK, "Created: %s", dst)
				}
			}
		}

		// Copy folders
		folderList = []string{"static", "templates"}
		goPath := os.Getenv("GOPATH")
		if goPath == "" {
			if runtime.GOOS == "windows" {
				goPath = filepath.Join(os.Getenv("USERPROFILE"), "go")
			} else {
				goPath = filepath.Join(os.Getenv("HOME"), "go")
			}
			utils.Trail(utils.INFO, "Your GOPATH environment variable is not set. Using the default path: %s", goPath)
		}

		// The path from where to copy static files and templates will depend on
		// where uadmin folder is located:
		// pre 1.16 with no modules: $GOPATH/src/github.com/uadmin/uadmin
		// 1.16 and above: $GOPATH/pkg/mod/github.com/uadmin/uadmin@$uadmin.Version
		// where uadmin.Verion is the installed version of uAdmin
		uadminPathSrc := []string{goPath, "src", "github.com", "uadmin", "uadmin"}
		uadminPathMod := []string{goPath, "pkg", "mod", "github.com", "uadmin", "uadmin@v" + preloaded.Version}

		if _, err := os.Stat("go.mod"); err == nil {
			// check if there is a go.mod file and the version from that
			buf, _ := ioutil.ReadFile("go.mod")
			fs, err := modfile.Parse("go.mod", buf, nil)
			if err == nil {
				for i := range fs.Require {
					if fs.Require[i].Mod.Path == "github.com/uadmin/uadmin" {
						uadminPathMod[len(uadminPathMod)-1] = "uadmin@v" + fs.Require[i].Mod.Version
						break
					}
				}

				// Search for replace
				for i := range fs.Replace {
					if fs.Replace[i].Old.Path == "github.com/uadmin/uadmin" {
						// Check if new if a new is a file system path or module path
						if strings.HasPrefix(fs.Replace[i].New.Path, "./") ||
							strings.HasPrefix(fs.Replace[i].New.Path, "/") ||
							(len(fs.Replace[i].New.Path) > 2 && fs.Replace[i].New.Path[1] == ':') {
							uadminPathMod = []string{fs.Replace[i].New.Path}
						} else {
							uadminPathMod = append([]string{goPath, "pkg", "mod"}, strings.Split(fs.Replace[i].New.Path+"@"+fs.Replace[i].New.Version, "/")...)
						}
						break
					}
				}
			}
		}

		// By default, we will use the module version unless the command
		// was passed with --src parameter
		uadminPath := filepath.Join(uadminPathMod...)
		if len(args) > 2 && args[2] == "--src" {
			uadminPath = filepath.Join(uadminPathSrc...)
		}

		utils.Trail(utils.INFO, "Copying static/templates from: %s", uadminPath)

		for _, v := range folderList {
			msg := "Updated"
			if _, err = os.Stat(filepath.Join(ex, v)); os.IsNotExist(err) {
				msg = "Created"
			}
			dst = filepath.Join(ex, v)
			src = filepath.Join(uadminPath, v)
			err := copy.Copy(src, dst)
			if err != nil {
				utils.Trail(utils.WARNING, "Unable to copy \"%s\" folder: %s", v, err)
			} else {
				utils.Trail(utils.OK, msg+": %s", dst)
			}
		}
		return
	} else if command == "version" {
		utils.Trail(utils.INFO, preloaded.Version)
		return
	}
}
