package config

type Location struct {
	Title           string  `fig:"title" yaml:"title"`
	Address         string  `fig:"address" yaml:"address"`
	Latitude        float32 `fig:"latitude" yaml:"latitude"`
	Longitude       float32 `fig:"longitude" yaml:"longitude"`
	FoursquareID    string  `fig:"foursquare_id" yaml:"foursquare_id"`
	FoursquareType  string  `fig:"foursquare_type" yaml:"foursquare_type"`
	GooglePlaceID   string  `fig:"google_place_id" yaml:"google_place_id"`
	GooglePlaceType string  `fig:"google_place_type" yaml:"google_place_type"`
}
