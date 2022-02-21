package main

import (
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
	initIPFS()
	app.Route("/", &messenger{})
	app.RunWhenOnBrowser()
	initServer()
}
