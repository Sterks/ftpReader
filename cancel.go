package main

import "encoding/xml"

// Export ...
type Export struct {
	XMLName               xml.Name `xml:"export"`
	Text                  string   `xml:",chardata"`
	Ns5                   string   `xml:"ns5,attr"`
	Xmlns                 string   `xml:"xmlns,attr"`
	Ns6                   string   `xml:"ns6,attr"`
	Ns7                   string   `xml:"ns7,attr"`
	Ns8                   string   `xml:"ns8,attr"`
	Ns9                   string   `xml:"ns9,attr"`
	Ns2                   string   `xml:"ns2,attr"`
	Ns3                   string   `xml:"ns3,attr"`
	Ns4                   string   `xml:"ns4,attr"`
	FcsNotificationCancel struct {
		Text           string `xml:",chardata"`
		SchemeVersion  string `xml:"schemeVersion,attr"`
		ID             string `xml:"id"`
		PurchaseNumber string `xml:"purchaseNumber"`
		DocNumber      string `xml:"docNumber"`
		DocDate        string `xml:"docDate"`
		DocPublishDate string `xml:"docPublishDate"`
		Href           string `xml:"href"`
		PrintForm      struct {
			Text      string `xml:",chardata"`
			URL       string `xml:"url"`
			Signature struct {
				Text string `xml:",chardata"`
				Type string `xml:"type,attr"`
			} `xml:"signature"`
		} `xml:"printForm"`
		CancelReason struct {
			Text                string `xml:",chardata"`
			ResponsibleDecision struct {
				Text         string `xml:",chardata"`
				DecisionDate string `xml:"decisionDate"`
			} `xml:"responsibleDecision"`
		} `xml:"cancelReason"`
		AddInfo     string `xml:"addInfo"`
		Attachments struct {
			Text       string `xml:",chardata"`
			Attachment struct {
				Text               string `xml:",chardata"`
				PublishedContentId string `xml:"publishedContentId"`
				FileName           string `xml:"fileName"`
				FileSize           string `xml:"fileSize"`
				DocDescription     string `xml:"docDescription"`
				URL                string `xml:"url"`
				CryptoSigns        struct {
					Text      string `xml:",chardata"`
					Signature struct {
						Text string `xml:",chardata"`
						Type string `xml:"type,attr"`
					} `xml:"signature"`
				} `xml:"cryptoSigns"`
			} `xml:"attachment"`
		} `xml:"attachments"`
	} `xml:"fcsNotificationCancel"`
}
