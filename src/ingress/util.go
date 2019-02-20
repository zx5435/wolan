package ingress

import (
	"os"
	"path/filepath"
	"html/template"
	"os/exec"
	"flag"
	"github.com/zx5435/wolan/src/log"
	"fmt"
	"io/ioutil"
	"errors"
)

func MkDirAll(dir string, perm os.FileMode) {
	log.Change("mkdir", dir)

	if err := os.MkdirAll(dir, perm); err != nil {
		log.Fatal(err, dir)
	}
}

func sameDir(filename string, perm os.FileMode) error {
	dir := filepath.Dir(filename)
	MkDirAll(dir, perm)

	return nil
}

func NginxReload() error {
	cmd := exec.Command("/bin/sh", "-c", "nginx -s reload")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// res
func getFile(filename string) string {
	file, _ := os.Open(tplDir + filename)

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Debugf("load file: %s, size: %d", filename, len(bytes))
	return string(bytes)
}
func fetchResource(filename string) ([]byte, error) {
	if filename == "" {
		return nil, errors.New("ngx: empty resource name")
	}

	file, _ := os.Open(tplDir + filename)

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func writeResource(filename string) error {
	if filename == "" {
		return errors.New("ngx: empty resource name")
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := sameDir(filename, 0700); err != nil {
			return err
		}
		fileCtn := getFile(filepath.Base(filename))
		return ioutil.WriteFile(filename, []byte(fileCtn), 0600)
	}

	return nil
}

// tpl file
func fileCreate(tpl *template.Template, fp string, data interface{}) error {
	log.Change("fileCreate", fp)
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		fn, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer fn.Close()

		return tpl.Execute(fn, data)
	}

	return os.ErrExist
}

func fileEdit(tpl *template.Template, fp string, data interface{}) error {
	fn, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer fn.Close()

	return tpl.Execute(fn, data)
}

func UsageAndExit(msg string) {
	if msg != "" {
		log.Error(msg)
	}
	fmt.Println()
	flag.Usage()
	fmt.Println(`Demo:
  wolan-ingress -s new -d www.test.com
  wolan-ingress -env=prod -s=new -d zx5435.com`)
	os.Exit(1)
}
