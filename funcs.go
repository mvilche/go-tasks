package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/gomail.v2"
)

type Config struct {
	Mailpassword string
	Mailhost     string
	Mailuser     string
	Mailport     int
	Maildestinos string
	Ambiente     string
}

func fileGetAbsolutePath(path string) (string, os.FileInfo) {
	ret, err := filepath.Abs(path)
	if err != nil {
		LoggerError.Fatalf("Invalid file: %v", err)
	}

	f, err := os.Lstat(ret)
	if err != nil {
		LoggerError.Fatalf("File stats failed: %v", err)
	}

	return ret, f
}

func ReadConfig(configFile []string) Config {
	var c = strings.Join(opts.ConfigFile, "")
	var config Config
	if _, err := toml.DecodeFile(c, &config); err != nil {
		LoggerError.Println("Archivo contiene errores, se descarta su uso")

	}
	return config
}

func NowDate() string {

	current := time.Now()
	format := current.Format("02-01-2006")
	return format
}

func NowHora() string {

	current := time.Now()
	format := current.Format("15:04:05")
	return format
}

func SendEmail(config Config, estado string, body string) {

	m := gomail.NewMessage()
	m.SetHeader("From", config.Mailuser)
	m.SetHeader("To", config.Maildestinos)
	m.SetHeader("Subject", "Task ejecutado - "+estado+" - Ambiente: "+config.Ambiente)
	m.SetBody("text/html", body)
	//m.Attach("/home/Alex/lolcat.jpg")
	d := gomail.NewDialer(config.Mailhost, config.Mailport, config.Mailuser, config.Mailpassword)
	// disable ssl check
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		LoggerError.Println("Error al enviar el email", err)
	} else {
		LoggerInfo.Println("Notificacion de estado enviada por email")
	}
}

func checkIfFileExists(path string) bool {
	f, err := os.Stat(path)

	// check if path exists
	if err != nil {
		return false
	}

	return (!f.IsDir()) && f.Mode().IsRegular()
}

func checkIfDirectoryExists(path string) bool {
	f, err := os.Stat(path)

	// check if path exists
	if err != nil {
		return false
	}

	return f.IsDir()
}

func checkIfFileExistsAndOwnedByRoot(path string) bool {
	f, err := os.Stat(path)

	// check if path exists
	if err != nil {
		return false
	}

	// check if it is not a file
	if !f.Mode().IsRegular() {
		return false
	}

	uidS := fmt.Sprint(f.Sys().(*syscall.Stat_t).Uid)
	uid, err := strconv.Atoi(uidS)
	if err != nil {
		return false
	}

	if uid != 0 {
		return false
	}

	return true
}

func checkIfFileIsValid(f os.FileInfo, path string) bool {
	if f.IsDir() {
		return false
	}

	if f.Mode().IsRegular() {
		if f.Mode().Perm()&0022 == 0 {
			return true
		} else {
			return true
		}
	} else {
		LoggerInfo.Printf("Ignoring non regular file %s\n", path)
	}

	return false
}
