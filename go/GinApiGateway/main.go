package main

import (
	"GinApiGateway/models"
	"GinApiGateway/routers"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "GinApiGateway/docs"
)
// @title           Prizoo Api Document
// @version         1.0
// @description     包含所有 api 的用法.

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
func main() {
	// //	Create log file
	// file, err := os.Create("gin.log")
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()

	// //	Set log file write in file and console
	// gin.DefaultWriter = io.MultiWriter(file, os.Stdout)

	router := gin.Default()

	//	Set trust proxy ip
	router.SetTrustedProxies([]string{"127.0.0.1"})

	//	Initial authentication
	models.InitFirebaseAuth()
	models.InitFirebaseStorage()

	//	Initial routers
	routers.InitApiRouters(router)

	//	Graceful shutdown
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	//	Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)		//	Make a quit channel to get os signal
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)	//	If receive shutdown signal, put this signal into quit channel
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)	//	Set timeout context with 5 seconds
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown: ", err)
	}
	log.Println("Server exiting")
}
