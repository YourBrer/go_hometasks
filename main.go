package main

import (
	"fmt"
	"net/http"
	"project/configs"
	"project/internal/auth"
	"project/internal/verify"
)

func main() {
	appConfig := configs.GetAppConfig()
	router := http.NewServeMux()
	auth.NewAuthHandler(router, &appConfig.Auth)
	verify.NewVerifyHandler(router, &appConfig.Mail)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", appConfig.ServerPort),
		Handler: router,
	}

	fmt.Printf("Server started on port %d\n", appConfig.ServerPort)

	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
