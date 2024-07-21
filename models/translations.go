package models

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

var causeTranslations = map[string]string{
	"door geplande werkzaamheden":                                        "due to planned engineering work",
	"door werkzaamheden":                                                 "due to engineering work",
	"door werkzaamheden elders":                                          "due to engineering work elsewhere",
	"door werkzaamheden aan de Hogesnelheidslijn":                        "due to engineering work on the HSL",
	"door onverwachte werkzaamheden":                                     "due to unexpected engineering work",
	"door uitgelopen werkzaamheden":                                      "due to over-running engineering works",
	"door uitloop van werkzaamheden":                                     "due to over-running engineering works",
	"door de aanleg van een nieuw spoor":                                 "due to construction of a new track",
	"door een spoedreparatie aan het spoor":                              "due to emergency repairs",
	"door een aangepaste dienstregeling":                                 "due to an amended timetable",
	"door te grote vertraging":                                           "due to large delay",
	"door te hoog opgelopen vertraging":                                  "due to excessive delay",
	"door te veel vertraging in het buitenland":                          "due to excessive delay abroad",
	"door een eerdere verstoring":                                        "due to an earlier disruption",
	"door herstelwerkzaamheden":                                          "due to reparation works",
	"door een seinstoring":                                               "due to signal failure",
	"door een sein- en wisselstoring":                                    "due to signalling and points failure",
	"door een sein-en wisselstoring":                                     "due to signalling and points failure",
	"door een grote sein- en wisselstoring":                              "due to a large signalling failure",
	"door een storing aan bediensysteem seinen en wissels":               "due to a control system failure",
	"door een storing in de bediening van seinen":                        "due to a control system failure",
	"door defect materieel":                                              "due to a broken down train",
	"door een defecte trein":                                             "due to a broken down train",
	"door defecte treinen":                                               "due to broken down trains",
	"door een ontspoorde trein":                                          "due to a derailed train",
	"door een gestrande trein":                                           "due to a stranded train",
	"door een defecte spoorbrug":                                         "due to a defective railway bridge",
	"door een beschadigd spoorviaduct":                                   "due to a damaged railway bridge",
	"door een beschadigde spoorbrug":                                     "due to a damaged railway bridge",
	"door beperkingen in de materieelinzet":                              "due to rolling stock problems",
	"door beperkingen in het buitenland":                                 "due to restrictions abroad",
	"door acties van het personeel":                                      "due to staff strike",
	"door acties in het buitenland":                                      "due to staff strike abroad",
	"door een wisselstoring":                                             "due to points failure",
	"door een defect wissel":                                             "due to a defective switch",
	"door veel defect materieel":                                         "due to numerous broken down trains",
	"door een overwegstoring":                                            "due to level crossing failure",
	"door overwegstoringen":                                              "due to level crossing failures",
	"door een aanrijding met een persoon":                                "due to a person hit by a train",
	"door een aanrijding":                                                "due to a collision",
	"door een aanrijding met een dier":                                   "due to a collision with an animal",
	"door een aanrijding met een voertuig":                               "due to a collision with a vehicle",
	"door een auto op het spoor":                                         "due to a car on the track",
	"door mensen op het spoor":                                           "due to persons on the track",
	"door mensen langs het spoor":                                        "due to people along the track",
	"door een dier op het spoor":                                         "due to an animal on the track",
	"door een boom op het spoor":                                         "due to a tree on the track",
	"door een verstoring elders":                                         "due to a disruption elsewhere",
	"door een persoon op het spoor":                                      "due to a trespassing incident",
	"door een persoon langs het spoor":                                   "due to a trespassing incident",
	"door een defect spoor":                                              "due to a defective rail",
	"door een defect aan het spoor":                                      "due to a defective rail",
	"door gladde sporen":                                                 "due to slippery rail",
	"door een defecte bovenleiding":                                      "due to overhead wire problems",
	"door een beschadigde bovenleiding":                                  "due to a damaged overhead wire",
	"door een beschadigde overweg":                                       "due to a damaged level crossing",
	"door een defecte overweg":                                           "due to a defective level crossing",
	"door een versperring":                                               "due to an obstruction on the line",
	"door een versperring van het spoor":                                 "due to an obstruction on the line",
	"door inzet van de brandweer":                                        "due to deployment of the fire brigade",
	"door inzet van de politie":                                          "due to police action",
	"door brand in een trein":                                            "due to fire in a train",
	"op last van de politie":                                             "due to restrictions imposed by the police",
	"op last van de brandweer":                                           "due to restrictions imposed by the fire brigade",
	"door politieonderzoek":                                              "due to police investigation",
	"door vandalisme":                                                    "due to vandalism",
	"door inzet van hulpdiensten":                                        "due to an emergency call",
	"door een stroomstoring":                                             "due to power disruption",
	"door stormschade":                                                   "due to storm damage",
	"door een bermbrand":                                                 "due to a lineside fire",
	"door diverse oorzaken":                                              "due to various reasons",
	"door meerdere verstoringen":                                         "due to multiple disruptions",
	"door koperdiefstal":                                                 "due to copper theft",
	"verwachte weersomstandigheden":                                      "due to expected weather conditions",
	"door de verwachte weersomstandigheden":                              "due to expected weather conditions",
	"door de weersomstandigheden":                                        "due to bad weather conditions",
	"sneeuw":                                                             "due to snow",
	"door rijp aan de bovenleiding":                                      "due to frost on the overhead wires",
	"door ijzelvorming aan de bovenleiding":                              "due to ice on the overhead wires",
	"door harde wind op de Hogesnelheidslijn":                            "due to strong winds on the high-speed line",
	"door het onschadelijk maken van een bom uit de Tweede Wereldoorlog": "due to defusing a bomb from World War II",
	"door het onschadelijk maken van een bom uit de 2e WO":               "due to defusing a bomb from World War II",
	"door een evenement":                                                 "due to an event",
	"door een sein-en overwegstoring":                                    "due to signalling failure and a level crossing failure",
	"door een sein- en overwegstoring":                                   "due to signalling failure and a level crossing failure",
	"door technisch onderzoek":                                           "due to technical inspection",
	"door een brandmelding":                                              "due to a fire alarm",
	"door een voorwerp in de bovenleiding":                               "due to an obstacle in the overhead wire",
	"door een voorwerp op het spoor":                                     "due to an obstacle on the track",
	"door rommel op het spoor":                                           "due to rubbish on the track",
	"door grote drukte":                                                  "due to large crowds",
	"door blikseminslag":                                                 "due to lightning",
	"door wateroverlast":                                                 "due to flooding",
	"door problemen op het spoor in het buitenland":                      "due to railway problems abroad",
	"door problemen in het buitenland":                                   "due to railway problems abroad",
	"door spoorproblemen in het buitenland":                              "due to railway problems abroad",
	"door een storing in een tunnel":                                     "due to a problem in a tunnel",
	"door hinder op het spoor":                                           "due to interference on the line",
	"veiligheidsredenen":                                                 "due to safety reasons",
	"door het onverwacht ontbreken van personeel":                        "due to missing crew",
	"door problemen met de personeelsinzet":                              "due to staffing problems",
	"door een vervangende trein":                                         "due to a replacement train",
	"door het vervangen van een spoorbrug":                               "due to replacement of a railway bridge",
	"door Koningsdag":                                                    "due to King's day",
	"door de Vierdaagse":                                                 "due to the Four Days Marches",
	"door een object op het spoor":                                       "due to an object on the track",
	"door een voertuig op het spoor":                                     "due to a vehicle on the track",
	"door nog onbekende oorzaak":                                         "due to a yet unknown reason",
	"door de inzet van ander materieel":                                  "due to a replacement train",
	"door personen op het spoor":                                         "due to persons on the railway",
	"door personen langs het spoor":                                      "due to persons along the railway",
	"door een staking":                                                   "due to a strike",
	"door een tekort aan beschikbaar personeel":                          "due to staff shortage",
	"door een versperring op het spoor":                                  "due to an obstruction on the line",
	"door logistieke beperkingen":                                        "due to logistical limitations",
	"tekort aan verkeersleiders":                                         "due to shortage of traffic controllers",
	"door een tekort aan verkeersleiders":                                "due to shortage of traffic controllers",
	"door een tekort aan personeel":                                      "due to staff shortage",
	"door acties personeel":                                              "due to staff strike",
	"door een systeemstoring":                                            "due to a system failure",
	"door noodzakelijke aanpassing in de dienstregeling":                 "due to a necessary timetable change",
	"om veiligheidsredenen":                                              "for safety reasons",
	"door een dienstregelingswijziging":                                  "due to a timetable change",
	"inzet veiligheidsmedewerkers":                                       "due to deployment of safety personnel",
}

// Translate returns the appropriate translation based on language
func Translate(remarkNL, remarkEN, language string) string {
	remark := remarkNL

	if language == "en" {
		remark = remarkEN
	}

	return remark
}

// TranslateStations returns the appropriate translation based on language
func TranslateStations(remarkNL, remarkEN string, stations []Station, language string) string {
	translation := Translate(remarkNL, remarkEN, language)

	return fmt.Sprintf(translation, stationsStringTranslated(stations, language))
}

// TranslateCause translates a Dutch cause (long version) to an English cause
func TranslateCause(causeLong string) string {
	translation, exists := causeTranslations[causeLong]

	if exists {
		return translation
	}

	// Unknown translation, return original
	log.WithField("cause", causeLong).Warnf("No translation for cause: %s", causeLong)
	return causeLong
}

func stationsStringTranslated(stations []Station, language string) string {
	stationsText := ""

	for index, station := range stations {
		if index > 0 {
			if index < len(stations)-1 {
				stationsText += ", "
			} else {
				if language == "en" {
					stationsText += " and "
				} else {
					stationsText += " en "
				}
			}
		}
		stationsText += station.NameLong
	}

	return stationsText
}
