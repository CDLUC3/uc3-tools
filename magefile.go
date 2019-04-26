// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// ------------------------------------------------------------
// Targets

// Build builds binaries for the current platform.
func Build() error {
	subprojects, err := Subprojects()
	if err != nil {
		return err
	}
	for _, sp := range subprojects {
		err = sp.Build()
		if err != nil {
			return err
		}
	}
	return nil
}

// BuildLinux builds linux-amd64 binaries (the most common cross-compile case).
func BuildLinux() error {
	subprojects, err := Subprojects()
	if err != nil {
		return err
	}
	opSys, arch := "linux", "amd64"
	for _, sp := range subprojects {
		err = sp.BuildFor(opSys, arch)
		if err != nil {
			return err
		}
	}
	return nil
}

// BuildAll builds binaries for each target platform.
func BuildAll() error {
	subprojects, err := Subprojects()
	if err != nil {
		return err
	}
	for _, sp := range subprojects {
		for _, opSys := range operatingSystems {
			for _, arch := range architectures {
				err = sp.BuildFor(opSys, arch)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Clean removes compiled binaries from the current working directory.
func Clean() error {
	subprojects, err := Subprojects()
	if err != nil {
		return err
	}
	for _, sp := range subprojects {
		err = sp.Clean()
		if err != nil {
			return err
		}
	}
	return nil
}

// Install installs in $GOPATH/bin.
func Install() error {
	subprojects, err := Subprojects()
	if err != nil {
		return err
	}
	for _, sp := range subprojects {
		err = sp.Install()
		if err != nil {
			return err
		}
	}
	return nil
}

// Platforms lists target platforms for buildAll.
func Platforms() {
	for _, opSys := range operatingSystems {
		for _, arch := range architectures {
			fmt.Printf("%s-%s\n", opSys, arch)
		}
	}
}

// ------------------------------------------------------------
// Unexported symbols

var goCmd = mg.GoCmd()
var operatingSystems = []string{"darwin", "linux", "windows"}
var architectures = []string{"amd64"}
var projectRoot = getProjectRoot()
var moduleDeclarationRe = regexp.MustCompile("(?m:^module ([[:^space:]]+)$)")

func getProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	wdAbs, err := filepath.Abs(wd)
	if err != nil {
		panic(err)
	}
	return wdAbs
}

func reportBuilt(binName string) {
	binAbs, err := filepath.Abs(binName)
	if err != nil {
		panic(err)
	}
	binRel, err := filepath.Rel(projectRoot, binAbs)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stderr, "%v\n", binRel)
}

// -----------------------------
// Subproject

type Subproject interface {
	Name() string
	Build() error
	BuildFor(opSys, arch string) error
	Clean() error
	Install() error
}

func Subprojects() ([]Subproject, error) {
	var subprojects []Subproject
	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if info.Name() == "go.mod" {
			parentDir := filepath.Dir(path)
			if parentDir == projectRoot {
				return nil
			}
			parentDirRelative := filepath.Base(parentDir)
			sp := subproject{name: parentDirRelative}
			subprojects = append(subprojects, &sp)
		}
		return nil
	})
	return subprojects, err
}

type subproject struct {
	name          string
	projectDirAbs string
	rootPkg       string
	cmdPkg        string
}

func (p *subproject) Name() string {
	return p.name
}

func (p *subproject) ProjectDirAbs() string {
	if p.projectDirAbs == "" {
		relpath := path.Join(projectRoot, p.name)
		resolved, err := filepath.EvalSymlinks(relpath)
		if err != nil {
			panic(err)
		}
		abspath, err := filepath.Abs(resolved)
		if err != nil {
			panic(err)
		}
		p.projectDirAbs = abspath
	}
	return p.projectDirAbs
}

func (p *subproject) RootPkg() string {
	if p.rootPkg == "" {
		gomod := path.Join(p.ProjectDirAbs(), "go.mod")
		data, err := ioutil.ReadFile(gomod)
		if err != nil {
			panic(err)
		}
		matches := moduleDeclarationRe.FindStringSubmatch(string(data))
		if len(matches) < 2 {
			panic(fmt.Errorf("no module declaration found in %v:\n%v\n", gomod, string(data)))
		}
		p.rootPkg = matches[1]
	}
	return p.rootPkg
}

func (p *subproject) CmdPkg() string {
	if p.cmdPkg == "" {
		cmdPkg := ""
		projectDir := p.ProjectDirAbs()
		err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
			if info.Name() == "cmd" && info.IsDir() {
				relPath, err := filepath.Rel(projectDir, path)
				if err != nil {
					return err
				}
				cmdPkg = filepath.Join(p.RootPkg(), relPath)
				return io.EOF
			}
			return nil
		})
		if err != nil && err != io.EOF {
			panic(err)
		}
		if cmdPkg == "" {
			panic(fmt.Errorf("no cmd package found under %v", projectDir))
		}
		p.cmdPkg = cmdPkg
	}
	return p.cmdPkg
}

func (p *subproject) LdFlags() (string, error) {
	commitHash, err := sh.Output("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return "", err
	}
	tag, err := sh.Output("git", "describe", "--tags")
	if err != nil {
		return "", err
	}
	timestamp := time.Now().Format(time.RFC3339)

	flagVals := map[string]string{
		"commitHash": commitHash,
		"tag":        tag,
		"timestamp":  timestamp,
	}

	var flags []string
	for k, v := range flagVals {
		flag := fmt.Sprintf("-X %s.%s=%s", p.CmdPkg(), k, v)
		flags = append(flags, flag)
	}
	return strings.Join(flags, " "), nil
}

func (p *subproject) BinNameFor(opSys string, arch string) string {
	binName := p.name
	binName = fmt.Sprintf("%s-%s-%s", binName, opSys, arch)
	if opSys == "windows" {
		binName += ".exe"
	}
	return binName
}

// TODO: extract shared build code

func (p *subproject) Install() error {
	err := os.Chdir(p.ProjectDirAbs())
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Chdir(projectRoot)
	}()
	binName := p.name
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	gopath, err := sh.Output(goCmd, "env", "GOPATH")
	if err != nil {
		return fmt.Errorf("error determining GOPATH: %v", err)
	}
	binDir := filepath.Join(gopath, "bin")
	binPath := filepath.Join(binDir, binName)

	flags, err := p.LdFlags()
	if err != nil {
		return fmt.Errorf("error determining ldflags: %v", err)
	}
	err = sh.RunV(goCmd, "build", "-o", binPath, "-ldflags", flags)
	if err == nil {
		reportBuilt(binName)
	}
	return err
}

func (p *subproject) Build() error {
	err := os.Chdir(p.ProjectDirAbs())
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Chdir(projectRoot)
	}()

	binName := p.name
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	flags, err := p.LdFlags()
	if err != nil {
		return fmt.Errorf("error determining ldflags: %v", err)
	}
	err = sh.RunV(goCmd, "build", "-o", binName, "-ldflags", flags)
	if err == nil {
		reportBuilt(binName)
	}
	return err
}

func (p *subproject) BuildFor(opSys, arch string) error {
	err := os.Chdir(p.ProjectDirAbs())
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Chdir(projectRoot)
	}()

	binName := p.BinNameFor(opSys, arch)

	flags, err := p.LdFlags()
	if err != nil {
		return fmt.Errorf("error determining ldflags: %v", err)
	}

	env := map[string]string{
		"GOOS":   opSys,
		"GOARCH": arch,
	}
	err = sh.RunWith(env, goCmd, "build", "-o", binName, "-ldflags", flags)
	if err == nil {
		reportBuilt(binName)
	}
	return err
}

func (p *subproject) Clean() error {
	err := os.Chdir(p.ProjectDirAbs())
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Chdir(projectRoot)
	}()
	var binRe = regexp.MustCompile("^" + p.name + "(-[a-zA-Z0-9]+-[a-zA-Z0-9]+)?(.exe)?$")

	rmcmd := "rm"
	if runtime.GOOS == "windows" {
		rmcmd = "del"
	}

	files, err := ioutil.ReadDir("./")
	if err != nil {
		return err
	}

	for _, f := range files {
		mode := f.Mode()
		isPlainFile := mode.IsRegular() && mode&os.ModeSymlink == 0
		isExecutable := mode&0111 != 0
		if isPlainFile && isExecutable {
			name := f.Name()
			if binRe.MatchString(name) {
				err := sh.RunV(rmcmd, name)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
