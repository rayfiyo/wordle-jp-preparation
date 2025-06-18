# 長音をいい感じにするやつ、途中

import csv
import sys


def replace_long_marks(input_path: str) -> None:
    with open(input_path, newline="", encoding="utf-8") as infile:
        reader = csv.reader(infile)
        writer = csv.writer(sys.stdout)

        for row in reader:
            # 1列目と2列目を取得
            first, second = row[0], row[1]

            # 長音記号「ー」が含まれるなら位置対応で1列目を置換
            if "ー" in second:
                # 文字数が異なる場合はそのままスキップ or パディングするようにアレンジ可
                # if len(first) != len(second):
                # print(
                # f"Warning: length mismatch in row {row}",
                # file=sys.stderr,
                # )
                # min長だけループして置換

                new_first = "".join(
                    (
                        "ー"
                        if i < len(second) and second[i] == "ー"
                        else first[i] if i < len(first) else ""
                    )
                    for i in range(max(len(first), len(second)))
                )
                row[0] = new_first

            writer.writerow(row)


if __name__ == "__main__":
    replace_long_marks("for_wordle_jp.csv")
