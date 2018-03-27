package Model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type LCRData struct {
	CarrierID bson.ObjectId `json:"carrierid" bson:"_carrierid"`
	Active    bool          `json:"active" bson:"active"`
	Number    Number        `json:"number" bson:"number"`
}
type Carrier struct {
	Name        string    `json:"name" bson:"name"`
	DispatcerID int       `json:"dispatcer" bson:"dis[atcerid"`
	Active      bool      `json:"active" bson:"active"`
	CreateDate  time.Time `json:"createdate" bson:"createdate"`
	UpdateDate  time.Time `json:"updatedate" bson:"updatedate"`
}

type LCRDataa struct {
	ID        bson.ObjectId `json:"ID" bson:"_id"`
	CarrierID bson.ObjectId `json:"carrierid " bson:"_carrierid"`
	Active    bool          `json:"active" bson:"active"`
	Number    Number        `json:"number" bson:"number"`
}

type Number struct {
	Value int  `json:"value" bson:"value"`
	Rates Rate `json:"rates" bson:"rates"`
}

type Rate struct {
	Interstate    float64 `json:"interstate" bson:"interstate"`
	Intrastate    float64 `json:"intrastate" bson:"intrastate"`
	Indeterminate float64 `json:"indeterminate" bson:"indeterminate"`
}
