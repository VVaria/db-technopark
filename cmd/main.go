package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/VVaria/db-technopark/configs"
	"github.com/VVaria/db-technopark/internal/app/tools/databases"
	"github.com/VVaria/db-technopark/internal/app/tools/middleware"

	forumHandler "github.com/VVaria/db-technopark/internal/app/forum/delivery/http"
	forumRepo "github.com/VVaria/db-technopark/internal/app/forum/repository/postgres"
	forumUsecase "github.com/VVaria/db-technopark/internal/app/forum/usecase"

	postHandler "github.com/VVaria/db-technopark/internal/app/post/delivery/http"
	postRepo "github.com/VVaria/db-technopark/internal/app/post/repository/postgres"
	postUsecase "github.com/VVaria/db-technopark/internal/app/post/usecase"

	serviceHandler "github.com/VVaria/db-technopark/internal/app/service/delivery/http"
	serviceRepo "github.com/VVaria/db-technopark/internal/app/service/repository/postgres"
	serviceUsecase "github.com/VVaria/db-technopark/internal/app/service/usecase"

	threadHandler "github.com/VVaria/db-technopark/internal/app/thread/delivery/http"
	threadRepo "github.com/VVaria/db-technopark/internal/app/thread/repository/postgres"
	threadUsecase "github.com/VVaria/db-technopark/internal/app/thread/usecase"

	userHandler "github.com/VVaria/db-technopark/internal/app/user/delivery/http"
	userRepo "github.com/VVaria/db-technopark/internal/app/user/repository/postgres"
	userUsecase "github.com/VVaria/db-technopark/internal/app/user/usecase"
)

func main() {
	postgresDB, err := databases.NewPostgres(configs.Configs.GetPostgresConfig())
	if err != nil {
		log.Fatal(err)
	}

	forumRepo := forumRepo.NewForumRepository(postgresDB.GetDatabase())
	postRepo := postRepo.NewPostRepository(postgresDB.GetDatabase())
	serviceRepo := serviceRepo.NewServiceRepository(postgresDB.GetDatabase())
	threadRepo := threadRepo.NewThreadRepository(postgresDB.GetDatabase())
	userRepo := userRepo.NewUserRepository(postgresDB.GetDatabase())

	forumUsecase := forumUsecase.NewForumUsecase(forumRepo)
	postUsecase := postUsecase.NewPostUsecase(postRepo, threadRepo)
	serviceUsecase := serviceUsecase.NewServiceUsecase(serviceRepo)
	threadUsecase := threadUsecase.NewThreadUsecase(threadRepo, forumRepo)
	userUsecase := userUsecase.NewUserUsecase(userRepo)

	forumHandler := forumHandler.NewForumHandler(forumUsecase, threadUsecase)
	postHandler := postHandler.NewPostHandler(postUsecase)
	serviceHandler := serviceHandler.NewServiceHandler(serviceUsecase)
	threadHandler := threadHandler.NewThreadHandler(threadUsecase, postUsecase)
	userHandler := userHandler.NewUserHandler(userUsecase)

	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AppJSONMiddleware)

	forumHandler.Configure(api)
	postHandler.Configure(api)
	serviceHandler.Configure(api)
	threadHandler.Configure(api)
	userHandler.Configure(api)

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
