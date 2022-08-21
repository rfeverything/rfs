build:
	GO_CFLAGS="-I/home/yuuin/rocksdb/include" \
CGO_LDFLAGS="-L/home/yuuin/rocksdb -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd" \
go build