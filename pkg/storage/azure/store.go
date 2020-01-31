package azure

import (
	"encoding/base64"

	azure "github.com/Azure/azure-sdk-for-go/storage"
	"github.com/ne2blink/antenna/pkg/storage"
	uuid "github.com/satori/go.uuid"
)

func (s *store) CreateApp(app storage.App) (string, error) {
	id := base64.RawURLEncoding.EncodeToString(uuid.NewV4().Bytes())

	entity := s.getAppReference(id)
	entity.Properties = map[string]interface{}{
		"name":   app.Name,
		"secret": app.Secret,
	}
	if err := entity.Insert(azure.EmptyPayload, nil); err != nil {
		return "", err
	}

	return id, nil
}

func (s *store) UpdateApp(app storage.App) error {
	entity := s.getAppReference(app.ID)
	entity.Properties = map[string]interface{}{
		"name":   app.Name,
		"secret": app.Secret,
	}
	return entity.Update(true, nil)
}

func (s *store) GetApp(id string) (storage.App, error) {
	entity := s.getAppReference(id)
	if err := entity.Get(defaultTimeout, azure.MinimalMetadata, nil); err != nil {
		return storage.App{}, err
	}
	return appFromEntity(entity), nil
}

func (s *store) DeleteApp(id string) error {
	entity := s.getAppReference(id)
	return entity.Delete(true, nil)
}

func (s *store) ListApps() ([]storage.App, error) {
	entities, err := queryEntities(s.app)
	if err != nil {
		return nil, err
	}

	var apps []storage.App
	for _, entity := range entities {
		apps = append(apps, appFromEntity(entity))
	}

	return apps, nil
}

func (s *store) ListSubscribers(ID string) ([]int64, error) {
	panic("not implemented") // TODO: Implement
}

func (s *store) ListSubscribedApps(ChatID int64) ([]storage.App, error) {
	panic("not implemented") // TODO: Implement
}

func (s *store) Subscribe(ChatID int64, AppID string) error {
	panic("not implemented") // TODO: Implement
}

func (s *store) Unsubscribe(ChatID int64, AppID string) error {
	panic("not implemented") // TODO: Implement
}

func (s *store) UnsubscribeAll(ChatID int64) error {
	panic("not implemented") // TODO: Implement
}

func (s *store) Close() error {
	return nil
}

func (s *store) getAppReference(id string) *azure.Entity {
	return s.app.GetEntityReference(id, "-")
}
