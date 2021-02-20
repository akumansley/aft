package server

type Option func(*config)

type config struct {
	ServePort   string
	ServeDir    string
	CatalogPort string
	Authed      bool
	DBLogPath   string
	TLSCert     string
	TLSKey      string
}

func newConfig() *config {
	return &config{
		ServePort:   "8081",
		ServeDir:    "",
		CatalogPort: "8080",
		Authed:      true,
		DBLogPath:   "",
	}
}

func ServeDir(dir string) Option {
	return func(c *config) {
		c.ServeDir = dir
	}
}

func ServePort(port string) Option {
	return func(c *config) {
		c.ServePort = port
	}
}

func CatalogPort(port string) Option {
	return func(c *config) {
		c.CatalogPort = port
	}
}

func DBLogPath(path string) Option {
	return func(c *config) {
		c.DBLogPath = path
	}
}

func Authed(authed bool) Option {
	return func(c *config) {
		c.Authed = authed
	}
}

func TLSCert(path string) Option {
	return func(c *config) {
		c.TLSCert = path
	}
}

func TLSKey(path string) Option {
	return func(c *config) {
		c.TLSKey = path
	}
}
