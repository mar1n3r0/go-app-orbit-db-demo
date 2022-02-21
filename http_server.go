//go:build !wasm

// file: http_server.go

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	orbitdb "berty.tech/go-orbit-db"
	"berty.tech/go-orbit-db/iface"
	"github.com/NYTimes/gziphandler"
	"github.com/maxence-charriere/go-app/v9/pkg/app"

	config "github.com/ipfs/go-ipfs-config"
	icore "github.com/ipfs/interface-go-ipfs-core"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader" // This package is needed so that all the preloaded plugins are loaded automatically
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

var KvStore iface.KeyValueStore

func initServer() {
	withGz := gziphandler.GzipHandler(&app.Handler{
		Name:        "messenger",
		Description: "A orbit-db example with ipfs",
		Styles:      []string{"https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css"},
	})
	http.Handle("/", withGz)

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}

func setupPlugins(externalPluginsPath string) error {
	// Load any external plugins if available on externalPluginsPath
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	// Load preloaded and external plugins
	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}

func createTempRepo() (string, error) {
	repoPath, err := ioutil.TempDir("", "ipfs-shell")
	if err != nil {
		return "", fmt.Errorf("failed to get temp dir: %s", err)
	}

	// Create a config with default options and a 2048 bit key
	cfg, err := config.Init(ioutil.Discard, 2048)
	if err != nil {
		return "", err
	}

	// Create the repo with the config
	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to init ephemeral node: %s", err)
	}

	return repoPath, nil
}

/// ------ Spawning the node

// Creates an IPFS node and returns its coreAPI
func createNode(ctx context.Context, repoPath string) (icore.CoreAPI, error) {
	// Open the repo
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	// Construct the node

	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
		// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		Repo: repo,
		ExtraOpts: map[string]bool{
			"pubsub": true,
		},
	}

	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, err
	}

	// Attach the Core API to the constructed node
	return coreapi.NewCoreAPI(node)
}

// Spawns a node to be used just for this run (i.e. creates a tmp repo)
func spawnEphemeral(ctx context.Context) (icore.CoreAPI, error) {
	if err := setupPlugins(""); err != nil {
		return nil, err
	}

	// Create a Temporary Repo
	repoPath, err := createTempRepo()
	if err != nil {
		return nil, fmt.Errorf("failed to create temp repo: %s", err)
	}

	// Spawning an ephemeral IPFS node
	return createNode(ctx, repoPath)
}

func initIPFS() {
	/// ------ Setting up the IPFS Repo

	/// --- Part I: Getting a IPFS node running

	fmt.Println("-- Getting an IPFS node running -- ")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Spawn a node using a temporary path, creating a temporary repo for the run
	fmt.Println("Spawning node on a temporary repo")
	ipfs, err := spawnEphemeral(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to spawn ephemeral node: %s", err))
	}

	fmt.Println("IPFS node is running")

	db, err := orbitdb.NewOrbitDB(ctx, ipfs, &orbitdb.NewOrbitDBOptions{})
	if err != nil {
		panic(fmt.Errorf("failed to create new orbitdb: %s", err))
	}

	dbStore, err := db.Create(ctx, "test", "keyvalue", &orbitdb.CreateDBOptions{})
	if err != nil {
		panic(fmt.Errorf("failed to create new db store: %s", err))
	}

	fmt.Printf("dbStore address: %s\n", dbStore.Address())

	KvStore, err = db.KeyValue(ctx, dbStore.Address().String(), &orbitdb.CreateDBOptions{})
	if err != nil {
		panic(fmt.Errorf("failed to get kv store: %s", err))
	}

	kvStorePut(ctx, "message", []byte("test1"))

	v := kvStoreGet(ctx, "message")

	fmt.Printf("kv store value: %s\n", v)
}

func kvStorePut(ctx context.Context, key string, val []byte) {
	_, err := KvStore.Put(ctx, "message", []byte("test1"))
	if err != nil {
		panic(fmt.Errorf("failed to put in kv store: %s", err))
	}
}

func kvStoreGet(ctx context.Context, key string) []byte {
	v, err := KvStore.Get(ctx, "message")
	if err != nil {
		panic(fmt.Errorf("failed to get value from kv store: %s", err))
	}

	return v
}
