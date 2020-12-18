setlocal
call :"%1"
endlocal
exit /b

:""
:"all"
    go fmt
    go build
    exit /b

:"test"
    cd unicodetest && go fmt
    go run unicodetest/main.go
    exit /b

:"test-linux"
    set "GOOS=linux"
    pushd unicodetest
    go build
    "%windir%\System32\bash" -c "./unicodetest"
    popd
    exit /b

:"get"
    go get -u all
    go mod tidy
    exit /b
