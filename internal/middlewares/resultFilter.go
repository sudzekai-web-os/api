package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/sudzekai/web-os-api/internal/objects/responses"
	"github.com/sudzekai/web-os-api/packages/abstractions"
	"github.com/sudzekai/web-os-api/packages/types"
)

type ResultFilter struct {
	logger abstractions.ILogger
}

func NewResultFilter(loggerFactory abstractions.ILoggerFactory) *ResultFilter {
	return &ResultFilter{
		logger: loggerFactory.NewLogger("result-filter"),
	}
}

func (rf *ResultFilter) FilterFunc(w http.ResponseWriter, r *http.Request, result types.HandlerResult) {
	response := responses.ResponseEnvelope{}

	response.IsSuccess = result.Error == nil

	response.Data = result.Data

	if result.Error != nil {
		response.Error = &responses.Error{
			StatusCode: result.StatusCode,
			Message:    result.Error.Error(),
		}
	} else {
		response.Error = nil
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(result.StatusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		rf.logger.LogError("Ошибка сериализации ответа на %s %s: %s",
			r.Method,
			r.URL.Path,
			err,
		)
	}
}
