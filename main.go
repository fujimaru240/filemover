package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type MoveRule struct {
	Pattern string
	DestDir string
}

func main() {
	// コマンドライン引数の定義
	sourceDir := flag.String("source", "", "ファイル群の格納先ディレクトリ")
	destBase := flag.String("dest", "", "移動先のベースディレクトリ")
	csvFile := flag.String("config", "config.csv", "設定CSVファイルのパス")
	dryRun := flag.Bool("dry-run", false, "実際には移動せず、動作を表示するのみ")
	flag.Parse()

	// 必須引数のチェック
	if *sourceDir == "" || *destBase == "" {
		fmt.Println("使用方法:")
		flag.PrintDefaults()
		fmt.Println("\n例: go run main.go -source ./source -dest ./destination -config config.csv")
		os.Exit(1)
	}

	// CSVファイルの読み込み
	rules, err := loadRules(*csvFile)
	if err != nil {
		fmt.Printf("CSVファイルの読み込みエラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("設定を読み込みました: %d件のルール\n", len(rules))

	// ソースディレクトリ内のファイルを取得
	files, err := os.ReadDir(*sourceDir)
	if err != nil {
		fmt.Printf("ソースディレクトリの読み込みエラー: %v\n", err)
		os.Exit(1)
	}

	movedCount := 0
	skippedCount := 0

	// 各ファイルをチェックして移動
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		matched := false

		// ルールに従ってマッチング
		for _, rule := range rules {
			re, _ := regexp.Compile(rule.Pattern)
			if re.MatchString(fileName) {
				// 移動先ディレクトリのパスを構築
				destDir := filepath.Join(*destBase, rule.DestDir)
				destPath := filepath.Join(destDir, fileName)
				sourcePath := filepath.Join(*sourceDir, fileName)

				if *dryRun {
					fmt.Printf("[DRY-RUN] %s -> %s\n", sourcePath, destPath)
				} else {
					// 移動先ディレクトリを作成
					if err := os.MkdirAll(destDir, 0755); err != nil {
						fmt.Printf("ディレクトリ作成エラー (%s): %v\n", destDir, err)
						skippedCount++
						matched = true
						break
					}

					// ファイルを移動
					if err := moveFile(sourcePath, destPath); err != nil {
						fmt.Printf("ファイル移動エラー (%s): %v\n", fileName, err)
						skippedCount++
					} else {
						fmt.Printf("移動完了: %s -> %s\n", fileName, destPath)
						movedCount++
					}
				}
				matched = true
				break
			}
		}

		if !matched {
			fmt.Printf("スキップ (パターン未一致): %s\n", fileName)
			skippedCount++
		}
	}

	// 結果サマリー
	fmt.Printf("\n=== 処理結果 ===\n")
	if *dryRun {
		fmt.Println("(DRY-RUNモード)")
	}
	fmt.Printf("移動: %d件\n", movedCount)
	fmt.Printf("スキップ: %d件\n", skippedCount)
}

// loadRules CSVファイルから移動ルールを読み込む
func loadRules(csvPath string) ([]MoveRule, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	var rules []MoveRule
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// 空行やヘッダーをスキップ
		if len(record) < 2 || record[0] == "" {
			continue
		}

		rules = append(rules, MoveRule{
			Pattern: strings.TrimSpace(record[0]),
			DestDir: strings.TrimSpace(record[1]),
		})
	}

	return rules, nil
}

// moveFile ファイルを移動する
func moveFile(source, dest string) error {
	// 同名ファイルが存在する場合はエラー
	if _, err := os.Stat(dest); err == nil {
		return fmt.Errorf("移動先に同名ファイルが既に存在します")
	}

	// os.Renameを試す（同じファイルシステム内の場合）
	err := os.Rename(source, dest)
	if err == nil {
		return nil
	}

	// os.Renameが失敗した場合はコピー&削除
	if err := copyFile(source, dest); err != nil {
		return err
	}

	return os.Remove(source)
}

// copyFile ファイルをコピーする
func copyFile(source, dest string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// パーミッションをコピー
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}
	return os.Chmod(dest, sourceInfo.Mode())
}
