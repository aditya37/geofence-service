package util

import "strings"

func ReplaceGeojsonSingleQuote(geojson string) string {
	if isContainSingleQuote := strings.Contains(geojson, "'"); isContainSingleQuote {
		return strings.Replace(geojson, "'", `"`, -1)
	} else if isPetik := strings.Contains(geojson, `"`); isPetik {
		return strings.Replace(geojson, "'", `"`, -1)
	} else {
		return geojson
	}
}
