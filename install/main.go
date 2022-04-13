package main

import (
	"flag"
	"github.com/beego/beego/v2/server/web"
	"github.com/sondertara/tara-wiki/install/storage"
	"log"
	"os"
	"path/filepath"
)

// install

var (
	port = flag.String("port", "8091", "please input listen port")
)

func main() {
	flag.Parse()

	_, err := os.Stat(filepath.Join(storage.RootDir, "./install.lock"))
	if err == nil || !os.IsNotExist(err) {
		log.Println("MM-Wiki already installed!")
		os.Exit(1)
	}

	//web.BConfig.RunMode = "prod"
	web.Run(":" + *port)
}
