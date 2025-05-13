package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type StandardResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// WriteJSON sends a standardised success response as JSON.

func WriteJSON(w http.ResponseWriter, statusCode int, message string, data interface{}) error {
	resp := StandardResponse{
		Status:  http.StatusText(statusCode),
		Message: message,
		Data:    data,
	}

	js, err := json.Marshal(resp)
	if err != nil {
		// Log and fallback to plain-text error response
		fmt.Println("❌ WriteJSON marshal error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(js); err != nil {
		fmt.Println("❌ WriteJSON write error:", err)
		return err
	}

	fmt.Println("✅ WriteJSON sent:", string(js))
	return nil
}

// DecodeJSON parses the request body into the dst struct.
// Disallows unknown fields and multiple JSON values.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxErr):
			return fmt.Errorf("malformed JSON at position %d", syntaxErr.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("malformed JSON")
		case errors.As(err, &unmarshalTypeErr):
			return fmt.Errorf("wrong type for field %q", unmarshalTypeErr.Field)
		case errors.Is(err, io.EOF):
			return errors.New("request body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("unknown field %s", field)
		case err.Error() == "http: request body too large":
			return errors.New("request body must not be larger than 1MB")
		default:
			return err
		}
	}

	// Check for more than one JSON object
	if dec.More() {
		return errors.New("only one JSON object allowed")
	}

	return nil
}
