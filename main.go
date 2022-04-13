package main

import (
	"github.com/beego/beego/v2/server/web"
	_ "github.com/beego/beego/v2/server/web/session/memcache"
	_ "github.com/beego/beego/v2/server/web/session/redis"
	_ "github.com/beego/beego/v2/server/web/session/redis_cluster"
	_ "github.com/sondertara/tara-wiki/app"
)

func main() {
	web.Run()
}
