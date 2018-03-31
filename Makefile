run:IPAddressTool.exe
	IPAddressTool.exe
IPAddressTool.exe:main.go
	go build -ldflags="-H windowsgui"