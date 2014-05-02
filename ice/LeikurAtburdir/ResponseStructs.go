package leikuratburdir

import "encoding/xml"

type Response struct{
	XMLName xml.Name `xml:"Envelope"`
	Body *ResponseBody
}

type ResponseBody struct{
	XMLName xml.Name `xml:"Body"`
	LeikurAtburdirResponse *ResponseLeikurAtburdir
}

type ResponseLeikurAtburdir struct{
	XMLName xml.Name `xml:"LeikurAtburdirResponse"`
	LeikurAtburdirSvar *ResponseLeikurAtburdirSvar
}

type ResponseLeikurAtburdirSvar struct{
	XMLName xml.Name `xml:"LeikurAtburdurSvar"`
	ArrayLeikurAtburdir *ResponseArrayLeikurAtburdir
}

type ResponseArrayLeikurAtburdir struct {
	XMLName xml.Name `xml:"ArrayLeikurAtburdir"`
	LeikurAtburdir []ResponseLeikurAtburdur
}

type ResponseLeikurAtburdur struct{
	XMLName xml.Name `xml:"LeikurAtburdir"`
	LeikmadurNumer int
	LeikmadurNafn string 
	FelagNafn string
	TreyjuNumer uint8
	AtburdurMinuta uint8
	AtburdurNumer uint8
	AtburdurNafn string
}