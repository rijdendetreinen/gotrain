package parsers

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/beevik/etree"
	"github.com/rijdendetreinen/gotrain/models"
)

// ParseDvsMessage parses a DVS XML message to a Departure object
func ParseDvsMessage(reader io.Reader) (departure models.Departure, err error) {
	doc := etree.NewDocument()

	if _, err := doc.ReadFrom(reader); err != nil {
		return departure, err
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Parser error: %+v", r)
		}
	}()

	product := doc.SelectElement("PutReisInformatieBoodschapIn").SelectElement("ReisInformatieProductDVS")

	if product == nil {
		err = errors.New("Missing DVS element")
		return
	}

	productAdministration := product.SelectElement("RIPAdministratie")
	infoProduct := product.SelectElement("DynamischeVertrekStaat")
	trainProduct := infoProduct.SelectElement("Trein")

	departure.Timestamp = ParseInfoPlusDateTime(productAdministration.SelectElement("ReisInformatieTijdstip"))
	departure.ProductID = productAdministration.SelectElement("ReisInformatieProductID").Text()

	departure.ServiceID = infoProduct.SelectElement("RitId").Text()
	departure.ServiceDate = infoProduct.SelectElement("RitDatum").Text()
	departure.Station = ParseInfoPlusStation(infoProduct.SelectElement("RitStation"))
	departure.ID = departure.ServiceDate + "-" + departure.ServiceID + "-" + departure.Station.Code

	departure.ServiceNumber = trainProduct.SelectElement("TreinNummer").Text()
	departure.ServiceType = trainProduct.SelectElement("TreinSoort").Text()
	departure.ServiceTypeCode = trainProduct.SelectElement("TreinSoort").SelectAttrValue("Code", "")
	departure.Company = trainProduct.SelectElement("Vervoerder").Text()
	departure.Status, _ = strconv.Atoi(trainProduct.SelectElement("TreinStatus").Text())

	// Train name, e.g. special trains like the museum train
	nameNode := trainProduct.SelectElement("TreinNaam")
	if nameNode != nil {
		departure.ServiceName = nameNode.Text()
	}

	departure.DepartureTime = ParseInfoPlusDateTime(ParseWhenAttribute(trainProduct, "VertrekTijd", "InfoStatus", "Gepland"))
	departure.Delay = ParseInfoPlusDuration(trainProduct.SelectElement("ExacteVertrekVertraging"))

	departure.DestinationActual = ParseInfoPlusStations(ParseWhenAttributeMulti(trainProduct, "TreinEindBestemming", "InfoStatus", "Actueel"))
	departure.DestinationPlanned = ParseInfoPlusStations(ParseWhenAttributeMulti(trainProduct, "TreinEindBestemming", "InfoStatus", "Gepland"))

	departure.PlatformActual = ParseInfoPlusPlatform(ParseWhenAttributeMulti(trainProduct, "VertrekSpoor", "InfoStatus", "Actueel"))
	departure.PlatformPlanned = ParseInfoPlusPlatform(ParseWhenAttributeMulti(trainProduct, "VertrekSpoor", "InfoStatus", "Gepland"))

	departure.ReservationRequired = ParseInfoPlusBoolean(trainProduct.SelectElement("Reserveren"))
	departure.WithSupplement = ParseInfoPlusBoolean(trainProduct.SelectElement("Toeslag"))
	departure.SpecialTicket = ParseInfoPlusBoolean(trainProduct.SelectElement("SpeciaalKaartje"))
	departure.RearPartRemains = ParseInfoPlusBoolean(trainProduct.SelectElement("AchterBlijvenAchtersteTreinDeel"))
	departure.DoNotBoard = ParseInfoPlusBoolean(trainProduct.SelectElement("NietInstappen"))

	viaNodeActual := ParseWhenAttribute(trainProduct, "VerkorteRoute", "InfoStatus", "Actueel")
	viaNodePlanned := ParseWhenAttribute(trainProduct, "VerkorteRoute", "InfoStatus", "Actueel")

	if viaNodeActual != nil {
		departure.ViaActual = ParseInfoPlusStations(viaNodeActual.SelectElements("Station"))
	}
	if viaNodePlanned != nil {
		departure.ViaPlanned = ParseInfoPlusStations(viaNodePlanned.SelectElements("Station"))
	}

	// TODO:
	// Tips: InstapTip, ReisTip, OverstapTip
	// Wings

	departure.Modifications = ParseInfoPlusModifications(trainProduct)

	return
}
