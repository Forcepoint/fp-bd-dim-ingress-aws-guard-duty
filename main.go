package main

import (
	"fmt"
	"fp-dim-aws-guard-duty-ingress/api"
	"fp-dim-aws-guard-duty-ingress/internal"
	"fp-dim-aws-guard-duty-ingress/internal/config"
	"fp-dim-aws-guard-duty-ingress/internal/hooks"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

func main() {
	InitViper("config", "./config")

	InitLogrus()
	internal.Register()
	config.CreateAndSetUrlAuthToken()

	router := mux.NewRouter()

	router.HandleFunc("/run", api.HandleIncomingData).Methods("POST", "OPTIONS")
	router.HandleFunc("/health", api.SendHealth).Methods("GET", "OPTIONS")
	router.HandleFunc("/config", api.ConfigEndpoint).Methods("GET", "OPTIONS")
	router.HandleFunc("/icon", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, "./icon/aws.png")
	})

	modulePort := os.Getenv("LOCAL_PORT")

	http.ListenAndServe(fmt.Sprintf(":%s", modulePort), router)
}

func InitLogrus() {
	logrus.SetLevel(logrus.InfoLevel)
	// Show where error was logged, function, line number, etc.
	logrus.SetReportCaller(true)

	// Output to stdout and logfile
	logrus.SetOutput(os.Stdout)

	logrus.AddHook(&hooks.LoggingHook{})
}

func InitViper(filename, configLocation string) {
	if _, err := os.Stat(fmt.Sprintf("%s/%s.yml", configLocation, filename)); err != nil {
		createConfigFile(configLocation, filename)
	}

	// Set up viper config library
	viper.SetConfigName(filename)
	viper.AddConfigPath(configLocation)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func createConfigFile(configLocation, filename string) {
	f, err := os.Create(fmt.Sprintf("%s/%s.yml", configLocation, filename))

	if err != nil {
		panic(err)
	}

	defer f.Close()
}
