package router

import (
	"github.com/Zhang-Yu-Bo/curly-garbanzo/controller"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	mRouter := mux.NewRouter()

	mRouter.HandleFunc("/", controller.HomePage).Methods("GET")
	mRouter.HandleFunc("/{username}", controller.ShowUserInfo).Methods("GET")
	mRouter.HandleFunc("/eventsub/{userID}", controller.EventSub).Methods("POST")

	mRouter.HandleFunc("/test", controller.TestPage).Methods("GET")

	return mRouter
}
