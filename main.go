package main

import (
	"context"
	"fmt"
	"net/http"

	orbitdb "berty.tech/go-orbit-db"
	client "github.com/ipfs/go-ipfs-http-client"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// messenger is a component that does a simple messenger on ipfs. A component is a
// customizable, independent, and reusable UI element. It is created by
// embedding app.Compo into a struct.
type messenger struct {
	app.Compo
	message   string
	dbMessage string
}

func (m *messenger) OnMount(app.Context) {
	initIPFS()
}

// The Render method is where the component appearance is defined. Here, a
// "messenger World!" is displayed as a heading.
func (m *messenger) Render() app.UI {
	return app.Div().Class("input-group my-3").Body(
		app.Input().Type("text").Class("form-control").Value(m.message).Placeholder("Store message").OnInput(m.OnStoreMessage),
		app.Input().Type("text").Class("form-control").Value(m.dbMessage).Placeholder("DB message"),
		app.Div().Class("input-group-append").Body(
			app.Button().Class("btn btn-primary").Body(app.Text("Store Message")).OnClick(m.storeMessage),
		),
	)
}

func (m *messenger) OnStoreMessage(ctx app.Context, e app.Event) {
	m.message = e.Get("target").Get("value").String()

	m.Update()
}

func (m *messenger) storeMessage(ctx app.Context, e app.Event) {

}

func main() {
	app.Route("/", &messenger{})
	app.RunWhenOnBrowser()
	initServer()
}

func initIPFS() {
	c, err := client.NewURLApiWithClient("localhost:5001", &http.Client{})
	if err != nil {
		panic(fmt.Errorf("failed to connect to local api: %s", err))
	}

	fmt.Println(c)

	/// ------ Setting up the IPFS Repo

	/// --- Part I: Getting a IPFS node running

	// fmt.Println("-- Getting an IPFS node running -- ")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// // Spawn a node using a temporary path, creating a temporary repo for the run
	// fmt.Println("Spawning node on a temporary repo")
	// ipfs, err := spawnEphemeral(ctx)
	// if err != nil {
	// 	panic(fmt.Errorf("failed to spawn ephemeral node: %s", err))
	// }

	// fmt.Println("IPFS node is running")

	db, err := orbitdb.NewOrbitDB(ctx, c, &orbitdb.NewOrbitDBOptions{})
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
