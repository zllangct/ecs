go build -o bin\rockmem-ecs.exe .\cmd\rockmem-ecs\
bin\rockmem-ecs.exe generate -o . ecs.rm
bin\rockmem-ecs.exe generate -o .\test\cmp\testdata\ .\test\cmp\testdata\components.rm
