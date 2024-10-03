package internal

import "embed"

var (
	MaxRetries   = 3
	TemplatePath = "user_invitations.tmpl"
)

//go:embed templates
var FS embed.FS

type Client interface {
	Send(templateFile string, toEmail []string, data any) error
}
