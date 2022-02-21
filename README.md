# go-app-orbit-db-demo
Experimenting with go-app and orbit-db on IPFS

# requirements

https://github.com/maxence-charriere/go-app

https://github.com/berty/go-orbit-db

https://ipfs.io/

# what works

We can spawn an ephemeral node on the server part for the current run and store and retrieve data in db

# what doesn't work

We can't connect to orbit-db from the wasm itself since IPFS needs access to the host file system and wasm doesn't support host bindings yet

# ideas

It would make sense for IPFS to include orbit-db implementation and expose it via the http-client-api. That way we can call the local node via http from the wasm itself and use the local daemon rather than programattically spawning temporary nodes from within the client app.
