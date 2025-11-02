package web_api

import (
	"github.com/andreychh/coopera/internal/adapter/controller/web_api/middleware"
	chimw "github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"

	"github.com/andreychh/coopera/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	userController       *UserController
	teamController       *TeamController
	membershipController *MembershipController
}

func NewRouter(
	userUseCase usecase.UserUseCase,
	teamUseCase usecase.TeamUseCase,
	membershipUseCase usecase.MembershipUseCase,
) *Router {
	return &Router{
		userController:       NewUserController(userUseCase),
		teamController:       NewTeamController(teamUseCase),
		membershipController: NewMembershipController(membershipUseCase),
	}
}

func (r *Router) SetupRoutes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.TimeoutMiddleware(5*time.Second), chimw.Recoverer)

	router.Route("/api/v1", func(api chi.Router) {
		api.Route("/users", func(users chi.Router) {
			users.Post("/", r.userController.Create)
			users.Get("/", r.userController.Get)
		})

		api.Route("/teams", func(teams chi.Router) {
			teams.Post("/", r.teamController.Create)
			teams.Get("/", r.teamController.Get)
		})

		api.Route("/memberships", func(members chi.Router) {
			members.Post("/", r.membershipController.AddMember)
		})
	})

	return router
}
