# go-app-orbit-db-demo
Experimenting with go-app and orbit-db on IPFS

# requirements

https://github.com/maxence-charriere/go-app

https://github.com/berty/go-orbit-db

https://github.com/ipfs/go-ipfs-http-client

# what works
We can connect to the local IPFS daemon via the http client

# what doesn't work

We can't connect to orbit-db from the wasm itself since it requires access to the host file system and wasm doesn't support that(wasi does that).

# github.com/syndtr/goleveldb/leveldb/storage
../../go/pkg/mod/github.com/syndtr/goleveldb@v1.0.0/leveldb/storage/file_storage.go:107:16: undefined: newFileLock
../../go/pkg/mod/github.com/syndtr/goleveldb@v1.0.0/leveldb/storage/file_storage.go:123:13: cannot assign error to err in multiple assignment
../../go/pkg/mod/github.com/syndtr/goleveldb@v1.0.0/leveldb/storage/file_storage.go:127:16: cannot assign error to err in multiple assignment
../../go/pkg/mod/github.com/syndtr/goleveldb@v1.0.0/leveldb/storage/file_storage.go:192:3: undefined: rename
../../go/pkg/mod/github.com/syndtr/goleveldb@v1.0.0/leveldb/storage/file_storage.go:267:12: undefined: rename
../../go/pkg/mod/github.com/syndtr/goleveldb@v1.0.0/leveldb/storage/file_storage.go:272:12: undefined: syncDir
../../go/pkg/mod/github.com/syndtr/goleveldb@v1.0.0/leveldb/storage/file_storage.go:555:9: undefined: rename
../../go/pkg/mod/github.com/syndtr/goleveldb@v1.0.0/leveldb/storage/file_storage.go:591:13: undefined: syncDir

# ideas

It would make sense for IPFS to include orbit-db implementation as a user data store and expose it via the http client. 

WASI could be an alternative not as good as direct IPFS integration though since it requires additional runtime outside of the browser and adds an additional layer to deal with.
