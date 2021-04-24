module github.com/slyerr/econfig/server

go 1.16

require (
	github.com/slyerr/econfig/core v0.0.0
	github.com/slyerr/verifier v0.1.0
	go.etcd.io/bbolt v1.3.5
	go.uber.org/zap v1.16.0
	golang.org/x/sys v0.0.0-20210415045647-66c3f260301c // indirect
)

replace github.com/slyerr/econfig/core => ../core
