package midtrans

import (
	"ceo-suite-go/configs"
	"errors"
	"fmt"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type ChargeResponse = coreapi.ChargeResponse

type MidtransService interface {
	GenerateTransaction(result int, paymentType string, invoice string, custDetails *midtrans.CustomerDetails, itemDetails []midtrans.ItemDetails) (*ChargeResponse, map[string]interface{}, error)
	GenerateTransactionSnap(result int, paymentType string, invoice string, custDetails *midtrans.CustomerDetails, itemDetails []midtrans.ItemDetails) (*snap.Response, error)
	TransactionStatus(notificationPayload map[string]interface{}) (int, string, error)
}

type midtransService struct {
	core coreapi.Client
}

func InitMidtrans(c configs.ProgrammingConfig) MidtransService {
	var core coreapi.Client
	var envi midtrans.EnvironmentType
	if c.MidtransEnvironment == "production" {
		envi = midtrans.Production
	} else {
		envi = midtrans.Sandbox
	}

	core.New(c.MidtransServerKey, envi)

	return &midtransService{
		core: core,
	}
}

func (ms *midtransService) GenerateTransaction(result int, paymentType string, invoice string, custDetails *midtrans.CustomerDetails, itemDetails []midtrans.ItemDetails) (*ChargeResponse, map[string]interface{}, error) {
	fmt.Println("Get data: ", result, invoice, paymentType)

	var chargeReq *coreapi.ChargeReq
	response := map[string]any{}

	// if paymentType == "credit_card" {
	// 	chargeReq = &coreapi.ChargeReq{
	// 		PaymentType: "credit_card",
	// 		TransactionDetails: midtrans.TransactionDetails{
	// 			OrderID:  invoice,
	// 			GrossAmt: int64(result),
	// 		},
	// 		CustomerDetails: custDetails,
	// 		Items: &itemDetails,
	// 		CreditCard: &coreapi.CreditCardDetails{

	// 		},
	// 	}
	// }

	if paymentType == "qris" {
		chargeReq = &coreapi.ChargeReq{
			PaymentType: "qris",
			Gopay: &coreapi.GopayDetails{
				EnableCallback: true,
			},
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  invoice,
				GrossAmt: int64(result),
			},
			CustomerDetails: custDetails,
			Items:           &itemDetails,
		}
	}

	if paymentType == "gopay" {
		chargeReq = &coreapi.ChargeReq{
			Gopay: &coreapi.GopayDetails{
				EnableCallback: true,
			},
			PaymentType: coreapi.PaymentTypeGopay,
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  invoice,
				GrossAmt: int64(result),
			},
			CustomerDetails: custDetails,
			Items:           &itemDetails,
		}
	}

	if paymentType == "bca" || paymentType == "bni" || paymentType == "bri" {
		var midtransBank midtrans.Bank

		switch paymentType {
		case "bca":
			midtransBank = midtrans.BankBca
		case "bri":
			midtransBank = midtrans.BankBri
		case "bni":
			midtransBank = midtrans.BankBni
		default:
			midtransBank = midtrans.BankBca
		}

		chargeReq = &coreapi.ChargeReq{
			PaymentType:  coreapi.PaymentTypeBankTransfer,
			BankTransfer: &coreapi.BankTransferDetails{Bank: midtransBank},
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  invoice,
				GrossAmt: int64(result),
			},
			CustomerDetails: custDetails,
			Items:           &itemDetails,
		}
	}

	// fmt.Println("Charge req:", chargeReq)
	chargeResp, err := ms.core.ChargeTransaction(chargeReq)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, nil, err
	}

	if paymentType == "qris" || paymentType == "gopay" {
		fmt.Println("Nyampe sini")
		if len(chargeResp.Actions) > 0 {
			for _, action := range chargeResp.Actions {
				switch action.Name {
				case "generate-qr-code":
					deepLinkURL := action.URL
					response["callback_url"] = deepLinkURL
					response["payment_type"] = "qris"
				case "deeplink-redirect":
					deepLinkURL := action.URL
					response["callback_url"] = deepLinkURL
					response["payment_type"] = "gopay"
				}
			}
		}
	}

	if paymentType == "bca" || paymentType == "bni" || paymentType == "bri" {
		var vaAccount string

		for _, va := range chargeResp.VaNumbers {
			if va.Bank == paymentType {
				vaAccount = va.VANumber
				break
			}
		}
		response["invoice_id"] = chargeResp.OrderID
		response["payment_type"] = paymentType
		response["va_account"] = vaAccount
	}

	return chargeResp, response, nil
}

func (ms *midtransService) GenerateTransactionSnap(result int, paymentType string, invoice string, custDetails *midtrans.CustomerDetails, itemDetails []midtrans.ItemDetails) (*snap.Response, error) {
	config := configs.InitConfig()
	if config == nil {
		return nil, errors.New("failed to load configuration")
	}

	midtrans.ServerKey = config.MidtransServerKey
	if config.MidtransEnvironment == "production" {
		midtrans.Environment = midtrans.Production
	} else {
		midtrans.Environment = midtrans.Sandbox
	}

	// itemDetail []midtrans.ItemDetails
	// allName    string
	// allPrice   int64

	// for i, prod := range cartData {
	// 	var totalPerProd int64

	// 	allName += prod.ProductDetail.Name
	// 	totalPerProd += int64(prod.ProductDetail.Price.PriceSRPAfterDiscount)

	// 	for _, namePU := range prod.UpgradeProduct {
	// 		totalPerProd += int64(namePU.ProductUpgradeDetail.Price.PriceSRPAfterDiscount)
	// 	}

	// 	for _, namePA := range prod.UpgradeAccessories {
	// 		totalPerProd += int64(namePA.ProductAccessoriesDetail.Price.PriceSRPAfterDiscount)
	// 	}

	// 	item := midtrans.ItemDetails{
	// 		ID:           fmt.Sprintf("%v", i),
	// 		Name:         allName,
	// 		Price:        totalPerProd,
	// 		Qty:          int32(prod.Quantity),
	// 		Brand:        fmt.Sprintf("Zyrex %v", cartData[i].ProductDetail.Name),
	// 		Category:     fmt.Sprintf("%v", cartData[i].ProductDetail.ProductCategory),
	// 		MerchantName: "Zyrex",
	// 	}

	// 	allPrice += (totalPerProd * int64(prod.Quantity))
	// 	itemDetail = append(itemDetail, item)
	// }

	// pajak := midtrans.ItemDetails{
	// 	ID:           fmt.Sprintf("%v", "Pajak"),
	// 	Name:         "Tax 10%",
	// 	Price:        (10 * allPrice) / 100,
	// 	Qty:          1,
	// 	Brand:        "Tax 10%",
	// 	Category:     "Tax",
	// 	MerchantName: "Pajak Zyrex",
	// }

	// itemDetail = append(itemDetail, pajak)

	// if shippingPrice != 0 {

	// 	shipping := midtrans.ItemDetails{
	// 		ID:           fmt.Sprintf("%v", "Shipping"),
	// 		Name:         "Shipping",
	// 		Price:        int64(shippingPrice),
	// 		Qty:          1,
	// 		Brand:        "Shipping",
	// 		Category:     "Shipping",
	// 		MerchantName: "Shipping Zyrex",
	// 	}

	// 	itemDetail = append(itemDetail, shipping)
	// }
	// var snapReq *snap.Request

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  invoice,
			GrossAmt: int64(result),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},

		CustomerDetail: custDetails,
		Items:          &itemDetails,
		Gopay: &snap.GopayDetails{
			EnableCallback: true,
		},
	}

	snapResp, err := snap.CreateTransaction(snapReq)
	if err != nil {
		return nil, err
	}

	return snapResp, nil
}

func (ms *midtransService) TransactionStatus(notificationPayload map[string]interface{}) (int, string, error) {
	var paymentStatus int
	orderId, exists := notificationPayload["order_id"].(string)
	if !exists {
		return 0, "", errors.New("Order ID Not Found")
	}

	transactionStatusResp, e := ms.core.CheckTransaction(orderId)
	if e != nil {
		return 0, "", errors.New(e.GetMessage())
	} else {
		if transactionStatusResp != nil {
			if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "challenge" {
					fmt.Println("Payment status challenged")
					paymentStatus = 1
					return paymentStatus, transactionStatusResp.OrderID, nil
				} else if transactionStatusResp.FraudStatus == "accept" {
					fmt.Println("Payment received")
					paymentStatus = 2
					return paymentStatus, transactionStatusResp.OrderID, nil
				}
			} else if transactionStatusResp.TransactionStatus == "settlement" {
				fmt.Println("Payment status settlement")
				paymentStatus = 2
				return paymentStatus, transactionStatusResp.OrderID, nil
			} else if transactionStatusResp.TransactionStatus == "deny" {
				fmt.Println("Payment status denied")
				paymentStatus = 3
				return paymentStatus, transactionStatusResp.OrderID, nil
			} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
				fmt.Println("Payment status failure")
				paymentStatus = 4
				return paymentStatus, transactionStatusResp.OrderID, nil
			} else if transactionStatusResp.TransactionStatus == "pending" {
				fmt.Println("Payment status pending")
				paymentStatus = 5
				return paymentStatus, transactionStatusResp.OrderID, nil
			}
		}
	}

	return 0, "", nil
}
