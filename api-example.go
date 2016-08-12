package main

import (
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
)

type Bill struct {
	Id, Name string
}

type BillService struct {
	bills map[string]Bill
}

func (u BillService) Register() {
	ws := new(restful.WebService)
	ws.
		Path("/bills").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) 

	ws.Route(ws.GET("/").To(u.findAllBills).
		// docs
		Doc("get all bills").
		Operation("findAllBills").
		Returns(200, "OK", []Bill{}))

	ws.Route(ws.GET("/{bill-id}").To(u.findBill).
		// docs
		Doc("get a bill").
		Operation("findBill").
		Param(ws.PathParameter("bill-id", "identifier of the bill").DataType("string")).
		Writes(Bill{})) // on the response

	ws.Route(ws.PUT("/{bill-id}").To(u.updateBill).
		Doc("update a bill").
		Operation("updateBill").
		Param(ws.PathParameter("bill-id", "identifier of the bill").DataType("string")).
		Reads(Bill{})) 

	ws.Route(ws.PUT("").To(u.createBill).
		Doc("create a bill").
		Operation("createBill").
		Reads(Bill{})) 

	ws.Route(ws.DELETE("/{bill-id}").To(u.removeBill).
		Doc("delete a bill").
		Operation("removeBill").
		Param(ws.PathParameter("bill-id", "identifier of the bill").DataType("string")))

	restful.Add(ws)
}

// GET http://localhost:8080/bills
//
func (u BillService) findAllBills(request *restful.Request, response *restful.Response) {
	response.WriteEntity(u.bills)
}

// GET http://localhost:8080/bills/1
//
func (u BillService) findBill(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("bill-id")
	usr := u.bills[id]
	if len(usr.Id) == 0 {
		response.WriteErrorString(http.StatusNotFound, "Bill could not be found.")
	} else {
		response.WriteEntity(usr)
	}
}

// PUT http://localhost:8080/bills/1
// <Bill><Id>1</Id><Name>Melissa Raspberry</Name></Bill>
//
func (u *BillService) updateBill(request *restful.Request, response *restful.Response) {
	usr := new(Bill)
	err := request.ReadEntity(&usr)
	if err == nil {
		u.bills[usr.Id] = *usr
		response.WriteEntity(usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// PUT http://localhost:8080/bills/1
// <Bill><Id>1</Id><Name>Melissa</Name></Bill>
//
func (u *BillService) createBill(request *restful.Request, response *restful.Response) {
	usr := Bill{Id: request.PathParameter("bill-id")}
	err := request.ReadEntity(&usr)
	if err == nil {
		u.bills[usr.Id] = usr
		response.WriteHeaderAndEntity(http.StatusCreated, usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// DELETE http://localhost:8080/bills/1
//
func (u *BillService) removeBill(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("bill-id")
	delete(u.bills, id)
}

func main() {
	u := BillService{map[string]Bill{}}
	u.Register()

	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open http://localhost:8080/apidocs and enter http://localhost:8080/apidocs.json in the api input field.
	config := swagger.Config{
		WebServices:    restful.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: "http://localhost:8080",
		ApiPath:        "/apidocs.json",

		// Optionally, specifiy where the UI is located
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "/Bills/emicklei/Projects/swagger-ui/dist"}
	swagger.InstallSwaggerService(config)

	log.Printf("start listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
