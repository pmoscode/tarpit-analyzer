package geo_location

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

var debug = false

func DoLocalization(pathSource string, pathTarget string, batchSize int, geoLocationLongitude string, geoLocationLatitude string, debugParam bool) error {
	debug = debugParam
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	file, err := os.Open(pathSource)
	if err != nil {
		log.Errorln(err)
		return err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)

	ips := make([]string, 0)

	// Read through 'tokens' until an EOF is encountered.
	for sc.Scan() {
		text := sc.Text()
		ip, err := processLine(text)
		if err != nil {
			log.Warningln(err)
		} else {
			if ip != "" {
				ips = append(ips, ip)
			}
		}
	}

	if err := sc.Err(); err != nil {
		log.Errorln(err)
	}

	ips = uniqueNonEmptyElementsOf(ips)

	log.Debugln("Unique ips found: ", ips)
	ipLocations, err := processIps(ips, batchSize)
	log.Debugln("Locations for ips: ", ipLocations)

	errOutput := writeConvertedDataToFile(ipLocations, pathTarget, geoLocationLongitude, geoLocationLatitude)
	if errOutput != nil {
		return errOutput
	}

	return nil
}

func processLine(line string) (string, error) {
	chunks := strings.Split(line, ",")

	ip := chunks[1]

	return ip, nil
}

func processIps(ips []string, batchSize int) ([]IpLocation, error) {
	ipsArraySize := len(ips)
	ipLocations := make([]IpLocation, 0)

	for i := 0; i < ipsArraySize; i += batchSize {
		if i+batchSize >= ipsArraySize {
			batchSize = ipsArraySize - i
		}

		ipBatch := ips[i : i+batchSize]
		resolved, _ := DoQuery(ipBatch)
		ipLocations = append(ipLocations, resolved...)
	}

	return ipLocations, nil
}

func uniqueNonEmptyElementsOf(s []string) []string {
	unique := make(map[string]bool, len(s))
	us := make([]string, len(unique))
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				us = append(us, elem)
				unique[elem] = true
			}
		}
	}

	return us
}

func writeConvertedDataToFile(ips []IpLocation, path string, geoLocationLongitude string, geoLocationLatitude string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Errorln(err)
		return err
	}

	dataWriter := bufio.NewWriter(file)

	_, _ = dataWriter.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
	<Document>
		<Style id="transBluePoly">
			<LineStyle>
				<width>1.5</width>
				<color>501400E6</color>
			</LineStyle>
		</Style>`)

	for _, data := range ips {
		if data.Status == "success" {
			_, _ = dataWriter.WriteString(`
		<Placemark>
			<name>` + data.Country + `</name>
			<extrude>1</extrude>
			<tessellate>1</tessellate>
			<styleUrl>#transBluePoly</styleUrl>
			<LineString>
				<coordinates> 
					` + fmt.Sprintf("%f", data.Lon) + `, ` + fmt.Sprintf("%f", data.Lat) + `
					` + geoLocationLongitude + `,` + geoLocationLatitude + `					
				</coordinates>
			</LineString>
		</Placemark>`)
		}
	}

	_, _ = dataWriter.WriteString(`
	</Document>
</kml>`)

	errWriter := dataWriter.Flush()
	if errWriter != nil {
		return errWriter
	}

	errFile := file.Close()
	if errFile != nil {
		return errFile
	}

	return nil
}

/*
<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<kml xmlns=\"http://www.opengis.net/kml/2.2\">
	<Document>
		<Style id=\"transBluePoly\">
			<LineStyle>
				<width>1.5</width>
				<color>501400E6</color>
			</LineStyle>
		</Style>

		<Placemark>
			<name> wert6 </name>
			<extrude>1</extrude>
			<tessellate>1</tessellate>
			<styleUrl>#transBluePoly</styleUrl>
			<LineString>
				<coordinates> wert8 , wert7 + System.getProperty("line.separator")  + longlat </coordinates>
			</LineString>
		</Placemark>

	</Document>
</kml>
*/
