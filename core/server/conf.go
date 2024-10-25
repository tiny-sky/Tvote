package server

type Settings struct {
	ListenOn       string `yaml:"listenOn"`
	Interval       int    `yaml:"interval"`
	MaxTicketUsage int    `yaml:"maxTicketUsage"`
}
