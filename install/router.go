package main

import (
	"github.com/beego/beego/v2/server/web"
	"net/http"
	"os"
	"path/filepath"
	"github.com/sondertara/tara-wiki/app/utils"
	"github.com/sondertara/tara-wiki/install/controllers"
	"github.com/sondertara/tara-wiki/install/storage"
)

func init() {

	storage.InstallDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	storage.RootDir = filepath.Join(storage.InstallDir, "../")

	web.AppConfig.Set("sys.name", "mm-wiki-installer")
	web.BConfig.AppName, _ = web.AppConfig.String("sys.name")
	web.BConfig.ServerName, _ = web.AppConfig.String("sys.name")

	// set static path
	web.SetStaticPath("/static/", filepath.Join(storage.InstallDir, "../static"))
	// views path
	web.BConfig.WebConfig.ViewsPath = filepath.Join(storage.InstallDir, "../views/")

	// session
	web.BConfig.WebConfig.Session.SessionName = "mmwikiinstallssid"
	web.BConfig.WebConfig.Session.SessionOn = true

	// router
	web.BConfig.WebConfig.AutoRender = false
	web.BConfig.RouterCaseSensitive = false

	web.AutoRouter(&controllers.InstallController{})
	web.Router("/", &controllers.InstallController{}, "*:Index")
	web.ErrorHandler("404", http_404)
	web.ErrorHandler("500", http_500)

	// add template func
	web.AddFuncMap("dateFormat", utils.NewDate().Format)

}

func http_404(rs http.ResponseWriter, req *http.Request) {
	rs.Write([]byte("404 not found!"))
}

func http_500(rs http.ResponseWriter, req *http.Request) {
	rs.Write([]byte("500 server error!"))
}
