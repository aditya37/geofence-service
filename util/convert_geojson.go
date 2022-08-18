package util

import "strings"

func ReplaceGeojsonSingleQuote(geojson string) string {
	if isContainSingleQuote := strings.Contains(geojson, "'"); isContainSingleQuote {
		return strings.Replace(geojson, "'", `"`, -1)
	}
	return geojson
}
