package main

import (
	"os"
	"os/exec"
	"path"
)

func addSrc(src string) error {
	cmd := exec.Command("git", "clone", "--depth=1", src)
	cmd.Dir = src
	return cmd.Run()
}

func listSrc() []string {
	d, err := os.Open(src)
	if err != nil {
		return []string{}
	}
	defer d.Close()
	n, err := d.Readdirnames(0)
	if err != nil {
		return []string{}
	}
	return n
}

func makePkg(pkgName string) error {
	cmd := exec.Command("git", "clone", "--depth=1", path.Join(src, pkgName))
	cmd.Dir = build
	err := cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("makepkg")
	cmd.Dir = path.Join(build, pkgName)
	err = cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("cp", "*.pkg.tar.xz", repo)
	cmd.Dir = path.Join(build, pkgName)
	err = cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("repo-add", repo, repoDB)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil

}
