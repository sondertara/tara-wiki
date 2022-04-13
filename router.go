package main

import (
	"github.com/beego/beego/v2/server/web"
	"html/template"
	"net/http"
	"github.com/sondertara/tara-wiki/app"
	"github.com/sondertara/tara-wiki/app/controllers"
	systemControllers "github.com/sondertara/tara-wiki/app/modules/system/controllers"
	"github.com/sondertara/tara-wiki/app/utils"
)

func init() {
	initRouter()
}

func initRouter() {
	// router
	web.BConfig.WebConfig.AutoRender = false
	web.BConfig.RouterCaseSensitive = false

	web.Router("/", &controllers.MainController{}, "*:Index")
	web.Router("/author", &controllers.AuthorController{}, "*:Index")
	web.AutoRouter(&controllers.AuthorController{})
	web.AutoRouter(&controllers.MainController{})
	web.AutoRouter(&controllers.SpaceController{})
	web.AutoRouter(&controllers.CollectionController{})
	web.AutoRouter(&controllers.FollowController{})
	web.AutoRouter(&controllers.UserController{})
	web.AutoRouter(&controllers.DocumentController{})
	web.AutoRouter(&controllers.PageController{})
	web.AutoRouter(&controllers.ImageController{})
	web.AutoRouter(&controllers.AttachmentController{})

	systemNamespace := web.NewNamespace("/system",
		web.NSAutoRouter(&systemControllers.MainController{}),
		web.NSAutoRouter(&systemControllers.ProfileController{}),
		web.NSAutoRouter(&systemControllers.UserController{}),
		web.NSAutoRouter(&systemControllers.RoleController{}),
		web.NSAutoRouter(&systemControllers.PrivilegeController{}),
		web.NSAutoRouter(&systemControllers.SpaceController{}),
		web.NSAutoRouter(&systemControllers.Space_UserController{}),
		web.NSAutoRouter(&systemControllers.LogController{}),
		web.NSAutoRouter(&systemControllers.EmailController{}),
		web.NSAutoRouter(&systemControllers.LinkController{}),
		web.NSAutoRouter(&systemControllers.AuthController{}),
		web.NSAutoRouter(&systemControllers.ConfigController{}),
		web.NSAutoRouter(&systemControllers.ContactController{}),
		web.NSAutoRouter(&systemControllers.StaticController{}),
	)
	web.AddNamespace(systemNamespace)

	web.ErrorHandler("404", http_404)
	web.ErrorHandler("500", http_500)

	// add template func
	web.AddFuncMap("dateFormat", utils.Date.Format)
}

func http_404(rw http.ResponseWriter, req *http.Request) {
	t, _ := template.New("404.html").ParseFiles(web.BConfig.WebConfig.ViewsPath + "/error/404.html")
	data := make(map[string]interface{})
	data["content"] = "page not found"
	data["copyright"] = app.CopyRight
	t.Execute(rw, data)
}

func http_500(rw http.ResponseWriter, req *http.Request) {
	t, _ := template.New("500.html").ParseFiles(web.BConfig.WebConfig.ViewsPath + "/error/500.html")
	data := make(map[string]interface{})
	data["content"] = "Server Error"
	data["copyright"] = app.CopyRight
	t.Execute(rw, data)
}
