module example/client

go 1.16

require (
	github.com/slyerr/econfig/client v0.0.0
	github.com/slyerr/econfig/core v0.0.0
	go.uber.org/zap v1.16.0
)

replace github.com/slyerr/econfig/client => ../../client

replace github.com/slyerr/econfig/core => ../../core
