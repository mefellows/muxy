#!/bin/sh
#go test -v --race $(find . -maxdepth 10 -not -path './.git*' -not -path '*/_*' -type d | grep -v examples | grep -v pkg | grep -v bin | grep -v scripts  | xargs -I{} echo "{}/")

# Get Test dependencies
go get github.com/axw/gocov/gocov
go get github.com/mattn/goveralls
go get golang.org/x/tools/cmd/cover
go get github.com/modocache/gover

# Run test coverage on each subdirectories and merge the coverage profile.
echo "mode: count" > profile.cov

# Standard go tooling behavior is to ignore dirs with leading underscors
for dir in $(find . -maxdepth 10 -not -path './.git*' -not -path '*/_*' -type d | grep -v examples); do
  if ls $dir/*.go &> /dev/null; then
    go test -covermode=count -coverprofile=$dir/profile.tmp $dir
    if [ -f $dir/profile.tmp ]; then
    	cat $dir/profile.tmp | tail -n +2 >> profile.cov
    	rm $dir/profile.tmp
    fi
  fi
done

go tool cover -func profile.cov