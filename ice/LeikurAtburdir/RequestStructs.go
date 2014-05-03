package leikuratburdir

import "encoding/xml"

type RequestEnvelope struct {
	XMLName xml.Name `xml:"http://www.w3.org/2001/XMLSchema-instance soap:Envelope"`
	Xmlns   string   `xml:"xmlns:soap,attr"`
	Soap    *RequestSoapBody
}
type RequestSoapBody struct {
	XMLName xml.Name               `xml:"soap:Body"`
	Xmlns   string                 `xml:"xmlns:xsi,attr"`
	Content *RequestLeikurAtburdir `xml:"LeikurAtburdir"`
}

type RequestLeikurAtburdir struct {
	LeikurNumer string
	Xmlns       string `xml:"xmlns,attr"`
}
