package web_api

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func InitValidator(v *validator.Validate) {
	validate = v
}

// BindRequest парсит запрос в dst (body JSON или query) и валидирует
func BindRequest(r *http.Request, dst any) error {
	// Body для POST, PUT, PATCH
	if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
		if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
			return err
		}
	} else { // Query для GET, DELETE и т.д.
		if err := decodeQuery(r, dst); err != nil {
			return err
		}
	}

	// Валидация
	if err := validate.Struct(dst); err != nil {
		return err
	}
	return nil
}

// decodeQuery с конвертацией типов
func decodeQuery(r *http.Request, dst any) error {
	values := r.URL.Query()
	v := reflect.ValueOf(dst).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("form")
		if tag == "" {
			continue
		}

		param := values.Get(tag)
		if param == "" {
			continue
		}

		f := v.Field(i)
		if !f.CanSet() {
			continue
		}

		switch f.Kind() {
		case reflect.String:
			f.SetString(param)
		case reflect.Int, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(param, 10, 64)
			if err != nil {
				return err
			}
			f.SetInt(n)
		case reflect.Uint, reflect.Uint32, reflect.Uint64:
			n, err := strconv.ParseUint(param, 10, 64)
			if err != nil {
				return err
			}
			f.SetUint(n)
		case reflect.Bool:
			b, err := strconv.ParseBool(param)
			if err != nil {
				return err
			}
			f.SetBool(b)
		}
	}

	return nil
}

// helper для возврата ошибок в JSON
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeValidationError(w http.ResponseWriter, err error) {
	if ve, ok := err.(validator.ValidationErrors); ok {
		errors := make([]string, 0, len(ve))
		for _, e := range ve {
			errors = append(errors, e.Error())
		}
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error":   "Validation failed",
			"details": errors,
		})
	} else {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
}
