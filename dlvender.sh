pkgs="github.com/huandu/xstrings github.com/lytics/base62 github.com/pkg/errors github.com/emirpasic/gods github.com/ianlancetaylor/demangle github.com/PuerkitoBio/goquery github.com/andybalholm/cascadia github.com/Workiva/go-datastructures github.com/bitly/go-simplejson github.com/cheekybits/genny github.com/kitech/colog"
# github.com/go-clang/v3.9

commits="gods=79df803e554cab06f835fb729064282b2fd4cc99 goquery=70a02e53e345395ce667b8749b04dbb89098186e go-funk=d866354f4b2438943bec6be044dfa312d3720020 demangle=28f6c0f3b63983aaa99575ca3b693afff7996387"

# go version: go1.11.x
# https://dl.google.com/go/go1.11.13.linux-amd64.tar.gz
# maybe 202006

gopath=$GOPATH
depsrc=$GOPATH/src

oldpwd=$PWD
set -x
for pkg in $pkgs ; do
    echo "\nprocess $pkg ..."
    cd $depsrc
    dir=$(echo $pkg|awk -F'/' '{print $1 "/" $2}')
    name=$(basename $pkg)
    echo "$dir -- $name"
    cd $depsrc
    mkdir -pv $dir
    cd $depsrc/$dir
    if [ ! -d "$name" ]; then
        pxyrun.sh 12334 git clone "https://${pkg}.git"
    fi
    cd $oldpwd
done


mkdir -pv $depsrc/golang.org/x
cd $depsrc/golang.org/x
if [ ! -d "sys" ]; then
    pxyrun.sh.sh 12334 git clone https://go.googlesource.com/sys.git
    # 2d18734c6014bdfc4d7fd7cfc8e0608c9684967b
fi
if [ ! -d "net" ]; then
    pxyrun.sh.sh 12334 git clone https://go.googlesource.com/net.git
    # 6772e930b67bb09bf22262c7378e7d2f67cf59d1
fi

cd $oldpwd
cd $depsrc

# github.com/go-clang/v3.9/clang
mkdir -pv github.com/go-clang
cd github.com/go-clang
rm -rf v3.9
# cp -a ~/aprog/go-clang-v3.9 v3.9
# fork and modified, cannot use upstream one
git clone https://github.com/kitech/go-clang-v3.9.git v3.9
cd $depsrc
if [ ! -d "gopp" ]; then
    pxyrun.sh 12334 git clone https://github.com/kitech/goplusplus.git gopp
fi
cd $oldpwd
