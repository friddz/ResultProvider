package motleikir

import "encoding/xml"

type Response struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    *ResponseBody
}

type ResponseBody struct {
	XMLName           xml.Name `xml:"Body"`
	MotLeikirResponse *ResponseMotLeikir
}

type ResponseMotLeikir struct {
	XMLName       xml.Name `xml:"MotLeikirResponse"`
	MotLeikirSvar *ResponseMotLeikirSvar
}

type ResponseMotLeikirSvar struct {
	XMLName        xml.Name `xml:"MotLeikirSvar"`
	ArrayMotLeikir *ResponseArrayMotLeikir
}

type ResponseArrayMotLeikir struct {
	XMLName   xml.Name `xml:"ArrayMotLeikir"`
	MotLeikur []ResponseMotLeikur
}

type ResponseMotLeikur struct {
	XMLName                 xml.Name `xml:"MotLeikur"`
	LeikurNumer             string
	UmferdNumer             string
	LeikDagur               string
	FelagHeimaNafn          string
	FelagUtiNafn            string
	UrslitHeima             string
	UrslitUti               string
	StadaFyrriHalfleikHeima string
	StadaFyrriHalfleikUti   string
}
