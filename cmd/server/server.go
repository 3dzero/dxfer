package main

import (
	"embed"
	"fmt"
	"github.com/3dzero/dxfer/internal/handle"
	"github.com/alexsuslov/godotenv"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"syscall"
)

var version = "developer preview"

//go:embed static/*
var Static embed.FS

//go:embed templates/*
var Templates embed.FS



func main() {
	log.Printf("Starting " + getMessage())

	// load env
	if err := godotenv.Load(".env"); err != nil {
		log.Println("warrning load env", err)
	}

	//parse templates
	Templates := template.Must(template.ParseFS(Templates, "templates/*.tmpl"))

	//router
	r := mux.NewRouter()

	r.HandleFunc("/", handle.Home(Templates)).Methods("GET")


	//r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(Static))))

	//data, err := Static.ReadDir("static/css/")
	//if err != nil{
	//	panic(err)
	//}
	//
	//for _,v :=range data{
	//	logrus.WithField("ls", v.Name()).Info("debug")
	//}


	r.PathPrefix("/static/").Handler(http.FileServer(http.FS(Static)))



	//staticFS := http.FS(Static)
	//fs := http.FileServer(staticFS)
	//http.Handle("/static/", fs)

	host:=Env("HTTP_HOST", "0.0.0.0")
	port:=Env("HTTP_PORT", "8080")

	httpAddr := fmt.Sprintf("%s:%s",
		host,
		port)
	log.Println("listen", httpAddr)

	server := http.Server{Addr: httpAddr, Handler: r}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	//d := dxf.NewDrawing()
	//d.Header().LtScale = 100.0
	//d.AddLayer("sq", dxf.DefaultColor, dxf.DefaultLineType, true)
	//h := 100.0
	////w:=100.0
	//d.ChangeLayer("sq")
	//d.Line(0, 0, 0, 0, h, 0)
	//err := d.SaveAs("test.dxf")
	//if err != nil {
	//	log.Fatal(err)
	//}

}

func getMessage() string {
	return os.Getenv("MESSAGE") + "(" + version + ")"
}

func Env(key string, def string) string {
	v, _ := syscall.Getenv(key)
	if v == "" {
		return def
	}
	return v

}
