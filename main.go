package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/gorilla/mux"
)

const (
	src    = "/home/souenzzo/.local/src"
	repo   = "/home/souenzzo/.local/pkg"
	repoDB = "/home/souenzzo/.local/pkg/pacgo.db.tar.gz"
	build  = "/tmp/pacgo"
)

type Origin struct {
	Src string `json:"src"`
}

func addSrcHandler(w http.ResponseWriter, r *http.Request) {
	var O Origin
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &O)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = addSrc(O.Src)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func addSrc(src string) error {
	cmd := exec.Command("git", "clone", "--depth=1", src)
	cmd.Dir = src
	return cmd.Run()
}

func listSrcHandler(w http.ResponseWriter, r *http.Request) {
	names := listSrc()
	n, err := json.Marshal(&names)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "%s", n)
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

func makePkgHandler(w http.ResponseWriter, r *http.Request) {
	pkgName := mux.Vars(r)["pkg"]
	err := makePkg(pkgName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	fmt.Fprintln(w, "Ok!")
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

func web(addr string) error {
	r := mux.NewRouter()
	r.HandleFunc("/src", addSrcHandler).Methods("POST")
	r.HandleFunc("/src", listSrcHandler).Methods("GET")
	r.HandleFunc("/make/{pkg}", makePkgHandler).Methods("POST")
	s := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	return s.ListenAndServe()
}
