package internal

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type ErrToJson struct {
	logger *zap.SugaredLogger
}

func NewErrToJson(l *zap.SugaredLogger) *ErrToJson {
	return &ErrToJson{l}
}

func (e ErrToJson) ErrToJson(w http.ResponseWriter, externalError error) {
	errrr := ErrOut{externalError.Error()}
	marshaled, err := json.Marshal(errrr)
	if err != nil {
		e.logger.Error(err)
		return
	}
	_, err = w.Write(marshaled)
	if err != nil {
		e.logger.Error(err)
		return
	}
}
