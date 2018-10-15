package common

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"time"

	pi "github.com/christoofar/gowebapi"
)

var cfg = pi.NewConfiguration()
var client *pi.APIClient
var auth context.Context

func Init() {
	cfg.BasePath = "https://piwebserver.yourcompany.com/piwebapi"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	auth = context.WithValue(context.Background(), pi.ContextBasicAuth, pi.BasicAuth{
		UserName: "login",
		Password: "password",
	})

	client = pi.NewAPIClient(cfg)
	cfg.HTTPClient.Transport = tr

}

func UploadData(data HealthData, attribute string, position int) {

	heartRateWebId := pi.EncodeWebID(pi.NewAFAttributeWebID("CLSAF\\HealthData\\"+data.Device+"|"+attribute, pi.IS_AF_ELEMENT))

	var values []pi.TimedValue

	for i := 0; i < len(data.Readings); i++ {
		var val interface{}

		switch position {
		case 0:
			val = data.Readings[i].Value
			break
		case 1:
			val = data.Readings[i].Value2
			break
		case 2:
			val = data.Readings[i].Value3
			break
		case 3:
			val = data.Readings[i].Value3
			break
		}

		var tv pi.TimedValue

		tv.Timestamp = data.Readings[i].EndTime
		tv.Value = &val

		values = append(values, tv)
	}

	low := 0  // lowest index
	high := 0 // highest index

	if len(values) > 10000 {
		high = 9999
	} else {
		high = len(values) - 1
	}

	for low < len(values)-1 {

		optionals := make(map[string]interface{})
		optionals["updateOption"] = "NoReplace"

		posted, resp, err := client.StreamApi.StreamUpdateValues(auth, heartRateWebId, values[low:high], optionals)

		if err != nil {
			log.Println(err)
		}

		if resp != nil {
			log.Printf("Response code is %d\n", resp.StatusCode)
			log.Printf("# of values received: %d  posted: %d\n", len(values), len(posted.Items))
		}

		if (high + 10000) > len(values) {
			low += 10000
			high = len(values) - 1
		} else {
			low += 10000
			high += 10000
		}

	}
}

// Create the Device element in AF if it's not already there
func CheckPhoneElement(deviceName string) bool {
	phoneWebId := pi.EncodeWebID(pi.NewAFElementWebID("CLSAF\\HealthData\\" + deviceName))

	// Look and see if we already know about this device and this reading type in AF.
	phoneElement, _, err := client.ElementApi.ElementGet(auth, phoneWebId, nil)
	if err != nil {
		log.Println(err) // :(
	} else {
		return true
	}

	if phoneElement.Id == "" {
		// We need to create the element
		var element pi.Element

		element.Name = deviceName
		element.TemplateName = "Health"

		dbname := pi.EncodeWebID(pi.NewAFDatabaseWebID("CLSAF\\HealthData"))

		response, err := client.AssetDatabaseApi.AssetDatabaseCreateElement(auth, dbname, element, nil)
		if err != nil {
			log.Println(err)
			return false
		}
		if response.StatusCode == 201 {
			log.Println("Created phone element")
			client.ElementApi.ElementCreateConfig(auth, phoneWebId, nil)
		} else {
			return false
		}
	}

	return true
}

// Writes received data to AF
func PostValues(data HealthData, attribute string, position int) {

	Init()
	if CheckPhoneElement(data.Device) {
		// Push data up.  To do this, we should send it as a batch.
		UploadData(data, attribute, position)
	}

}

// This is a bundle to send a PI Web API batch request.
type BatchRequest struct {
	Resource string `json:"Resource,omitempty"`
	Method   string `json:"Method,omitempty"`
	Content  string `json:"Content,omitempty"`
}

// This is the message that is passed back after the phone calls this service.  It returns
// the number of items received, the time and any error messages.
type HealthResponse struct {
	CachedItems int       `json:"cached_items"`
	Message     string    `json:"message"`
	Time        time.Time `json:"time"`
}

// A basic health value reading.  For things like steps where there are multiple measurements, you
// place them in Value2, Value3 and Value4 and read them back out.
type HealthReading struct {
	Created   time.Time `json:"created_time"`
	BeginTime time.Time `json:"begin_time"`
	EndTime   time.Time `json:"end_time"`
	Value     float64   `json:"value"`
	Value2    float64   `json:"value2"`
	Value3    float64   `json:"value3"`
	Value4    float64   `json:"value4"`
}

// A collection of health value readings
type HealthData struct {
	DataType string          `json:"data_type"`
	Device   string          `json:"device"`
	Readings []HealthReading `json:"readings"`
}
