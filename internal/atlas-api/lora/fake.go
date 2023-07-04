package lora

import (
	"github.com/thingspect/atlas/pkg/alog"
	"go.uber.org/mock/gomock"
)

// NewFake builds a new Loraer using a mock and returns it. It may be used in
// place of a Lora implementation, but should be accompanied by a warning.
func NewFake() Loraer {
	// Controller.Finish() is not called because usage is expected to be
	// long-lived.
	loraer := NewMockLoraer(gomock.NewController(alog.Default()))
	loraer.EXPECT().CreateGateway(gomock.Any(), gomock.Any()).Return(nil).
		AnyTimes()
	loraer.EXPECT().DeleteGateway(gomock.Any(), gomock.Any()).Return(nil).
		AnyTimes()
	loraer.EXPECT().CreateDevice(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).AnyTimes()
	loraer.EXPECT().DeleteDevice(gomock.Any(), gomock.Any()).Return(nil).
		AnyTimes()

	return loraer
}
