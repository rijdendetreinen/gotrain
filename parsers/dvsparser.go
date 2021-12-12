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
	departure.GenerateID()

	departure.ServiceNumber = trainProduct.SelectElement("TreinNummer").Text()
	departure.ServiceType = trainProduct.SelectElement("TreinSoort").Text()
	departure.ServiceTypeCode = trainProduct.SelectElement("TreinSoort").SelectAttrValue("Code", "")
	departure.Company = trainProduct.SelectElement("Vervoerder").Text()
	departure.Status, _ = strconv.Atoi(trainProduct.SelectElement("TreinStatus").Text())

	lineNumberNode := trainProduct.SelectElement("LijnNummer")
	if lineNumberNode != nil {
		departure.LineNumber = lineNumberNode.Text()
	}

	// Train name, e.g. special trains like the museum train
	nameNode := trainProduct.SelectElement("TreinNaam")
	if nameNode != nil {
		departure.ServiceName = nameNode.Text()
	}

	departure.DepartureTime = ParseInfoPlusDateTime(ParseWhenAttribute(trainProduct, "VertrekTijd", "InfoStatus", "Gepland"))
	departure.Delay = ParseInfoPlusDuration(trainProduct.SelectElement("ExacteVertrekVertraging"))

	departure.DestinationActual = ParseInfoPlusStations(ParseWhenAttributeMulti(trainProduct, "TreinEindBestemming", "InfoStatus", "Actueel"))
	departure.DestinationPlanned = ParseInfoPlusStations(ParseWhenAttributeMulti(trainProduct, "TreinEindBestemming", "InfoStatus", "Gepland"))

	departure.PlatformActual = ParseInfoPlusPlatform(ParseWhenAttributeMulti(trainProduct, "TreinVertrekSpoor", "InfoStatus", "Actueel"))
	departure.PlatformPlanned = ParseInfoPlusPlatform(ParseWhenAttributeMulti(trainProduct, "TreinVertrekSpoor", "InfoStatus", "Gepland"))

	departure.ReservationRequired = ParseInfoPlusBoolean(trainProduct.SelectElement("Reserveren"))
	departure.WithSupplement = ParseInfoPlusBoolean(trainProduct.SelectElement("Toeslag"))
	departure.SpecialTicket = ParseInfoPlusBoolean(trainProduct.SelectElement("SpeciaalKaartje"))
	departure.RearPartRemains = ParseInfoPlusBoolean(trainProduct.SelectElement("AchterBlijvenAchtersteTreinDeel"))
	departure.DoNotBoard = ParseInfoPlusBoolean(trainProduct.SelectElement("NietInstappen"))

	viaNodeActual := ParseWhenAttribute(trainProduct, "VerkorteRoute", "InfoStatus", "Actueel")
	viaNodePlanned := ParseWhenAttribute(trainProduct, "VerkorteRoute", "InfoStatus", "Gepland")

	if viaNodeActual != nil {
		departure.ViaActual = ParseInfoPlusStations(viaNodeActual.SelectElements("Station"))
	}
	if viaNodePlanned != nil {
		departure.ViaPlanned = ParseInfoPlusStations(viaNodePlanned.SelectElements("Station"))
	}

	departure.Modifications = ParseInfoPlusModifications(trainProduct)

	boardingTipNodes := trainProduct.SelectElements("InstapTip")
	for _, boardingTipNode := range boardingTipNodes {
		var boardingTip models.BoardingTip

		boardingTip.ExitStation = ParseInfoPlusStation(boardingTipNode.SelectElement("InstapTipUitstapStation"))
		boardingTip.Destination = ParseInfoPlusStation(boardingTipNode.SelectElement("InstapTipTreinEindBestemming"))

		boardingTip.TrainType = boardingTipNode.SelectElement("InstapTipTreinSoort").Text()
		boardingTip.TrainTypeCode = boardingTipNode.SelectElement("InstapTipTreinSoort").SelectAttrValue("Code", "")

		boardingTip.DeparturePlatform = ParseInfoPlusPlatform(boardingTipNode.SelectElements("InstapTipVertrekSpoor"))
		boardingTip.DepartureTime = ParseInfoPlusDateTime(boardingTipNode.SelectElement("InstapTipVertrekTijd"))

		departure.BoardingTips = append(departure.BoardingTips, boardingTip)
	}

	travelTipNodes := trainProduct.SelectElements("ReisTip")
	for _, travelTipNode := range travelTipNodes {
		var travelTip models.TravelTip

		travelTip.TipCode = travelTipNode.SelectElement("ReisTipCode").Text()
		travelTip.Stations = ParseInfoPlusStations(travelTipNode.SelectElements("ReisTipStation"))

		departure.TravelTips = append(departure.TravelTips, travelTip)
	}

	changeTipNodes := trainProduct.SelectElements("OverstapTip")
	for _, changeTipNode := range changeTipNodes {
		var changeTip models.ChangeTip

		changeTip.ChangeStation = ParseInfoPlusStation(changeTipNode.SelectElement("OverstapTipOverstapStation"))
		changeTip.Destination = ParseInfoPlusStation(changeTipNode.SelectElement("OverstapTipBestemming"))

		departure.ChangeTips = append(departure.ChangeTips, changeTip)
	}

	for _, wingInfo := range trainProduct.SelectElements("TreinVleugel") {
		var trainWing models.TrainWing

		trainWing.DestinationActual = ParseInfoPlusStations(ParseWhenAttributeMulti(wingInfo, "TreinVleugelEindBestemming", "InfoStatus", "Actueel"))
		trainWing.DestinationPlanned = ParseInfoPlusStations(ParseWhenAttributeMulti(wingInfo, "TreinVleugelEindBestemming", "InfoStatus", "Gepland"))
		trainWing.Modifications = ParseInfoPlusModifications(wingInfo)

		stationsNode := ParseWhenAttribute(wingInfo, "StopStations", "InfoStatus", "Actueel")
		stationsNodePlanned := ParseWhenAttribute(wingInfo, "StopStations", "InfoStatus", "Gepland")

		if stationsNode != nil {
			for _, stationInfo := range stationsNode.SelectElements("Station") {
				station := ParseInfoPlusStation(stationInfo)

				trainWing.Stations = append(trainWing.Stations, station)
			}
		}

		if stationsNodePlanned != nil {
			for _, stationInfo := range stationsNodePlanned.SelectElements("Station") {
				station := ParseInfoPlusStation(stationInfo)

				trainWing.StationsPlanned = append(trainWing.StationsPlanned, station)
			}
		}

		for _, materialInfo := range wingInfo.SelectElements("MaterieelDeelDVS") {
			var material models.Material

			material.NaterialType = materialInfo.SelectElement("MaterieelSoort").Text() + "-" + materialInfo.SelectElement("MaterieelAanduiding").Text()

			materialNumberNode := materialInfo.SelectElement("MaterieelNummer")
			materialPositionNode := materialInfo.SelectElement("MaterieelDeelVolgordeVertrek")

			if materialNumberNode != nil {
				material.Number = materialNumberNode.Text()
			}

			if materialPositionNode != nil {
				material.Position, _ = strconv.Atoi(materialPositionNode.Text())
			}

			material.DestinationActual = ParseInfoPlusStation(ParseWhenAttribute(materialInfo, "MaterieelDeelEindBestemming", "InfoStatus", "Actueel"))
			material.DestinationPlanned = ParseInfoPlusStation(ParseWhenAttribute(materialInfo, "MaterieelDeelEindBestemming", "InfoStatus", "Gepland"))

			trainWing.Material = append(trainWing.Material, material)
		}

		departure.TrainWings = append(departure.TrainWings, trainWing)
	}

	// Check for flags that may be set:
	for _, modification := range departure.Modifications {
		switch modification.ModificationType {
		case models.ModificationCancelledDeparture:
			departure.Cancelled = true
		case models.ModificationNotActual:
			departure.NotRealTime = true
		}
	}

	return
}
