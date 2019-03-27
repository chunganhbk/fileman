package main

import (
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
	"net"
	"net/http"
)
var (
	host     = flag.String("host", "127.0.0.1", "TCP host to listen to")
	port     = flag.String("port", "8081", "TCP port to listen to")
	homeDir     = flag.String("dir", "", "Set folder's user have been uploads")
	baseURL      = flag.String("baseurl", "", "Directory to serve static files from")
)

func main() {
	flag.Parse();
	//folder home os
	home, err := homedir.Dir();
	//home file man
	homeFileMan := home + "/fileman/";
	checkErr(err);
	DatabasePath := homeFileMan + "fileman.db"
	appfs := afero.NewOsFs();
	//create forder fileman
	appfs.MkdirAll(*homeDir, 0755);
	db, err := storm.Open(DatabasePath);
	checkErr(err)
	defer db.Close()
	store, err := bolt.NewStorage(db);
	checkErr(err)
	cfServer, err := store.Config.GetServer();
	if err != nil{
		forderUpload := homeFileMan + "/uploads";
		cfServer = &config.Server{
			RootPath: forderUpload,
			JWTKEY: "fileman@2019",
			BaseURL: *baseURL,
		}
		appfs.MkdirAll(forderUpload, 755);
		store.Config.SaveServer(cfServer);
		go createUserDefault(store);
	}
	if(*homeDir != ""){
		//Check forder permission
		if err = appfs.MkdirAll(*homeDir, 755); err == nil {
			cfServer.RootPath = *homeDir;
			store.Config.SaveServer(cfServer);
			fmt.Println("Forder uploads was changed!");
		}else{
			fmt.Println("Error:", err);
		}
	}else {
		var listener net.Listener
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", *host, *port))
		checkErr(err)
		handler, _ := ctl.NewHandler(store, cfServer)
		log.Println("Listening on", listener.Addr().String())
		if err := http.Serve(listener, handler); err != nil {
			log.Fatal(err)
		}
	}
}
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func createUserDefault(store *storage.Storage) error {
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
