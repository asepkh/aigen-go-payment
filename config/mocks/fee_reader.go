// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	config "github.com/imrenagi/go-payment/config"
	mock "github.com/stretchr/testify/mock"

	payment "github.com/imrenagi/go-payment"

	time "time"
)

// FeeConfigReader is an autogenerated mock type for the FeeConfigReader type
type FeeConfigReader struct {
	mock.Mock
}

// GetAdminFeeConfig provides a mock function with given fields: currency
func (_m *FeeConfigReader) GetAdminFeeConfig(currency string) *config.Fee {
	ret := _m.Called(currency)

	var r0 *config.Fee
	if rf, ok := ret.Get(0).(func(string) *config.Fee); ok {
		r0 = rf(currency)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.Fee)
		}
	}

	return r0
}

// GetGateway provides a mock function with given fields:
func (_m *FeeConfigReader) GetGateway() payment.Gateway {
	ret := _m.Called()

	var r0 payment.Gateway
	if rf, ok := ret.Get(0).(func() payment.Gateway); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(payment.Gateway)
	}

	return r0
}

// GetInstallmentFeeConfig provides a mock function with given fields: currency
func (_m *FeeConfigReader) GetInstallmentFeeConfig(currency string) *config.Fee {
	ret := _m.Called(currency)

	var r0 *config.Fee
	if rf, ok := ret.Get(0).(func(string) *config.Fee); ok {
		r0 = rf(currency)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.Fee)
		}
	}

	return r0
}

// GetPaymentWaitingTime provides a mock function with given fields:
func (_m *FeeConfigReader) GetPaymentWaitingTime() *time.Duration {
	ret := _m.Called()

	var r0 *time.Duration
	if rf, ok := ret.Get(0).(func() *time.Duration); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*time.Duration)
		}
	}

	return r0
}
