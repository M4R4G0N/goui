package components

import (
	"net/http"
)

// DateValue representa a parsed date value do frontend.
type DateValue struct {
	Single string   // Preenchido se for data única
	Start  string   // Preenchido se for range
	End    string   // Preenchido se for range
	IsRange bool     // Indica se é um intervalo
}

// ParseDate extrai de forma genérica uma data ou range de um request HTTP.
func ParseDate(r *http.Request, name string) DateValue {
	// r.FormValue já chama ParseForm ou ParseMultipartForm internamente.
	// No entanto, para pegar múltiplos valores (slice), precisamos do r.Form.
	r.ParseForm()
	if r.MultipartForm != nil {
		// Se for multipart, os valores ficam no r.MultipartForm.Value
		values := r.MultipartForm.Value[name]
		if len(values) >= 2 {
			return DateValue{Start: values[0], End: values[1], IsRange: true}
		}
		if len(values) == 1 {
			return DateValue{Single: values[0], IsRange: false}
		}
	}
	
	// Caso padrão (urlencoded)
	values := r.Form[name]
	if len(values) >= 2 {
		return DateValue{Start: values[0], End: values[1], IsRange: true}
	}
	if len(values) == 1 {
		return DateValue{Single: values[0], IsRange: false}
	}
	
	return DateValue{}
}
