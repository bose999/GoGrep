GoGrep
======

# 処理機能
GoGrep はgrepコマンドの 
grep 検索文字列 ファイル名 
という実行機能をGo言語で実装してみたものです。 

# 追加の機能
grep 検索文字列 ファイル名 処理開始行数 
として文字列チェックを開始する行数を指定する事を可能にしています。 
この機能を使った場合は処理終了時に 
GoGrep finished. All Line: 1048577 Check Start Line: 204800 
のように全行数と処理開始行を出力します。

# 内部処理概要
ファイルの先頭からのループ処理をして文字列のチェックは 
ゴルーチンで別処理にして少しだけ並列処理にしています。

# 性能チェック
## 環境
Mac book pro 15 retina 2013(Corei7/Mem16GB/SSD 512GB)

## 検索対象のファイルを生成したシェル
    #!/bin/bash
    for i in `seq 1 1048576`
    do
      pwgen -0A 1024 1 >> /Users/matakeda/Documents/a.txt
    done

## GoGrepとOSのgrepを実行したシェル
    #!/bin/bash
    echo "GoGrep & Grep 実行"
    GG_START=`date +'%s'`
    /Users/matakeda/Documents/git/GoGrep/bin/gogrep "doshuc" /Users/matakeda/Documents/a.txt > ./gogrep.out
    GG_END=`date +'%s'`
    
    G_START=`date +'%s'`
    /usr/bin/grep "doshuc" /Users/matakeda/Documents/a.txt > ./grep.out
    G_END=`date +'%s'`
    
    echo "GoGrep"
    echo `expr $GG_END - $GG_START`
    echo "grep"
    echo `expr $G_END - $G_START`
    echo "diff"
    echo `diff ./gogrep.out ./grep.out`
    echo "finish"
    
## 3回実行した結果
grep 約14秒〜15秒、goGrep 約3秒。OS標準のgrepより早く処理しています。 
    % ./grep-time.sh
    GoGrep & Grep 実行
    GoGrep
    3
    grep
    14
    diff
    
    finish
    % ./grep-time.sh
    GoGrep & Grep 実行
    GoGrep
    3
    grep
    15
    diff
    
    finish
    % ./grep-time.sh
    GoGrep & Grep 実行 
    GoGrep
    3
    grep
    15
    diff
    
    finish