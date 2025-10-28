### ビルド

```bash
GOOS=windows GOARCH=amd64 go build -o filemover.exe main.go
GOOS=darwin GOARCH=amd64 go build -o filemover-mac-intel main.go
GOOS=darwin GOARCH=arm64 go build -o filemover-mac-arm main.go
GOOS=linux GOARCH=amd64 go build -o filemover-linux main.go
```

### 実行

```bash
./filemover.exe -config config.csv -source source/ -dest destination/
./filemover-linux -config config.csv -source source/ -dest destination/
```
