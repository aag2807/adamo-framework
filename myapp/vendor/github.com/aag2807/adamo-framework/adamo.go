package adamo

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/aag2807/adamo-framework/render"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

const version = "0.0.1"

type Adamo struct {
	AppName  string
	Debug    bool
	Version  string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	RootPath string
	config   config
	Routes   *chi.Mux
	Render   *render.Render
	JetViews *jet.Set
}

type config struct {
	port     string
	renderer string
}

func (a *Adamo) New(rootPath string) error {
	pathConfig := initPaths{
		RootPath: rootPath,
		folderNames: []string{
			"config",
			"handlers",
			"migrations",
			"views",
			"data",
			"public",
			"tmp",
			"logs",
			"middleware",
		},
	}

	err := a.Init(pathConfig)
	if err != nil {
		return err
	}

	err = a.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	//read .dotenv
	err = godotenv.Load(fmt.Sprintf("%s/.env", rootPath))
	if err != nil {
		return err
	}

	//create loggers
	infoLog, errorLog := a.startLoggers()
	a.InfoLog = infoLog
	a.ErrorLog = errorLog

	//gets debug from .env
	a.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))

	a.Version = version
	a.RootPath = rootPath
	a.Routes = a.routes().(*chi.Mux)

	//set config values
	a.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
	}

	// configure jet set renderer
	jetViewPath := fmt.Sprintf("%s/views", rootPath)
	views := jet.NewSet(
		jet.NewOSFileSystemLoader(jetViewPath),
		jet.InDevelopmentMode(),
	)

	a.JetViews = views

	// config renderer
	a.createRenderer()

	return nil
}

func (a *Adamo) Init(p initPaths) error {
	root := p.RootPath

	for _, path := range p.folderNames {
		folderCreationDir := fmt.Sprintf("%s/%s", root, path) // equals root + "/" + path
		err := a.CreateDirNoNotExist(folderCreationDir)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Adamo) ListenAndServe() {
	timeoutSpan := 30 * time.Second

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", a.config.port),
		ErrorLog:     a.ErrorLog,
		Handler:      a.Routes,
		IdleTimeout:  timeoutSpan,
		ReadTimeout:  timeoutSpan,
		WriteTimeout: 600 * time.Second,
	}

	a.InfoLog.Printf("Listening on port %s", a.config.port)
	err := srv.ListenAndServe()
	a.ErrorLog.Fatal(err)
}

func (a *Adamo) checkDotEnv(path string) error {
	dotEnvPath := fmt.Sprintf("%s/.env", path) // equivalent to path + "/.env"
	err := a.CreateFileNoExists(dotEnvPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *Adamo) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

func (a *Adamo) createRenderer() {
	myRenderer := render.Render{
		Renderer: a.config.renderer,
		RootPath: a.RootPath,
		Port:     a.config.port,
		JetViews: a.JetViews,
	}

	a.Render = &myRenderer
}
