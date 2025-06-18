package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	// フラグ定義
	lines := flag.Int("n", 12, "選択する行数")
	filePath := flag.String("f", "for_wordle_jp.csv", "入力CSVファイル名")
	savePath := flag.String("o", "./assets/", "保存ディレクトリ（ファイル名は指定不可）")
	blank := flag.String("b", "＠", "伏せ字の文字")
	flag.Parse()

	// CSV読み込み
	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "入力ファイルを開けません: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	r := csv.NewReader(file)

	// レコード収集
	var records [][]string
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "CSV読み込みエラー: %v\n", err)
			os.Exit(1)
		}
		records = append(records, rec)
	}

	// 選択対象数
	maxRow := len(records)
	maxRow = min(maxRow, 19786)
	if *lines > maxRow {
		fmt.Fprintf(os.Stderr,
			"-n の値が大きすぎます: 最大 %d 行まで選択可能です\n", maxRow)
		os.Exit(1)
	}

	// ランダムに重複なし行を選択
	perm := rand.Perm(maxRow)
	chosenIdx := perm[:*lines]

	// 出力用レコード
	out := make([][]string, 0, *lines)
	runedBlank := []rune(*blank)[0]
	for _, idx := range chosenIdx {
		row := records[idx]
		first := row[0]
		// rune スライスに変換
		runes := []rune(first)
		length := len(runes)
		// マスク処理: 重複を許して3回ランダム置換
		for range 3 {
			rIdx := rand.N(length)
			runes[rIdx] = runedBlank
		}
		newFirst := string(runes)
		out = append(out, []string{newFirst, row[1]})
	}

	// 出力ディレクトリ準備
	if err := os.MkdirAll(*savePath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "ディレクトリ作成エラー: %v\n", err)
		os.Exit(1)
	}

	// 既存ファイル調査
	files, err := filepath.Glob(filepath.Join(*savePath, "*.csv"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ファイル一覧取得エラー: %v\n", err)
		os.Exit(1)
	}
	maxNum := 0
	for _, fn := range files {
		base := filepath.Base(fn)
		numStr := base[:len(base)-len(filepath.Ext(base))]
		// ファイル名末尾に日付が含まれている場合はトリム
		if parts := strings.Split(numStr, "_"); len(parts) > 1 {
			numStr = parts[0]
		}
		if num, err := strconv.Atoi(numStr); err == nil {
			maxNum = max(num, maxNum)
		}
	}

	// 新ファイル名生成
	next := maxNum + 1
	// 東京時間を明示的に使用
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		loc = time.Local
	}
	dateStr := time.Now().In(loc).Format("20060102")
	outName := fmt.Sprintf("%06d_%s.csv", next, dateStr)
	outPath := filepath.Join(*savePath, outName)

	// CSV書き出し
	outFile, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "出力ファイル作成エラー: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	// ヘッダー行を書き込む（すべて二重引用符で囲む）
	outFile.WriteString(`"Question","Answer"` + "\n")

	// レコードを書き込む（全フィールドを二重引用符で囲む）
	for _, rec := range out {
		var fields []string
		for _, v := range rec {
			escaped := strings.ReplaceAll(v, `"`, `""`)
			fields = append(fields, `"`+escaped+`"`)
		}
		line := strings.Join(fields, ",") + "\n"
		outFile.WriteString(line)
	}

	fmt.Printf("出力完了: %s\n", outPath)
}
