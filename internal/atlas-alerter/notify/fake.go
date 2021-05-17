package notify

import (
	"github.com/golang/mock/gomock"
	"github.com/thingspect/atlas/pkg/alog"
)

// NewFake builds a new Notifier using a mock and returns it. It may be used in
// place of a Notify implementation, but should be accompanied by a warning.
func NewFake() Notifier {
	// Controller.Finish() is not called because usage is expected to be
	// long-lived.
	notifier := NewMockNotifier(gomock.NewController(alog.Default()))
	notifier.EXPECT().App(gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any()).Return(nil).AnyTimes()
	notifier.EXPECT().VaildateSMS(gomock.Any(), gomock.Any()).Return(nil).
		AnyTimes()
	notifier.EXPECT().SMS(gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any()).Return(nil).AnyTimes()
	notifier.EXPECT().Email(gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	return notifier
}
