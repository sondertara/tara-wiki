package app

import (
	"flag"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/fatih/color"
	"github.com/go-ego/riot/types"
	"github.com/sondertara/go-activerecord/mysql"
	"log"
	"os"
	"path"
	"path/filepath"
	"github.com/sondertara/tara-wiki/app/models"
	"github.com/sondertara/tara-wiki/app/utils"
	"github.com/sondertara/tara-wiki/app/work"
	"github.com/sondertara/tara-wiki/global"
	"time"
)

var (
	defaultConf = "conf/mm-wiki.conf"

	confPath = flag.String("conf", "", "please set mm-wiki conf path")

	version = flag.Bool("version", false, "mm-wiki version")

	upgrade = flag.Bool("upgrade", false, "mm-wiki upgrade")

	Version = global.SYSTEM_VERSION

	CopyRight = web.Str2html(global.SYSTEM_COPYRIGHT)

	StartTime = int64(0)

	RootDir = ""

	DocumentAbsDir = ""

	MarkdownAbsDir = ""

	ImageAbsDir = ""

	AttachmentAbsDir = ""

	SearchIndexAbsDir = ""
)

func init() {
	initFlag()
	poster()
	initConfig()
	initDB()
	checkUpgrade()
	initDocumentDir()
	//initSearch()
	//initWork()
	StartTime = time.Now().Unix()
}

// init flag
func initFlag() {
	flag.Parse()
	// --version
	if *version == true {
		fmt.Printf(Version)
		os.Exit(0)
	}
}

// poster logo
func poster() {
	fg := color.New(color.FgBlue)
	logo := `
                                            _   _      _ 
 _ __ ___    _ __ ___           __      __ (_) | | __ (_)
| '_ ' _ \  | '_ ' _ \   _____  \ \ /\ / / | | | |/ / | |
| | | | | | | | | | | | |_____|  \ V  V /  | | |   <  | |
|_| |_| |_| |_| |_| |_|           \_/\_/   |_| |_|\_\ |_|
` +
		"Author: phachon\r\n" +
		"Version: " + Version + "\r\n" +
		"Link: https://github.com/phachon/mm-wiki"
	fg.Println(logo)
}

// init beego config
func initConfig() {

	RootDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println("init config error: " + err.Error())
		os.Exit(1)
	}
	confFile := *confPath
	if *confPath == "" {
		confFile = filepath.Join(RootDir, defaultConf)
	}
	ok, _ := utils.NewFile().PathIsExists(confFile)
	if ok == false {
		log.Println("conf file " + confFile + " not exists!")
		os.Exit(1)
	}
	// init config file
	web.LoadAppConfig("ini", confFile)

	// init name
	web.AppConfig.Set("sys.name", "mm-wiki")
	web.BConfig.AppName, _ = web.AppConfig.String("sys.name")
	web.BConfig.ServerName, _ = web.AppConfig.String("sys.name")

	// set static path
	web.SetStaticPath("/static/", filepath.Join(RootDir, "./static"))
	// views path
	web.BConfig.WebConfig.ViewsPath = filepath.Join(RootDir, "./views/")

	// session
	//web.BConfig.WebConfig.Session.SessionProvider = "memory"
	//web.BConfig.WebConfig.Session.SessionProviderConfig = ".session"
	//web.BConfig.WebConfig.Session.SessionName = "mmwikissid"
	//web.BConfig.WebConfig.Session.SessionOn = true

	// log
	logConfigs, err := web.AppConfig.GetSection("log")
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	for adapter, config := range logConfigs {
		logs.SetLogger(adapter, config)
	}
	logs.SetLogFuncCall(true)
}

//init db
func initDB() {
	host, _ := web.AppConfig.String("db::host")
	port, _ := web.AppConfig.Int("db::port")
	user, _ := web.AppConfig.String("db::user")
	pass, _ := web.AppConfig.String("db::pass")
	dbname, _ := web.AppConfig.String("db::name")
	dbTablePrefix, _ := web.AppConfig.String("db::table_prefix")
	maxIdle, _ := web.AppConfig.Int("db::conn_max_idle")
	maxConn, _ := web.AppConfig.Int("db::conn_max_connection")
	models.G = mysql.NewDBGroup("default")
	cfg := mysql.NewDBConfigWith(host, port, dbname, user, pass)
	cfg.MaxIdleConns = maxIdle
	cfg.MaxOpenConns = maxConn
	cfg.TablePrefix = dbTablePrefix
	cfg.TablePrefixSqlIdentifier = "__PREFIX__"
	err := models.G.Regist("default", cfg)
	if err != nil {
		logs.Error(fmt.Errorf("database error:%s,with config : %v", err, cfg))
		os.Exit(1)
	}
}

// init document dir
func initDocumentDir() {
	docRootDir, _ := web.AppConfig.String("document::root_dir")
	if docRootDir == "" {
		logs.Error("document root dir " + docRootDir + " is not empty!")
		os.Exit(1)
	}
	ok, _ := utils.File.PathIsExists(docRootDir)
	if !ok {
		logs.Error("document root dir " + docRootDir + " is not exists!")
		os.Exit(1)
	}

	documentAbsDir, err := filepath.Abs(docRootDir)
	if err != nil {
		logs.Error("document root dir " + docRootDir + " is error!")
		os.Exit(1)
	}

	DocumentAbsDir = documentAbsDir

	// markdown save dir
	markDownAbsDir := path.Join(documentAbsDir, "markdowns")
	// image save dir
	imagesAbsDir := path.Join(documentAbsDir, "images")
	// attachment save dir
	attachmentAbsDir := path.Join(documentAbsDir, "attachment")
	// search index dir
	searchIndexAbsDir := path.Join(documentAbsDir, "search-index")

	MarkdownAbsDir = markDownAbsDir
	ImageAbsDir = imagesAbsDir
	AttachmentAbsDir = attachmentAbsDir
	SearchIndexAbsDir = searchIndexAbsDir

	dirList := []string{MarkdownAbsDir, ImageAbsDir, AttachmentAbsDir, SearchIndexAbsDir}
	// create dir
	for _, dir := range dirList {
		ok, _ = utils.File.PathIsExists(dir)
		if !ok {
			err := os.Mkdir(dir, 0777)
			if err != nil {
				logs.Error("create document dir "+dir+" error=%s", err.Error())
				os.Exit(1)
			}
		}
	}

	// utils document
	utils.Document.MarkdownAbsDir = markDownAbsDir
	utils.Document.DocumentAbsDir = documentAbsDir

	web.SetStaticPath("/images/", ImageAbsDir)
	// todo
	web.SetStaticPath("/images/:space_id/:document_id/", ImageAbsDir)
}

// check upgrade
func checkUpgrade() {
	if *upgrade == true {
		logs.Info("Start checking whether MM-Wiki needs upgrading.")
		var versionDb = "v0.0.0"
		versionConf := models.ConfigModel.GetConfigValueByKey(models.ConfigKeySystemVersion, "v0.0.0")
		if versionConf != "" {
			versionDb = versionConf
		}
		logs.Info("MM-Wiki Database version：" + versionDb)
		logs.Info("MM-Wiki Now version: " + Version)

		if versionDb == Version {
			logs.Info("MM-Wiki does not need updating.")
		} else {
			logs.Info("MM-Wiki start upgrading.")
			err := models.UpgradeModel.Start(versionDb)
			if err != nil {
				logs.Error("MM-Wiki upgrade failed.")
				os.Exit(1)
			}
			logs.Info("MM-Wiki upgrade finish.")
		}
		os.Exit(0)
	}
}

func initSearch() {

	gseFile := filepath.Join(RootDir, "docs/search_dict/dictionary.txt")
	stopFile := filepath.Join(RootDir, "docs/search_dict/stop_tokens.txt")
	ok, _ := utils.File.PathIsExists(gseFile)
	if !ok {
		logs.Error("search dict file " + gseFile + " is not exists!")
		os.Exit(1)
	}
	ok, _ = utils.File.PathIsExists(stopFile)
	if !ok {
		logs.Error("search stop dict file " + stopFile + " is not exists!")
		os.Exit(1)
	}
	global.DocSearcher.Init(types.EngineOpts{
		UseStore:    true,
		StoreFolder: SearchIndexAbsDir,
		Using:       3,
		//GseDict:       "zh",
		GseDict:       gseFile,
		StopTokenFile: stopFile,
		IndexerOpts: &types.IndexerOpts{
			IndexType: types.LocsIndex,
		},
	})
}

func initWork() {
	work.DocSearchWorker.Start()
}
