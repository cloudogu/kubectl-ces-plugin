package portforward

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"net/url"
	"testing"
)

func Test_executivePortForwarder_ExecuteWithPortForward(t *testing.T) {
	type fields struct {
		dialerFactory    dialerFactory
		forwarderFactory portForwarderFactory
	}
	tests := []struct {
		name    string
		fields  fields
		payload func() error
		wantErr func(t *testing.T, err error)
	}{
		{
			name: "should fail to create dialer",
			fields: fields{
				dialerFactory: failingDialerFactory(t),
			},
			wantErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorContains(t, err, "failed to create dialer for port-forward")
			},
		},
		{
			name: "should fail to create kubernetes port forwarder",
			fields: fields{
				dialerFactory:    succeedingDialerFactory(t),
				forwarderFactory: failingPortForwarderFactory(t),
			},
			wantErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorContains(t, err, "failed to create kubernetes port forwarder")
			},
		},
		{
			name: "should fail to create port forward",
			fields: fields{
				dialerFactory:    succeedingDialerFactory(t),
				forwarderFactory: portForwarderFactoryWithFailingPortForwarder(t),
			},
			wantErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorContains(t, err, "could not create port-forward")
			},
		},
		{
			name: "should fail during payload execution",
			fields: fields{
				dialerFactory:    succeedingDialerFactory(t),
				forwarderFactory: portForwarderFactoryWithSucceedingPortForwarder(t),
			},
			payload: func() error {
				return assert.AnError
			},
			wantErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorContains(t, err, "encountered error during port-forward")
			},
		},
		{
			name: "should succeed",
			fields: fields{
				dialerFactory:    succeedingDialerFactory(t),
				forwarderFactory: portForwarderFactoryWithSucceedingPortForwarder(t),
			},
			payload: func() error {
				return nil
			},
			wantErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			sut := &executivePortForwarder{
				dialerFactory:    tt.fields.dialerFactory,
				forwarderFactory: tt.fields.forwarderFactory,
			}

			// when
			err := sut.ExecuteWithPortForward(tt.payload)

			// then
			tt.wantErr(t, err)
		})
	}
}

func failingDialerFactory(t *testing.T) dialerFactory {
	t.Helper()

	factory := newMockDialerFactory(t)
	factory.EXPECT().create().Return(nil, assert.AnError)

	return factory
}

func succeedingDialerFactory(t *testing.T) dialerFactory {
	t.Helper()

	factory := newMockDialerFactory(t)
	factory.EXPECT().create().Return(newMockDialer(t), nil)

	return factory
}

func failingPortForwarderFactory(t *testing.T) portForwarderFactory {
	t.Helper()

	factory := newMockPortForwarderFactory(t)
	factory.EXPECT().create(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError)

	return factory
}

func portForwarderFactoryWithFailingPortForwarder(t *testing.T) portForwarderFactory {
	t.Helper()

	forwarder := newMockPortForwarder(t)
	forwarder.EXPECT().ForwardPorts().Return(assert.AnError)

	factory := newMockPortForwarderFactory(t)
	factory.EXPECT().create(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(forwarder, nil)

	return factory
}

func portForwarderFactoryWithSucceedingPortForwarder(t *testing.T) portForwarderFactory {
	t.Helper()

	factory := newMockPortForwarderFactory(t)
	factory.EXPECT().create(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).RunAndReturn(
		func(dialer dialer, stopCh chan struct{}, readyCh chan struct{}, writer io.Writer, writer2 io.Writer) (portForwarder, error) {
			forwarder := newMockPortForwarder(t)
			forwarder.EXPECT().ForwardPorts().RunAndReturn(func() error {
				readyCh <- struct{}{}
				return nil
			})
			return forwarder, nil
		},
	)

	return factory
}

func Test_createApiUrl(t *testing.T) {
	type args struct {
		host string
		pod  types.NamespacedName
	}
	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr func(t *testing.T, err error)
	}{
		{
			name: "should fail to create url",
			args: args{
				host: "\x00",
				pod: types.NamespacedName{
					Namespace: "test-namespace",
					Name:      "test-name",
				},
			},
			wantErr: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.ErrorContains(t, err, "net/url: invalid control character in URL")
			},
		},
		{
			name: "should create url correctly",
			args: args{
				host: "http://localhost:8080",
				pod: types.NamespacedName{
					Namespace: "test-namespace",
					Name:      "test-name",
				},
			},
			want: &url.URL{
				Scheme: "http",
				Host:   "localhost:8080",
				Path:   "/api/v1/namespaces/test-namespace/pods/test-name/portforward",
			},
			wantErr: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			got, err := createApiUrl(tt.args.host, tt.args.pod)

			// then
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNew(t *testing.T) {
	t.Run("should create executive port-forwarder", func(t *testing.T) {
		// given
		restConfig := &rest.Config{}
		pod := types.NamespacedName{
			Namespace: "test-ns",
			Name:      "test",
		}
		localPort, clusterPort := 45678, 1234

		// when
		actual := New(restConfig, pod, localPort, clusterPort)

		// then
		assert.NotNil(t, actual)
	})
}
