### ビルド

```bash
GOOS=windows GOARCH=amd64 go build -o filemover.exe main.go
GOOS=darwin GOARCH=amd64 go build -o filemover-mac-intel main.go
GOOS=darwin GOARCH=arm64 go build -o filemover-mac-arm main.go
GOOS=linux GOARCH=amd64 go build -o filemover-linux main.go
```

### 実行

- 使用方法:
  - `-config string`
    - 設定CSVファイルのパス (default "config.csv")
  - `-dest string`
    - 移動先のベースディレクトリ
  - `-dry-run`
    - 実際には移動せず、動作を表示するのみ
  - `-source string`
    - ファイル群の格納先ディレクトリ

```bash
./filemover.exe -config config.csv -source source/ -dest destination/
./filemover-linux -config config.csv -source source/ -dest destination/
```
