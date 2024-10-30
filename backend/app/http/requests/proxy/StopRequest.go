package proxyrequest

type StopProxy struct {
	Proxy ProxyPayload `json:"proxy"`
	Error string       `json:"error"`
}
type ProxyPayload struct {
	ID       string `json:"id" bson:"id"`
	Server   string `json:"server" bson:"server"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}
