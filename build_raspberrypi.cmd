setlocal 
  set GOOS=linux&& set GOARCH=arm&& set GOARM=5&& go build -o gover -v -mod vendor
endlocal