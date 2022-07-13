package commands

import (
	"gopkg.in/telebot.v3"
)

func (h Commands) Address(tx telebot.Context) error {
	location := h.cfg.Location

	venue := &telebot.Venue{
		Title:           location.Title,
		Address:         location.Address,
		FoursquareID:    location.FoursquareID,
		FoursquareType:  location.FoursquareType,
		GooglePlaceID:   location.GooglePlaceID,
		GooglePlaceType: location.GooglePlaceType,
		Location: telebot.Location{
			Lat: location.Latitude,
			Lng: location.Longitude,
		},
	}

	return tx.Reply(venue, h.menu(tx))
}
