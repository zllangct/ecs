module github.com/zllangct/ecs/example/example2

go 1.24

require (
	github.com/zllangct/ecs v0.0.0
	github.com/zllangct/rockmem v0.0.0
)

replace (
	github.com/zllangct/ecs => ../..
	github.com/zllangct/rockmem => ../../../rockmem
)
