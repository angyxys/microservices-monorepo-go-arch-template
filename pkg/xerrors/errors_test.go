package xerrors

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestCustomErrorHandler_StandardgRPCError(t *testing.T) {
	grpcErr := status.Error(codes.Unauthenticated, "missing token")

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/v1/users", nil)
	mux := runtime.NewServeMux()

	marshaler := &runtime.JSONBuiltin{}

	CustomErrorHandler(context.Background(), mux, marshaler, w, r, grpcErr)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var body CustomErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, body.Code)
	assert.Equal(t, "missing token", body.Message)
	assert.Equal(t, "Unauthenticated", body.Reason)
	assert.Equal(t, "error during processing the request", body.LocalizedMessage)
	assert.Empty(t, body.Details)
}

func TestCustomErrorHandler_ValidationViolations(t *testing.T) {
	protoViolations := &validate.Violations{
		Violations: []*validate.Violation{
			{
				Message: proto.String("value length must be at least 6 runes"),
				Field: &validate.FieldPath{
					Elements: []*validate.FieldPathElement{
						{
							FieldName: proto.String("password"),
						},
					},
				},
			},
		},
	}

	baseStatus := status.New(codes.InvalidArgument, "field validations")
	statusWithDetails, err := baseStatus.WithDetails(protoViolations)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/v1/users", nil)
	mux := runtime.NewServeMux()
	marshaler := &runtime.JSONBuiltin{}

	CustomErrorHandler(context.Background(), mux, marshaler, w, r, statusWithDetails.Err())

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var body CustomErrorResponse
	err = json.Unmarshal(w.Body.Bytes(), &body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, body.Code)
	assert.Equal(t, "field validations", body.Message)
	assert.Equal(t, "InvalidArgument", body.Reason)
	assert.Equal(t, "validate, body request", body.LocalizedMessage)

	assert.Len(t, body.Details, 1)
	assert.Equal(t, "password", body.Details[0].Field)
	assert.Equal(t, "value length must be at least 6 runes", body.Details[0].Description)
}
