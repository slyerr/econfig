module github.com/slyerr/econfig/client

go 1.16

require (
	github.com/pkg/errors v0.9.1
	github.com/slyerr/econfig/core v0.0.0
	github.com/slyerr/verifier v0.1.0
	go.uber.org/zap v1.16.0
)

replace github.com/slyerr/econfig/core => ../core
