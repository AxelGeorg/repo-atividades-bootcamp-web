package errors

import (
	"fmt"
	"net/http"
)

type CustomError struct {
	StatusHttp int    // Tipo do erro
	Message    string // Mensagem do erro
}

// Implementação da interface error
func (e *CustomError) Error() string {
	return fmt.Sprintf("%s: %s", http.StatusText(e.StatusHttp), e.Message)
}

// Funções para criar novos erros personalizados
func NewBadRequestError(message string) error {
	return &CustomError{StatusHttp: http.StatusBadRequest, Message: message}
}

func NewConflictError(message string) error {
	return &CustomError{StatusHttp: http.StatusConflict, Message: message}
}
