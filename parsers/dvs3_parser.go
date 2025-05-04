package parsers

import (
	"errors"
	"fmt"
	"io"

	"github.com/beevik/etree"
	"github.com/rijdendetreinen/gotrain/models"
)

func ParseDvs3Message(reader io.Reader) (departure models.Departure, err error) {
	doc := etree.NewDocument()

	if _, err := doc.ReadFrom(reader); err != nil {
		return departure, err
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parser error: %+v", r)
			// log.Error().Err(err).Msg("Recovered from panic in ParseDvs3Message")
		}
	}()

	dvs3_product := doc.SelectElement("PutReisInformatieBoodschapIn").SelectElement("reisInformatieProductDVS")

	if dvs3_product != nil {
		if dvs3_product.NamespaceURI() != "urn:ns:cdm:reisinformatie:data:dvs:3" {
			dvs3_product = nil
		} else {
			return parseDvs3Product(dvs3_product)
		}
	}

	err = errors.New("missing DVS3 element")
	// log.Error().Err(err).Msg("Failed to find DVS3 element in XML")

	return
}

// parseDvs3Product parses the DVS3 product element and returns a Departure object.
func parseDvs3Product(product *etree.Element) (departure models.Departure, err error) {
	departure.DvsVersion = 3

	productAdministration := product.SelectElement("ripAdministratie")
	infoProduct := product.SelectElement("dynamischeVertrekStaat")
	trainProduct := infoProduct.SelectElement("trein")

	departure.Timestamp = ParseIsoTime(product.SelectAttrValue("timestamp", ""))
	departure.ProductID = productAdministration.SelectElement("reisInformatieProductID").Text()

	departure.ServiceID = infoProduct.SelectElement("ritNummer").Text()
	departure.ServiceDate = infoProduct.SelectElement("ritDatum").Text()
	departure.Station = ParseInfoPlusStation2024(infoProduct.SelectElement("vertrekStation"))
	departure.GenerateID()

	departure.ServiceNumber = trainProduct.SelectElement("nummer").Text()
	departure.ServiceType = trainProduct.SelectElement("soort").SelectElement("presentatieTekstPerTaal").SelectElement("tekst").Text()
	departure.ServiceTypeCode = trainProduct.SelectElement("soort").SelectElement("code").Text()
	departure.Company = trainProduct.SelectElement("vervoerder").Text()

	departureStatus := trainProduct.SelectElement("status").Text()

	switch departureStatus {
	case "ONBEKEND":
		departure.Status = models.DepartureStatusUnknown
	case "NADERT":
		departure.Status = models.DepartureStatusApproaching
	case "BINNENKOMST":
		departure.Status = models.DepartureStatusArriving
	case "VERTROKKEN":
		departure.Status = models.DepartureStatusDeparted

	}

	lineNumberNode := trainProduct.SelectElement("lijnNummer")
	if lineNumberNode != nil {
		departure.LineNumber = lineNumberNode.Text()
	}

	// Train name, e.g. special trains like the museum train
	nameNode := trainProduct.SelectElement("naam")
	if nameNode != nil {
		departure.ServiceName = nameNode.Text()
	}

	// Departure time:
	departureTimes := trainProduct.SelectElements("vertrekTijd")

	for _, departureTime := range departureTimes {
		if departureTime.SelectAttrValue("infoStatus", "") == "GEPLAND" {
			departure.DepartureTime = ParseIsoTime(departureTime.Text())
		}
	}

	delayNode := trainProduct.SelectElement("vertraging")
	if delayNode != nil {
		departure.Delay = ParseInfoPlusDuration(delayNode.SelectElement("exact"))
	}
	return
}
