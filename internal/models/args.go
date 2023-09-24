package models

type Args struct {
	File           string
	Out            string
	Style          string
	Watch          bool
	HttpServer     bool
	Debug          bool
	ServerPort     int
	ServerHostname string
}
