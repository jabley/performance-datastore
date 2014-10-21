package main

import (
	"github.com/alext/tablecloth"
	"github.com/jabley/performance-datastore/pkg/config_api"
	"github.com/jabley/performance-datastore/pkg/handlers"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func main() {
	if wd := os.Getenv("GOVUK_APP_ROOT"); wd != "" {
		tablecloth.WorkingDir = wd
	}

	var (
		port         = getEnvDefault("HTTP_PORT", "8080")
		databaseName = getEnvDefault("DBNAME", "backdrop")
		mongoURL     = getEnvDefault("MONGO_URL", "localhost")
		bearerToken  = getEnvDefault("BEARER_TOKEN", "EMPTY")
		configAPIURL = getEnvDefault("CONFIG_API_URL", "https://stagecraft.production.performance.service.gov.uk/")
		maxGzipBody  = getEnvDefault("MAX_GZIP_SIZE", "10000000")
	)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	handlers.ConfigAPIClient = config_api.NewClient(configAPIURL, bearerToken)
	handlers.DataSetStorage = handlers.NewMongoStorage(mongoURL, databaseName)

	maxBody, err := strconv.Atoi(maxGzipBody)

	if err != nil {
		panic(err)
	}

	go serve(":"+port, handlers.NewHandler(maxBody), wg)
	wg.Wait()
}

func serve(addr string, handler http.Handler, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Fatal(tablecloth.ListenAndServe(addr, handler))
}

func getEnvDefault(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	return val
}
