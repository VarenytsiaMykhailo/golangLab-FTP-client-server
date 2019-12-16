package main

import (
	"flag"
	filedriver "github.com/goftp/file-driver"
	"github.com/goftp/server"
	"log"
)

type myAuth struct{}

func (a *myAuth) CheckPasswd(user, passwd string) (bool, error) { //корявая бд
	if user == "admin" && passwd == "12345" {
		return true, nil
	} else if user == "dragon" && passwd == "qwerty" {
		return true, nil
	}
	return false, nil
}

func main() {
	var (
		root = flag.String("root", "./server/root", "Root directory to server")
		//user = flag.String("user", "admin", "Username for login")
		//pass = flag.String("pass", "admin123", "Password for login")
		port = flag.Int("port", 9876, "Port")
		host = flag.String("host", "localhost", "Host") //вместо localhost 185.20.227.83
	)
	flag.Parse()
	if *root == "" {
		log.Fatalf("Please set a root to serve with -root")
	}

	factory := &filedriver.FileDriverFactory{
		RootPath: *root,
		Perm:     server.NewSimplePerm("user", "group"),
	}

	opts := &server.ServerOpts{
		Factory:  factory,
		Port:     *port,
		Hostname: *host,
		//Auth:     &server.SimpleAuth{Name: *user, Password: *pass},
		Auth: &myAuth{},
	}

	log.Printf("Starting ftp server on %v:%v", opts.Hostname, opts.Port)
	//log.Printf("Username %v, Password %v", *user, *pass)
	server := server.NewServer(opts)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
