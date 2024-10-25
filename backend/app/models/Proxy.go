package models

type Proxy struct {
	Server   string `json:"server" bson:"server"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Status   bool   `json:"status" bson:"status"`
}

func (c *Proxy) GetTableName() string {
	return "proxies"
}
