package inits

import (
	"errors"
	"flygoose/configs"
	"flygoose/datasource"
	"flygoose/pkg/models"
	"flygoose/pkg/tlog"
	"flygoose/pkg/tools"
	"flygoose/web/controllers/flygoose"
	"fmt"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"go.uber.org/zap"
	"path/filepath"
)

type FlygooseApp struct {
	Cfg    *configs.Config
	Engine *iris.Application
}

func NewFlygooseApp(cfg *configs.Config) *FlygooseApp {
	app := new(FlygooseApp)
	app.Cfg = cfg
	app.Engine = iris.Default()
	return app
}

func (m *FlygooseApp) Start() {
	m.InitDir()
	m.initLog()
	m.initDB()
	m.initRouter()
	m.run()
}

func (m *FlygooseApp) run() {
	err := m.Engine.Listen(fmt.Sprintf(":%d", m.Cfg.Http.Port), iris.WithoutBodyConsumptionOnUnmarshal, iris.WithOptimizations)
	if err != nil {
		tlog.Error2("FlygooseApp run 启动出错:", err)
		panic(err)
	}
}

func (m *FlygooseApp) InitDir() {
	//初始化可执行文件所在目录
	m.Cfg.ExecuteDir = tools.GetExecuteDir()
	m.Cfg.StaticDir = "/static"
	m.Cfg.StaticImgDir = "/static/img"

	//创建相应的文件目录
	tools.CreateDir(filepath.Join(m.Cfg.ExecuteDir, m.Cfg.StaticDir))
	tools.CreateDir(filepath.Join(m.Cfg.ExecuteDir, m.Cfg.StaticImgDir))

	var abcStaticDir = filepath.Join(m.Cfg.ExecuteDir, m.Cfg.StaticDir)
	m.Engine.HandleDir(configs.Flygoose_Url_Prefix+"/static", abcStaticDir)
}

func (m *FlygooseApp) initLog() {
	tlog.InitLog()
}

func (m *FlygooseApp) initDB() {
	var cfg = m.Cfg
	if m.Cfg.Database.Driver == configs.DbDriverMySQL {
		datasource.InitMySql(cfg)
	} else if m.Cfg.Database.Driver == configs.DbDriverPostgreSQL {
		datasource.InitPostgreSQL(cfg)
	} else {
		panic(errors.New("暂不支持其它数据库"))
	}

	m.initAutoMigrate()
}

func (m *FlygooseApp) initAutoMigrate() {
	//if m.Cfg.EnvMode != configs.EnvModeTypeDevelop {
	//	return
	//}

	err := datasource.GetMasterDB().AutoMigrate(
		&models.Admin{},
		&models.Blog{},
		&models.Category{},
		&models.Link{},
		&models.Notice{},
		&models.Section{},
		&models.Site{},
		&models.Special{},
		&models.Webmaster{},
		&models.Banner{},
	)

	if err != nil {
		tlog.Error("映射models出错:", zap.Error(err))
	}
}

func (m *FlygooseApp) initRouter() {
	// 跨域
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // allows everything, use that to change the hosts.
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"*"},
		//在这里需要加Authorization字段，不然js无法添加自定义header
		ExposedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		AllowCredentials: true,
	})

	v1 := m.Engine.Party(configs.Flygoose_Url_Prefix, crs).AllowMethods(iris.MethodOptions)
	{
		mvc.New(v1.Party("/site")).Handle(flygoose.NewSiteController())
		mvc.New(v1.Party("/blog")).Handle(flygoose.NewBlogController())
		mvc.New(v1.Party("/special")).Handle(flygoose.NewSpecialController())
		mvc.New(v1.Party("/health")).Handle(flygoose.NewHealthController())
	}
}
