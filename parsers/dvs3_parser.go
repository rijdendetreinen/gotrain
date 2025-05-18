package parsers

import (
	"errors"
	"fmt"
	"io"
	"strings"

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

	// Departure time and delay
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

	// Departure platform
	departurePlatformNode := trainProduct.SelectElement("vertrekSpoor")
	if departurePlatformNode != nil {
		departurePlatformTrackNode := departurePlatformNode.SelectElement("spoor")

		trackNumberNode := departurePlatformTrackNode.SelectElement("nummer")
		trackPhaseNode := departurePlatformTrackNode.SelectElement("fase")

		if trackNumberNode == nil || trackPhaseNode == nil {
			departure.PlatformActual = ""
		} else if trackPhaseNode == nil {
			departure.PlatformActual = trackNumberNode.Text()
		} else {
			departure.PlatformActual = trackNumberNode.Text() + strings.ToLower(trackPhaseNode.Text())
		}
	}

	// Parse modifications
	// modifications := trainProduct.SelectElements("wijziging")
	// for _, modification := range modifications {
	// 	modificationType, _ := strconv.Atoi(modification.SelectElement("code").Text())

	// 	departure.Modifications = append(departure.Modifications, models.Modification{
	// 		ModificationType: modificationType,
	// 	})
	// }

	// Parse wings
	wings := trainProduct.SelectElements("vleugel")

	for _, wing := range wings {
		wing := parseDvs3TrainWing(wing, &departure)

		departure.TrainWings = append(departure.TrainWings, wing)

		departure.DestinationActual = append(departure.DestinationActual, wing.DestinationActual...)
		departure.DestinationPlanned = append(departure.DestinationPlanned, wing.DestinationPlanned...)
	}

	return
}

// parseDvs3TrainWing parses a DVS3 train wing element and returns a TrainWing object.
func parseDvs3TrainWing(wing *etree.Element, departure *models.Departure) models.TrainWing {
	wingDeparture := models.TrainWing{}

	// Destinations:
	destinations := wing.SelectElements("bestemming")
	for _, destination := range destinations {
		destinationStation := ParseInfoPlusStation2024(destination)
		wingDeparture.DestinationActual = append(wingDeparture.DestinationActual, destinationStation)
	}

	return wingDeparture
}
