package parsers

import (
	"errors"
	"fmt"
	"io"

	"github.com/beevik/etree"
	"github.com/rijdendetreinen/gotrain/models"
	"github.com/rs/zerolog/log"
)

func ParseDvs3Message(reader io.Reader) (departure models.Departure, err error) {
	doc := etree.NewDocument()

	if _, err := doc.ReadFrom(reader); err != nil {
		return departure, err
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parser error: %+v", r)
			log.Error().Err(err).Msg("Recovered from panic in ParseDvsMessage")
		}
	}()

	dvs3_product := doc.SelectElement("PutReisInformatieBoodschapIn").SelectElement("reisInformatieProductDVS")

	if dvs3_product != nil {
		if dvs3_product.NamespaceURI() != "urn:ns:cdm:reisinformatie:data:dvs:3" {
			dvs3_product = nil
		}

		return parseDvs3Product(dvs3_product)
	}

	err = errors.New("missing DVS3 element")
	log.Error().Err(err).Msg("Failed to find DVS3 element in XML")

	return
}

func parseDvs3Product(product *etree.Element) (departure models.Departure, err error) {
	departure.DvsVersion = 3

	return departure, errors.New("not implemented")
}
