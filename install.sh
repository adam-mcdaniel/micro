cd ~

go get -d github.com/blang/semver github.com/dustin/go-humanize \
       github.com/flynn/json5 github.com/go-errors/errors \
       github.com/mattn/go-isatty github.com/mitchellh/go-homedir \
       github.com/sergi/go-diff/diffmatchpatch github.com/yuin/gopher-lua \
       github.com/zyedidia/clipboard github.com/zyedidia/glob \
       github.com/zyedidia/tcell github.com/zyedidia/tcell/encoding \
       github.com/zyedidia/terminal layeh.com/gopher-luar \
       github.com/adam-mcdaniel/micro/cmd/micro


if [ -z "$GOPATH" ]; then DIR=~/go; else DIR=$GOPATH; fi

cd $DIR/src/github.com/adam-mcdaniel/micro

make >/dev/null || make

echo "Add \"$DIR/bin\" to your path"