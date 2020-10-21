package objectstore

import (
	"context"
	"fmt"

	crd "cloud.redhat.com/clowder/v2/apis/cloud.redhat.com/v1alpha1"
	"cloud.redhat.com/clowder/v2/controllers/cloud.redhat.com/config"
	"cloud.redhat.com/clowder/v2/controllers/cloud.redhat.com/errors"
	"cloud.redhat.com/clowder/v2/controllers/cloud.redhat.com/providers"
	p "cloud.redhat.com/clowder/v2/controllers/cloud.redhat.com/providers"
)

// ObjectStoreProvider is the interface for apps to use to configure object
// stores
type ObjectStoreProvider interface {
	p.Configurable
	CreateBuckets(app *crd.ClowdApp) error
}

func GetObjectStore(c *p.Provider) (ObjectStoreProvider, error) {
	objectStoreMode := c.Env.Spec.Providers.ObjectStore.Mode
	switch objectStoreMode {
	case "minio":
		return NewMinIO(c)
	case "app-interface":
		return &AppInterfaceObjectstoreProvider{Provider: *c}, nil
	default:
		errStr := fmt.Sprintf("No matching object store mode for %s", objectStoreMode)
		return nil, errors.New(errStr)
	}
}

func RunAppProvider(provider providers.Provider, c *config.AppConfig, app *crd.ClowdApp) error {
	objectStoreProvider, err := GetObjectStore(&provider)

	if err != nil {
		return err
	}

	err = objectStoreProvider.CreateBuckets(app)

	if err != nil {
		return err
	}

	objectStoreProvider.Configure(c)
	return nil
}

func RunEnvProvider(provider providers.Provider) error {
	_, err := GetObjectStore(&provider)

	if err != nil {
		return err
	}

	return nil
}

type mockBucket struct {
	Name        string
	Exists      bool
	CreateError error
	ExistsError error
}

type mockBucketHandler struct {
	hostname              string
	port                  int
	accessKey             *string
	secretKey             *string
	wantCreateClientError bool
	ExistsCalls           []string
	MakeCalls             []string
	MockBuckets           []mockBucket
}

func (c *mockBucketHandler) Exists(ctx context.Context, bucketName string) (bool, error) {
	// track the calls to this mock func
	c.ExistsCalls = append(c.ExistsCalls, bucketName)

	for _, mockBucket := range c.MockBuckets {
		if mockBucket.Name == bucketName {
			if mockBucket.ExistsError == nil {
				return mockBucket.Exists, nil
			}
			return mockBucket.Exists, mockBucket.ExistsError
		}
	}
	// todo: really we should error out of the test here if there's no MockBuckets
	return false, nil
}

func (c *mockBucketHandler) Make(ctx context.Context, bucketName string) (err error) {
	// track the calls to this mock func
	c.MakeCalls = append(c.MakeCalls, bucketName)

	for _, mockBucket := range c.MockBuckets {
		if mockBucket.Name == bucketName {
			if mockBucket.CreateError == nil {
				return nil
			}
			return mockBucket.CreateError
		}
	}
	// todo: really we should error out of the test here if there's no MockBuckets
	return nil
}

func (c *mockBucketHandler) CreateClient(
	hostname string, port int, accessKey *string, secretKey *string,
) error {
	if c.wantCreateClientError == true {
		return errors.New("create client error")
	}
	c.hostname = hostname
	c.port = port
	c.accessKey = accessKey
	c.secretKey = secretKey
	return nil
}
