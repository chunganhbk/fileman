package main

import (
	"encoding/base64"
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

func main() {

	home, err := homedir.Dir();
	checkErr(err);
	//home file man
	//windows
	homeFileMan := home + "/fileman/";
	appfs := afero.NewOsFs();
	//create forder fileman
	appfs.MkdirAll(homeFileMan, 0755);
	appfs.MkdirAll(homeFileMan + "/api/resources/", 0775)
	config := config.ConfigState{
		DatabasePath: homeFileMan + "fileman.db",
		Port: "4000",
		RootPath: homeFileMan,
		JWTKEY: "fileman@2019",
	}
	db, err := storm.Open(config.DatabasePath)
	checkErr(err)
	defer db.Close()
	store, err := bolt.NewStorage(db);
	go createAdminDemo(store)
	adr := "127.0.0.1:"   + config.Port
	var listener net.Listener
	listener, err = net.Listen("tcp", adr)
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
	fmt.Println(err)
	return  err

}
func randStr(len int) string {
	buff := make([]byte, len)
	rand.Read(buff)
	str := base64.StdEncoding.EncodeToString(buff)
	// Base 64 can be longer than len
	return str[:len]
}
