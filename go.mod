module github.com/iburinoc/wg-docker-net

go 1.15

require (
	github.com/coreos/go-iptables v0.5.0
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-plugins-helpers v0.0.0-20200102110956-c9a8a2d92ccc
	github.com/hashicorp/go-multierror v1.1.0
	github.com/vishvananda/netlink v1.1.0
	github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	gopkg.in/go-ini/ini.v1 v1.62.0
)

replace github.com/docker/go-plugins-helpers => github.com/iburinoc/go-plugins-helpers v0.0.2
