package new

import (
	"errors"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"

	"github.com/rickywei/sparrow/project"
)

var CmdNew = &cobra.Command{
	Use:   "new",
	Short: "Create a new project",
	Long:  "Create a new project",
	Run:   run,
}

func run(_ *cobra.Command, args []string) {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	project := ""
	if len(args) == 0 {
		prompt := &survey.Input{
			Message: "What is the project name?",
			Help:    "Please input the project name. eg. github.com/rickywei/sparrow or sparrow.",
		}
		if err = survey.AskOne(prompt, &project); err != nil || project == "" {
			return
		}
	} else {
		project = args[0]
	}

	p := &newProject{
		project: project,
		name:    filepath.Base(project),
		dir:     workDir,
	}
	if err = p.New(); err != nil {
		log.Printf("failed to new project %v\n", err)
	}

}

type newProject struct {
	project string
	name    string
	dir     string
}

func (p *newProject) New() (err error) {
	dst := filepath.Join(p.dir, p.name)
	_, err = os.Stat(dst)
	if !os.IsNotExist(err) {
		log.Printf("%s already exists\n", p.name)
		prompt := &survey.Confirm{
			Message: "Do you want to override the folder?",
			Help:    "Delete the existing folder and create the project.",
		}
		confirm := false
		if err = survey.AskOne(prompt, &confirm); err != nil {
			return
		}
		if !confirm {
			return
		}
		// if err = os.RemoveAll(dst); err != nil {
		// 	return
		// }
	}

	err = fs.WalkDir(project.EmbedFS, ".", func(path string, d fs.DirEntry, e error) (err error) {
		if e != nil {
			return e
		}

		to := filepath.Join(dst, strings.Replace(path, "project", "", 1))
		log.Println(to)

		if d.IsDir() {
			// return os.Mkdir(to, 0755)
			err = os.Mkdir(to, 0755)
			if errors.Is(err, os.ErrExist) {
				err = nil
			}
			return
		}
		data, err := fs.ReadFile(project.EmbedFS, path)
		if err != nil {
			return
		}
		if filepath.Ext(path) == ".go" {
			if filepath.Base(path) == "template.go" {
				return
			}
			if data, err = p.replaceImport(string(data)); err != nil {
				return
			}
		} else if filepath.Base(path) == "go.mod" {
			var f *modfile.File
			f, err = modfile.Parse(path, data, nil)
			if err != nil {
				return
			}
			if err = f.AddModuleStmt(p.project); err != nil {
				return
			}
			if data, err = f.Format(); err != nil {
				return
			}
		}
		err = os.WriteFile(to, data, 0644)
		return
	})

	return
}

func (p *newProject) replaceImport(src string) (data []byte, err error) {
	f, err := parser.ParseFile(token.NewFileSet(), "", src, parser.ImportsOnly)
	if err != nil {
		return
	}

	last, err := lo.Last(f.Imports)
	if err != nil {
		data, err = []byte(src), nil
		return
	}
	packagePart := src[:f.Name.End()]
	importPart := src[f.Name.End():last.End()]
	importPart = strings.ReplaceAll(importPart, "github.com/rickywei/sparrow/project", p.project)
	otherPart := src[last.End():]

	data = []byte(packagePart + importPart + otherPart)

	return
}
