package web_api

import (
	"github.com/andreychh/coopera-backend/internal/adapter/controller/web_api/middleware"
	chimw "github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"

	"github.com/andreychh/coopera-backend/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	userController       *UserController
	teamController       *TeamController
	taskController       *TaskController
	membershipController *MembershipController
}

func NewRouter(
	userUseCase usecase.UserUseCase,
	teamUseCase usecase.TeamUseCase,
	taskUseCase usecase.TaskUseCase,
	membershipUseCase usecase.MembershipUseCase,
) *Router {
	return &Router{
		userController:       NewUserController(userUseCase),
		teamController:       NewTeamController(teamUseCase),
		taskController:       NewTaskController(taskUseCase),
		membershipController: NewMembershipController(membershipUseCase),
	}
}

func (r *Router) SetupRoutes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.TimeoutMiddleware(5*time.Second), chimw.Recoverer)

	router.Route("/api/v1", func(api chi.Router) {
		api.Route("/users", func(users chi.Router) {
			users.Post("/", middleware.ErrorHandler(r.userController.Create))
			users.Get("/", middleware.ErrorHandler(r.userController.Get))
		})

		api.Route("/teams", func(teams chi.Router) {
			teams.Post("/", middleware.ErrorHandler(r.teamController.Create))
			teams.Get("/", middleware.ErrorHandler(r.teamController.Get))
			teams.Delete("/", middleware.ErrorHandler(r.teamController.Delete))
		})

		api.Route("/memberships", func(members chi.Router) {
			members.Post("/", middleware.ErrorHandler(r.membershipController.AddMember))
			members.Delete("/", middleware.ErrorHandler(r.membershipController.DeleteMember))
		})

		api.Route("/tasks", func(tasks chi.Router) {
			tasks.Post("/", middleware.ErrorHandler(r.taskController.Create))
			//tasks.Get("/", middleware.ErrorHandler(r.taskController.Get))
			//tasks.Delete("/", middleware.ErrorHandler(r.taskController.Delete))
		})
	})

	return router
}
