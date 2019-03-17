package parsers

import (
	"io"

	"github.com/beevik/etree"
	"github.com/rijdendetreinen/gotrain/models"
)

// ParseDasMessage parses a DVS XML message to an Arrival object
func ParseDasMessage(reader io.Reader) models.Arrival {
	doc := etree.NewDocument()

	if _, err := doc.ReadFrom(reader); err != nil {
		panic(err)
	}

	product := doc.SelectElement("PutReisInformatieBoodschapIn").SelectElement("ReisInformatieProductDAS")
	productAdministration := product.SelectElement("RIPAdministratie")
	infoProduct := product.SelectElement("DynamischeAankomstStaat")
	trainProduct := infoProduct.SelectElement("TreinAankomst")

	var arrival models.Arrival

	arrival.Timestamp = ParseInfoPlusDateTime(productAdministration.SelectElement("ReisInformatieTijdstip"))
	arrival.ProductID = productAdministration.SelectElement("ReisInformatieProductID").Text()

	arrival.ServiceID = infoProduct.SelectElement("RitId").Text()
	arrival.ServiceDate = infoProduct.SelectElement("RitDatum").Text()
	arrival.Station = ParseInfoPlusStation(infoProduct.SelectElement("RitStation"))
	arrival.ID = arrival.ServiceDate + "-" + arrival.ServiceID + "-" + arrival.Station.Code

	arrival.ServiceNumber = trainProduct.SelectElement("TreinNummer").Text()
	arrival.ServiceType = trainProduct.SelectElement("TreinSoort").Text()
	arrival.ServiceTypeCode = trainProduct.SelectElement("TreinSoort").SelectAttrValue("Code", "")
	arrival.Company = trainProduct.SelectElement("Vervoerder").Text()

	arrival.ArrivalTime = ParseInfoPlusDateTime(ParseWhenAttribute(trainProduct, "AankomstTijd", "InfoStatus", "Gepland"))

	// TODO: Parse other fields

	return arrival
}
