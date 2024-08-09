package config

import "github.com/midtrans/midtrans-go"

func SetupMidtrans() {
	midtrans.ServerKey = Env.MidtransServerKey
	midtrans.Environment = Env.MidtransEnvironment
}
