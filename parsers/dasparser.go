package parsers

import (
	"fmt"
	"io"
	"strconv"

	"github.com/beevik/etree"
	"github.com/rijdendetreinen/gotrain/models"
)

// ParseDasMessage parses a DAS XML message to an Arrival object
func ParseDasMessage(reader io.Reader) (arrival models.Arrival, err error) {
	doc := etree.NewDocument()

	if _, err := doc.ReadFrom(reader); err != nil {
		return arrival, err
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Parser error: %+v", r)
		}
	}()

	product := doc.SelectElement("PutReisInformatieBoodschapIn").SelectElement("ReisInformatieProductDAS")
	productAdministration := product.SelectElement("RIPAdministratie")
	infoProduct := product.SelectElement("DynamischeAankomstStaat")
	trainProduct := infoProduct.SelectElement("TreinAankomst")

	arrival.Timestamp = ParseInfoPlusDateTime(productAdministration.SelectElement("ReisInformatieTijdstip"))
	arrival.ProductID = productAdministration.SelectElement("ReisInformatieProductID").Text()

	arrival.ServiceID = infoProduct.SelectElement("RitId").Text()
	arrival.ServiceDate = infoProduct.SelectElement("RitDatum").Text()
	arrival.Station = ParseInfoPlusStation(infoProduct.SelectElement("RitStation"))
	arrival.GenerateID()

	arrival.ServiceNumber = trainProduct.SelectElement("TreinNummer").Text()
	arrival.ServiceType = trainProduct.SelectElement("TreinSoort").Text()
	arrival.ServiceTypeCode = trainProduct.SelectElement("TreinSoort").SelectAttrValue("Code", "")
	arrival.Company = trainProduct.SelectElement("Vervoerder").Text()
	arrival.Status, _ = strconv.Atoi(trainProduct.SelectElement("TreinStatus").Text())

	lineNumberNode := trainProduct.SelectElement("LijnNummer")
	if lineNumberNode != nil {
		arrival.LineNumber = lineNumberNode.Text()
	}

	// Train name, e.g. special trains like the museum train
	nameNode := trainProduct.SelectElement("TreinNaam")
	if nameNode != nil {
		arrival.ServiceName = nameNode.Text()
	}

	arrival.ArrivalTime = ParseInfoPlusDateTime(ParseWhenAttribute(trainProduct, "AankomstTijd", "InfoStatus", "Gepland"))
	arrival.Delay = ParseInfoPlusDuration(trainProduct.SelectElement("ExacteAankomstVertraging"))

	arrival.OriginActual = ParseInfoPlusStations(ParseWhenAttributeMulti(trainProduct, "TreinHerkomst", "InfoStatus", "Actueel"))
	arrival.OriginPlanned = ParseInfoPlusStations(ParseWhenAttributeMulti(trainProduct, "TreinHerkomst", "InfoStatus", "Gepland"))

	// Workaround for DAS bug:
	if arrival.OriginActual[0].Code == arrival.Station.Code {
		arrival.OriginActual = arrival.OriginPlanned
	}

	arrival.PlatformActual = ParseInfoPlusPlatform(ParseWhenAttributeMulti(trainProduct, "TreinAankomstSpoor", "InfoStatus", "Actueel"))
	arrival.PlatformPlanned = ParseInfoPlusPlatform(ParseWhenAttributeMulti(trainProduct, "TreinAankomstSpoor", "InfoStatus", "Gepland"))

	viaNodeActual := ParseWhenAttribute(trainProduct, "VerkorteRouteHerkomst", "InfoStatus", "Actueel")
	viaNodePlanned := ParseWhenAttribute(trainProduct, "VerkorteRouteHerkomst", "InfoStatus", "Gepland")

	if viaNodeActual != nil {
		arrival.ViaActual = ParseInfoPlusStations(viaNodeActual.SelectElements("Station"))
	}
	if viaNodePlanned != nil {
		arrival.ViaPlanned = ParseInfoPlusStations(viaNodePlanned.SelectElements("Station"))
	}

	arrival.Modifications = ParseInfoPlusModificationsByElement(trainProduct, "WijzigingHerkomst")

	// Check for flags that may be set:
	for _, modification := range arrival.Modifications {
		switch modification.ModificationType {
		case models.ModificationCancelledArrival:
			arrival.Cancelled = true
		case models.ModificationNotActual:
			arrival.NotRealTime = true
		}
	}

	return
}
