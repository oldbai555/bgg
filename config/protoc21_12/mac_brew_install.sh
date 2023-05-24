brew unlink protobuf
brew install protobuf.rb
cd $GOPATH/bin
rm -rf protoc
brew link protobuf
ln -s $(brew list protobuf |grep '/bin/protoc' |head -n 1) protoc
protoc --version