package main

import (
	"fmt"
	"github.com/huge-kumo/dw-10th/pkg/cors"
	recover2 "github.com/huge-kumo/dw-10th/pkg/recover"
	"github.com/kataras/iris/v12"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	registerSignal()
	app := iris.New()

	app.Use(recover2.New())
	_registerHandler(app)
}

func registerSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP:
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("退出", s)
				os.Exit(0)
			default:
				fmt.Println("other", s)
			}
		}
	}()

	f, err := os.Create("server_pid.txt")
	if err != nil {
		panic(err)
	}

	if _, err = f.WriteString(strconv.Itoa(os.Getpid())); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	if err = f.Close(); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

}

func _registerHandler(app *iris.Application) {
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
	})

	m := app.Party("/", crs).AllowMethods(iris.MethodOptions)
	{
		m.Handle("GET", "/", func(ctx iris.Context) {
			ctx.WriteString("ok")
		})
		m.HandleDir("/", "../assets")

	}
	app.Run(iris.Addr("0.0.0.0:8030"), iris.WithoutServerError(iris.ErrServerClosed))
}
