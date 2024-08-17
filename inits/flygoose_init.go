package inits

import (
	"errors"
	"flygoose/configs"
	"flygoose/datasource"
	"flygoose/pkg/models"
	"flygoose/pkg/tlog"
	"flygoose/pkg/tools"
	"flygoose/web/controllers/admin"
	"flygoose/web/controllers/flygoose"
	"flygoose/web/daos"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/cors"
	"github.com/kataras/iris/v12/mvc"
	"go.uber.org/zap"
	"path/filepath"
	"time"
)

type FlygooseApp struct {
	Cfg    *configs.Config
	Engine *iris.Application
}

func NewFlygooseApp(cfg *configs.Config) *FlygooseApp {
	app := new(FlygooseApp)
	app.Cfg = cfg
	//app.Engine = iris.Default()
	app.Engine = iris.New()
	app.Engine.UseRouter(cors.New().
		ExtractOriginFunc(cors.DefaultOriginExtractor).
		ReferrerPolicy(cors.NoReferrerWhenDowngrade).
		AllowOriginFunc(cors.AllowAnyOrigin).
		AllowHeaders("token", "content-type", "Authorization", "x-requested-with").
		ExposeHeaders("token", "content-type", "Authorization", "x-requested-with").
		Handler())
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
	m.Engine.HandleDir("/static", abcStaticDir) //http://192.168.1.6:29090/img/aa.jpg
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

	ad := daos.NewAccessDao()
	cnt, err := ad.CountUsername(configs.Flygoose_Admin_Phone)
	if err != nil {
		fmt.Printf("查询数据库中是否存在默认账户错误. err: %w\n", err)
		panic(errors.New("查询数据库中是否存在默认账户错误"))
	}

	if cnt == 0 {
		initAccountUser := models.Admin{
			Phone:      configs.Flygoose_Admin_Phone,
			Password:   "21232f297a57a5a743894a0e4a801fc3",
			Nicker:     "admin",
			Avatar:     "https://img-hello-world.oss-cn-beijing.aliyuncs.com/imgs/b3e9e8fb50b3eba780178256a21234ec.jpg",
			CreateTime: time.Now(),
			ValidTime:  time.Now(),
			LoginTime:  time.Now(),
			Status:     models.AdminStatusNormal,
		}
		if ad.Create(&initAccountUser) != nil {
			fmt.Printf("创建数默认账户错误. err: %w\n", err)
			panic(errors.New("创建数默认账户错误"))
		}
	}
}

func (m *FlygooseApp) initAutoMigrate() {

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
	v1 := m.Engine.Party(configs.Flygoose_Url_Prefix)
	{
		mvc.New(v1.Party("/site")).Handle(flygoose.NewSiteController())
		mvc.New(v1.Party("/blog")).Handle(flygoose.NewBlogController())
		mvc.New(v1.Party("/special")).Handle(flygoose.NewSpecialController())

		mvc.New(v1.Party("/health")).Handle(admin.NewHealthController())                 //健康检查
		mvc.New(v1.Party("/admin/access")).Handle(admin.NewAccessController())           //访问相关接口
		mvc.New(v1.Party("/admin/blog")).Handle(admin.NewBlogController())               //博客相关接口
		mvc.New(v1.Party("/admin/link")).Handle(admin.NewLinkController())               //友链相关接口
		mvc.New(v1.Party("/admin/site")).Handle(admin.NewSiteController())               //网站信息相关接口
		mvc.New(v1.Party("/admin/category")).Handle(admin.NewCategoryController())       //博客分类相关接口
		mvc.New(v1.Party("/admin/notice")).Handle(admin.NewNoticeController())           //公告分类相关接口
		mvc.New(v1.Party("/admin/special")).Handle(admin.NewSpecialController())         //专栏相关接口
		mvc.New(v1.Party("/admin/file")).Handle(admin.NewFileController())               //文件相关接口
		mvc.New(v1.Party("/admin/banner")).Handle(admin.NewBannerController())           //轮播图相关接口
		mvc.New(v1.Party("/admin/workStation")).Handle(admin.NewWorkStationController()) //统计数据
	}
}
