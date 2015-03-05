@echo off

if exist build.bat goto ok
echo build.bat must be run from its folder
goto end

:ok

cd main

@if exist "main-res.syso" (
    @del "main-res.syso"
)

@if exist "%~dp0bin\goapp.exe" (
    @del "%~dp0bin\goapp.exe"
)

windres -o main-res.syso main.rc


IF "%1"=="noconsole" (
    go build -ldflags="-H windowsgui" -o ../bin/goapp.exe
    rem @if %ERRORLEVEL% neq 0 goto end
) else (
    go build -o ../bin/goapp.exe
    rem @if %ERRORLEVEL% neq 0 goto end
)

cd ..\browser
go build -o ../bin/browser.exe

cd ..\appserver
go build -o ../bin/appserver.exe

cd ../bin
goapp.exe
cd ..

:end

pause
