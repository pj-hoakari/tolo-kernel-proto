# tolo-kernel-proto

共有カーネル`tolo/kernel/v1/kernel.proto`の原本と、その配布物を置くリポジトリ  
hoakari-rootにある定義は参照用の写しとして扱う

1回のリリース（`vX.Y.Z`）で、次の2つを同じ版で公開する

| 配布物 | 場所 | 消費側 |
| -- | -- | -- |
| Goモジュール | `github.com/pj-hoakari/tolo-kernel-proto` | Goのサービス |
| ORASの成果物（`.proto`） | `ghcr.io/pj-hoakari/tolo-kernel-proto` | Goを使わない消費側（Pythonのサービス、スタッフアプリ） |

Goの生成コードのimportパスは`github.com/pj-hoakari/tolo-kernel-proto/gen/tolo/kernel/v1`（package `kernelv1`）  
`kernel.proto`をimportするGoのサービスは、kernelのGoコードを自前で生成せず、このパッケージを参照する  
同じprotoファイルを2つのGoパッケージが登録すると、プロセスの起動時にpanicする

## Goの消費側

CLIを`go.mod`のtoolに登録する

```sh
go get -tool github.com/pj-hoakari/tolo-kernel-proto/cmd/tolo-kernel-proto@v0.1.0
```

モジュールに埋め込んだ`kernel.proto`を、bufの入力ディレクトリへ書き出す

```sh
go tool tolo-kernel-proto --write-proto proto
# proto/tolo/kernel/v1/kernel.proto が作成される
```

書き出すファイルの版はgo.modで決まるので、`.gitignore`に入れてよい  
`buf.gen.go.yaml`ではkernelを生成対象から外し、`go_package`の上書きを無効にする

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

## Goを使わない消費側

同じ版の`.proto`をORASで取得し、各自で生成する  
各サービスが公開する`.proto`の成果物には`kernel.proto`を含めない。サービスの`.proto`がkernelをimportする場合は、この成果物も同じディレクトリに取得する  
版は、取得するサービスの`go.mod`が参照する`github.com/pj-hoakari/tolo-kernel-proto`の版のうち、最も新しいものにする

```sh
oras pull ghcr.io/pj-hoakari/tolo-kernel-proto:0.1.0 -o proto
```

## 開発

ツールはmiseで入れる

```sh
mise install
```

| タスク | 内容 |
| -- | -- |
| `task proto` | `buf lint`とGoコードの再生成 |
| `task proto:breaking` | 最新のリリースタグと比べる`buf breaking` |
| `task release VERSION=vX.Y.Z` | GitHubリリースの作成 |

`kernel.proto`を変更したら、`task proto`で`gen/`を更新してコミットする  
CIは、生成コードの差分、`buf lint`、`buf breaking`、`go test`、golangci-lint、actionlintを検査する

## リリース

```sh
task release VERSION=v0.1.0
```

リリースが公開されると、`publish-proto.yml`が同じ版のORASの成果物を公開する  
Goモジュールの版には、リリースのタグがそのまま使われる

破壊的変更を含む版をリリースしたら、`kernel.proto`をimportする全サービスの`go.mod`を同じ時期にその版へ上げる  
Goを使わない消費側は、サービスの版のうち最も新しいものに合わせて`kernel.proto`を取得するので、上げていないサービスの`.proto`と食い違う
