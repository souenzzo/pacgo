package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	src    = "/home/souenzzo/.local/src"
	repo   = "/home/souenzzo/.local/pkg"
	repoDB = "/home/souenzzo/.local/pkg/pacgo.db.tar.gz"
	build  = "/tmp/pacgo"
)

func main() {
	os.MkdirAll(src, os.ModePerm)
	os.MkdirAll(repo, os.ModePerm)
	os.MkdirAll(build, os.ModePerm)

	newSrc := flag.String("add", "", "Indica a origem que deseja adicionar")
	addr := flag.String("http", "", "Interface web")
	list := flag.Bool("list", false, "Lista pacotes fontes")
	make := flag.String("make", "", "Compila pacote")
	flag.Parse()
	switch {
	case *addr != "":
		err := web(*addr)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	case *make != "":
		err := makePkg(*make)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	case *newSrc != "":
		err := addSrc(*newSrc)
		if err != nil {
			fmt.Println(err.Error())
		}
		return
	case *list:
		fmt.Println(listSrc())
		return
	default:
		fmt.Println("-help")
	}
}
