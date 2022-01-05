package pkg

import (
	"net/http"

	"github.com/justinas/alice"
)

func initRoutes(app *Application) http.Handler {
	mux := http.NewServeMux()
	//staticFileServer := http.FileServer(http.Dir("./ui/static"))
	//mux.Handle("/static/", http.StripPrefix("/static", staticFileServer))

	mux.Handle(HomeRoute, home(app))
	mux.Handle(DashboardRoute, dashboard(app))
	mux.Handle(SignUpRoute, signUp(app))
	mux.Handle(AllUsersRoute, allUsers(app))
	mux.Handle(AddBill, addBill(app))

	middlewareBeforeRouting := alice.New(recoverPanic, logRequest, secureHeaders, makeSwaggerCompatible)
	return middlewareBeforeRouting.Then(mux)
}
