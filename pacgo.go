package main

import (
	"net/url"
	"os"
	"os/exec"
	"path"
)

func src2local(src string) (string, error) {
	u, err := url.Parse(src)
	if err != nil {
		return "", err
	}
	b := os.Getenv("GOPATH")
	if p := u.Path; p[len(p)-4:] == ".git" {
		u.Path = p[:len(p)-4]
	}
	return path.Join(b, "src", u.Host, u.Path), nil
}

func addSrc(src string) error {
	dir, err := src2local(src)
	if err != nil {
		return err
	}
	cmd := exec.Command("git", "clone", src, dir)
	return cmd.Run()
}

// LIST STUFF
func getSub(dir string) []string {
	d, err := os.Open(dir)
	if err != nil {
		return []string{}
	}
	defer d.Close()
	ds, err := d.Readdirnames(0)
	if err != nil {
		return []string{}
	}
	return ds
}

func has(elem string, array []string) bool {
	for i := range array {
		if elem == array[i] {
			return true
		}
	}
	return false
}

func listProj(dir string) []string {
	projs := []string{}
	dirs := getSub(dir)
	if has(".git", dirs) {
		return []string{dir}
	} else {
		for _, v := range dirs {
			projs = append(projs, listProj(path.Join(dir, v))...)
		}
		return projs
	}
}

func listSrc() []string {
	src := path.Join(os.Getenv("GOPATH"), "src")
	fprojs := listProj(src)
	var projs []string
	k := len(src) + 1
	for _, v := range fprojs {
		projs = append(projs, v[k:])
	}
	return projs
}

// END OF LIST STUFF. PLEASE DO BETTER!!!!

func newBuildDir() (string, error) {
	dir := "/tmp/pacgo"
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", err
	}
	return dir, nil
}

func makePkg(src string) (err error) {
	repo := path.Join(os.Getenv("GOPATH"), "repo") // output dir
	db := path.Join(repo, "pacgo.db.tar.gz")       // db file
	orig, err := src2local(src)                    // local origin
	if err != nil {
		return err
	}
	dir, err := newBuildDir() // Volative build dir (future: chroot)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	cmd := exec.Command("git", "clone", orig, dir)
	err = cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("makepkg")
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("cp", "*.pkg.tar.xz", repo)
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("repo-add", repo, db)
	return cmd.Run()
}

func printConf() {
	println("[pacgo]")
	println("SigLevel = Optional TrustAll")
	println("Server = file://" + repo)
}
