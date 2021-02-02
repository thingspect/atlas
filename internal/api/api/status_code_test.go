// +build !integration

package api

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

const testTimeout = 2 * time.Second

func TestStatusCode(t *testing.T) {
	t.Parallel()

	t.Run("Modify status code", func(t *testing.T) {
		t.Parallel()

		mdHeader := metadata.MD{"atlas-status-code": []string{"201"}}
		wHeader := http.Header{"Grpc-Metadata-Atlas-Status-Code": []string{
			"201"}}
		t.Logf("mdHeader, wHeader: %+v, %+v", mdHeader, wHeader)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		respWriter := NewMockResponseWriter(ctrl)
		respWriter.EXPECT().Header().Return(wHeader).Times(1)
		respWriter.EXPECT().WriteHeader(201).Times(1)

		ctx, cancel := context.WithTimeout(runtime.NewServerMetadataContext(
			context.Background(), runtime.ServerMetadata{HeaderMD: mdHeader}),
			testTimeout)
		defer cancel()

		err := statusCode(ctx, respWriter, nil)
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		t.Logf("mdHeader, wHeader: %+v, %+v", mdHeader, wHeader)
		require.Empty(t, mdHeader)
		require.Empty(t, wHeader)
	})

	t.Run("Pass through status code without header", func(t *testing.T) {
		t.Parallel()

		mdHeader := metadata.MD{}
		wHeader := http.Header{}
		t.Logf("mdHeader, wHeader: %+v, %+v", mdHeader, wHeader)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		respWriter := NewMockResponseWriter(ctrl)

		ctx, cancel := context.WithTimeout(runtime.NewServerMetadataContext(
			context.Background(), runtime.ServerMetadata{HeaderMD: mdHeader}),
			testTimeout)
		defer cancel()

		err := statusCode(ctx, respWriter, nil)
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		t.Logf("mdHeader, wHeader: %+v, %+v", mdHeader, wHeader)
		require.Empty(t, mdHeader)
		require.Empty(t, wHeader)
	})

	t.Run("Pass through status code without metadata", func(t *testing.T) {
		t.Parallel()

		mdHeader := metadata.MD{}
		wHeader := http.Header{}
		t.Logf("mdHeader, wHeader: %+v, %+v", mdHeader, wHeader)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		respWriter := NewMockResponseWriter(ctrl)

		ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
		defer cancel()

		err := statusCode(ctx, respWriter, nil)
		t.Logf("err: %#v", err)
		require.NoError(t, err)

		t.Logf("mdHeader, wHeader: %+v, %+v", mdHeader, wHeader)
		require.Empty(t, mdHeader)
		require.Empty(t, wHeader)
	})

	t.Run("Don't modify status code with invalid metadata", func(t *testing.T) {
		t.Parallel()

		mdHeader := metadata.MD{"atlas-status-code": []string{"aaa"}}
		wHeader := http.Header{"Grpc-Metadata-Atlas-Status-Code": []string{
			"201"}}
		t.Logf("mdHeader, wHeader: %+v, %+v", mdHeader, wHeader)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		respWriter := NewMockResponseWriter(ctrl)

		ctx, cancel := context.WithTimeout(runtime.NewServerMetadataContext(
			context.Background(), runtime.ServerMetadata{HeaderMD: mdHeader}),
			testTimeout)
		defer cancel()

		err := statusCode(ctx, respWriter, nil)
		t.Logf("err: %#v", err)
		require.ErrorIs(t, err, strconv.ErrSyntax)

		t.Logf("mdHeader, wHeader: %+v, %+v", mdHeader, wHeader)
		require.Equal(t, metadata.MD{"atlas-status-code": []string{"aaa"}},
			mdHeader)
		require.Equal(t, http.Header{
			"Grpc-Metadata-Atlas-Status-Code": []string{"201"}}, wHeader)
	})
}
