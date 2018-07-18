//go:generate go run ../../../../TIBCOSoftware/flogo-lib/flogo/gen/gen.go $GOPATH
package main

import (
	"context"
	"os"
	"strconv"

	"github.com/TIBCOSoftware/flogo-contrib/trigger/rest"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/TIBCOSoftware/flogo-lib/flogo"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/retgits/flogo-components/activity/addtodate"
	"github.com/retgits/flogo-components/activity/randomnumber"
)

var (
	httpport = os.Getenv("HTTPPORT")
)

func main() {
	// Create a new Flogo app
	app := flogo.NewApp()

	// Convert the HTTPPort to an integer
	port, err := strconv.Atoi(httpport)

	// Register the HTTP trigger
	trg := app.NewTrigger(&rest.RestTrigger{}, map[string]interface{}{"port": port})
	trg.NewFuncHandler(map[string]interface{}{"method": "GET", "path": "/api/expected-date/:invoiceId"}, handler)

	e, err := flogo.NewEngine(app)

	if err != nil {
		logger.Error(err)
		return
	}

	engine.RunEngine(e)
}

func handler(ctx context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error) {
	// Get the ID from the path
	id := inputs["pathParams"].Value().(map[string]string)["invoiceId"]

	// Generate a random number for the amount
	// There are definitely better ways to do this with Go, but this keeps the flow consistent with the UI version
	in := map[string]interface{}{"min": 0, "max": 10}
	out, err := flogo.EvalActivity(&randomnumber.MyActivity{}, in)
	if err != nil {
		return nil, err
	}
	datediff := out["result"].Value().(int)

	// Generate a new date
	// There are definitely better ways to do this with Go, but this keeps the flow consistent with the UI version
	in = map[string]interface{}{"number": datediff, "units": "days"}
	out, err = flogo.EvalActivity(&addtodate.MyActivity{}, in)
	if err != nil {
		return nil, err
	}
	expectedPaymentDate := out["result"].Value().(string)

	// The return message is a map[string]*data.Attribute which we'll have to construct
	response := make(map[string]*data.Attribute)

	attr, _ := data.NewAttribute("id", data.TypeString, id)
	response["id"] = attr
	attr, _ = data.NewAttribute("expectedDate", data.TypeString, expectedPaymentDate)
	response["expectedDate"] = attr

	return response, nil
}