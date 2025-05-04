package parsers

import (
	"errors"
	"fmt"
	"io"

	"github.com/beevik/etree"
	"github.com/rijdendetreinen/gotrain/models"
	"github.com/rs/zerolog/log"
)

// ParseDvsMessage parses a DVS XML message to a Departure object
func ParseDvsMessage(reader io.Reader) (departure models.Departure, err error) {
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

	// try to find DVS2 product first
	dvs2_product := doc.SelectElement("PutReisInformatieBoodschapIn").SelectElement("ReisInformatieProductDVS")
	if dvs2_product != nil {
		if dvs2_product.NamespaceURI() != "urn:ndov:cdm:trein:reisinformatie:data:4" {
			dvs2_product = nil
		} else {
			return parseDvs2Product(dvs2_product)
		}
	}

	// try to find DVS3 product
	dvs3_product := doc.SelectElement("PutReisInformatieBoodschapIn").SelectElement("reisInformatieProductDVS")

	if dvs3_product != nil {
		if dvs3_product.NamespaceURI() != "urn:ns:cdm:reisinformatie:data:dvs:3" {
			dvs3_product = nil
		} else {
			return parseDvs3Product(dvs3_product)
		}
	}

	// if neither DVS2 nor DVS3 product is found, return an error
	err = errors.New("missing DVS element (neither DVS2 nor DVS3)")
	log.Error().Err(err).Msg("Failed to find DVS element in XML")

	return
}
