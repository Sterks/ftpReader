package ftp

import "encoding/xml"

type Notification struct {
	XMLName           xml.Name `xml:"export"`
	Text              string   `xml:",chardata"`
	Ns5               string   `xml:"ns5,attr"`
	Xmlns             string   `xml:"xmlns,attr"`
	Ns6               string   `xml:"ns6,attr"`
	Ns7               string   `xml:"ns7,attr"`
	Ns8               string   `xml:"ns8,attr"`
	Ns9               string   `xml:"ns9,attr"`
	Ns2               string   `xml:"ns2,attr"`
	Ns3               string   `xml:"ns3,attr"`
	Ns4               string   `xml:"ns4,attr"`
	FcsNotificationEF struct {
		Text           string `xml:",chardata"`
		SchemeVersion  string `xml:"schemeVersion,attr"`
		ID             string `xml:"id"`
		PurchaseNumber string `xml:"purchaseNumber"`
		DirectDate     string `xml:"directDate"`
		DocPublishDate string `xml:"docPublishDate"`
		DocNumber      string `xml:"docNumber"`
		Href           string `xml:"href"`
		PrintForm      struct {
			Text      string `xml:",chardata"`
			URL       string `xml:"url"`
			Signature struct {
				Text string `xml:",chardata"`
				Type string `xml:"type,attr"`
			} `xml:"signature"`
		} `xml:"printForm"`
		PurchaseObjectInfo  string `xml:"purchaseObjectInfo"`
		PurchaseResponsible struct {
			Text           string `xml:",chardata"`
			ResponsibleOrg struct {
				Text            string `xml:",chardata"`
				RegNum          string `xml:"regNum"`
				ConsRegistryNum string `xml:"consRegistryNum"`
				FullName        string `xml:"fullName"`
				PostAddress     string `xml:"postAddress"`
				FactAddress     string `xml:"factAddress"`
				INN             string `xml:"INN"`
				KPP             string `xml:"KPP"`
			} `xml:"responsibleOrg"`
			ResponsibleRole string `xml:"responsibleRole"`
			ResponsibleInfo struct {
				Text           string `xml:",chardata"`
				OrgPostAddress string `xml:"orgPostAddress"`
				OrgFactAddress string `xml:"orgFactAddress"`
				ContactPerson  struct {
					Text       string `xml:",chardata"`
					LastName   string `xml:"lastName"`
					FirstName  string `xml:"firstName"`
					MiddleName string `xml:"middleName"`
				} `xml:"contactPerson"`
				ContactEMail string `xml:"contactEMail"`
				ContactPhone string `xml:"contactPhone"`
				ContactFax   string `xml:"contactFax"`
			} `xml:"responsibleInfo"`
		} `xml:"purchaseResponsible"`
		PlacingWay struct {
			Text string `xml:",chardata"`
			Code string `xml:"code"`
			Name string `xml:"name"`
		} `xml:"placingWay"`
		ContractConclusionOnSt83Ch2 string `xml:"contractConclusionOnSt83Ch2"`
		ETP                         struct {
			Text string `xml:",chardata"`
			Code string `xml:"code"`
			Name string `xml:"name"`
			URL  string `xml:"url"`
		} `xml:"ETP"`
		ProcedureInfo struct {
			Text       string `xml:",chardata"`
			Collecting struct {
				Text      string `xml:",chardata"`
				StartDate string `xml:"startDate"`
				Place     string `xml:"place"`
				Order     string `xml:"order"`
				EndDate   string `xml:"endDate"`
			} `xml:"collecting"`
			Scoring struct {
				Text string `xml:",chardata"`
				Date string `xml:"date"`
			} `xml:"scoring"`
			Bidding struct {
				Text string `xml:",chardata"`
				Date string `xml:"date"`
			} `xml:"bidding"`
		} `xml:"procedureInfo"`
		Lot struct {
			Text     string `xml:",chardata"`
			MaxPrice string `xml:"maxPrice"`
			Currency struct {
				Text string `xml:",chardata"`
				Code string `xml:"code"`
				Name string `xml:"name"`
			} `xml:"currency"`
			FinanceSource        string `xml:"financeSource"`
			QuantityUndefined    string `xml:"quantityUndefined"`
			CustomerRequirements struct {
				Text                string `xml:",chardata"`
				CustomerRequirement struct {
					Text     string `xml:",chardata"`
					Customer struct {
						Text            string `xml:",chardata"`
						RegNum          string `xml:"regNum"`
						ConsRegistryNum string `xml:"consRegistryNum"`
						FullName        string `xml:"fullName"`
					} `xml:"customer"`
					MaxPrice    string `xml:"maxPrice"`
					KladrPlaces struct {
						Text       string `xml:",chardata"`
						KladrPlace struct {
							Text          string `xml:",chardata"`
							DeliveryPlace string `xml:"deliveryPlace"`
						} `xml:"kladrPlace"`
					} `xml:"kladrPlaces"`
					DeliveryTerm      string `xml:"deliveryTerm"`
					ContractGuarantee struct {
						Text              string `xml:",chardata"`
						Amount            string `xml:"amount"`
						ProcedureInfo     string `xml:"procedureInfo"`
						SettlementAccount string `xml:"settlementAccount"`
						PersonalAccount   string `xml:"personalAccount"`
						Bik               string `xml:"bik"`
					} `xml:"contractGuarantee"`
					PurchaseCode   string `xml:"purchaseCode"`
					TenderPlanInfo struct {
						Text               string `xml:",chardata"`
						Plan2017Number     string `xml:"plan2017Number"`
						Position2017Number string `xml:"position2017Number"`
					} `xml:"tenderPlanInfo"`
					PurchaseObjectDescription string `xml:"purchaseObjectDescription"`
				} `xml:"customerRequirement"`
			} `xml:"customerRequirements"`
			PurchaseObjects struct {
				Text           string `xml:",chardata"`
				PurchaseObject []struct {
					Text  string `xml:",chardata"`
					OKPD2 struct {
						Text string `xml:",chardata"`
						Code string `xml:"code"`
						Name string `xml:"name"`
					} `xml:"OKPD2"`
					Name string `xml:"name"`
					OKEI struct {
						Text         string `xml:",chardata"`
						Code         string `xml:"code"`
						NationalCode string `xml:"nationalCode"`
						FullName     string `xml:"fullName"`
					} `xml:"OKEI"`
					CustomerQuantities struct {
						Text             string `xml:",chardata"`
						CustomerQuantity struct {
							Text     string `xml:",chardata"`
							Customer struct {
								Text            string `xml:",chardata"`
								RegNum          string `xml:"regNum"`
								ConsRegistryNum string `xml:"consRegistryNum"`
								FullName        string `xml:"fullName"`
							} `xml:"customer"`
							Quantity string `xml:"quantity"`
						} `xml:"customerQuantity"`
					} `xml:"customerQuantities"`
					Price    string `xml:"price"`
					Quantity struct {
						Text  string `xml:",chardata"`
						Value string `xml:"value"`
					} `xml:"quantity"`
					Sum string `xml:"sum"`
				} `xml:"purchaseObject"`
				TotalSum string `xml:"totalSum"`
			} `xml:"purchaseObjects"`
			Requirements struct {
				Text        string `xml:",chardata"`
				Requirement []struct {
					Text      string `xml:",chardata"`
					ShortName string `xml:"shortName"`
					Name      string `xml:"name"`
					Content   string `xml:"content"`
				} `xml:"requirement"`
			} `xml:"requirements"`
			Restrictions struct {
				Text        string `xml:",chardata"`
				Restriction struct {
					Text      string `xml:",chardata"`
					ShortName string `xml:"shortName"`
					Name      string `xml:"name"`
					Content   string `xml:"content"`
				} `xml:"restriction"`
			} `xml:"restrictions"`
			MustPublicDiscussion string `xml:"mustPublicDiscussion"`
		} `xml:"lot"`
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
		Modification struct {
			Text               string `xml:",chardata"`
			ModificationNumber string `xml:"modificationNumber"`
			Info               string `xml:"info"`
			Reason             struct {
				Text                string `xml:",chardata"`
				ResponsibleDecision struct {
					Text         string `xml:",chardata"`
					DecisionDate string `xml:"decisionDate"`
				} `xml:"responsibleDecision"`
			} `xml:"reason"`
		} `xml:"modification"`
		Documentation struct {
			Text                   string `xml:",chardata"`
			PurchaseObjectsCh9St37 string `xml:"purchaseObjectsCh9St37"`
			Modifiable             string `xml:"modifiable"`
			ClarificationInfo      struct {
				Text                    string `xml:",chardata"`
				StartDate               string `xml:"startDate"`
				FilledManuallyStartDate string `xml:"filledManuallyStartDate"`
				EndDate                 string `xml:"endDate"`
				DeliveryProcedure       string `xml:"deliveryProcedure"`
			} `xml:"clarificationInfo"`
			OnesideRejectionCh9St95 string `xml:"onesideRejectionCh9St95"`
			PrintFormInfo           struct {
				Text      string `xml:",chardata"`
				URL       string `xml:"url"`
				Signature struct {
					Text string `xml:",chardata"`
					Type string `xml:"type,attr"`
				} `xml:"signature"`
			} `xml:"printFormInfo"`
		} `xml:"documentation"`
	} `xml:"fcsNotificationEF"`
}
