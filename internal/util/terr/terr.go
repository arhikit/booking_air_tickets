package terr

import (
	"encoding/json"
	"net/http"
)

// Error defines model for Error.
type Error struct {
	// Код состояния HTTP
	HTTPStatusCode int `json:"-"`

	// Код ошибки
	Code string `json:"code"`

	// Сообщение об ошибке
	Message string `json:"message"`
}

// From returns Error structure from all kind of error types.
func From(err error) *Error {
	e, ok := err.(*Error)
	if !ok {
		e = InternalServerError("UNKNOWN_ERROR", err.Error())
	}
	return e
}

// Equal returns true if errors have same code.
func Equal(e1, e2 error) bool {
	if e1 == nil && e2 == nil {
		return true
	}
	if (e1 == nil && e2 != nil) || (e1 != nil && e2 == nil) {
		return false
	}

	err1, ok := e1.(*Error)
	if !ok {
		return e1.Error() == e2.Error()
	}

	err2, ok := e2.(*Error)
	if !ok {
		return e1.Error() == e2.Error()
	}

	return err1.Code == err2.Code
}

// Error реализация интерфейса ошибки.
func (e *Error) Error() string {
	return e.Message
}

func WriteError(w http.ResponseWriter, err *Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.HTTPStatusCode)
	_ = json.NewEncoder(w).Encode(&err)
}
