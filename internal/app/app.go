package app

import (
	"github.com/gorilla/mux"
	sql "github.com/jmoiron/sqlx"
	deliveryForum "github.com/just4n4cc/tp-sem2-db/internal/service/forum/delivery"
	repositoryForum "github.com/just4n4cc/tp-sem2-db/internal/service/forum/repository"
	usecaseForum "github.com/just4n4cc/tp-sem2-db/internal/service/forum/usecase"
	deliveryPost "github.com/just4n4cc/tp-sem2-db/internal/service/post/delivery"
	repositoryPost "github.com/just4n4cc/tp-sem2-db/internal/service/post/repository"
	usecasePost "github.com/just4n4cc/tp-sem2-db/internal/service/post/usecase"
	deliveryService "github.com/just4n4cc/tp-sem2-db/internal/service/service/delivery"
	repositoryService "github.com/just4n4cc/tp-sem2-db/internal/service/service/repository"
	usecaseService "github.com/just4n4cc/tp-sem2-db/internal/service/service/usecase"
	deliveryThread "github.com/just4n4cc/tp-sem2-db/internal/service/thread/delivery"
	repositoryThread "github.com/just4n4cc/tp-sem2-db/internal/service/thread/repository"
	usecaseThread "github.com/just4n4cc/tp-sem2-db/internal/service/thread/usecase"
	deliveryUser "github.com/just4n4cc/tp-sem2-db/internal/service/user/delivery"
	repositoryUser "github.com/just4n4cc/tp-sem2-db/internal/service/user/repository"
	usecaseUser "github.com/just4n4cc/tp-sem2-db/internal/service/user/usecase"
	"github.com/just4n4cc/tp-sem2-db/internal/utils"
	"github.com/just4n4cc/tp-sem2-db/pkg/logger"
	"net/http"
	"strings"
)

const logMessage = "app:"

type App struct {
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
		//logger.Error(message, err)
		return nil, err
	}

	repoUser := repositoryUser.NewRepository(db)
	repoPost := repositoryPost.NewRepository(db)
	repoForum := repositoryForum.NewRepository(db)
	repoThread := repositoryThread.NewRepository(db)
	repoService := repositoryService.NewRepository(db)

	useUser := usecaseUser.NewUseCase(repoUser)
	usePost := usecasePost.NewUseCase(repoPost, repoForum, repoThread, repoUser)
	useThread := usecaseThread.NewUseCase(repoThread, repoPost)
	useForum := usecaseForum.NewUseCase(repoForum, repoThread)
	useService := usecaseService.NewUseCase(repoService)

	return &App{
		deliveryForum:   deliveryForum.NewDelivery(useForum),
		deliveryPost:    deliveryPost.NewDelivery(usePost),
		deliveryService: deliveryService.NewDelivery(useService),
		deliveryThread:  deliveryThread.NewDelivery(useThread),
		deliveryUser:    deliveryUser.NewDelivery(useUser),
	}, nil
}

func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func urlPrintMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()
		if strings.Contains(url, "service") {
			logger.Info("URL: " + url)
		}
		next.ServeHTTP(w, r)
	})
}

func newRouter(app *App) *mux.Router {
	r := mux.NewRouter()
	r.Use(contentTypeMiddleware)
	//r.Use(urlPrintMiddleware)
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

	//message := logMessage + "Run:"
	//logger.Info(message + "start")

	port := ":5000"
	r := newRouter(app)
	s := &http.Server{
		Addr:    port,
		Handler: r,
		//ReadTimeout:  10 * time.Second,
		//WriteTimeout: 10 * time.Second,
	}
	err := s.ListenAndServe()
	if err != nil {
		//logger.Error(message, err)
		return err
	}
	return nil
}
