@echo off
cd /d %~dp0

REM 生成组件代码
echo Generating component code...
..\..\bin\rockmem-ecs.exe generate components.rm -o ./components -l golang

if errorlevel 1 (
    echo Failed to generate component code
    exit /b 1
)

echo Component code generated successfully!
