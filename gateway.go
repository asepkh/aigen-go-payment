package payment

import (
	"fmt"
	"strings"
)

// Gateway represent payment gateway used
type Gateway int

const (
	// UnknownGateway gateway is unknown
	UnknownGateway Gateway = iota
	// GatewayMidtrans is midtrans payment gateway
	GatewayMidtrans
	// GatewayXendit is xendit payment gateway
	GatewayXendit
	// GatewayFinpay is finpay payment gateway
	GatewayFinpay
)

func (g Gateway) String() string {
	return []string{"unkown", "midtrans", "xendit", "finpay"}[g]
}

// UnmarshalYAML convert string to Gateway enum
func (g *Gateway) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var n string
	if err := unmarshal(&n); err != nil {
		return err
	}

	*g = NewGateway(n)

	if *g == UnknownGateway {
		return fmt.Errorf("payment gateway is not recognized")
	}

	return nil
}

// NewGateway return a gateway for its string name
func NewGateway(name string) Gateway {
	var g Gateway
	switch strings.ToLower(name) {
	case "midtrans":
		g = GatewayMidtrans
	case "xendit":
		g = GatewayXendit
	case "finpay":
		g = GatewayFinpay
	default:
		g = UnknownGateway
	}
	return g
}
