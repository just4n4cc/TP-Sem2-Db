package app

import (
	"github.com/gorilla/mux"
	sql "github.com/jmoiron/sqlx"
	deliveryForum "github.com/just4n4cc/tp-sem2-db/internal/service/forum/delivery"
	deliveryPost "github.com/just4n4cc/tp-sem2-db/internal/service/post/delivery"
	deliveryService "github.com/just4n4cc/tp-sem2-db/internal/service/service/delivery"
	deliveryThread "github.com/just4n4cc/tp-sem2-db/internal/service/thread/delivery"
	deliveryUser "github.com/just4n4cc/tp-sem2-db/internal/service/user/delivery"
	"net/http"
	"time"

	//repositoryForum "github.com/just4n4cc/tp-sem2-db/internal/service/forum/repository"
	//repositoryPost "github.com/just4n4cc/tp-sem2-db/internal/service/post/repository"
	//repositoryService "github.com/just4n4cc/tp-sem2-db/internal/service/service/repository"
	//repositoryThread "github.com/just4n4cc/tp-sem2-db/internal/service/thread/repository"
	repositoryUser "github.com/just4n4cc/tp-sem2-db/internal/service/user/repository"

	//usecaseForum "github.com/just4n4cc/tp-sem2-db/internal/service/forum/repository"
	//usecasePost "github.com/just4n4cc/tp-sem2-db/internal/service/post/repository"
	//usecaseService "github.com/just4n4cc/tp-sem2-db/internal/service/service/repository"
	//usecaseThread "github.com/just4n4cc/tp-sem2-db/internal/service/thread/repository"
	usecaseUser "github.com/just4n4cc/tp-sem2-db/internal/service/user/usecase"
	"github.com/just4n4cc/tp-sem2-db/internal/utils"
	"github.com/just4n4cc/tp-sem2-db/pkg/logger"
)

const logMessage = "app:"

type App struct {
	// options
	deliveryForum   *deliveryForum.Delivery
	deliveryPost    *deliveryPost.Delivery
	deliveryService *deliveryService.Delivery
	deliveryThread  *deliveryThread.Delivery
	deliveryUser    *deliveryUser.Delivery
	db              *sql.DB
}

func NewApp() (*App, error) {
	const message = logMessage + "NewApp"
	db, err := utils.InitDb()
	if err != nil {
		logger.Error(message, err)
		return nil, err
	}

	repoUser := repositoryUser.NewRepository(db)
	useUser := usecaseUser.NewUseCase(repoUser)
	//deliveryUser := deliveryUser.NewDelivery(usecaseUser)

	return &App{
		//deliveryForum: deliveryUser.NewDelivery(useForum),
		//deliveryPost: deliveryUser.NewDelivery(usePost),
		//deliveryService: deliveryUser.NewDelivery(useService),
		//deliveryThread: deliveryUser.NewDelivery(useThread),
		deliveryUser: deliveryUser.NewDelivery(useUser),
	}, nil
}

func newRouter(app *App) *mux.Router {
	r := mux.NewRouter()
	rApi := r.PathPrefix("/api").Subrouter()

	rApi.HandleFunc("/forum/create", app.deliveryForum.ForumCreate).Methods("POST")
	rApi.HandleFunc("/forum/{slug}/details", app.deliveryForum.ForumGet).Methods("GET")
	rApi.HandleFunc("/forum/{slug}/create", app.deliveryForum.ForumThreadCreate).Methods("POST")
	rApi.HandleFunc("/forum/{slug}/users", app.deliveryForum.ForumUsers).Methods("GET")
	rApi.HandleFunc("/forum/{slug}/threads", app.deliveryForum.ForumThreads).Methods("GET")

	rApi.HandleFunc("/post/{id}/details", app.deliveryPost.PostGet).Methods("GET")
	rApi.HandleFunc("/post/{id}/details", app.deliveryPost.PostUpdate).Methods("POST")

	rApi.HandleFunc("/service/clear", app.deliveryService.ServiceClear).Methods("POST")
	rApi.HandleFunc("/service/status", app.deliveryService.ServiceStatus).Methods("GET")

	rApi.HandleFunc("/thread/{slug_or_id}/create", app.deliveryThread.ThreadCreate).Methods("POST")
	rApi.HandleFunc("/thread/{slug_or_id}/details", app.deliveryThread.ThreadGet).Methods("GET")
	rApi.HandleFunc("/thread/{slug_or_id}/details", app.deliveryThread.ThreadUpdate).Methods("POST")
	rApi.HandleFunc("/thread/{slug_or_id}/posts", app.deliveryThread.ThreadPosts).Methods("GET")
	rApi.HandleFunc("/thread/{slug_or_id}/vote", app.deliveryThread.ThreadVote).Methods("POST")

	rApi.HandleFunc("/user/{nickname}/create", app.deliveryUser.UserCreate).Methods("POST")
	rApi.HandleFunc("/user/{nickname}/profile", app.deliveryUser.UserProfileGet).Methods("GET")
	rApi.HandleFunc("/user/{nickname}/profile", app.deliveryUser.UserProfileUpdate).Methods("POST")

	return r
}

func (app *App) Run() error {
	if app.db != nil {
		defer app.db.Close()
	}

	message := logMessage + "Run:"
	logger.Info(message + "start")

	port := ":5050"
	r := newRouter(app)
	s := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := s.ListenAndServe()
	if err != nil {
		logger.Error(message, err)
		return err
	}
	return nil
}
