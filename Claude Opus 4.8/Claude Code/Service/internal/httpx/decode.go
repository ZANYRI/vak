package httpx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// DecodeAndValidate reads a JSON body into dst and runs struct validation.
// It returns an *APIError on malformed input so handlers can return it directly.
func DecodeAndValidate(r *http.Request, dst any) error {
	dec := json.NewDecoder(io.LimitReader(r.Body, 1<<20)) // 1 MiB cap
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		if err == io.EOF {
			return ErrValidation("request body is required")
		}
		return ErrValidation(fmt.Sprintf("invalid JSON body: %s", err.Error()))
	}

	if err := validate.Struct(dst); err != nil {
		var verrs validator.ValidationErrors
		if ok := asValidationErrors(err, &verrs); ok && len(verrs) > 0 {
			first := verrs[0]
			return ErrValidation("request validation failed").WithDetails(map[string]any{
				"field":  first.Field(),
				"reason": fmt.Sprintf("failed rule '%s'", first.Tag()),
			})
		}
		return ErrValidation(err.Error())
	}
	return nil
}

func asValidationErrors(err error, target *validator.ValidationErrors) bool {
	if ve, ok := err.(validator.ValidationErrors); ok {
		*target = ve
		return true
	}
	return false
}
