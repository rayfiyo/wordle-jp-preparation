# wordle-jp-preparation

- Wordle Japanese の対策
- Wordle Japanese Preparation

## more

- mottox2 さんの [🔴Wordle Japanese](https://wordle.mottox2.com) （以降、WordleJP）対策。
- WordleJP は CC BY-NC-SA 3.0 で公開されている、[国立国語研究所(2004)『分類語彙表増補改訂版データベース』(ver.1.0) (masayu-a/WLSP)](https://github.com/masayu-a/WLSP) を元データとして利用しているので、このリポジトリもお世話になった。(2025-06-18 アクセス)
- また、Go のプログラム `main.go` は単語帳アプリ [MyFlasher](https://apps.apple.com/jp/app/id644699476) で使える練習問題（１～３文字の伏せ字）を生成する

## dir

- `bunruidb_utf-8.csv`
  - masayu-a/WLSP の`bunruidb.csv`を、BOM なしUTF-8 にしたもの
- `five_character_only.csv`
  - Wordle jp 用に抽出した５文字の単語
  - `ひらがな,漢字（最適な表現）` の形式にしている
  - 抽出方法は次の標準出力をリダイレクトした（`> five_character_only.csv`）
    ```bash
    awk -F, 'length($14)==5 { print $14,$12 }' bunruidb_utf-8.csv | sed -E 's/ /,/g'
    ```
- `for_wordle_jp.csv`

  - `for_wordle_jp.csv` の読み（１列目）を Wordle で出題される形式に変換したもの
  - 掲載の単語と Wordle JP で入力できる形式になれるという目的で、
    単語帳アプリ等で使えるよう `ひらがな（Wordleで入力できる形式）,漢字（最適な表現）`
    の形式にしている
  - 抽出方法は次の標準出力をリダイレクトした（`> for_wordle_jp.csv`）
    ```bash
    awk -F',' '{
      cmd = "echo " $1 " | sed -e \
      \"s/が/か/g; s/ぎ/き/g; s/ぐ/く/g; s/げ/け/g; s/ご/こ/g; \
        s/ざ/さ/g; s/じ/し/g; s/ず/す/g; s/ぜ/せ/g; s/ぞ/そ/g; \
        s/だ/た/g; s/ぢ/ち/g; s/づ/つ/g; s/で/て/g; s/ど/と/g; \
        s/ば/は/g; s/び/ひ/g; s/ぶ/ふ/g; s/べ/へ/g; s/ぼ/ほ/g; \
        s/ぱ/は/g; s/ぴ/ひ/g; s/ぷ/ふ/g; s/ぺ/へ/g; s/ぽ/ほ/g; \
        s/ぁ/あ/g; s/ぃ/い/g; s/ぅ/う/g; s/ぇ/え/g; s/ぉ/お/g; \
        s/ゃ/や/g; s/ゅ/ゆ/g; s/ょ/よ/g; s/っ/つ/g; s/ゎ/わ/g; s/ゕ/か/g; s/ゖ/け/g\""
        cmd | getline new_col1; close(cmd)
        print new_col1 "," $2
    }' five_character_only.csv
    ```
