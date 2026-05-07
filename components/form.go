package components

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ─── CSRF ──────────────────────────────────────────────────────────────────

var csrfSecret = []byte("goui-csrf-default-secret-change-this")

// SetCSRFSecret sets the HMAC key used to sign CSRF tokens.
// Call this once at startup with a random, persistent secret.
func SetCSRFSecret(secret string) {
	csrfSecret = []byte(secret)
}

// NewCSRFToken generates a signed token, persists it in an HttpOnly cookie,
// and returns the raw token to be embedded in the form via CSRF(token).
func NewCSRFToken(w http.ResponseWriter) string {
	raw := make([]byte, 32)
	rand.Read(raw)
	token := base64.URLEncoding.EncodeToString(raw)

	mac := hmac.New(sha256.New, csrfSecret)
	mac.Write([]byte(token))
	sig := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	http.SetCookie(w, &http.Cookie{
		Name:     "_goui_csrf",
		Value:    sig + "|" + token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(1 * time.Hour),
	})
	return token
}

// ValidateCSRF returns true when the "_csrf" form field matches the signed cookie.
// Call this at the top of every POST handler that uses Form with CSRF protection.
func ValidateCSRF(r *http.Request) bool {
	cookie, err := r.Cookie("_goui_csrf")
	if err != nil {
		return false
	}
	parts := strings.SplitN(cookie.Value, "|", 2)
	if len(parts) != 2 {
		return false
	}
	storedSig, token := parts[0], parts[1]

	mac := hmac.New(sha256.New, csrfSecret)
	mac.Write([]byte(token))
	expectedSig := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(storedSig), []byte(expectedSig)) {
		return false
	}
	return hmac.Equal([]byte(r.FormValue("_csrf")), []byte(token))
}

// ─── FORM COMPONENT ────────────────────────────────────────────────────────

// csrf is an unexported sentinel type for CSRF tokens.
type csrf string

// CSRF sets the CSRF token that will be embedded as a hidden field in the form.
func CSRF(token string) csrf { return csrf(token) }

// FormComponent renders an HTML <form> element.
type FormComponent struct {
	Action    string
	Method    string
	Children  []Component
	Class     string
	ID        string
	CSRFToken string
	Attrs     Attr
}

func (f *FormComponent) GetID() string {
	if f.ID == "" {
		f.ID = AutoID()
	}
	return f.ID
}

// Form creates a <form> with optional CSRF protection.
//
// Example (with CSRF):
//
//	// GET handler:
//	token := components.NewCSRFToken(w)
//	page := MyPage(token)
//
//	// POST handler:
//	if !components.ValidateCSRF(r) { http.Error(w, "CSRF inválido", 403); return }
//
//	// Page function:
//	components.Form("/submit", "POST",
//	    components.CSRF(token),
//	    components.FormField("Nome", components.Input(components.Name("name")), ""),
//	    components.Button("Enviar"),
//	)
func Form(action, method string, opts ...any) *FormComponent {
	f := &FormComponent{Action: action, Method: method}
	for _, opt := range opts {
		switch v := opt.(type) {
		case csrf:
			f.CSRFToken = string(v)
		case Class:
			f.Class = string(v)
		case ID:
			f.ID = string(v)
		case Attr:
			f.Attrs = v
		case Component:
			f.Children = append(f.Children, v)
		case []Component:
			f.Children = append(f.Children, v...)
		}
	}
	return f
}

func (f *FormComponent) Render() string {
	var b strings.Builder
	cls := "goui-form"
	if f.Class != "" {
		cls += " " + f.Class
	}
	fmt.Fprintf(&b, `<form id="%s" class="%s" action="%s" method="%s"%s>`,
		f.GetID(), cls, f.Action, f.Method, renderAttrs(f.Attrs))

	if f.CSRFToken != "" {
		fmt.Fprintf(&b, `<input type="hidden" name="_csrf" value="%s">`, f.CSRFToken)
	}
	for _, child := range f.Children {
		b.WriteString(child.Render())
	}
	b.WriteString(`</form>`)
	return b.String()
}

// ─── CLIENT-SIDE VALIDATION ────────────────────────────────────────────────

// Validation holds declarative validation rules used by Input and Textarea.
// Add it as an option: Input(Name("email"), Validation{Required: true, ...})
type Validation struct {
	Required    bool
	RequiredMsg string
	MinLen      int
	MinLenMsg   string
	MaxLen      int
	MaxLenMsg   string
	MinNum      *float64
	MinNumMsg   string
	MaxNum      *float64
	MaxNumMsg   string
	Pattern     string
	PatternMsg  string
}

// renderValidationAttrs translates a Validation struct into HTML attribute strings.
func renderValidationAttrs(v *Validation) []string {
	var attrs []string
	if v.Required {
		attrs = append(attrs, "required")
	}
	if v.MinLen > 0 {
		attrs = append(attrs, fmt.Sprintf(`minlength="%d"`, v.MinLen))
	}
	if v.MaxLen > 0 {
		attrs = append(attrs, fmt.Sprintf(`maxlength="%d"`, v.MaxLen))
	}
	if v.MinNum != nil {
		attrs = append(attrs, fmt.Sprintf(`min="%s"`, strconv.FormatFloat(*v.MinNum, 'f', -1, 64)))
	}
	if v.MaxNum != nil {
		attrs = append(attrs, fmt.Sprintf(`max="%s"`, strconv.FormatFloat(*v.MaxNum, 'f', -1, 64)))
	}
	if v.Pattern != "" {
		attrs = append(attrs, fmt.Sprintf(`pattern="%s"`, v.Pattern))
	}
	return attrs
}

// renderValidationScript emits a <script> that sets custom browser validity messages.
func renderValidationScript(id string, v *Validation) string {
	if v == nil {
		return ""
	}
	var checks []string
	if v.Required && v.RequiredMsg != "" {
		checks = append(checks, fmt.Sprintf(`if(el.validity.valueMissing){el.setCustomValidity('%s');return;}`, v.RequiredMsg))
	}
	if v.MinLen > 0 && v.MinLenMsg != "" {
		checks = append(checks, fmt.Sprintf(`if(el.validity.tooShort){el.setCustomValidity('%s');return;}`, v.MinLenMsg))
	}
	if v.MaxLen > 0 && v.MaxLenMsg != "" {
		checks = append(checks, fmt.Sprintf(`if(el.validity.tooLong){el.setCustomValidity('%s');return;}`, v.MaxLenMsg))
	}
	if v.Pattern != "" && v.PatternMsg != "" {
		checks = append(checks, fmt.Sprintf(`if(el.validity.patternMismatch){el.setCustomValidity('%s');return;}`, v.PatternMsg))
	}
	if len(checks) == 0 {
		return ""
	}
	return fmt.Sprintf(
		`<script>(function(){var el=document.getElementById('%s');if(!el)return;el.addEventListener('input',function(){el.setCustomValidity('');%s});})();</script>`,
		id, strings.Join(checks, ""),
	)
}

// ─── SERVER-SIDE VALIDATION ─────────────────────────────────────────────────

// FormErrors maps field names to their error messages.
// An empty map means the form is valid.
type FormErrors map[string]string

// FieldRule defines server-side validation rules for a single form field.
type FieldRule struct {
	Required  bool
	MinLen    int
	MaxLen    int
	Pattern   string
	MinNum    *float64
	MaxNum    *float64
	CustomMsg string
}

// ValidateForm validates form values from r against the provided rules.
// Returns a FormErrors map; if empty, the form passed validation.
//
// Example:
//
//	errs := components.ValidateForm(r, map[string]components.FieldRule{
//	    "email": {Required: true, Pattern: `^[^@]+@[^@]+\.[^@]+$`, CustomMsg: "E-mail inválido."},
//	    "name":  {Required: true, MinLen: 2},
//	})
//	if len(errs) > 0 { /* render form with errors */ }
func ValidateForm(r *http.Request, rules map[string]FieldRule) FormErrors {
	_ = r.ParseForm()
	errs := make(FormErrors)

	for field, rule := range rules {
		val := strings.TrimSpace(r.FormValue(field))
		msg := rule.CustomMsg

		if rule.Required && val == "" {
			if msg == "" {
				msg = fmt.Sprintf("O campo %q é obrigatório.", field)
			}
			errs[field] = msg
			continue
		}
		if val == "" {
			continue
		}
		if rule.MinLen > 0 && len([]rune(val)) < rule.MinLen {
			if msg == "" {
				msg = fmt.Sprintf("Mínimo de %d caracteres.", rule.MinLen)
			}
			errs[field] = msg
			continue
		}
		if rule.MaxLen > 0 && len([]rune(val)) > rule.MaxLen {
			if msg == "" {
				msg = fmt.Sprintf("Máximo de %d caracteres.", rule.MaxLen)
			}
			errs[field] = msg
			continue
		}
		if rule.Pattern != "" {
			rx, err := regexp.Compile(rule.Pattern)
			if err == nil && !rx.MatchString(val) {
				if msg == "" {
					msg = fmt.Sprintf("O campo %q tem formato inválido.", field)
				}
				errs[field] = msg
				continue
			}
		}
		if rule.MinNum != nil || rule.MaxNum != nil {
			num, err := strconv.ParseFloat(val, 64)
			if err != nil {
				errs[field] = fmt.Sprintf("O campo %q deve ser numérico.", field)
				continue
			}
			if rule.MinNum != nil && num < *rule.MinNum {
				if msg == "" {
					msg = fmt.Sprintf("Valor mínimo: %g.", *rule.MinNum)
				}
				errs[field] = msg
				continue
			}
			if rule.MaxNum != nil && num > *rule.MaxNum {
				if msg == "" {
					msg = fmt.Sprintf("Valor máximo: %g.", *rule.MaxNum)
				}
				errs[field] = msg
				continue
			}
		}
	}
	return errs
}

// FieldError renders a validation error message below an input.
// Returns an empty Component when msg is blank (no-op in the layout).
func FieldError(msg string) Component {
	if msg == "" {
		return HTML("")
	}
	return HTML(fmt.Sprintf(`<span class="goui-field-error">%s</span>`, msg))
}
