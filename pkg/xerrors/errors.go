package xerrors

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

type ErrorDetail struct {
	Field       string `json:"field,omitempty"`
	Description string `json:"description"`
}

type CustomErrorResponse struct {
	Code             int           `json:"code"`
	Message          string        `json:"message"`
	Reason           string        `json:"reason"`
	LocalizedMessage string        `json:"localized_message"`
	Details          []ErrorDetail `json:"details,omitempty"`
}

func CustomErrorHandler(_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	st := status.Convert(err)
	httpCode := runtime.HTTPStatusFromCode(st.Code())

	response := CustomErrorResponse{
		Code:             httpCode,
		Message:          st.Message(),
		Reason:           st.Code().String(),
		LocalizedMessage: "error during processing the request",
	}

	for _, detail := range st.Details() {
		if violations, ok := detail.(*validate.Violations); ok {
			response.LocalizedMessage = "validate, body request"

			for _, violation := range violations.GetViolations() {
				fieldName := ""
				protoField := violation.GetField()
				if protoField != nil {
					elements := protoField.GetElements()
					var parts []string
					for _, el := range elements {
						if el.GetFieldName() != "" {
							parts = append(parts, el.GetFieldName())
						}
					}
					fieldName = strings.Join(parts, ".")
				}

				response.Details = append(response.Details, ErrorDetail{
					Field:       fieldName,
					Description: violation.GetMessage(),
				})
			}
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	_ = json.NewEncoder(w).Encode(response)
}
