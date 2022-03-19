# ecs
A Go-implementation of the ECS (Entity-Component-System).

ECS engine for RockGO

Need go1.18 for generics code

## Install

how to get ecs?
```
$: go get -u github.com/zllangct/ecs
```

## Benchmark

```azure
goos: darwin
goarch: amd64
pkg: test_ecs_m_d
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkNormal
BenchmarkNormal-12             	      27	  39827654 ns/op
BenchmarkNormalParallel
BenchmarkNormalParallel-12     	      20	  54198601 ns/op
BenchmarkEcsCollectionV1
BenchmarkEcsCollectionV1-12    	     453	   2687113 ns/op
BenchmarkEcsCollectionV2
BenchmarkEcsCollectionV2-12    	     442	   2813405 ns/op
```

## Example

- [benchmark-0](https://github.com/zllangct/ecs/tree/master/example/benchmark-0)
- [fake game server](https://github.com/zllangct/ecs/tree/master/example/fake-simple-game-server)
