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
	"gopkg.in/yaml.v2"
	"log"
	"math/rand"
	"net"
	"net/http"
	"path"
)

func main() {

	err := loadConfig();
	db, err := storm.Open(config.State.DatabasePath)
	checkErr(err)
	defer db.Close()
	store, err := bolt.NewStorage(db);
	adr := "127.0.0.1:"   + config.State.Port
	var listener net.Listener
	listener, err = net.Listen("tcp", adr)
	checkErr(err)
	handler, err := ctl.NewHandler(store, config.State)
	log.Println("Listening on", listener.Addr().String())
	if err := http.Serve(listener, handler); err != nil {
		log.Fatal(err)
	}
	go createAdminDemo(store)
}
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func loadConfig () error {

	home, err := homedir.Dir();
	checkErr(err);
	// CONFIG
	err = config.Load(home +"/fileman/config/config.yaml")
	if(err != nil){
		appfs := afero.NewOsFs();
		path := home + "/fileman/";
		//create forder fileman
		appfs.MkdirAll(path, 0755);

		dataConfig := fmt.Sprintf(`
			debug: true
			database_path: %s
			port: 4000
			JWTKEY: %s
		`, path, randStr(10))
		fmt.Println("dataconfig", dataConfig)
		err := yaml.Unmarshal([]byte(dataConfig), config.State);
		fmt.Println("err", err)
		return err;
	}
	return nil;
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
	root := config.State.RootPath;
	appfs := afero.NewOsFs()
	appfs.MkdirAll(path.Dir(root + "/api/resources/"), 0775)
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
