package components

import (
	"fmt"
)

type ToastType string

const (
	ToastSuccess ToastType = "success"
	ToastError   ToastType = "error"
	ToastInfo    ToastType = "info"
	ToastWarning ToastType = "warning"
)

// ToastContainer renders the hidden container where toasts will be injected.
// This should be included once in the main layout (usually via Headbar/AppLayout).
func ToastContainer() string {
	return `<div id="goui-toast-container" class="goui-toast-container"></div>`
}

// ShowToast returns a script that triggers a toast notification.
// Use this to send feedback from the backend during interactions.
func ShowToast(message string, t ToastType) string {
	return fmt.Sprintf(`<script>gouiToast("%s", "%s")</script>`, message, t)
}
