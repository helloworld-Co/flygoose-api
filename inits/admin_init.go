package inits

import (
	"errors"
	"flygoose/configs"
	"flygoose/datasource"
	"flygoose/pkg/models"
	"flygoose/pkg/tlog"
	"flygoose/pkg/tools"
	"flygoose/web/controllers/admin"
	"fmt"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"go.uber.org/zap"
	"path/filepath"
)

type AdminApp struct {
	Cfg    *configs.Config
	Engine *iris.Application
}

func NewAdminApp(cfg *configs.Config) *AdminApp {
	app := new(AdminApp)
	app.Cfg = cfg
	app.Engine = iris.Default()
	return app
}

func (m *AdminApp) Start() {
	m.InitDir()
	m.initDB()
	m.initRouter()
	m.run()
}

func (m *AdminApp) run() {
	var abcStaticDir = filepath.Join(m.Cfg.ExecuteDir, m.Cfg.StaticDir)
	m.Engine.HandleDir(configs.Admin_Url_Prefix+"/static", abcStaticDir)

	err := m.Engine.Listen(fmt.Sprintf(":%d", m.Cfg.Http.Port), iris.WithoutBodyConsumptionOnUnmarshal, iris.WithOptimizations)
	if err != nil {
		panic(err)
	}
}

func (m *AdminApp) InitDir() {
	//初始化可执行文件所在目录
	m.Cfg.ExecuteDir = tools.GetExecuteDir()
	m.Cfg.StaticDir = "/static"
	m.Cfg.StaticImgDir = "/static/img"

	//创建相应的文件目录
	tools.CreateDir(filepath.Join(m.Cfg.ExecuteDir, m.Cfg.StaticDir))
	tools.CreateDir(filepath.Join(m.Cfg.ExecuteDir, m.Cfg.StaticImgDir))
}

func (m *AdminApp) initDB() {
	var cfg = m.Cfg
	if m.Cfg.Database.Driver == configs.DbDriverMySQL {
		datasource.InitMySql(cfg)
	} else if m.Cfg.Database.Driver == configs.DbDriverPostgreSQL {
		datasource.InitPostgreSQL(cfg)
	} else {
		panic(errors.New("配置文件中database下driver字段配置错误"))
	}

	//自动将model映射成表
	m.initAutoMigrate()
}

func (m *AdminApp) initAutoMigrate() {

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

func (m *AdminApp) initRouter() {
	// 跨域
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // allows everything, use that to change the hosts.
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"*"},
		//在这里需要加Authorization字段，不然js无法添加自定义header
		ExposedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		AllowCredentials: true,
	})

	v8 := m.Engine.Party(configs.Admin_Url_Prefix, crs).AllowMethods(iris.MethodOptions)
	{
		mvc.New(v8.Party("/health")).Handle(admin.NewHealthController())           //健康检查
		mvc.New(v8.Party("/access")).Handle(admin.NewAccessController())           //访问相关接口
		mvc.New(v8.Party("/blog")).Handle(admin.NewBlogController())               //博客相关接口
		mvc.New(v8.Party("/link")).Handle(admin.NewLinkController())               //友链相关接口
		mvc.New(v8.Party("/site")).Handle(admin.NewSiteController())               //网站信息相关接口
		mvc.New(v8.Party("/category")).Handle(admin.NewCategoryController())       //博客分类相关接口
		mvc.New(v8.Party("/notice")).Handle(admin.NewNoticeController())           //公告分类相关接口
		mvc.New(v8.Party("/special")).Handle(admin.NewSpecialController())         //专栏相关接口
		mvc.New(v8.Party("/file")).Handle(admin.NewFileController(m.Cfg))          //文件相关接口
		mvc.New(v8.Party("/banner")).Handle(admin.NewBannerController())           //轮播图相关接口
		mvc.New(v8.Party("/workStation")).Handle(admin.NewWorkStationController()) //统计数据
	}
}
