- UNIX系プラットフォームでのインクリメンタルサーチで Backspace キー`\x7F` が直前の文字を削除しない不具合を修正 (Ctrl-H は機能していた)

v1.6.1
======
Nov 8, 2024

- UNIX系デスクトップのターミナルの機能で、複数行を貼り付けた時、一部のテキストをとりこぼす不具合を修正 ([hymkor/go-multiline-ny v0.16.2](https://github.com/hymkor/go-multiline-ny/releases/tag/v0.16.2) 向け修正)

v1.6.0
======
Nov 4, 2024

- 入力予想機能の予想関数を差し替えられるようにした(`(*Editor) Predictor = ...`)

v1.5.0
======
Oct 6, 2024

- PowerShell 7 風の入力予測機能の実装

次のように表示部分の表示開始・終了時のエスケープシーケンスを設定することで有効化される

```
editor := &readline.Editor{
    PredictColor: [...]string{"\x1B[3;22;34m", "\x1B[23;39m"},
}
```

v1.4.1
======
Oct 2, 2024

- WindowsTerminal 1.22 向けに、Zero Width Join Sequence の文字幅の計算関数を差し替えられるようにした(SetZWJSWidthGetter)

v1.4.0
======
Jun 16, 2024

- Ctrl-P/N: 履歴を切り替えるときに変更したエントリを保存し、(Enterが入力されるまでは)再度切り替えたときに復元するようにした

v1.3.1
======
Apr 21, 2024

- [#6] サブパッケージ completion で、空状態で補完するとクラッシュする不具合を修正  
  (Thanks [@glejeune])

[#6]: https://github.com/nyaosorg/go-readline-ny/pull/6
[@glejeune]: https://github.com/glejeune

v1.3.0
======
Apr 17, 2024

- 定数追加: keys.ShiftTAB = "\x1B[Z"
- 端末のインターフェイスと実装を簡素化 (`ITty` 周辺の互換性破壊)

v1.2.0
======
Feb 29, 2024

- "keys": キーの名前定数は型無しとした (元々は `keys.Code` だった)

v1.1.0
======
Feb 27, 2024

- "completion": 1候補に絞れた時に空白を追加していたが、空白のかわりに CmdCompletion.Postfix や CmdCompletionOrList.Postfix で指定できるようにした (デフォルトは0文字)

v1.0.1
======
Oct 08, 2023

- 合字や異体字など単純なコードポイントでない文字位置で色を変えられなかった問題を修正

v1.0.0
======
Oct 06, 2023

- バージョンを v1.0.0 にしました。

v0.15.2
=======
Oct 02, 2023

- カーソルが末尾にある時、`Coloring.Next` が CursorPositionDummyRune(U+E000) を受けとれない問題を修正
- `(ColorSequence) Chain` を追加 (二つのColorSequenceを連結する)

v0.15.1
=======
Oct 01, 2023

- `(ColorSequence) Add` を作成
- カーソル位置で、`Coloring.Next` は CursorPositionDummyRune(U+E000) を受けとるようにした

v0.15.0
=======
Sep 29, 2023

- v1.0.0 に向けて、非推奨としていたフィールド、メソッド、関数を削除しました
    - `KeyGoFuncT`. `GoCommand` をかわりに使ってください
    - `moji.GetCharWidth`
    - `GetFunc`. `nameutils.GetFunc` をかわりに使ってください
    - `(*Editor) LineFeed`. `(*Editor) LineFeedWriter` をかわりに使ってください
    - `(*Editor) Prompt`. `(*Editor) PromptWriter` をかわりに使ってください
    - `(*Editor) GetBindKey`
    - `(*KeyMap) BindKeyFunc`. `nameutils.BindKeyFunc()` をかわりに使ってください
    - `(*KeyMap) BindKeyClosure`
    - `(*KeyMap) GetBindKey`
    - `(*KeyMap) BindKeySymbol`. `nameutils.BindKeySymbol` をかわりに使ってください

v0.14.1
=======
Sep 13 2023

- 関数 `GetKey(tty) string` を公開
    - 内部的には前からあった関数でしたが、利用パッケージで同等の関数を再実装することが多かったので、ユーティリティー関数として公開しました。
    - GetKey は `(*TTY) ReadRune` と違って、`"[\x1B[A"` のような一連のキー入力シーケンスを１文字列といて返します
    - 引数 *tty* は [go-tty.TTY] もしくは、その互換型を想定。メソッド Raw(), ReadRune(), Buffered() を実装していなければいけません
- `Editor.Prompt` フィールドに **Deprecated** コメントをセットしました

[go-tty.TTY]: https://pkg.go.dev/github.com/mattn/go-tty#TTY

v0.14.0
=======
Aug 28 2023

- `(*Editor) PromptWriter` が `Ctrl-H` や `ESC]...\007` を出力した場合でも、プロンプトの幅を正しくカウントするようになった
- プロンプトの最後の行とユーザ入力テキストだけを再表示する `(*Buffer) RepaintLastLine()` を実装した。
    - プロンプト中の `\n` を `\r` に置換して出力する
    - 再表示の際に不必要に行が改められるのを避けるために実装した。
    - 利用するには PromptWriter フィールドが設定されていることが必要になります

v0.13.2
=======
Jul 29 2023

- keys/code.go で `\x` と書くべきものが `0x` になっていたリテラルを修正
- カーソルオフの状態で呼ばれると、そのままカーソルがオンにならず、表示されない不具合を修正

v0.13.1
=======
May 28 2023

- メソッド `(*KeyMap) Lookup` 追加

v0.13.0
=======
May 19 2023

- タブ文字を `^I` の代わりに４桁ごとの位置までの空白で表現できるようになった。

v0.12.3
=======
May 16 2023

- 新メソッド `(*Editor) Init` を追加
    - nil なフィールドをデフォルト値へ置き換える。`(*Editor) Readline` を呼び出す前に `Editor` のフィールドを参照する時、明示的に `(*Editor) Init` を呼ばなくてはいけない。主に go-readline-ny のアドオンパッケージ向け

v0.12.2
=======
May 15 2023

- パスがディレクトリをまったく含まないとき、`completion.File` が失敗する問題を修正

v0.12.1
=======
May 15 2023

-  `CmdCompletion` and `CmdCompletionOrList` で候補の絞り込みをするようにしたので、`Completion` interface 側で絞り込みをしなくてもよくなった
- `Completion.List` は basenames に nil を設定して省略できるようになった。
- `completion.File` で、ファイル名の中に `./` が含まれている時にひとつもマッチしなかった問題を修正（filepath パッケージが `./` を削除してしまうため）

v0.12.0
=======
May 13 2023

- [nyagos](https://github.com/nyaosorg/nyagos) の補完のサブセットを、サブパッケージ `completion` として逆輸入した。
- サブパッケージ `completion` を使用するよう `examples/example2.go` を書き換えた。
- LineFeed に io.Writer パラメーターを追加した新フィールド `LineFeedWriter` を追加。

v0.11.7
=======
May 12 2023

- グローバルな sync.Mutex 変数を `Editor` インスタンスの変数に変更した。キーに割り当てられたコマンド内で、新たな `Editor` インスタンスを生成して、`(*Editor) ReadLine` を呼ぶことができない問題があった。
- インスタンスのテーブル・グローバルなテーブルの双方から、キーに割り当てられたコマンドを探す `(*Editor) loolup(KEY)` を `LookUp` に改名（公開）
- `(*KeyMap) BindKey(KEY,nil)` で KEY に割り当てられた機能を削除できるようにした。
-  `b:=&strings.Builder{};c.Moji.WriteTo(b);b.String()` と等価な動作をする  `(c Cell) String` を実装
- `(*Buffer) GetKey` で `(*Buffer).Out.Flush` を呼ぶようにしたので、ユーザは明示的に `Flush` を呼ぶ必要はなくなった。

v0.11.6
=======
May 8 2023

- 関数 `StringToMoji` と `GetStringWidth` のメモリアローケーション回数を削減した
- 新関数 `MojiCountInString` を追加（  `len(StringToMoji(s))` よりコストが低い）

v0.11.5
=======
May 7 2023

+ TeraTerm で Ubuntu に接続して `make demo` を実行したときに色がおかしい不具合を修正
    + カラーサンプルの `coloring/vimbatch.go` が初期化する際に `ESC[0m` ではなく、`ESC[1;39;40m` を使っていた。
    + リセット用のカラーコード `ESC[0m` に相当する定数値 `readline.ColorReset` を用意

v0.11.4
=======
May 6 2023

- [Windows以外で桁数と行数が取り違えられる問題][issue-one] を修正するため、[go-tty] を [v0.0.5] へ更新しました。

[issue-one]: https://github.com/hymkor/go-multiline-ny/issues/1
[go-tty]: https://github.com/mattn/go-tty
[v0.0.5]: https://github.com/mattn/go-tty/releases/tag/v0.0.5

v0.11.3
=======
May 5 2023

+ 改行やエスケープシーケンス、絵文字が含まれていても出力する幅を正しく計測できるプロンプトフックとして、`Editor` 型に新しいフィールド `PromptWriter` を追加しました。

v0.11.2
=======
May 2 2023

無駄に複雑化している内部構造を簡素化しました。

### v0.11.2 (May 1 2023)

- （使われている場所がなかった）`(*Buffer) Write` メソッドを削除
- （外部からは使いどころがない）`(*Buffer) RefreshColor()` メソッドを非公開化
- サブパッケージの  `tty10` の goroutine リークを修正（デフォルトでは未使用）
- `AnonymousCommand`（匿名のキー割り当てコマンド型）と `SelfInserter` （文字挿入コマンド型）を追加

### v0.11.1 (Apr 28 2023)

+ キーコードを定義するサブパッケージ `keys` を作成
+ 型 `KeyGoFuncT` を `GoCommand` へ改名
+ 型 `KeyFuncT` を `Command` へ改名
+ 関数 `GetKey(ITty)` を隠蔽
+ キー名定数 `K_*****` を削除。かわりにキーコードそのものである `keys.*` を直接ご利用ください
+ コマンド名定数 `F_*****` を削除。かわりに `(Command) String()` メソッドの結果をご利用ください
+ （使われている場所がなかった） `(Result) String()` を削除

v0.11.0
=======
Apr 26 2023

- 以下のテストで動作確認ができており、fork 原因のWindowsTerminal の不具合が修正されているようなので、[go-tty] の独自forkバージョンを削除し、オリジナルの v0.0.4 を使用するようにした。
    - &#x2460; OK: CIRCLE DIGIT ONE: U+2460
    - &#x1F468;&#x200D;&#x1F33E; OK: FARMER: MAN(U+1F468)+ZERO WIDTH JOINER(U+200D)+EAR OF RICE(U+1F33E)
    - &#x908A;&#xE0104; OK: KANJI with VARIATION SELECTOR(U+908A U+E0104)
- [go-tty] のかわりに golang.org/[x/term] を使うための内部スイッチを用意（現行は従来どおり go-tty を使用）
- [go-runewidth] のかわりに golang.org/[x/text/width] を使うための内部スイッチを用意（現行は従来どおり go-runewidth を使用）
- 変数SurrogatePairOkを削除した。関数  EnableSurrogatePair() and IsSurrogatePairEnabled() を使用のこと
- 関数NewDefaultTtyを削除した。 golang.org/[x/term] もしくは [go-tty] を使用のこと

[go-tty]: https://github.com/mattn/go-tty
[go-runewidth]: https://github.com/mattn/go-runewidth
[x/term]: https://pkg.go.dev/golang.org/x/term
[x/text/width]: https://pkg.go.dev/golang.org/x/text/width

v0.10.1
=======
Apr 14 2023

- v0.10.0
    - `Coloring` interface を変更
        - `Init()` と `Next()` は int のかわりに `ColorSequence` (=int64) 型を返すものとした
        - `ESC[0m` で出力属性をリセットできるようになった
- v0.10.1
    - v0.10.0 で色用の定数が壊れていたのを修正

v0.9.1
=======
Apr 10 2023

- `Coloring.Init()` / `.Next()` が 0 を返した時、`m` というゴミ文字列が表示されるのを修正

v0.9.0
=======
Apr 9 2023

- 非公開の型 `_Cell` を `Cell` に改名（公開）。`Cell` のインスタンスは１文字を表すコードポイント群(=`Moji`)と非公開の色情報を含んでいます。
- `Moji` 型の配列を文字列に変換する非公開の関数 `string2moji` を `StringToMoji` に改名（公開）。`Moji` のインスタンスは１文字を表すコードポイント群を保持しています。

**例**
農夫の絵文字（👨‍🌾）は３つのコードポイント (U+0001F468, U+200D , and U+0001F33E)で表されますが、`Moji` は１インスタンスでそれらを格納します。

v0.8.5
=======
Apr 2 2023

- v0.8.4 の非互換性を修正しました。v0.8.4 は nyagos でリンクできなくなってしまっており、削除したメソッド `(*KeyMap)GetBindKey(string)`  を復元しました。


v0.8.4
=======
Mar 25 2023

- キーのマッピングを返す関数の `(*KeyMap) GetBindKey` が nil を返すのは都合が悪いので、 かわりに `(*Editor)GetBindKey` を作成し、`(*KeyMap) GetBindKey` は削除した

v0.8.3
=======
Sep 24 2022

- サンプルの色プラグイン `vimbatch` の文字色を `ESC[37m` (白) から `ESC[39m` (端末の標準文字色) へ変更

v0.8.2
=======
Aug 12 2022

- 最初の文字を表示する前に色をリセット (以前は省略されていた)

v0.8.1
=======
Aug 12 2022

- Ctrl-E がタイプされたとき、カーソル上に
- Fix: On Ctrl-E typed, sometimes non-space character remains on the cursor.

v0.8.0
=======
Apr 29 2022

- WezTerm で Surrogate-pair を有効にした。
- Contour Terminal で Surrogate-pair を有効にした。
- 出力にバックスペース(\b)を使わないようにした。

( 本バージョンは [nyagos 4.4.12\_0] で使用 )

[nyagos 4.4.12\_0]: https://github.com/nyaosorg/nyagos/releases/tag/4.4.12_0

v0.7.0
=======
Feb 26 2022

Windows Terminal での透明色に対応

- `Coloring.Init()` でデフォルト色を戻さなくてはいけなくなった (この interface の互換性は壊れた)
- デフォルト背景色で `ESC[40m` ではなく `ESC[49m` を使うようにした

v0.6.3
=======
Dec 29 2021

- [#2],[#3] フラグフィールド `Editor.HistoryCycling` を追加
- Color: Multi SelectGraphcReition Parameters をサポート(v0.6.2)
- 色設定のために小関数 `SGR1`, `SGR2`, `SGR3`, `SGR4` を実装

[#2]: https://github.com/nyaosorg/go-readline-ny/issues/2
[#3]: https://github.com/nyaosorg/go-readline-ny/issues/2

v0.6.1
=======
Dec 10 2021

カラー対応

![image](https://user-images.githubusercontent.com/3752189/145574044-9056cedb-731f-4ce1-b968-adb91da26432.png)

See [example2.go](https://github.com/nyaosorg/go-readline-ny/blob/master/example2.go) and [coloring/vimbatch.go](https://github.com/nyaosorg/go-readline-ny/blob/master/coloring/vimbatch.go)

v0.5.0
=======
Sep 12 2021

レポジトリオーナーを zetamatta から nyaosorg へ変更

v0.4.14
=======
Aug 27 2021

- Windows10 の、WindowsTerminal でない端末で、罫線キャラクターの幅が不正確になっていた  
 ( [East Asian Ambiguous Character · Issue #412 · zetamatta/nyagos](https://github.com/zetamatta/nyagos/issues/412) )

v0.4.13
=======
May 3 2021

- Windows Terminal で数学向けボールド文字 (U+1D400 - U+1D7FF) をサポートした。

v0.4.12
=======
May 3 2021

- Visual Studio Code のターミナルではサロゲートペアがサポートされていないので強制オフにするようにした。

この問題は VSCode を WindowsTerminal より起動すると表面化します。

v0.4.11
=======
Apr 14 2021

- Emoji Moifier Sequence (skin tone) をサポート: 🏻(U+1F3FB)～ 🏿 (U+1F3FF)

v0.4.10
=======
Apr 14 2021

- ReadLine メソッドを呼ばれる前から押されていて、呼ばれてからキーが離されたコードが入力される不具合を修正した

v0.4.9
======
Apr 14 2021

- RAINBOW FLAG をサポート:(U+1F3F3 U+200D U+1F308 🏳‍🌈)

v0.4.8
======
Apr 14 2021

- WAVING WHITE FLAG と異体字をサポート (U+1F3F3 & U+1F3F3 U+FE0F / 🏳 & 🏳️)
