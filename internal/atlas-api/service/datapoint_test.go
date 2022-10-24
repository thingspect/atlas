//go:build !integration

package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/api/go/message"
	"github.com/thingspect/atlas/internal/atlas-api/session"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/atlas/pkg/test/matcher"
	"github.com/thingspect/atlas/pkg/test/random"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPublishDataPoints(t *testing.T) {
	t.Parallel()

	t.Run("Publish valid data point", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()
		point := &common.DataPoint{
			UniqId: "api-point-" + random.String(16), Attr: "count",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts:       timestamppb.New(time.Now().Add(-15 * time.Minute)),
		}

		dpQueue := queue.NewFake()
		vInSub, err := dpQueue.Subscribe("")
		require.NoError(t, err)
		vInPubTopic := "topic-" + random.String(10)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_PUBLISHER,
			}), testTimeout)
		defer cancel()

		dpSvc := NewDataPoint(dpQueue, vInPubTopic, nil)
		_, err = dpSvc.PublishDataPoints(ctx, &api.PublishDataPointsRequest{
			Points: []*common.DataPoint{point},
		})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		select {
		case msg := <-vInSub.C():
			msg.Ack()
			t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(), msg.Payload())
			require.Equal(t, vInPubTopic, msg.Topic())

			vIn := &message.ValidatorIn{}
			require.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
			t.Logf("vIn: %+v", vIn)

			// Normalize generated trace ID.
			point.TraceId = vIn.Point.TraceId

			// Testify does not currently support protobuf equality:
			// https://github.com/stretchr/testify/issues/758
			if !proto.Equal(&message.ValidatorIn{
				Point: point, OrgId: orgID, SkipToken: true,
			}, vIn) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", &message.ValidatorIn{
					Point: point, OrgId: orgID, SkipToken: true,
				}, vIn)
			}
		case <-time.After(testTimeout):
			t.Fatal("Message timed out")
		}
	})

	t.Run("Publish valid data point without timestamp", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()
		point := &common.DataPoint{
			UniqId: "api-point-" + random.String(16), Attr: "count",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
		}

		pubQueue := queue.NewFake()
		pubSub, err := pubQueue.Subscribe("")
		require.NoError(t, err)
		pubTopic := "topic-" + random.String(10)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		dpSvc := NewDataPoint(pubQueue, pubTopic, nil)
		_, err = dpSvc.PublishDataPoints(ctx, &api.PublishDataPointsRequest{
			Points: []*common.DataPoint{point},
		})
		t.Logf("err: %v", err)
		require.NoError(t, err)

		select {
		case msg := <-pubSub.C():
			msg.Ack()
			t.Logf("msg.Topic, msg.Payload: %v, %s", msg.Topic(), msg.Payload())
			require.Equal(t, pubTopic, msg.Topic())

			vIn := &message.ValidatorIn{}
			require.NoError(t, proto.Unmarshal(msg.Payload(), vIn))
			t.Logf("vIn: %+v", vIn)

			// Normalize generated trace ID.
			point.TraceId = vIn.Point.TraceId
			// Normalize timestamp.
			require.WithinDuration(t, time.Now(), vIn.Point.Ts.AsTime(),
				2*time.Second)
			point.Ts = vIn.Point.Ts

			// Testify does not currently support protobuf equality:
			// https://github.com/stretchr/testify/issues/758
			if !proto.Equal(&message.ValidatorIn{
				Point: point, OrgId: orgID, SkipToken: true,
			}, vIn) {
				t.Fatalf("\nExpect: %+v\nActual: %+v", &message.ValidatorIn{
					Point: point, OrgId: orgID, SkipToken: true,
				}, vIn)
			}
		case <-time.After(testTimeout):
			t.Fatal("Message timed out")
		}
	})

	t.Run("Publish data point with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		dpSvc := NewDataPoint(nil, "", nil)
		_, err := dpSvc.PublishDataPoints(ctx, &api.PublishDataPointsRequest{})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(api.Role_PUBLISHER), err)
	})

	t.Run("Publish data point with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_VIEWER,
			}), testTimeout)
		defer cancel()

		dpSvc := NewDataPoint(nil, "", nil)
		_, err := dpSvc.PublishDataPoints(ctx, &api.PublishDataPointsRequest{})
		t.Logf("err: %v", err)
		require.Equal(t, errPerm(api.Role_PUBLISHER), err)
	})

	t.Run("Publish valid data point with bad queue", func(t *testing.T) {
		t.Parallel()

		orgID := uuid.NewString()
		point := &common.DataPoint{
			UniqId: "api-point-" + random.String(16), Attr: "count",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts:       timestamppb.New(time.Now().Add(-15 * time.Minute)),
		}
		vInPubTopic := "topic-" + random.String(10)

		queuer := queue.NewMockQueuer(gomock.NewController(t))
		queuer.EXPECT().Publish(vInPubTopic, gomock.Any()).
			Return(dao.ErrNotFound).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		dpSvc := NewDataPoint(queuer, vInPubTopic, nil)
		_, err := dpSvc.PublishDataPoints(ctx, &api.PublishDataPointsRequest{
			Points: []*common.DataPoint{point},
		})
		t.Logf("err: %v", err)
		require.Equal(t, status.Error(codes.Internal, "publish failure"), err)
	})
}

func TestListDataPoints(t *testing.T) {
	t.Parallel()

	t.Run("List data points by valid UniqID with ts", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{
			UniqId: "api-point-" + random.String(16), Attr: "count",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts:       timestamppb.Now(), TraceId: uuid.NewString(),
		}
		retPoint, _ := proto.Clone(point).(*common.DataPoint)
		orgID := uuid.NewString()
		end := time.Now().UTC()
		start := time.Now().UTC().Add(-15 * time.Minute)

		datapointer := NewMockDataPointer(gomock.NewController(t))
		datapointer.EXPECT().List(gomock.Any(), orgID, point.UniqId, "", "",
			end, start).Return([]*common.DataPoint{retPoint}, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		dpSvc := NewDataPoint(nil, "", datapointer)
		listPoints, err := dpSvc.ListDataPoints(ctx, &api.ListDataPointsRequest{
			IdOneof: &api.ListDataPointsRequest_UniqId{UniqId: point.UniqId},
			EndTime: timestamppb.New(end), StartTime: timestamppb.New(start),
		})
		t.Logf("point, listPoints, err: %+v, %+v, %v", point, listPoints, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListDataPointsResponse{
			Points: []*common.DataPoint{point},
		}, listPoints) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", &api.ListDataPointsResponse{
				Points: []*common.DataPoint{point},
			}, listPoints)
		}
	})

	t.Run("List data points by valid dev ID with attr", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{
			UniqId: "api-point-" + random.String(16), Attr: "count",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts:       timestamppb.Now(), TraceId: uuid.NewString(),
		}
		retPoint, _ := proto.Clone(point).(*common.DataPoint)
		orgID := uuid.NewString()
		devID := uuid.NewString()

		datapointer := NewMockDataPointer(gomock.NewController(t))
		datapointer.EXPECT().List(gomock.Any(), orgID, "", devID, point.Attr,
			matcher.NewRecentMatcher(2*time.Second),
			matcher.NewRecentMatcher(24*time.Hour+2*time.Second)).
			Return([]*common.DataPoint{retPoint}, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		dpSvc := NewDataPoint(nil, "", datapointer)
		listPoints, err := dpSvc.ListDataPoints(ctx, &api.ListDataPointsRequest{
			IdOneof: &api.ListDataPointsRequest_DeviceId{DeviceId: devID},
			Attr:    point.Attr,
		})
		t.Logf("point, listPoints, err: %+v, %+v, %v", point, listPoints, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.ListDataPointsResponse{
			Points: []*common.DataPoint{point},
		}, listPoints) {
			t.Fatalf("\nExpect: %+v\nActual: %+v", &api.ListDataPointsResponse{
				Points: []*common.DataPoint{point},
			}, listPoints)
		}
	})

	t.Run("List data points with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		dpSvc := NewDataPoint(nil, "", nil)
		listPoints, err := dpSvc.ListDataPoints(ctx,
			&api.ListDataPointsRequest{})
		t.Logf("listPoints, err: %+v, %v", listPoints, err)
		require.Nil(t, listPoints)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("List data points with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_CONTACT,
			}), testTimeout)
		defer cancel()

		dpSvc := NewDataPoint(nil, "", nil)
		listPoints, err := dpSvc.ListDataPoints(ctx,
			&api.ListDataPointsRequest{})
		t.Logf("listPoints, err: %+v, %v", listPoints, err)
		require.Nil(t, listPoints)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("List data points by invalid time range", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		dpSvc := NewDataPoint(nil, "", nil)
		listPoints, err := dpSvc.ListDataPoints(ctx, &api.ListDataPointsRequest{
			IdOneof: &api.ListDataPointsRequest_UniqId{UniqId: "api-point-" +
				random.String(16)}, EndTime: timestamppb.Now(),
			StartTime: timestamppb.New(time.Now().Add(-91 * 24 * time.Hour)),
		})
		t.Logf("listPoints, err: %+v, %v", listPoints, err)
		require.Nil(t, listPoints)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"maximum time range exceeded"), err)
	})

	t.Run("List data points by invalid org ID", func(t *testing.T) {
		t.Parallel()

		datapointer := NewMockDataPointer(gomock.NewController(t))
		datapointer.EXPECT().List(gomock.Any(), "aaa", gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil,
			dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: "aaa", Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		dpSvc := NewDataPoint(nil, "", datapointer)
		listPoints, err := dpSvc.ListDataPoints(ctx, &api.ListDataPointsRequest{
			IdOneof: &api.ListDataPointsRequest_UniqId{UniqId: "api-point-" +
				random.String(16)},
		})
		t.Logf("listPoints, err: %+v, %v", listPoints, err)
		require.Nil(t, listPoints)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}

func TestLatestDataPoints(t *testing.T) {
	t.Parallel()

	t.Run("Latest data points by valid UniqID", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{
			UniqId: "api-point-" + random.String(16), Attr: "count",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts:       timestamppb.Now(), TraceId: uuid.NewString(),
		}
		retPoint, _ := proto.Clone(point).(*common.DataPoint)
		orgID := uuid.NewString()

		datapointer := NewMockDataPointer(gomock.NewController(t))
		datapointer.EXPECT().Latest(gomock.Any(), orgID, point.UniqId, "",
			matcher.NewRecentMatcher(30*24*time.Hour+2*time.Second)).
			Return([]*common.DataPoint{retPoint}, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		dpSvc := NewDataPoint(nil, "", datapointer)
		latPoints, err := dpSvc.LatestDataPoints(ctx,
			&api.LatestDataPointsRequest{
				IdOneof: &api.LatestDataPointsRequest_UniqId{
					UniqId: point.UniqId,
				},
			})
		t.Logf("point, latPoints, err: %+v, %+v, %v", point, latPoints, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.LatestDataPointsResponse{
			Points: []*common.DataPoint{point},
		}, latPoints) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.LatestDataPointsResponse{
					Points: []*common.DataPoint{point},
				}, latPoints)
		}
	})

	t.Run("Latest data points by valid dev ID with ts", func(t *testing.T) {
		t.Parallel()

		point := &common.DataPoint{
			UniqId: "api-point-" + random.String(16), Attr: "count",
			ValOneof: &common.DataPoint_IntVal{IntVal: 123},
			Ts:       timestamppb.Now(), TraceId: uuid.NewString(),
		}
		retPoint, _ := proto.Clone(point).(*common.DataPoint)
		orgID := uuid.NewString()
		devID := uuid.NewString()
		start := time.Now().UTC().Add(-15 * time.Minute)

		datapointer := NewMockDataPointer(gomock.NewController(t))
		datapointer.EXPECT().Latest(gomock.Any(), orgID, "", devID, start).
			Return([]*common.DataPoint{retPoint}, nil).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: orgID, Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		dpSvc := NewDataPoint(nil, "", datapointer)
		latPoints, err := dpSvc.LatestDataPoints(ctx,
			&api.LatestDataPointsRequest{
				IdOneof: &api.LatestDataPointsRequest_DeviceId{
					DeviceId: devID,
				},
				StartTime: timestamppb.New(start),
			})
		t.Logf("point, latPoints, err: %+v, %+v, %v", point, latPoints, err)
		require.NoError(t, err)

		// Testify does not currently support protobuf equality:
		// https://github.com/stretchr/testify/issues/758
		if !proto.Equal(&api.LatestDataPointsResponse{
			Points: []*common.DataPoint{point},
		}, latPoints) {
			t.Fatalf("\nExpect: %+v\nActual: %+v",
				&api.LatestDataPointsResponse{
					Points: []*common.DataPoint{point},
				}, latPoints)
		}
	})

	t.Run("Latest data points with invalid session", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		dpSvc := NewDataPoint(nil, "", nil)
		latPoints, err := dpSvc.LatestDataPoints(ctx,
			&api.LatestDataPointsRequest{})
		t.Logf("latPoints, err: %+v, %v", latPoints, err)
		require.Nil(t, latPoints)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("Latest data points with insufficient role", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_CONTACT,
			}), testTimeout)
		defer cancel()

		dpSvc := NewDataPoint(nil, "", nil)
		latPoints, err := dpSvc.LatestDataPoints(ctx,
			&api.LatestDataPointsRequest{})
		t.Logf("latPoints, err: %+v, %v", latPoints, err)
		require.Nil(t, latPoints)
		require.Equal(t, errPerm(api.Role_VIEWER), err)
	})

	t.Run("Latest data points by invalid time range", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: uuid.NewString(), Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		dpSvc := NewDataPoint(nil, "", nil)
		latPoints, err := dpSvc.LatestDataPoints(ctx,
			&api.LatestDataPointsRequest{
				IdOneof: &api.LatestDataPointsRequest_UniqId{
					UniqId: "api-point-" + random.String(16),
				}, StartTime: timestamppb.New(
					time.Now().Add(-91 * 24 * time.Hour)),
			})
		t.Logf("latPoints, err: %+v, %v", latPoints, err)
		require.Nil(t, latPoints)
		require.Equal(t, status.Error(codes.InvalidArgument,
			"maximum time range exceeded"), err)
	})

	t.Run("Latest data points by invalid org ID", func(t *testing.T) {
		t.Parallel()

		datapointer := NewMockDataPointer(gomock.NewController(t))
		datapointer.EXPECT().Latest(gomock.Any(), "aaa", gomock.Any(),
			gomock.Any(),
			matcher.NewRecentMatcher(30*24*time.Hour+2*time.Second)).
			Return(nil, dao.ErrInvalidFormat).Times(1)

		ctx, cancel := context.WithTimeout(session.NewContext(
			context.Background(), &session.Session{
				OrgID: "aaa", Role: api.Role_ADMIN,
			}), testTimeout)
		defer cancel()

		dpSvc := NewDataPoint(nil, "", datapointer)
		latPoints, err := dpSvc.LatestDataPoints(ctx,
			&api.LatestDataPointsRequest{
				IdOneof: &api.LatestDataPointsRequest_UniqId{
					UniqId: "api-point-" + random.String(16),
				},
			})
		t.Logf("latPoints, err: %+v, %v", latPoints, err)
		require.Nil(t, latPoints)
		require.Equal(t, status.Error(codes.InvalidArgument, "invalid format"),
			err)
	})
}
