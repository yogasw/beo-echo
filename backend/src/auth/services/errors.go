package services

// OAuth Error types
const (
	// General OAuth errors
	ErrOAuthDisabled          = "oauth_service_disabled"
	ErrOAuthNotConfigured     = "oauth_not_configured"
	ErrOAuthConfigRetrieval   = "oauth_config_retrieval_failed"
	ErrOAuthUserInfoRetrieval = "oauth_user_info_retrieval_failed"

	// Google OAuth specific errors
	ErrGoogleOAuthDisabled      = "google_oauth_service_disabled"
	ErrGoogleOAuthNotConfigured = "google_oauth_not_configured"

	// Registration errors
	ErrAutoRegistrationDisabled = "auto_registration_disabled"

	// Domain errors
	ErrDomainNotAllowed = "email_domain_not_allowed"

	// Request validation errors
	ErrMissingRedirectURI = "missing_redirect_uri"
	ErrInvalidRequestBody = "invalid_request_body"
	ErrMissingCode        = "missing_code"
	ErrInvalidState       = "invalid_state"
)

// Error message constructor functions
func NewAutoRegistrationDisabledError() error {
	return &OAuthError{
		Type:    ErrAutoRegistrationDisabled,
		Message: "Auto-registration is disabled and user does not exist. Please contact your administrator.",
	}
}

func NewDomainNotAllowedError(domain string) error {
	return &OAuthError{
		Type:    ErrDomainNotAllowed,
		Message: "Your email domain '" + domain + "' is not allowed. Please contact your administrator.",
		Data: map[string]interface{}{
			"domain": domain,
		},
	}
}

func NewGoogleOAuthDisabledError() error {
	return &OAuthError{
		Type:    ErrGoogleOAuthDisabled,
		Message: "Google OAuth service is disabled. Please contact your administrator.",
	}
}

func NewGoogleOAuthNotConfiguredError() error {
	return &OAuthError{
		Type:    ErrGoogleOAuthNotConfigured,
		Message: "Google OAuth credentials are not configured. Please contact your administrator.",
	}
}

func NewOAuthUserInfoRetrievalError() error {
	return &OAuthError{
		Type:    ErrOAuthUserInfoRetrieval,
		Message: "failed to retrieve user information, please try again later",
	}
}

func NewOAuthConfigRetrievalError(detail string) error {
	message := "failed to retrieve OAuth configuration"
	if detail != "" {
		message += ": " + detail
	}

	return &OAuthError{
		Type:    ErrOAuthConfigRetrieval,
		Message: message,
	}
}

func NewInvalidRequestBodyError() error {
	return &OAuthError{
		Type:    ErrInvalidRequestBody,
		Message: "invalid request body",
	}
}

func NewMissingRedirectURIError() error {
	return &OAuthError{
		Type:    ErrMissingRedirectURI,
		Message: "redirect_uri parameter is required",
	}
}

func NewMissingCodeError() error {
	return &OAuthError{
		Type:    ErrMissingCode,
		Message: "authorization code is required",
	}
}

func NewInvalidStateError() error {
	return &OAuthError{
		Type:    ErrInvalidState,
		Message: "invalid state parameter",
	}
}

// OAuthError represents a structured OAuth error
type OAuthError struct {
	Type    string
	Message string
	Data    map[string]interface{}
}

// Error implements the error interface
func (e *OAuthError) Error() string {
	return e.Message
}

// GetType returns the error type
func (e *OAuthError) GetType() string {
	return e.Type
}

// GetData returns any additional error data
func (e *OAuthError) GetData() map[string]interface{} {
	return e.Data
}
