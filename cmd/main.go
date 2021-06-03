package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/VVaria/db-technopark/configs"
	"github.com/VVaria/db-technopark/internal/app/tools/middleware"
	"github.com/VVaria/db-technopark/internal/app/tools/databases"
)

func main() {
	postgresDB, err := database.NewPostgres(configs.Configs.GetPostgresConfig())
	if err != nil {
		log.Fatal(err)
	}

	//forumRepo := forumRepo.NewForumRepository(postgresDB.GetDatabase())
	//postRepo := postRepo.NewPostRepository(postgresDB.GetDatabase())
	//serviceRepo := serviceRepo.NewServiceRepository(postgresDB.GetDatabase())
	//threadRepo := threadRepo.NewThreadRepository(postgresDB.GetDatabase())
	//userRepo := userRepo.NewUserRepository(postgresDB.GetDatabase())
	//
	//forumUsecase := forumUsecase.NewForumUsecase(forumRepo)
	//postUsecase := postUsecase.NewForumUsecase(postRepo, threadRepo)
	//serviceUsecase := serviceUsecase.NewForumUsecase(serviceRepo)
	//threadUsecase := threadUsecase.NewForumUsecase(threadRepo, forumRepo)
	//userUsecase := userUsecase.NewForumUsecase(userRepo)
	//
	//forumHandler := forumHandler.NewForumHandler(forumUsecase, threadUsecase)
	//postHandler := postHandler.NewPostHandler(postUsecase)
	//serviceHandler := serviceHandler.NewServiceHandler(serviceUsecase)
	//threadHandler := threadHandler.NewThreadHandler(threadUsecase, postUsecase)
	//userHandler := userHandler.NewUserHandler(userUsecase)

	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.appJSONMiddleware)

	//forumHandler.Configure(api)
	//postHandler.Configure(api)
	//serviceHandler.Configure(api)
	//threadHandler.Configure(api)
	//userHandler.Configure(api)

	server := http.Server{
		Addr:         fmt.Sprint(":", configs.Configs.GetMainPort()),
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
