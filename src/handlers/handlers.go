package handlers

import (
	"encoding/xml"
	"fmt"

	"github.com/gin-gonic/gin"
)

func IndexHandler(ctx *gin.Context) {
	name := ctx.Params.ByName("name")

	ctx.JSON(200, gin.H{
		"message": fmt.Sprintf("hello %v", name),
	})
}

type Person struct {
	XMLName   xml.Name `xml:"person"`
	FirstName string   `xml:"firstName,attr"`
	LastName  string   `xml:"lastName,attr"`
}

func PersonHandler(ctx *gin.Context) {
	ctx.XML(200, Person{
		FirstName: "Tester",
		LastName:  "Testing",
	})
}
