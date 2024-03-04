package main

import (
	"embed"
	"flag"
	"fmt"
	"hmcalister/htmxTest/api"
	"io/fs"
	"net/http"
	"os"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	//go:embed static/css/output.css
	embedCSSFile []byte

	//go:embed static/htmx/htmx.js
	embedHTMXFile []byte

	//go:embed static/templates/*.html
	templatesFS embed.FS
)

func loggerSetup(logFilePath string, debugFlag bool) *os.File {
	var err error
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.
		With().Caller().Logger().
		With().Timestamp().Logger()

	logFileHandle, err := os.Create(logFilePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}

	log.Logger = log.Output(logFileHandle)
	if debugFlag {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
		multiWriter := zerolog.MultiLevelWriter(consoleWriter, logFileHandle)
		log.Logger = log.Output(multiWriter)
	}

	return logFileHandle
}

func main() {
	var err error

	port := flag.Int("port", 8080, "The port to run the application on.")
	debugFlag := flag.Bool("debug", false, "Flag for debug level with console log outputs.")
	logFilePath := flag.String("logFilePath", "log", "File to write logs to. If nil, logs written to os.Stdout.")
	flag.Parse()

	logFile := loggerSetup(*logFilePath, *debugFlag)
	defer logFile.Close()

	applicationState := api.NewApplicationState()

	// Parse templates from embedded file system --------------------------------------------------

	templatesFS, err := fs.Sub(templatesFS, "static/templates")
	if err != nil {
		log.Fatal().Err(err).Msg("error during embedded file system")
	}
	indexTemplate, err := template.ParseFS(templatesFS, "index.html")
	if err != nil {
		log.Fatal().Err(err).Msg("error parsing template")
	}
	cardTemplate, err := template.ParseFS(templatesFS, "card.html")
	if err != nil {
		log.Fatal().Err(err).Msg("error parsing template")
	}

	router := chi.NewRouter()
	router.Use(zerologLogger)
	router.Use(recoverWithInternalServerError)

	// Add handlers for CSS and HTMX files --------------------------------------------------------

	router.Get("/css/output.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Write(embedCSSFile)
	})

	router.Get("/htmx/htmx.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")
		w.Write(embedHTMXFile)
	})

	// Add handlers for base routes, e.g. initial page --------------------------------------------
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		err = indexTemplate.Execute(w, nil)
		if err != nil {
			log.Fatal().Err(err).Msg("error during index template execute")
		}
	})

	// Add any API routes -------------------------------------------------------------------------
	router.Put("/api/addItem", func(w http.ResponseWriter, r *http.Request) {
		cardData := struct {
			ItemID int
		}{
			ItemID: applicationState.AddItem(),
		}
		err = cardTemplate.Execute(w, cardData)
		if err != nil {
			log.Fatal().Err(err).Msg("error during card template execute")
		}
	})

	router.Delete("/api/deleteItem", func(w http.ResponseWriter, r *http.Request) {
		applicationState.DeleteItem()
		w.Write(nil)
	})

	router.Delete("/api/deleteAll", func(w http.ResponseWriter, r *http.Request) {
		applicationState.DeleteAll()
		w.Write(nil)
	})

	// Start server -------------------------------------------------------------------------------

	fmt.Printf("Serving at http://localhost:%v/", *port)
	err = http.ListenAndServe(fmt.Sprintf(":%v", *port), router)
	if err != nil {
		log.Fatal().Err(err).Msg("error during http serving")
	}
}
