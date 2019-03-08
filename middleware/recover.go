package middleware

import (
	"net/http"

	"github.com/bingtianbaihua/goproxy/log"
)

type RecoverAdapter struct{}

func NewRecoverAdapter() (*RecoverAdapter, error) {
	return &RecoverAdapter{}, nil
}

// NewRecovery creates a new instance of Recovery
func (rv *RecoverAdapter) HandleTask(w http.ResponseWriter, r *http.Request, next func(http.ResponseWriter, *http.Request)) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Warn("panic with error: %v\n", err)
		}
	}()
	next(w, r)
}
