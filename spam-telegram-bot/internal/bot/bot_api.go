package bot

import (
	"net/http"
)

func SetRoutes(r *http.ServeMux, bot *SpamBot) {
	r.HandleFunc("/api/send", bot.SendHandler)
	r.HandleFunc("/healthz", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("OK"))
	})
}
