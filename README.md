# tolo-kernel-proto

共有カーネル `tolo/kernel/v1/kernel.proto` の正本と、その配布物を持つリポジトリ

1 つのリリース（`vX.Y.Z`）から次の 2 つを同じ版で配布する

| 配布物 | 場所 | 消費側 |
| -- | -- | -- |
| Go モジュール | `github.com/pj-hoakari/tolo-kernel-proto` | Go のサービス |
| ORAS の成果物（`.proto`） | `ghcr.io/pj-hoakari/tolo-kernel-proto` | Python のサービス |

Go の生成コードの import パスは `github.com/pj-hoakari/tolo-kernel-proto/gen/tolo/kernel/v1`（package `kernelv1`）  
`kernel.proto` を import するサービスは、kernel の Go コードを自前で生成せず、このパッケージを参照する  
同じ proto ファイルが 2 つの Go パッケージから登録されると、プロセスの起動時に panic するため

## Go の消費側

`go.mod` に CLI を tool として登録する

```sh
go get -tool github.com/pj-hoakari/tolo-kernel-proto/cmd/tolo-kernel-proto@v0.1.0
```

埋め込まれた `kernel.proto` を buf の入力ディレクトリへ書き出す

```sh
go tool tolo-kernel-proto --write-proto proto
# proto/tolo/kernel/v1/kernel.proto が作成される
```

書き出したファイルは版が go.mod で決まるため、`.gitignore` に入れてよい  
`buf.gen.go.yaml` では kernel を生成対象から除外し、`go_package` の上書きを無効にする

```yaml
version: v2
managed:
    enabled: true
    override:
        - file_option: go_package_prefix
          value: github.com/pj-hoakari/<service>/gen
    disable:
        - file_option: go_package
          path: tolo/kernel/v1/kernel.proto
plugins:
    - local: ["go", "tool", "protoc-gen-go"]
      out: gen
      opt: paths=source_relative
inputs:
    - directory: proto
      exclude_paths:
        - proto/tolo/kernel
```

## Python の消費側

同じ版の `.proto` を ORAS で取得し、自前で生成する

```sh
oras pull ghcr.io/pj-hoakari/tolo-kernel-proto:0.1.0 -o proto
```

## 開発

ツールは mise で揃える

```sh
mise install
```

| タスク | 内容 |
| -- | -- |
| `task proto` | `buf lint` と Go コードの再生成 |
| `task proto:breaking` | 最新のリリースタグとの `buf breaking` |
| `task release VERSION=vX.Y.Z` | GitHub リリースの作成 |

`kernel.proto` を変えたら `task proto` で `gen/` を更新してコミットする  
CI は生成コードの差分、`buf lint`、`buf breaking`、`go test`、golangci-lint、actionlint を検査する

## リリース

```sh
task release VERSION=v0.1.0
```

リリースの公開を契機に `publish-proto.yml` が同じ版の ORAS の成果物を公開する  
Go モジュールの版はリリースのタグそのものになる
