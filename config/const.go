package config

const (
	AccountSID          = "ACc39725a654edb7264186eca22f221a47"
	AuthToken           = "0f95c1f4de45c8e8fed7a68302acea7d"
	TwilioPhoneNumber   = "+13026045203"
	RedisAddr           = "localhost:6379" // Redis server address
	CUSTOMER_ROLE       = "customer"
	ADMIN_ROLE          = "admin"
	ERR_INFORMATION     = "The server has received the request and is continuing the process"
	SUCCESS             = "The request was successful"
	ERR_REDIRECTION     = "You have been redirected and the completion of the request requires further action"
	ERR_BADREQUEST      = "Bad request"
	ERR_INTERNAL_SERVER = "While the request appears to be valid, the server could not complete the request"
)

var SignedKey = []byte("MGJd@Ro]yKoCc)mVY1^c:upz~4rn9Pt!hYd]>c8dt#+%")
