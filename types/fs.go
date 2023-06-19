package types

// Structs - the values we bring in from *app.json configuration file
type Tp_general struct {
	EchoConfig          int
	Hostname            string
	Debuglevel          int
	Echojson            int
	Testsize            int    // Used to limit number of records posted, over rided when reading test cases from input_source,
	Sleep               int    // sleep time between API post
	Httpposturl         string // FeatureSpace API URL
	Cert_dir            string
	Cert_file           string
	Cert_key            string
	Call_fs_api         int     // Do we call the API endpoint
	Eventtype           string  // paymentRT or paymentNRT
	Json_to_file        int     // Do we output JSON to file in output_path
	Output_path         string  // output location
	Json_from_file      int     // Do we read JSON from input_path directory and post to FS API endpoint
	Input_path          string  // Where are my scenario JSON files located
	MinTransactionValue float64 // Min value if the fake transaction
	MaxTransactionValue float64 // Max value of the fake transaction
	SeedFile            string  // Which seed file to read in
	EchoSeed            int     // 0/1 Echo the seed data to terminal
	IntroBadEntity      int
	IntroBadPayer       int
	IntroBadAgent       int
}

// FS engineResponse components
type TAmount struct {
	BaseCurrency     string  `json:"baseCurrency,omitempty"`
	BaseValue        float64 `json:"baseValue,omitempty"`
	Currency         string  `json:"currency,omitempty"`
	Value            float64 `json:"value,omitempty"`
	MerchantCurrency string  `json:"merchantCurrency,omitempty"`
	MerchantValue    float64 `json:"merchantValue,omitempty"`
}

type TAddress struct {
	AddressLine1          string `json:"addressLine1,omitempty"`
	AddressLine2          string `json:"addressLine2,omitempty"`
	AddressLine3          string `json:"addressLine3,omitempty"`
	AddressType           string `json:"addressType,omitempty"`
	Country               string `json:"country,omitempty"`
	CountrySubDivision    string `json:"countrySubDivision,omitempty"`
	Latitude              string `json:"latitude,omitempty"`
	Longitude             string `json:"longitude,omitempty"`
	PostalCode            string `json:"postalCode,omitempty"`
	TownName              string `json:"townName,omitempty"`
	FullAddress           string `json:"fullAddress,omitempty"`
	ResidentAtAddressFrom string `json:"residentAtAddressFrom,omitempty"`
	ResidentAtAddressTo   string `json:"residentAtAddressTo,omitempty"`
	TimeAtAddress         string `json:"timeAtAddress,omitempty"`
}

type TDevice struct {
	AnonymizerInUseFlag    string  `json:"anonymizerInUseFlag,omitempty"`
	AreaCode               string  `json:"areaCode,omitempty"`
	BrowserType            string  `json:"browserType,omitempty"`
	BrowserVersion         string  `json:"browserVersion"`
	City                   string  `json:"city,omitempty"`
	ClientTimezone         string  `json:"clientTimezone,omitempty"`
	ContinentCode          string  `json:"continentCode,omitempty"`
	CookieId               string  `json:"cookieId"`
	CountryCode            string  `json:"countryCode,omitempty"`
	CountryName            string  `json:"countryName,omitempty"`
	DeviceFingerprint      string  `json:"deviceFingerprint,omitempty"`
	DeviceIMEI             string  `json:"deviceIMEI,omitempty"`
	DeviceName             string  `json:"deviceName,omitempty"`
	FlashPluginPresent     string  `json:"flashPluginPresent,omitempty"`
	HttpHeader             string  `json:"httpHeader,omitempty"`
	IpAddress              string  `json:"ipAddress,omitempty"`
	IpAddressV4            string  `json:"ipAddressV4,omitempty"`
	IpAddressV6            string  `json:"ipAddressV6,omitempty"`
	MetroCode              string  `json:"metroCode,omitempty"`
	MimeTypesPresent       string  `json:"mimeTypesPresent,omitempty"`
	MobileNumberDeviceLink string  `json:"mobileNumberDeviceLink,omitempty"`
	NetworkCarrier         string  `json:"networkCarrier,omitempty"`
	OS                     string  `json:"oS,omitempty"`
	PostalCode             string  `json:"postalCode,omitempty"`
	ProxyDescription       string  `json:"proxyDescription,omitempty"`
	ProxyType              string  `json:"proxyType,omitempty"`
	Region                 string  `json:"region,omitempty"`
	ScreenResolution       string  `json:"screenResolution,omitempty"`
	SessionLatitude        float64 `json:"sessionLatitude,omitempty"`
	SessionLongitude       float64 `json:"sessionLongitude,omitempty"`
	Timestamp              string  `json:"timestamp,omitempty"`
	Type                   string  `json:"type,omitempty"`
	UserAgentString        string  `json:"userAgentString,omitempty"`
}

type TVerificationType struct {
	Aa                           string `json:"aa,omitempty"`
	AccountDigitalSignature      string `json:"accountDigitalSignature,omitempty"`
	AuthenticationToken          string `json:"authenticationToken,omitempty"`
	Avs                          string `json:"avs,omitempty"`
	Biometry                     string `json:"biometry,omitempty"`
	CardholderIdentificationData string `json:"cardholderIdentificationData,omitempty"`
	CryptogramVerification       string `json:"cryptogramVerification,omitempty"`
	CscVerification              string `json:"cscVerification,omitempty"`
	Cvv                          string `json:"cvv,omitempty"`
	OfflinePIN                   string `json:"offlinePIN,omitempty"`
	OneTimePassword              string `json:"oneTimePassword,omitempty"`
	OnlinePIN                    string `json:"onlinePIN,omitempty"`
	Other                        string `json:"other,omitempty"`
	PaperSignature               string `json:"paperSignature,omitempty"`
	PassiveAuthentication        string `json:"passiveAuthentication,omitempty"`
	Password                     string `json:"password,omitempty"`
	ThreeDS                      string `json:"threeDS,omitempty"`
	TokenAuthentication          string `json:"tokenAuthentication,omitempty"`
}

type TName struct {
	FullName     string `json:"fullName,omitempty"`
	GivenName    string `json:"givenName,omitempty"`
	MiddleName   string `json:"middleName,omitempty"`
	NamePrefix   string `json:"namePrefix,omitempty"`
	PreviousName string `json:"previousName,omitempty"`
	Surname      string `json:"surname,omitempty"`
}

type TDecoration struct {
	Decoration []string `json:"decoration,omitempty"`
}

type TTenant struct {
	Name             string `json:"name"`
	TenantId         string `json:"tenantid"`
	BranchRangeStart int    `json:"branchrangestart"`
	BranchRangeEnd   int    `json:"branchrangeend"`
	Bicfi            string `json:"bicfi,omitempty"`
}

type TAccount struct {
	Id            string   `json:"id"`
	Name          TName    `json:"name"`
	TenantId      string   `json:"tenantid,omitempty"`
	AccountNumber string   `json:"accountnumber,omitempty"`
	Address       TAddress `json:"address,omitempty"`
}

type TPSeed struct {
	Decoration               []string `json:"decoration,omitempty"`
	Direction                []string `json:"direction,omitempty"`
	TransactionTypesNrt      []string `json:"transactionTypesNrt,omitempty"`
	TransactionTypesRt       []string `json:"transactionTypesRt,omitempty"`
	ChargeBearers            []string `json:"chargeBearers,omitempty"`
	RemittanceLocationMethod []string `json:"remittanceLocationMethod,omitempty"`
	SettlementMethod         []string `json:"settlementMethod,omitempty"`
	VerificationResult       []string `json:"verificationResult,omitempty"`
	PaymentFrequency         []string `json:"paymentFrequency,omitempty"`

	Tenants TTenants `json:"tenants,omitempty"`
	//	TenantsRT  []TTenant `json:"tenantsrt,omitempty"`
	//	TenantsNRT []TTenant `json:"tenantsnrt,omitempty"`

	Accounts TAccounts `json:"accounts,omitempty"`
	// GoodAccounts []TAccount `json:"goodAccounts,omitempty"`
	// BadAccounts  []TAccount `json:"badAccounts,omitempty"`

}

type TTenants struct {
	RT  []TTenant `json:"rt,omitempty"`
	NRT []TTenant `json:"nrt,omitempty"`
}

type TAccounts struct {
	Good []TAccount `json:"good,omitempty"`
	Bad  []TAccount `json:"bad,omitempty"`
}

type TPaymentNRT = struct {
	AccountAgentId                 string  `json:"accountAgentId,omitempty"`
	AccountEntityId                string  `json:"accountEntityId,omitempty"`
	AccountAgentName               string  `json:"accountAgentName,omitempty"`
	AccountId                      string  `json:"accountId,omitempty"`
	Amount                         TAmount `json:"amount,omitempty"`
	ChargeBearer                   string  `json:"chargeBearer,omitempty"` // hardcode SLEV
	CounterpartyAgentId            string  `json:"counterpartyAgentId,omitempty"`
	CounterpartyAgentName          string  `json:"counterpartyAgentName,omitempty"`
	CounterpartyEntityId           string  `json:"counterpartyEntityId,omitempty"`
	CounterpartyId                 string  `json:"counterpartyId,omitempty"`
	CreationDate                   string  `json:"creationDate,omitempty"`
	DestinationCountry             string  `json:"destinationCountry,omitempty"` // hardcode ZAF
	Direction                      string  `json:"direction,omitempty"`
	EventId                        string  `json:"eventId,omitempty"`
	EventTime                      string  `json:"eventTime,omitempty"`
	EventType                      string  `json:"eventType,omitempty"`
	FromFIBranchId                 string  `json:"fromFIBranchId,omitempty"`
	FromId                         string  `json:"fromId,omitempty"`
	LocalInstrument                string  `json:"localInstrument,omitempty"`
	MsgStatus                      string  `json:"msgStatus,omitempty"` // hardcode Success
	MsgType                        string  `json:"msgType,omitempty"`   // hardcode RCCT
	NumberOfTransactions           int     `json:"numberOfTransactions,omitempty"`
	PaymentClearingSystemReference string  `json:"paymentClearingSystemReference,omitempty"`
	PaymentMethod                  string  `json:"paymentMethod,omitempty"` // hardcode TRF
	PaymentReference               string  `json:"paymentReference,omitempty"`
	RemittanceId                   string  `json:"remittanceId,omitempty"`
	RequestExecutionDate           string  `json:"requestExecutionDate,omitempty"`
	SchemaVersion                  int     `json:"schemaVersion,omitempty"`
	SettlementClearingSystemCode   string  `json:"settlementClearingSystemCode,omitempty"` // hardcode RTC
	SettlementDate                 string  `json:"settlementDate,omitempty"`
	SettlementMethod               string  `json:"settlementMethod,omitempty"` // hardcode CLRG
	TenantId                       string  `json:"tenantId,omitempty"`
	ToFIBranchId                   string  `json:"toFIBranchId,omitempty"`
	ToId                           string  `json:"toId,omitempty"`
	TotalAmount                    TAmount `json:"totalAmount,omitempty"`
	TransactionId                  string  `json:"transactionId,omitempty"`
	TransactionType                string  `json:"transactionType,omitempty"`
	// FS Modifications required for these fields
	CounterpartyIDaccounttype  string `json:"counterpartyIDaccounttype,omitempty"`
	AccountIDaccounttype       string `json:"accountIDaccounttype,omitempty"`
	Usercode                   string `json:"usercode,omitempty"` // Hardcode RTC0000
	CounterpartyIDSuspenseflag string `json:"counterpartyIDSuspenseflag,omitempty"`
	AccountIDSuspenseflag      string `json:"accountIDSuspenseflag,omitempty"`
	EntryClass                 string `json:"entry,omitempty"` // hardcode RTC42
	UnpaidReasonCode           string `json:"unpaidReasonCode,omitempty"`
}

type TPaymentRT struct {
	AccountAddress                      TAddress          `json:"accountAddress,omitempty"`
	AccountAgentAddress                 TAddress          `json:"accountAgentAddress,omitempty"`
	AccountAgentId                      string            `json:"accountAgentId,omitempty"`
	AccountAgentName                    string            `json:"accountAgent,omitempty"`
	AccountBalanceAfter                 TAmount           `json:"accountBalanceAfter,omitempty"`
	AccountBICFI                        string            `json:"accountBICFI,omitempty"`
	AccountCustomerEntityId             string            `json:"accountCustomerEntityId,omitempty"`
	AccountCustomerId                   string            `json:"accountCustomerId,omitempty"`
	AccountDomain                       string            `json:"accountDomain,omitempty"`
	AccountEntityId                     string            `json:"accountEntity"`
	AccountId                           string            `json:"accountId,omitempty"`
	AccountName                         TName             `json:"accountName,omitempty"`
	AccountProxyType                    string            `json:"accountProxyType,omitempty"`
	AccountProxyId                      string            `json:"accountProxyId,omitempty"`
	AccountProxyEntityId                string            `json:"accountProxyEntityId,omitempty"`
	Amount                              TAmount           `json:"amount,omitempty"`
	CardEntityId                        string            `json:"cardEntityId,omitempty"`
	CardId                              string            `json:"cardId,omitempty"`
	Channel                             string            `json:"channel,omitempty"`
	ChargeBearer                        string            `json:"chargeBearer,omitempty"`
	CounterpartyAddress                 TAddress          `json:"CounterpartyAddress,omitempty"`
	CounterpartyAgentAddress            TAddress          `json:"counterpartyAgentAddress,omitempty"`
	CounterpartyAgentId                 string            `json:"counterpartyAgentId,omitempty"`
	CounterpartyAgentName               string            `json:"counterPartyAgentName"`
	CounterpartyEntityId                string            `json:"counterPartyEntityId,omitempty"`
	CounterpartyBICFI                   string            `json:"counterpartyBICFI,omitempty"`
	CounterpartyCustomerId              string            `json:"counterpartyCustomerId,omitempty"`
	CounterpartyCustomerEntityId        string            `json:"counterpartyCustomerEntityId,omitempty"`
	CounterpartyDomain                  string            `json:"counterpartyDomain,omitempty"`
	CounterpartyId                      string            `json:"counterPartyId,omitempty"`
	CounterpartyName                    TName             `json:"counterPartyName,omitempty"`
	CounterpartyNumber                  string            `json:"counterpartyNumber,omitempty"`
	CounterpartyProxyId                 string            `json:"counterpartyProxyId,omitempty"`
	CounterpartyProxyEntityId           string            `json:"counterpartyProxyEntityId,omitempty"`
	CounterpartyProxyType               string            `json:"counterpartyProxyType,omitempty"`
	CounterpartyType                    string            `json:"counterpartyType,omitempty"`
	CreationDate                        string            `json:"creationDate"`
	CustomerEntityId                    string            `json:"customerEntityId,omitempty"`
	CustomerId                          string            `json:"customerId,omitempty"`
	CustomerType                        string            `json:"customerType,omitempty"`
	DecorationId                        TDecoration       `json:"decorationId,omitempty"`
	DestinationCountry                  string            `json:"destinationCountry,omitempty"`
	Device                              TDevice           `json:"device,omitempty"`
	DeviceEntityId                      string            `json:"deviceEntityId,omitempty"`
	DeviceId                            string            `json:"deviceId,omitempty"`
	Direction                           string            `json:"direction,omitempty"`
	EventId                             string            `json:"eventId,omitempty"`
	EventTime                           string            `json:"eventTime,omitempty"`
	EventType                           string            `json:"eventType,omitempty"`
	FinalPaymentDate                    string            `json:"finalPaymentDate,omitempty"`
	FromFIBranchId                      string            `json:"fromFIBranchId,omitempty"`
	FromId                              string            `json:"fromId,omitempty"`
	InstructedAgentAddress              TAddress          `json:"instructedAgentAddress,omitempty"`
	InstructedAgentId                   string            `json:"instructedAgentId,omitempty"`
	InstructedAgentName                 string            `json:"instructedAgentName,omitempty"`
	InstructingAgentAddress             TAddress          `json:"instructingAgentAddress,omitempty"`
	InstructingAgentId                  string            `json:"instructingAgentId,omitempty"`
	InstructingAgentName                string            `json:"instructingAgentName,omitempty"`
	IntermediaryAgent1AccountId         string            `json:"intermediaryAgent1AccountId,omitempty"`
	IntermediaryAgent1Address           TAddress          `json:"intermeduartAgent1Address,omitempty"`
	IntermediaryAgent1Id                string            `json:"intermediaryAgent1Id,omitempty"`
	IntermediaryAgent1Name              string            `json:"intermediaryAgent1Name,omitempty"`
	IntermediaryAgent2AccountId         string            `json:"intermediaryAgent2AccountId,omitempty"`
	IntermediaryAgent2Address           TAddress          `json:"intermeduartAgent2Address,omitempty"`
	IntermediaryAgent2Id                string            `json:"intermediaryAgent2Id,omitempty"`
	IntermediaryAgent2Name              string            `json:"intermediaryAgent2Name,omitempty"`
	IntermediaryAgent3AccountId         string            `json:"intermediaryAgent3AccountId,omitempty"`
	IntermediaryAgent3Address           TAddress          `json:"intermeduartAgent3Address,omitempty"`
	IntermediaryAgent3Id                string            `json:"intermediaryAgent3Id,omitempty"`
	IntermediaryAgent3Name              string            `json:"intermediaryAgent3Name,omitempty"`
	LocalInstrument                     string            `json:"localInstrument,omitempty"`
	MsgStatus                           string            `json:"msgStatus,omitempty"`
	MsgStatusReason                     string            `json:"msgStatusReason,omitempty"`
	MsgType                             string            `json:"msgType,omitempty"`
	NumberOfTransactions                int               `json:"numberOfTransactions,omitempty"`
	PaymentClearingSystemReference      string            `json:"paymentClearingSystemReference,omitempty"`
	PaymentFrequency                    string            `json:"paymentFrequency,omitempty"`
	PaymentMethod                       string            `json:"paymentMethod,omitempty"`
	PaymentReference                    string            `json:"paymentReference,omitempty"`
	RemittanceId                        string            `json:"payRemittanceIdmentMethod,omitempty"`
	RemittanceLocationElectronicAddress string            `json:"remittanceLocationElectronicAddress,omitempty"`
	RemittanceLocationMethod            string            `json:"remittanceLocationMethod,omitempty"`
	RequestExecutionDate                string            `json:"paymentMRequestExecutionDateethod,omitempty"`
	SchemaVersion                       int               `json:"schemaVersion,omitempty"`
	ServiceLevelCode                    string            `json:"serviceLevelCode,omitempty"`
	SettlementClearingSystemCode        string            `json:"settlementClearingSystemCode,omitempty"`
	SettlementDate                      string            `json:"settlementDate,omitempty"`
	SettlementMethod                    string            `json:"settlementMethod,omitempty"`
	TenantId                            string            `json:"tenantId,omitempty"`
	ToFIBranchId                        string            `json:"toFIBranchId,omitempty"`
	ToId                                string            `json:"toId,omitempty"`
	TotalAmount                         TAmount           `json:"totalAmount,omitempty"`
	TransactionId                       string            `json:"transactionId,omitempty"`
	TransactionType                     string            `json:"transactionType,omitempty"`
	UltimateAccountAddress              TAddress          `json:"ultimateAccountAddress,omitempty"`
	UltimateAccountId                   string            `json:"ultimateAccountId,omitempty"`
	UltimateAccountName                 TName             `json:"ultimateAccountName,omitempty"`
	UltimateCounterpartyAddress         TAddress          `json:"ultimateCounterpartyAddress,omitempty"`
	UltimateCounterpartyId              string            `json:"ultimateCounterpartyId,omitempty"`
	UltimateCounterpartyName            TName             `json:"ultimateCounterpartyName,omitempty"`
	UnstructuredRemittanceInformation   string            `json:"unstructuredRemittanceInformation,omitempty"`
	VerificationResult                  string            `json:"verificationResult,omitempty"`
	VerificationType                    TVerificationType `json:"verificationType,omitempty"`
}
