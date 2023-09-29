/*****************************************************************************
*
*	File			: producer.go
*
* 	Created			: 27 Aug 2021
*
*	Description		: Creates Fake JSON Payload, originally posted onto Kafka topic, modified to post onto FeatureSpace event API endpoint
*
*	Modified		: 27 Aug 2021	- Start
*					: 24 Jan 2023   - Mod for applab sandbox
*					: 20 Feb 2023	- repackaged for TFM 2.0 load generation, we're creating fake FS Payment events messages onto
*					:				- FeatureSpace API endpoint
*
*					: 09 May 2023	- Enhancement: we can now read a directory of files (JSON structured), each as a paymentEvent and
*					:				- post the contents onto the API endpoint, allowing for pre designed schenarios to be used/tested
*					:				- NOTE, Original idea/usage was posting payloads onto Kafka topic, thus the fake data gen,
*					:				- With the new usage of reading scenario files allot of the bits below can be removed...
*					:				- Also removed the Prometheus instrumentation form this version as it will mostly be used to input/post a
*					: 				- coupld of files, not a big batch that needs to be timed/measure.
*
*					: 12 May 2023 	- Moved all environment variables from .exps environment export file to *_app.json file, this works better with a App
*									- destined for a desktop vs a app for a K8S server which prefers environment vars.
*					:				- https://onexlab-io.medium.com/golang-config-file-best-practise-d27d6a97a65a
*
*					: 24 May 2023	- Moved the seed data to a json structure seed.json thats ready in and then utilised instead of the seed package
*					:				- Modifying the payment structure to be aligned with Kiveshan's excell spread sheet, Makes for better down the line
*					:				- Fake data generation.
*					:				- also introduced the min and max transaction values.
*					:				- This required that I split the paymentNRT andpaymentRT into 2 dif functions as they are VERY different.
*
*					: 15 Jun 2023	- Expanded the tenant type to include BranchRangeStart/BranchRangeEnd and Bicfi, updated seed file, using
*					:				- possible brang number for toFIBrance and fromFIBrance
*
*					: 18 Jun 2023	- refactured so that we first pick a random payer and payerr account and then use the
*					:				- acocunts TenantId to fetch the tenant related information.
*					:				- staggered/split the tenants into RT and NRT as sub section under Tenant.
*					:				- staggered/split the accounts into Good: [] and Bad: [] as sub sections under Accounts
*					:				- Added additonal fields as determined through schema workshops to TenantRT message structure
*					:				- Removed Good and BadEntities
*					:
*					: 21 Aug 2023	- Change the fake data generation to be financial transaction based. if RT mode then it will generate a RT and NRT
*					: (2.0)			- if non RT then it creates 2 x NRT based payments.
*					:				- Add capability on the file input functionality to read/post AddProxy files.
*
*					: 28 Sept 2023	- For processing from file, collapsing inbound and outbound event into a single file (2 json payloads inside a array),
*					: 				- to enable us to refresh the trasanctionid so that both events match.
*
*
*	By				: George Leonard (georgelza@gmail.com)
*
*	jsonformatter 	: https://jsonformatter.curiousconcept.com/#
*
*****************************************************************************/

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/TylerBrock/colorjson"
	"github.com/tkanos/gonfig"
	glog "google.golang.org/grpc/grpclog"

	// My Types/Structs/functions
	"cmd/types"

	"crypto/tls"
	"crypto/x509"
	// Filter JSON array
)

var (
	grpcLog  glog.LoggerV2
	validate = validator.New()
	varSeed  types.TPSeed
	vGeneral types.Tp_general
	pathSep  = string(os.PathSeparator)
)

func init() {

	// Keeping it very simple
	grpcLog = glog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)

	grpcLog.Infoln("###############################################################")
	grpcLog.Infoln("#")
	grpcLog.Infoln("#   Project   : TFM 2.1")
	grpcLog.Infoln("#")
	grpcLog.Infoln("#   Comment   : FeatureSpace Scenario Publisher / Fake Data Generator")
	grpcLog.Infoln("#             : To be Event/Alert publisher")
	grpcLog.Infoln("#")
	grpcLog.Infoln("#   By        : George Leonard (georgelza@gmail.com)")
	grpcLog.Infoln("#")
	grpcLog.Infoln("#   Date/Time :", time.Now().Format("2006-01-02 15:04:05"))
	grpcLog.Infoln("#")
	grpcLog.Infoln("###############################################################")
	grpcLog.Infoln("")
	grpcLog.Infoln("")

}

func loadConfig(params ...string) types.Tp_general {

	var err error

	vGeneral := types.Tp_general{}
	env := "dev"
	if len(params) > 0 { // Input environment was specified, so lets use it
		env = params[0]
		grpcLog.Info("*")
		grpcLog.Info("* Called with ", env)
		grpcLog.Info("*")

	}

	vGeneral.CurrentPath, err = os.Getwd()
	if err != nil {
		grpcLog.Fatalln("Problem retrieving current path: %s", err)

	}

	vGeneral.OSName = runtime.GOOS

	fileName := fmt.Sprintf("%s%s%s_app.json", vGeneral.CurrentPath, pathSep, env)
	err = gonfig.GetConf(fileName, &vGeneral)
	if err != nil {
		grpcLog.Fatalln("Error Reading Config File: ", err)

	} else {

		vHostname, err := os.Hostname()
		if err != nil {
			grpcLog.Fatalln("Can't retrieve hostname %s", err)

		}
		vGeneral.Hostname = vHostname

		vGeneral.Cert_file = fmt.Sprintf("%s%s%s%s%s", vGeneral.CurrentPath, pathSep, vGeneral.Cert_dir, pathSep, vGeneral.Cert_file)
		vGeneral.Cert_key = fmt.Sprintf("%s%s%s%s%s", vGeneral.CurrentPath, pathSep, vGeneral.Cert_dir, pathSep, vGeneral.Cert_key)

		if vGeneral.Json_to_file == 1 {
			vGeneral.Output_path = fmt.Sprintf("%s%s%s", vGeneral.CurrentPath, pathSep, vGeneral.Output_path)

		} else {
			vGeneral.Output_path = ""

		}

		if vGeneral.Json_from_file == 1 {
			vGeneral.Input_path = fmt.Sprintf("%s%s%s", vGeneral.CurrentPath, pathSep, vGeneral.Input_path)

		} else {
			vGeneral.Input_path = ""

		}

		vGeneral.SeedFile = fmt.Sprintf("%s%s%s", vGeneral.CurrentPath, pathSep, vGeneral.SeedFile)

	}

	if vGeneral.EchoConfig == 1 {
		printConfig(vGeneral)
	}

	if vGeneral.Debuglevel > 0 {
		grpcLog.Infoln("*")
		grpcLog.Infoln("* Config:")
		grpcLog.Infoln("* Current path:", vGeneral.CurrentPath)
		grpcLog.Infoln("* Config File :", fileName)
		grpcLog.Infoln("*")

	}

	return vGeneral
}

func loadSeed(fileName string) types.TPSeed {

	var vSeed types.TPSeed

	err := gonfig.GetConf(fileName, &vSeed)
	if err != nil {
		grpcLog.Fatalln("Error Reading Seed File: ", err)

	}

	v, err := json.Marshal(vSeed)
	if err != nil {
		grpcLog.Fatalln("Marchalling error: ", err)
	}

	if vGeneral.EchoSeed == 1 {
		prettyJSON(string(v))

	}

	if vGeneral.Debuglevel > 0 {
		grpcLog.Infoln("*")
		grpcLog.Infoln("* Seed :")
		grpcLog.Infoln("* Current path:", vGeneral.CurrentPath)
		grpcLog.Infoln("* Seed File :", vGeneral.SeedFile)
		grpcLog.Infoln("*")

	}

	return vSeed
}

func printConfig(vGeneral types.Tp_general) {

	grpcLog.Info("****** General Parameters *****")
	grpcLog.Info("*")
	grpcLog.Info("* Hostname is\t\t\t", vGeneral.Hostname)
	grpcLog.Info("* OS is \t\t\t", vGeneral.OSName)
	grpcLog.Info("*")
	grpcLog.Info("* Debug Level is\t\t", vGeneral.Debuglevel)
	grpcLog.Info("*")
	grpcLog.Info("* Sleep Duration is\t\t", vGeneral.Sleep)
	grpcLog.Info("* Test Batch Size is\t\t", vGeneral.Testsize)
	grpcLog.Info("* Seed File is\t\t\t", vGeneral.SeedFile)
	grpcLog.Info("* Echo Seed is\t\t", vGeneral.EchoSeed)
	grpcLog.Info("* Echo JSON is\t\t", vGeneral.Echojson)
	grpcLog.Info("*")

	grpcLog.Info("* Current Path is \t\t", vGeneral.CurrentPath)
	grpcLog.Info("* Cert file is\t\t", vGeneral.Cert_file)
	grpcLog.Info("* Cert key is\t\t\t", vGeneral.Cert_key)

	grpcLog.Info("* HTTP JSON POST URL is\t", vGeneral.Httpposturl)

	grpcLog.Info("* Read JSON from file is\t", vGeneral.Json_from_file) // if 0 then we create fake data else
	grpcLog.Info("* Input path is\t\t", vGeneral.Input_path)            // if 1 then read files from input_path
	grpcLog.Info("* Data Gen Mode is\t\t", vGeneral.Datamode)           // if we're creating fake data then who's the input system

	grpcLog.Info("* Output JSON to file is\t", vGeneral.Json_to_file)
	grpcLog.Info("* Output path is\t\t", vGeneral.Output_path)

	grpcLog.Info("* MinTransactionValue is\tR ", vGeneral.MinTransactionValue)
	grpcLog.Info("* MaxTransactionValue is\tR ", vGeneral.MaxTransactionValue)

	grpcLog.Info("*")
	grpcLog.Info("*******************************")

	grpcLog.Info("")

}

// Pretty Print JSON string
func prettyJSON(ms string) {

	var obj map[string]interface{}

	json.Unmarshal([]byte(ms), &obj)

	// Make a custom formatter with indent set
	f := colorjson.NewFormatter()
	f.Indent = 4

	// Marshall the Colorized JSON
	result, _ := f.Marshal(obj)
	fmt.Println(string(result))

}

// We're driving this from an account perspective. so lets find the tenants information as per the account def
// Helper Func - Find the tenandId for the Bank sending or receiving the funds.
func findTenant(tenants []types.TTenant, filter string) (ret types.TTenant, err error) {

	for _, item := range tenants {
		if item.TenantId == filter {
			ret = item
			return ret, nil
		}
	}
	err = errors.New("tenant not found")

	return ret, err
}

// - FAKE Data generation
// - Build a fin transaction.
// 1. are we doing "hist" or rpp
// 2. Step 1 select debitor
// 3. Step 2 find creditor
// 4. build outbound record
// 5. build inbound record
// 6. send
// 7. print to file
// 8. print to screen

func constructFinTransaction() (t_OutboundPayment map[string]interface{}, t_InboundPayment map[string]interface{}) {

	var paymentStream string
	var TransactionTypeNRT string
	var localInstrument string
	var jDebtorAccount types.TAccount
	var jCreditorAccount types.TAccount
	var jDebtorBank types.TTenant
	var jCreditorBank types.TTenant
	var DebtorFIBranchId string
	var CreditorFIBranchId string

	var txnId string
	var eventTime string
	var requestExecutionDate string
	var settlementDate string
	var paymentClearingSystemReference string
	var paymentRef string
	var remittanceId string
	var msgType string

	// We just using gofakeit to pad the json document size a bit.
	//
	// https://github.com/brianvoe/gofakeit
	// https://pkg.go.dev/github.com/brianvoe/gofakeit

	gofakeit.Seed(0)

	nAmount := gofakeit.Price(vGeneral.MinTransactionValue, vGeneral.MaxTransactionValue)
	t_amount := &types.TAmount{
		BaseCurrency: "zar",
		BaseValue:    nAmount,
		Currency:     "zar",
		Value:        nAmount,
	}

	// Determine how many accounts we have in seed file, it's all about them...
	// and build the 2 structures from that viewpoint
	accountsCount := len(varSeed.Accounts.Good) - 1
	nDebtorAccount := gofakeit.Number(0, accountsCount)
	nCreditorAccount := gofakeit.Number(0, accountsCount)

	// check to make sure the 2 are not the same, if they are, redo...
	if nDebtorAccount == nCreditorAccount {
		nCreditorAccount = gofakeit.Number(0, accountsCount)
	}

	// find the debtor and creditor record to use
	jDebtorAccount = varSeed.Accounts.Good[nDebtorAccount]
	jCreditorAccount = varSeed.Accounts.Good[nCreditorAccount]

	if vGeneral.Datamode == "hist" {
		paymentStreamCount := len(varSeed.PaymentStreamNRT) - 1
		nPaymentStream := gofakeit.Number(0, paymentStreamCount)
		paymentStream = varSeed.PaymentStreamNRT[nPaymentStream]

		if paymentStream == "EFT" {
			TransactionTypesCount := len(varSeed.TransactionTypeNRT.EFT) - 1
			nTransactionTypesCount := gofakeit.Number(0, TransactionTypesCount)
			TransactionTypeNRT = varSeed.TransactionTypeNRT.EFT[nTransactionTypesCount].Name
			msgType = "EFT"

		} else if paymentStream == "ACD" {
			TransactionTypesCount := len(varSeed.TransactionTypeNRT.ACD) - 1
			nTransactionTypesCount := gofakeit.Number(0, TransactionTypesCount)
			TransactionTypeNRT = varSeed.TransactionTypeNRT.ACD[nTransactionTypesCount].Name
			msgType = "900000"

		} else if paymentStream == "RTC" {
			TransactionTypesCount := len(varSeed.TransactionTypeNRT.RTC) - 1
			nTransactionTypesCount := gofakeit.Number(0, TransactionTypesCount)
			TransactionTypeNRT = varSeed.TransactionTypeNRT.RTC[nTransactionTypesCount].Name
			msgType = "RTCCT"

		}
		jDebtorBank, _ = findTenant(varSeed.Tenants.NRT, jDebtorAccount.TenantId)
		jCreditorBank, _ = findTenant(varSeed.Tenants.NRT, jCreditorAccount.TenantId)

		localInstrumentCount := len(varSeed.LocalInstrument.HIST) - 1
		nlocalInstrumentCount := gofakeit.Number(0, localInstrumentCount)
		localInstrument = varSeed.LocalInstrument.HIST[nlocalInstrumentCount].Name

	} else { // RPP
		jDebtorBank, _ = findTenant(varSeed.Tenants.RT, jDebtorAccount.TenantId)
		jCreditorBank, _ = findTenant(varSeed.Tenants.RT, jCreditorAccount.TenantId)

		localInstrumentCount := len(varSeed.LocalInstrument.RPP) - 1
		nlocalInstrumentCount := gofakeit.Number(0, localInstrumentCount)
		localInstrument = varSeed.LocalInstrument.RPP[nlocalInstrumentCount].Value

	}

	// find FIB Id for the debtor and creditor bank
	DebtorFIBranchId = strconv.Itoa(gofakeit.Number(jDebtorBank.BranchRangeStart, jDebtorBank.BranchRangeEnd))
	CreditorFIBranchId = strconv.Itoa(gofakeit.Number(jCreditorBank.BranchRangeStart, jCreditorBank.BranchRangeEnd))

	txnId = uuid.New().String()
	eventTime = time.Now().Format("2006-01-02T15:04:05")

	requestExecutionDate = time.Now().Format("2006-01-02")
	settlementDate = time.Now().Format("2006-01-02")
	paymentClearingSystemReference = uuid.New().String()
	paymentRef = paymentClearingSystemReference
	remittanceId = paymentClearingSystemReference

	if vGeneral.Datamode == "hist" {

		// 2 x NRT records/events

		t_OutboundPayment = map[string]interface{}{
			"accountAgentId":                 jDebtorBank.TenantId,
			"accountId":                      jDebtorAccount.AccountNumber,
			"accountIdCode":                  jDebtorAccount.AccountIDCode, // Type of Account
			"accountNumber":                  jDebtorAccount.AccountNumber,
			"amount":                         t_amount,
			"chargeBearer":                   "SLEV",
			"counterpartyAgentId":            jCreditorBank.TenantId,
			"counterpartyId":                 jCreditorAccount.AccountNumber,
			"counterpartyIdCode":             jCreditorAccount.AccountIDCode, // Type of Account
			"counterpartyNumber":             jCreditorAccount.AccountNumber,
			"creationDate":                   time.Now().Format("2006-01-02T15:04:05"),
			"destinationCountry":             "ZAF",
			"direction":                      "outbound",
			"eventId":                        uuid.New().String(),
			"eventTime":                      eventTime,
			"eventType":                      "paymentNRT",
			"fromFIBranchId":                 DebtorFIBranchId,
			"fromId":                         jDebtorBank.TenantId,
			"localInstrument":                localInstrument, // aka Record ID's
			"msgStatus":                      "Settlement",
			"msgType":                        msgType,
			"msgStatusReason":                "JNL_ACQ.responseCode",
			"numberOfTransactions":           1,
			"paymentClearingSystemReference": paymentClearingSystemReference,
			"paymentMethod":                  "TRF",
			"paymentReference":               paymentRef,
			"remittanceId":                   remittanceId,
			"requestExecutionDate":           requestExecutionDate,
			"schemaVersion":                  1,
			"settlementClearingSystemCode":   paymentStream,
			"settlementDate":                 settlementDate,
			"settlementMethod":               "CLRG",
			"tenantId":                       jDebtorAccount.TenantId,
			"toFIBranchId":                   CreditorFIBranchId,
			"toId":                           jCreditorBank.TenantId,
			"totalAmount":                    t_amount,
			"transactionId":                  txnId,
			"transactionType":                TransactionTypeNRT,
			"verificationResult":             "SUCC",
			"usercode":                       "0000",
		}

		t_InboundPayment = map[string]interface{}{
			"accountAgentId":                 jCreditorBank.TenantId,
			"accountId":                      jCreditorAccount.AccountNumber,
			"accountIdCode":                  jCreditorAccount.AccountIDCode, // Type of Account
			"accountNumber":                  jCreditorAccount.AccountNumber,
			"amount":                         t_amount,
			"chargeBearer":                   "SLEV",
			"counterpartyAgentId":            jDebtorBank.TenantId,
			"counterpartyId":                 jDebtorAccount.AccountNumber,
			"counterpartyIdCode":             jDebtorAccount.AccountIDCode, // Type of Account
			"counterpartyNumber":             jDebtorAccount.AccountNumber,
			"creationDate":                   time.Now().Format("2006-01-02T15:04:05"),
			"destinationCountry":             "ZAF",
			"direction":                      "inbound",
			"eventId":                        uuid.New().String(),
			"eventTime":                      eventTime,
			"eventType":                      "paymentNRT",
			"fromFIBranchId":                 DebtorFIBranchId,
			"fromId":                         jDebtorBank.TenantId,
			"localInstrument":                localInstrument, // aka Record ID's
			"msgStatus":                      "Settlement",
			"msgType":                        msgType,
			"msgStatusReason":                "JNL_ACQ.responseCode",
			"numberOfTransactions":           1,
			"paymentClearingSystemReference": paymentClearingSystemReference,
			"paymentMethod":                  "TRF",
			"paymentReference":               paymentRef,
			"remittanceId":                   remittanceId,
			"requestExecutionDate":           requestExecutionDate,
			"schemaVersion":                  1,
			"settlementClearingSystemCode":   paymentStream,
			"settlementDate":                 settlementDate,
			"settlementMethod":               "CLRG",
			"tenantId":                       jCreditorBank.TenantId,
			"toFIBranchId":                   CreditorFIBranchId,
			"toId":                           jCreditorBank.TenantId,
			"totalAmount":                    t_amount,
			"transactionId":                  txnId,
			"transactionType":                TransactionTypeNRT,
			"verificationResult":             "SUCC",
			"usercode":                       "0000",
		}

	} else { // RPP
		// 1 NRT and 1 RT record/event

		// Outbound => paymentNRT
		// Inbound => paymentRT

		// Account => Debtor
		// CounterParty => Creditor

		t_OutboundPayment = map[string]interface{}{
			"accountAgentId":                    jDebtorBank.TenantId,
			"accountId":                         jDebtorAccount.AccountNumber,
			"accountIdCode":                     jDebtorAccount.AccountIDCode, // Type of Account
			"accountNumber":                     jDebtorAccount.AccountNumber,
			"amount":                            t_amount,
			"chargeBearer":                      "SLEV",
			"counterpartyAgentId":               jCreditorBank.TenantId,
			"counterpartyId":                    jCreditorAccount.AccountNumber,
			"counterpartyIdCode":                jCreditorAccount.AccountIDCode, // Type of Account
			"counterpartyNumber":                jCreditorAccount.AccountNumber,
			"creationDate":                      time.Now().Format("2006-01-02T15:04:05"),
			"destinationCountry":                "ZAF",
			"direction":                         "outbound",
			"eventId":                           uuid.New().String(),
			"eventTime":                         eventTime,
			"eventType":                         "paymentNRT",
			"fromFIBranchId":                    DebtorFIBranchId,
			"fromId":                            jDebtorBank.TenantId,
			"localInstrument":                   localInstrument, // Record ID's
			"msgStatus":                         "New",
			"msgType":                           "CRTRF",
			"msgStatusReason":                   "JNL_ACQ.responseCode",
			"numberOfTransactions":              1,
			"paymentClearingSystemReference":    paymentClearingSystemReference,
			"paymentMethod":                     "TRF",
			"paymentReference":                  paymentRef,
			"remittanceId":                      remittanceId,
			"requestExecutionDate":              requestExecutionDate,
			"schemaVersion":                     1,
			"settlementClearingSystemCode":      "RPP",
			"settlementDate":                    settlementDate,
			"settlementMethod":                  "CLRG",
			"tenantId":                          jDebtorAccount.TenantId,
			"toFIBranchId":                      CreditorFIBranchId,
			"toId":                              jCreditorBank.TenantId,
			"totalAmount":                       t_amount,
			"transactionId":                     txnId,
			"transactionType":                   "MTUP",
			"verificationResult":                "SUCC",
			"instructedAgentId":                 jDebtorBank.Bicfi,
			"instructingAgentId":                jCreditorBank.Bicfi,
			"intermediaryAgent1Id":              jDebtorBank.Bicfi,
			"intermediaryAgent2Id":              jCreditorBank.Bicfi,
			"ultimateAccountName":               jCreditorAccount.Name,
			"ultimateCounterpartyName":          jDebtorAccount.Name,
			"unstructuredRemittanceInformation": paymentClearingSystemReference,
		}

		// Account => Creditor
		// CounterParty => Debtor

		chargeBearersCount := len(varSeed.ChargeBearers) - 1
		nChargeBearers := gofakeit.Number(0, chargeBearersCount)
		chargeBearers := varSeed.ChargeBearers[nChargeBearers]

		//settlementMethodCount := len(varSeed.SettlementMethod) - 1
		//nSettlementMethod := gofakeit.Number(0, settlementMethodCount)

		t_InboundPayment = map[string]interface{}{
			"accountAgentId":    jCreditorBank.Bicfi,            // Bank Bicfi
			"accountId":         jCreditorAccount.AccountNumber, // Bank Acc Number
			"accountIdCode":     jCreditorAccount.AccountIDCode,
			"accountNumber":     jCreditorAccount.AccountNumber,
			"accountBICFI":      jCreditorBank.Bicfi,
			"accountProxyId":    jCreditorAccount.ProxyId,
			"accountProxyType":  jCreditorAccount.ProxyType,
			"accountDomain":     jCreditorAccount.ProxyDomain,
			"accountCustomerId": jCreditorAccount.AccountNumber,
			"accountAddress": types.TAddress{ // Payer - Payee Account
				AddressLine1:       jCreditorAccount.Address.AddressLine1, // CTT.creditor.streetName
				AddressLine2:       jCreditorAccount.Address.AddressLine2, // CTT.creditor.buildingNumber and buildingName
				TownName:           jCreditorAccount.Address.TownName,
				CountrySubDivision: jCreditorAccount.Address.CountrySubDivision,
				Country:            jCreditorAccount.Address.Country,
				PostalCode:         jCreditorAccount.Address.PostalCode,
				FullAddress:        jCreditorAccount.Address.FullAddress,
			},
			"accountName": types.TName{
				FullName:   jCreditorAccount.Name.FullName, // CTT.Creditor or CTT.Debtor
				NamePrefix: jCreditorAccount.Name.NamePrefix,
				Surname:    jCreditorAccount.Name.Surname,
			},
			"amount":                 t_amount,
			"chargeBearer":           chargeBearers,
			"counterpartyAgentId":    jDebtorBank.Bicfi, // Counterparty Bank
			"counterpartyId":         jDebtorAccount.AccountNumber,
			"counterpartyIdCode":     jDebtorAccount.AccountIDCode,
			"counterpartyNumber":     jDebtorAccount.AccountNumber,
			"counterpartyBICFI":      jDebtorBank.Bicfi,
			"counterpartyProxyId":    jDebtorAccount.ProxyId,
			"counterpartyProxyType":  jDebtorAccount.ProxyType,
			"counterpartyDomain":     jDebtorAccount.ProxyDomain,
			"counterpartyCustomerId": jDebtorAccount.AccountNumber,
			"counterpartyAddress": types.TAddress{ // Payer - Payee
				AddressLine1:       jDebtorAccount.Address.AddressLine1, // CTT.creditor.streetName
				AddressLine2:       jDebtorAccount.Address.AddressLine2, // CTT.creditor.buildingNumber and buildingName
				TownName:           jDebtorAccount.Address.TownName,
				CountrySubDivision: jDebtorAccount.Address.CountrySubDivision,
				Country:            jDebtorAccount.Address.Country,
				PostalCode:         jDebtorAccount.Address.PostalCode,
				FullAddress:        jDebtorAccount.Address.FullAddress,
			},
			"counterpartyName": types.TName{
				FullName:   jDebtorAccount.Name.FullName, // CTT.Creditor or CTT.Debtor
				NamePrefix: jDebtorAccount.Name.NamePrefix,
				Surname:    jDebtorAccount.Name.Surname,
			},
			"creationDate": time.Now().Format("2006-01-02T15:04:05"),
			//"customerId":                        jCreditorAccount.AccountNumber,
			"destinationCountry":                "ZAF",
			"direction":                         "inbound",
			"eventId":                           uuid.New().String(),
			"eventTime":                         eventTime,
			"eventType":                         "paymentRT",
			"fromId":                            jDebtorBank.TenantId,
			"localInstrument":                   localInstrument, // pick by LocalInstrumentRt, if pay buy proxy then no creditor account number
			"msgStatus":                         "New",
			"msgType":                           "CRTRF",
			"numberOfTransactions":              1,
			"paymentClearingSystemReference":    paymentClearingSystemReference,
			"paymentMethod":                     "TRF", // Hard coded CHK Cheque / TRF Transfer
			"paymentReference":                  paymentRef,
			"requestExecutionDate":              requestExecutionDate,
			"schemaVersion":                     1,
			"settlementClearingSystemCode":      "RPP",
			"settlementDate":                    settlementDate,
			"settlementMethod":                  "CLRG", // hard coded for now... -> varSeed.SettlementMethod[nSettlementMethod],
			"tenantId":                          jCreditorBank.TenantId,
			"toId":                              jCreditorBank.TenantId,
			"transactionId":                     txnId,
			"transactionType":                   "MTUP",
			"verificationResult":                "SUCC",
			"instructedAgentId":                 jDebtorBank.Bicfi,
			"instructingAgentId":                jCreditorBank.Bicfi,
			"intermediaryAgent1Id":              jDebtorBank.Bicfi,
			"intermediaryAgent2Id":              jCreditorBank.Bicfi,
			"ultimateAccountName":               jCreditorAccount.Name,
			"ultimateCounterpartyName":          jDebtorAccount.Name,
			"unstructuredRemittanceInformation": paymentClearingSystemReference,
		}
	}

	return t_OutboundPayment, t_InboundPayment
}

func ReadJSONFile(varRec string) []byte {

	// Let's first read the `config.json` file
	content, err := os.ReadFile(varRec)
	if err != nil {
		grpcLog.Fatalln("Error when opening file: ", err)

	}
	return content
}

func isJSON(content []byte) (isJson bool) {

	var t_Payment interface{}

	err := json.Unmarshal(content, &t_Payment)
	if err != nil {
		isJson = false

	} else {
		isJson = true

	}

	return isJson
}

func contructEventFromJSONFile(varRec string) (t_OutboundPayment map[string]interface{}, t_InboundPayment map[string]interface{}) {

	var objs interface{}

	// read our opened jsonFile as a byte array.
	byteValue, err := os.ReadFile(varRec)
	if err != nil {
		grpcLog.Errorln("ReadFile error ", err)

	}

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal([]byte(byteValue), &objs)
	if err != nil {
		grpcLog.Errorln("Unmarshall error ", err)

	}

	// Ensure that it is an array of objects.
	objArr, ok := objs.([]interface{})
	if !ok {
		grpcLog.Errorln("expected an array of objects ")

	}

	// Handle each object as a map[string]interface{}.
	for i, obj := range objArr {
		obj, ok := obj.(map[string]interface{})
		if !ok {
			grpcLog.Errorln(fmt.Sprintf("expected type map[string]interface{}, got %s", reflect.TypeOf(objArr[i])))

		}
		if i == 0 {
			t_InboundPayment = obj

		}
		if i == 1 {
			t_OutboundPayment = obj

		}
	}

	txnId := uuid.New().String()
	eventTime := time.Now().Format("2006-01-02T15:04:05")
	creationDate := time.Now().Format("2006-01-02T15:04:05")
	requestExecutionDate := time.Now().Format("2006-01-02")
	settlementDate := time.Now().Format("2006-01-02")

	// we update/refresh the eventID & eventTime, to ensure we don't get duplicate (and make it a payment in) id's at POST time
	t_OutboundPayment["eventId"] = uuid.New().String()
	t_InboundPayment["eventId"] = uuid.New().String()

	t_OutboundPayment["transactionId"] = txnId
	t_InboundPayment["transactionId"] = txnId

	t_OutboundPayment["eventTime"] = eventTime
	t_InboundPayment["eventTime"] = eventTime

	t_InboundPayment["creationDate"] = creationDate
	t_OutboundPayment["creationDate"] = creationDate

	if t_InboundPayment["eventType"] == "paymentRT" || t_InboundPayment["eventType"] == "paymentNRT" || t_OutboundPayment["eventType"] == "paymentRT" || t_OutboundPayment["eventType"] == "paymentNRT" {

		t_InboundPayment["requestExecutionDate"] = requestExecutionDate
		t_InboundPayment["settlementDate"] = settlementDate

		t_OutboundPayment["requestExecutionDate"] = requestExecutionDate
		t_OutboundPayment["settlementDate"] = settlementDate
	}

	return t_OutboundPayment, t_InboundPayment
}

// Return list of files located in input_path to be repackaged as JSON payloads and posted onto the API endpoint
func fetchJSONRecords(input_path string) (records map[int]string, count int) {

	count = 0

	m := make(map[int]string)

	// https://yourbasic.org/golang/list-files-in-directory/
	// Use the os.ReadDir function in package os. It returns a sorted slice containing elements of type os.FileInfo.
	files, err := os.ReadDir(input_path)
	if err != nil {
		grpcLog.Errorln("Problem retrieving list of input files: %s", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			m[count] = file.Name()
			count++

		}
	}

	records = m

	return records, count
}

func formatJSON(data []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", "    ")
	if err == nil {
		return out.Bytes(), err
	}
	return data, nil
}
func runLoader(arg string) {

	// Initialize the vGeneral struct variable - This holds our configuration settings.
	vGeneral = loadConfig(arg)

	// Lets get Seed Data from the specified seed file
	varSeed = loadSeed(vGeneral.SeedFile)

	// Create client with Cert once
	// https://stackoverflow.com/questions/38822764/how-to-send-a-https-request-with-a-certificate-golang

	caCert, err := os.ReadFile(vGeneral.Cert_file)
	if err != nil {
		grpcLog.Fatalln("Problem reading (Exiting) :", vGeneral.Cert_file, " Error :", err)

	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(vGeneral.Cert_file, vGeneral.Cert_key)
	if err != nil {
		grpcLog.Fatalln("Problem with LoadX509KeyPair (Exiting) :", vGeneral.Cert_key, " Error :", err)
		os.Exit(1)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            caCertPool,
				InsecureSkipVerify: true, // Self Signed cert
				Certificates:       []tls.Certificate{cert},
			},
		},
	}

	if vGeneral.Debuglevel > 0 {
		grpcLog.Info("**** LETS GO Processing ****")
		grpcLog.Infoln("")

	}

	////////////////////////////////////////////////////////////////////////
	// Lets fecth the records that need to be pushed to the fs api end point
	var todo_count = 0
	var returnedRecs map[int]string
	if vGeneral.Json_from_file == 0 { // Build Fake Record - atm we're generating the data, eventually we might fetch via SQL

		// As we're faking it:
		todo_count = vGeneral.Testsize // this will be recplaced by the value of todo_count from above.

	} else { // Build Record set from data fetched from JSON files in input_path

		// this will return an map of files names, each being a JSON document
		returnedRecs, todo_count = fetchJSONRecords(vGeneral.Input_path)

		if vGeneral.Debuglevel > 1 {
			grpcLog.Infoln("Checking input event files (Making sure it's valid JSON)...")
			grpcLog.Infoln("")

		}

		var weFailed bool = false
		for count := 0; count < todo_count; count++ {

			// Build the entire JSON Payload document, either a fake record or from a input/scenario JSON file

			filename := fmt.Sprintf("%s%s%s", vGeneral.Input_path, pathSep, returnedRecs[count])

			contents := ReadJSONFile(filename)
			if !isJSON(contents) {
				weFailed = true
				grpcLog.Infoln(filename, "=> FAIL")

			} else {
				grpcLog.Infoln(filename, "=> Pass")

			}

		}
		// We're saying we want to use records from input files but not all the files are valid JSON, die die die... ;)
		if weFailed {
			os.Exit(1)
		}
		grpcLog.Infoln("")

	}

	if vGeneral.Debuglevel > 1 {
		grpcLog.Infoln("Number of records to Process", todo_count) // just doing this to prefer a unused error

	}

	// now we loop through the results, building a json document based on FS requirements and then post it, for this code I'm posting to
	// Confluent Kafka topic, but it's easy to change to have it post to a API endpoint.

	// this is to keep record of the total batch run time
	vStart := time.Now()

	for count := 0; count < todo_count; count++ {

		if vGeneral.Debuglevel > 1 {
			grpcLog.Infoln("")
			grpcLog.Infoln("Record                        :", count+1)

		}

		// We're going to time every record and push that to prometheus
		txnStart := time.Now()

		// We creating fake data so we will have 2 events to deal with
		var t_OutboundPayload map[string]interface{}
		var t_InboundPayload map[string]interface{}
		var OutboundBytes []byte
		var InboundBytes []byte
		var Outboundbody []byte
		var tOutboundBody map[string]interface{}
		var Inboundbody []byte
		var tInboundBody map[string]interface{}

		// Build the entire JSON Payload document, either a fake record or from a input/scenario JSON file
		if vGeneral.Json_from_file == 0 { // Build Fake Record

			// They are just to different to have kept in one function, so split them into 2 seperate specific use case functions.
			t_OutboundPayload, t_InboundPayload = constructFinTransaction()

		} else {
			// We're reading data from files, so simply post data per file/payload

			// returnedRecs is a map of file names, each filename is JSON document which contains a FS Payment event,
			// At this point we simply post the contents of the payment onto the end point, and record the response.

			filename := fmt.Sprintf("%s%s%s", vGeneral.Input_path, pathSep, returnedRecs[count])

			if vGeneral.Debuglevel > 2 {
				grpcLog.Infoln("Source Event                  :", filename)

			}
			t_OutboundPayload, t_InboundPayload = contructEventFromJSONFile(filename)
		}

		if vGeneral.Debuglevel > 1 {

			grpcLog.Infoln("transactionId assigned        :", t_OutboundPayload["transactionId"])
			grpcLog.Infoln("eventTime assigned            :", t_OutboundPayload["eventTime"])
			grpcLog.Infoln("creationDate assigned         :", t_OutboundPayload["creationDate"])

			grpcLog.Infoln("Outbound        :")
			grpcLog.Infoln("eventId assigned              :", t_OutboundPayload["eventId"])

			if t_OutboundPayload["eventType"] == "paymentNRT" {
				grpcLog.Infoln("requestExecutionDate assigned :", t_OutboundPayload["requestExecutionDate"])
				grpcLog.Infoln("settlementDate assigned       :", t_OutboundPayload["settlementDate"])

			}

			grpcLog.Infoln("Inbound         :")
			grpcLog.Infoln("eventId assigned              :", t_InboundPayload["eventId"])

			if t_InboundPayload["eventType"] == "paymentRT" || t_InboundPayload["eventType"] == "paymentNRT" {
				grpcLog.Infoln("requestExecutionDate assigned :", t_InboundPayload["requestExecutionDate"])
				grpcLog.Infoln("settlementDate assigned       :", t_InboundPayload["settlementDate"])

			}
		}
		// check to see if we're calling API, if yes then make the call

		OutboundBytes, err = json.Marshal(t_OutboundPayload)
		if err != nil {
			grpcLog.Errorln("Marchalling error: ", err)

		}

		InboundBytes, err = json.Marshal(t_InboundPayload)
		if err != nil {
			grpcLog.Errorln("Marchalling error: ", err)

		}

		if vGeneral.Debuglevel > 1 && vGeneral.Echojson == 1 {
			grpcLog.Infoln("Outbound Payload   	:")
			prettyJSON(string(OutboundBytes))

			grpcLog.Infoln("")

			grpcLog.Infoln("Inbound Payload   	:")
			prettyJSON(string(InboundBytes))
		}

		if vGeneral.Call_fs_api == 1 { // POST to API endpoint

			// We need to do 2 api calls, 1 each for outbound and inbound event.

			// Outbound

			apiOutboundStart := time.Now()

			// https://golangtutorial.dev/tips/http-post-json-go/
			OutboundRequest, err := http.NewRequest("POST", vGeneral.Httpposturl, bytes.NewBuffer(OutboundBytes))
			if err != nil {
				grpcLog.Errorln("http.NewRequest error: ", err)

			}

			OutboundRequest.Header.Set("Content-Type", "application/json; charset=UTF-8")

			OutboundResponse, err := client.Do(OutboundRequest)
			if err != nil {
				grpcLog.Errorln("client.Do error: ", err)

			}
			defer OutboundResponse.Body.Close()

			// Did we call the API, how long did it take, do this here before we write to a file that will impact this time
			if vGeneral.Debuglevel > 0 {
				grpcLog.Infoln("API Call Time         :", time.Since(apiOutboundStart).Seconds(), "Sec")

			}

			// Do something with all the output/response
			Outboundbody, err = io.ReadAll(OutboundResponse.Body)
			if err != nil {
				grpcLog.Errorln("Outbound Body -> io.ReadAll(OutboundResponse.Body) error: ", err)

			}
			if vGeneral.Debuglevel > 2 {
				grpcLog.Infoln("response Payload      :")
				grpcLog.Infoln("response Status       :", OutboundResponse.Status)
				grpcLog.Infoln("response Headers      :", OutboundResponse.Header)
			}

			if OutboundResponse.Status == "204 No Content" {
				// it's a paymentNRT - SUCCESS
				// lets build a body of the header and some additional information
				if vGeneral.Debuglevel > 2 {

					grpcLog.Infoln("response Body         	: <...>NRT")
				}
				tOutboundBody = map[string]interface{}{
					"transactionId":   t_OutboundPayload["transactionId"],
					"eventId":         t_OutboundPayload["eventId"],
					"eventType":       t_OutboundPayload["eventType"],
					"responseStatus":  OutboundResponse.Status,
					"responseHeaders": OutboundResponse.Header,
					"processTime":     time.Now().UTC(),
				}

			} else {
				// oh sh$t, its not a success so now to try and build a body to fault fix later
				if vGeneral.Debuglevel > 2 {

					grpcLog.Infoln("response Body         :", string(Outboundbody))
					grpcLog.Infoln("response Result       : FAILED POST")
				}
				tOutboundBody = map[string]interface{}{
					"transactionId":   t_OutboundPayload["transactionId"],
					"eventId":         t_OutboundPayload["eventId"],
					"eventType":       t_OutboundPayload["eventType"],
					"responseResult":  "FAILED POST",
					"responseBody":    string(Outboundbody),
					"responseStatus":  OutboundResponse.Status,
					"responseHeaders": OutboundResponse.Header,
					"processTime":     time.Now().UTC(),
				}
			}

			// Inbound
			apiInboundStart := time.Now()

			// https://golangtutorial.dev/tips/http-post-json-go/
			InboundRequest, err := http.NewRequest("POST", vGeneral.Httpposturl, bytes.NewBuffer(InboundBytes))
			if err != nil {
				grpcLog.Errorln("http.NewRequest error: ", err)

			}

			InboundRequest.Header.Set("Content-Type", "application/json; charset=UTF-8")

			InboundResponse, err := client.Do(InboundRequest)
			if err != nil {
				grpcLog.Errorln("client.Do error: ", err)

			}
			defer InboundResponse.Body.Close()

			// Did we call the API, how long did it take, do this here before we write to a file that will impact this time
			if vGeneral.Debuglevel > 0 {
				grpcLog.Infoln("API Call Time         :", time.Since(apiInboundStart).Seconds(), "Sec")

			}

			Inboundbody, err = io.ReadAll(InboundResponse.Body)
			if err != nil {
				grpcLog.Errorln("Inbound Body -> io.ReadAll(InboundResponse.Body) error: ", err)

			}

			var x_inboundbody map[string]interface{}
			_ = json.Unmarshal(Inboundbody, &x_inboundbody)

			if InboundResponse.Status == "200 OK" {

				// it's a paymentRT or addProxyRT - SUCCESS
				// lets build a body of the header and some additional information

				tInboundBody = map[string]interface{}{
					"transactionId":   t_InboundPayload["transactionId"],
					"eventId":         t_InboundPayload["eventId"],
					"eventType":       t_InboundPayload["eventType"],
					"responseStatus":  InboundResponse.Status,
					"responseHeaders": InboundResponse.Header,
					"responseBody":    x_inboundbody,
					"processTime":     time.Now().UTC(),
				}

			} else if InboundResponse.Status == "204 No Content" {

				// it's a paymentNRT or addProxyNRT - SUCCESS
				// lets build a body of the header and some additional information

				tInboundBody = map[string]interface{}{
					"transactionId":   t_InboundPayload["transactionId"],
					"eventId":         t_InboundPayload["eventId"],
					"eventType":       t_InboundPayload["eventType"],
					"responseStatus":  InboundResponse.Status,
					"responseHeaders": InboundResponse.Header,
					"response Body":   "<...>NRT",
					"processTime":     time.Now().UTC(),
				}

			} else {

				tInboundBody = map[string]interface{}{
					"transactionId":   t_InboundPayload["transactionId"],
					"eventId":         t_InboundPayload["eventId"],
					"eventType":       t_InboundPayload["eventType"],
					"responseResult":  "FAILED POST",
					"responseBody":    x_inboundbody,
					"responseStatus":  InboundResponse.Status,
					"responseHeaders": InboundResponse.Header,
					"processTime":     time.Now().UTC(),
				}

			}

			if vGeneral.Debuglevel > 2 {
				prettyJSON, _ := formatJSON(Inboundbody)
				grpcLog.Infoln("response Payload      :")
				grpcLog.Infoln("response Status       :", InboundResponse.Status)
				grpcLog.Infoln("response Headers      :", InboundResponse.Header)
				grpcLog.Infoln("response Body         :", string(prettyJSON))
			}

		}

		// Output Cycle
		//
		// event if we post to FS API or not, we want isolated control if we output to the json file.
		if vGeneral.Json_to_file == 1 {

			//...................................
			// Writing struct type to a JSON file
			//...................................
			// Writing
			// https://www.golangprograms.com/golang-writing-struct-to-json-file.html
			// https://www.developer.com/languages/json-files-golang/
			// Reading
			// https://medium.com/kanoteknologi/better-way-to-read-and-write-json-file-in-golang-9d575b7254f2

			// We need to do 2 api calls, 1 each for outbound and inbound event.

			// Outbound
			fileStart := time.Now()

			TransactionId := t_InboundPayload["transactionId"]

			OutboundTagId := t_OutboundPayload["eventId"]

			loc_in := fmt.Sprintf("%s%s%s-%s.json", vGeneral.Output_path, pathSep, TransactionId, OutboundTagId)
			if vGeneral.Debuglevel > 0 {
				grpcLog.Infoln("Outbound Output Event         :", loc_in)

			}

			fd, err := json.MarshalIndent(t_OutboundPayload, "", " ")
			if err != nil {
				grpcLog.Errorln("MarshalIndent error", err)

			}

			err = os.WriteFile(loc_in, fd, 0644)
			if err != nil {
				grpcLog.Errorln("os.WriteFile error", err)

			}

			// Did we call the API endpoint above... if yes then do these steps
			if vGeneral.Call_fs_api == 1 { // we need to call the API to get a output/response on paymentRT events

				if t_OutboundPayload["eventType"] == "paymentRT" || t_InboundPayload["eventType"] == "paymentRT" {

					loc_out := fmt.Sprintf("%s%s%s-%s-out.json", vGeneral.Output_path, pathSep, TransactionId, OutboundTagId)
					if vGeneral.Debuglevel > 0 {
						grpcLog.Infoln("engineResponse                :", loc_out)

					}

					fj, err := json.MarshalIndent(tOutboundBody, "", " ")
					if err != nil {
						grpcLog.Errorln("MarshalIndent error", err)

					}

					err = os.WriteFile(loc_out, fj, 0644)
					if err != nil {
						grpcLog.Errorln("os.WriteFile error", err)

					}
				} else {

				}
			}

			if vGeneral.Debuglevel > 0 {
				grpcLog.Infoln("Outbound JSON to File         :", time.Since(fileStart).Seconds(), "Sec")

			}

			// Inbound
			fileStart = time.Now()

			InboundTagId := t_InboundPayload["eventId"]

			loc_in = fmt.Sprintf("%s%s%s-%s.json", vGeneral.Output_path, pathSep, TransactionId, InboundTagId)
			if vGeneral.Debuglevel > 0 {
				grpcLog.Infoln("Inbound Output Event          :", loc_in)

			}

			fd, err = json.MarshalIndent(t_InboundPayload, "", " ")
			if err != nil {
				grpcLog.Errorln("MarshalIndent error", err)

			}

			err = os.WriteFile(loc_in, fd, 0644)
			if err != nil {
				grpcLog.Errorln("os.WriteFile error", err)

			}

			// Did we call the API endpoint above... if yes then do these steps
			if vGeneral.Call_fs_api == 1 { // we need to call the API to get a output/response on paymentRT events

				loc_out := fmt.Sprintf("%s%s%s-%s-out.json", vGeneral.Output_path, pathSep, TransactionId, InboundTagId)

				if vGeneral.Debuglevel > 0 {
					grpcLog.Infoln("engineResponse                :", loc_out)

				}

				fj, err := json.MarshalIndent(tInboundBody, "", " ")
				if err != nil {
					grpcLog.Errorln("MarshalIndent error", err)

				}

				err = os.WriteFile(loc_out, fj, 0644)
				if err != nil {
					grpcLog.Errorln("os.WriteFile error", err)

				}
			}

			if vGeneral.Debuglevel > 0 {
				grpcLog.Infoln("Inbound JSON to File          :", time.Since(fileStart).Seconds(), "Sec")

			}

		}

		if vGeneral.Debuglevel > 0 {
			grpcLog.Infoln("Total Time                    :", time.Since(txnStart).Seconds(), "Sec")

		}

		//////////////////////////////////////////////////
		//
		// THIS IS SLEEP BETWEEN RECORD POSTS
		//
		// if 0 then sleep is disabled otherwise
		//
		// lets get a random value 0 -> vGeneral.sleep, then delay/sleep as up to that fraction of a second.
		// this mimics someone thinking, as if this is being done by a human at a keyboard, for batch file processing we don't have this.
		// ie if the user said 200 then it implies a randam value from 0 -> 200 milliseconds.
		//
		// USED TO SLOW THINGS DOWN
		//
		//////////////////////////////////////////////////

		if vGeneral.Sleep != 0 {
			n := rand.Intn(vGeneral.Sleep) // if vGeneral.sleep = 1000, then n will be random value of 0 -> 1000  aka 0 and 1 second
			if vGeneral.Debuglevel >= 2 {
				grpcLog.Infof("Going to sleep for            : %d Milliseconds\n", n)

			}
			time.Sleep(time.Duration(n) * time.Millisecond)
		}

	}

	if vGeneral.Debuglevel > 0 {
		grpcLog.Infoln("")
		grpcLog.Infoln("**** DONE Processing ****")
		grpcLog.Infoln("")

		if vGeneral.Debuglevel >= 1 {
			vEnd := time.Now()
			grpcLog.Infoln("Start                         : ", vStart)
			grpcLog.Infoln("End                           : ", vEnd)
			grpcLog.Infoln("Duration                      : ", vEnd.Sub(vStart))
			grpcLog.Infoln("Records                       : ", todo_count)

			grpcLog.Infoln("")
		}
	}

} // runLoader()

func main() {

	var arg string

	grpcLog.Info("****** Starting           *****")

	arg = os.Args[1]

	runLoader(arg)

	grpcLog.Info("****** Completed          *****")

}
