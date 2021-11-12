package common

// GocannonPlugin is an interface that has to be satisfied by a custom gocannnon plugin
type GocannonPlugin interface {
	// function called on gocannon startup with a config passed to it
	Startup(cfg Config)
	// function called before each request is sent
	BeforeRequest(cid int) (target string, method string, body RawRequestBody, headers RequestHeaders)
	// function that returns the plugin's name
	GetName() string
}
