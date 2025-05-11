package logger

import (
	"log/slog"
	"net/http"

	"github.com/iankencruz/sabiflow/internal/response"
)

// WriteJSONError logs the error and sends a structured JSON response.
func WriteJSONError(w http.ResponseWriter, log *slog.Logger, status int, msg string, err error, context ...any) {
	fields := append(context, "error", err)
	log.Error(msg, fields...)

	if writeErr := response.WriteJSON(w, status, msg, nil); writeErr != nil {
		log.Error("failed to write JSON error response", "write_error", writeErr)
	}
}
