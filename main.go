package main

import (
	"os"
	"strconv"
	"time"

	"github.com/bitcoin-sv/go-paymail/logging"

	"github.com/bitcoin-sv/go-paymail/server"
)

func main() {
	logger := logging.GetDefaultLogger()

	sl := server.PaymailServiceLocator{}
	sl.RegisterPaymailService(new(opnsServiceProvider))
	sl.RegisterPikeContactService(new(opnsServiceProvider))
	sl.RegisterPikePaymentService(new(opnsServiceProvider))

	var err error
	port := 3000
	portEnv := os.Getenv("PORT")
	if portEnv != "" {
		if port, err = strconv.Atoi(portEnv); err != nil {
			logger.Fatal().Msg(err.Error())
		}
	}
	// Custom server with lots of customizable goodies
	config, err := server.NewConfig(
		&sl,
		server.WithBasicRoutes(),
		server.WithP2PCapabilities(),
		server.WithBeefCapabilities(),
		server.WithDomain("1sat.app"),
		server.WithDomain("opns-paymail-production.up.railway.app"),
		server.WithDomain("localhost:3000"),
		// server.WithGenericCapabilities(),
		server.WithPort(port),
		// server.WithServiceName("BsvAliasCustom"),
		server.WithTimeout(15*time.Second),
		// server.WithCapabilities(customCapabilities()),
	)
	config.Prefix = "https://" //normally paymail requires https, but for demo purposes we'll use http
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}

	// Create & start the server
	server.StartServer(server.CreateServer(config), config.Logger)
}

// func customCapabilities() map[string]any {
// 	exampleBrfcKey := "406cef0ae2d6"
// 	return map[string]any{
// 		"custom_static_boolean": false,
// 		"custom_static_int":     10,
// 		exampleBrfcKey:          true,
// 		"custom_callable_cap": server.CallableCapability{
// 			Path:   fmt.Sprintf("/display_paymail/%s", server.PaymailAddressTemplate),
// 			Method: http.MethodGet,
// 			Handler: func(c *gin.Context) {
// 				incomingPaymail := c.Param(server.PaymailAddressParamName)

// 				response := map[string]string{
// 					"paymail": incomingPaymail,
// 				}

// 				c.Header("Content-Type", "application/json")
// 				c.JSON(http.StatusOK, response)
// 			},
// 		},
// 	}
// }
