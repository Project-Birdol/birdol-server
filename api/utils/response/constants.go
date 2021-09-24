package response

const (
	// result strings
	ResultOK 				= "ok"
	ResultFail 				= "failed"
	ResultNeedTokenRefresh 	= "need_refresh"
	ResultRefreshSuccess 	= "refreshed"

	// error strings
	ErrInvalidType			= "invalid_content_type"
	ErrAuthorizationFail 	= "fail_authorization"
	ErrInvalidToken 		= "invalid_token"
	ErrInvalidDevice 		= "invalid_device"
	ErrFailParseXML 		= "fail_parsing_xml"
	ErrInvalidSignature 	= "invalid_signature"
	ErrNotLoggedIn			= "not_logged_in"
	ErrFailParseJSON		= "fail_parsing_json"
	ErrFailAccountCreation	= "fail_account_creation"
	ErrInvalidAccount		= "invalid_account"
	ErrPasswordExpire		= "password_expired"
	ErrInvalidPassword		= "invalid_password"
	ErrFailDataLink			= "fail_datalink"
	ErrFailSetPassword		= "fail_set_password"
	ErrFailUnlink			= "fail_unlink"
	ErrFailCreateSession	= "fail_create_session"
	ErrInvalidRefreshToken	= "invalid_refresh_token"
	ErrFailRefresh		= "fail_refresh"
)
