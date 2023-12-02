package models

type Args struct {
	File               string
	Out                string
	Style              string
	Watch              bool
	HttpServer         bool
	Open               bool
	Debug              bool
	ServerPort         int
	ServerHostname     string
	NoExternalLibs     bool
	NoServerHeaderWait bool
}
