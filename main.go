package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/asdine/storm"
	"github.com/mitchellh/go-homedir"
	"github.com/raedahgroup/fileman/config"
	ctl "github.com/raedahgroup/fileman/http"
	"github.com/raedahgroup/fileman/storage"
	"github.com/raedahgroup/fileman/storage/bolt"
	"github.com/raedahgroup/fileman/users"
	"github.com/spf13/afero"
	"log"
	"math/rand"
	"net"
	"net/http"
)
var (
	host     = flag.String("host", "127.0.0.1", "TCP host to listen to")
	port     = flag.String("port", "8081", "TCP port to listen to")
	homeDir     = flag.String("dir", "", "Folder's user have been uploads")
	baseURL      = flag.String("baseurl", "", "Directory to serve static files from")
)

func main() {
	flag.Parse();
	//folder home os
	home, err := homedir.Dir();
	//home file man
	homeFileMan := home + "/fileman/";
	checkErr(err);
	if(*homeDir == ""){
		*homeDir = homeFileMan + "/uploads";
	}

	appfs := afero.NewOsFs();
	//create forder fileman
	appfs.MkdirAll(*homeDir, 0755);
	config := config.ConfigState{
		DatabasePath: homeFileMan + "fileman.db",
		Port: "4000",
		RootPath: *homeDir,
		JWTKEY: "fileman@2019",
		BaseURL: *baseURL,
	}
	db, err := storm.Open(config.DatabasePath)
	checkErr(err)
	defer db.Close()
	store, err := bolt.NewStorage(db);
	go createAdminDemo(store)
	var listener net.Listener
	listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", *host, *port))
	checkErr(err)
	handler, err := ctl.NewHandler(store, config)
	log.Println("Listening on", listener.Addr().String())
	if err := http.Serve(listener, handler); err != nil {
		log.Fatal(err)
	}

}
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func createAdminDemo(store *storage.Storage) error {
	pwd, err := users.HashPwd("123456")
	if err != nil {
		fmt.Println("hash password", err)
	}
	user := &users.User{
		Username:     "admin",
		Password:     pwd,
		LockPassword: true,
		Scope: "/",
		Perm: users.Permissions{
			Admin:    true,
			Execute:  true,
			Create:   true,
			Rename:   true,
			Modify:   true,
			Delete:   true,
			Share:    true,
			Download: true,
		},
	}
	err = store.Users.Save(user);
	return  err

}
func randStr(len int) string {
	buff := make([]byte, len)
	rand.Read(buff)
	str := base64.StdEncoding.EncodeToString(buff)
	// Base 64 can be longer than len
	return str[:len]
}
