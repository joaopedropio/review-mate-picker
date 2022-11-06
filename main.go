package main

import (
	"fmt"
	"github.com/joaopedropio/review-mate-picker/handlers"
	"github.com/joaopedropio/review-mate-picker/infrastructe"
	"github.com/joho/godotenv"
	"net/http"
)

func init() {
	godotenv.Load(".env")
}

func main() {
	env, err := infrastructe.NewEnvironment()
	if err != nil {
		panic(fmt.Errorf("unable to initialize environment: %w", err))
	}

	configureEndpoints(env)

	fmt.Println(fmt.Sprintf("[INFO] Server listening on port %d", env.GetHttpPort()))
	err = http.ListenAndServe(fmt.Sprintf(":%d", env.GetHttpPort()), nil)
	fmt.Println(err.Error())
}

func configureEndpoints(env infrastructe.Environment) {
	container := infrastructe.NewDependencyInjectionContainer(env)

	healthCheckHandler := handlers.NewHealthCheckHandler()
	eventsHandler := handlers.NewEventsHandler(container.SlackClient, env.GetSlackSigningSecret(), container.MateService, container.TimestampCache)

	http.HandleFunc("/healthCheck", healthCheckHandler.Handle)
	http.HandleFunc("/events", eventsHandler.Handle)
}
