package azure

import (
	"encoding/base64"
	"strconv"

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

func (s *store) DeleteApp(appID string) error {
	entity := s.getAppReference(appID)
	if err := entity.Delete(true, nil); err != nil {
		return err
	}

	chatIDs, err := s.ListSubscribers(appID)
	if err != nil {
		return err
	}

	var eg errGroup
	for len(chatIDs) > 0 {
		var ids []int64
		ids, chatIDs = splitInt64Slice(chatIDs, 50)
		batch := s.subscription.NewBatch()
		for _, chatID := range ids {
			entities := s.getSubscriptionReference(chatID, appID)
			for _, entity := range entities {
				batch.DeleteEntity(entity, true)
			}
		}
		eg.Append(batch.ExecuteBatch())
	}
	return eg.Simplify()
}

func (s *store) ListApps() ([]storage.App, error) {
	entities, err := queryEntities(s.app, "")
	if err != nil {
		return nil, err
	}

	apps := make([]storage.App, 0, len(entities))
	for _, entity := range entities {
		apps = append(apps, appFromEntity(entity))
	}

	return apps, nil
}

func (s *store) ListSubscribers(id string) ([]int64, error) {
	filter := "PartitionKey eq 'app_" + id + "'"
	entities, err := queryEntities(s.app, filter)
	if err != nil {
		return nil, err
	}
	if len(entities) == 0 {
		return nil, nil
	}

	ids := make([]int64, 0, len(entities))
	for _, entity := range entities {
		chatID, err := strconv.ParseInt(entity.RowKey, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, chatID)
	}

	return ids, nil
}

func (s *store) ListSubscribedApps(chatID int64) ([]storage.App, error) {
	appIDs, err := s.listSubscribedAppID(chatID)
	if err != nil {
		return nil, err
	}
	if len(appIDs) == 0 {
		return nil, nil
	}

	apps := make([]storage.App, 0, len(appIDs))
	for _, appID := range appIDs {
		app, err := s.GetApp(appID)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func (s *store) listSubscribedAppID(chatID int64) ([]string, error) {
	filter := "PartitionKey eq 'chat_" + strconv.FormatInt(chatID, 10) + "'"
	entities, err := queryEntities(s.app, filter)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(entities))
	for _, entity := range entities {
		ids = append(ids, entity.RowKey)
	}

	return ids, nil
}

func (s *store) Subscribe(chatID int64, appID string) error {
	entities := s.getSubscriptionReference(chatID, appID)
	var eg errGroup
	for _, entity := range entities {
		eg.Append(entity.Insert(azure.EmptyPayload, nil))
	}
	return eg.Simplify()
}

func (s *store) Unsubscribe(chatID int64, appID string) error {
	entities := s.getSubscriptionReference(chatID, appID)
	var eg errGroup
	for _, entity := range entities {
		eg.Append(entity.Delete(true, nil))
	}
	return eg.Simplify()
}

func (s *store) UnsubscribeAll(chatID int64) error {
	appIDs, err := s.listSubscribedAppID(chatID)
	if err != nil {
		return err
	}

	var eg errGroup
	for len(appIDs) > 0 {
		var ids []string
		ids, appIDs = splitStringSlice(appIDs, 50)
		batch := s.subscription.NewBatch()
		for _, appID := range ids {
			entities := s.getSubscriptionReference(chatID, appID)
			for _, entity := range entities {
				batch.DeleteEntity(entity, true)
			}
		}
		eg.Append(batch.ExecuteBatch())
	}
	return eg.Simplify()
}

func (s *store) Close() error {
	return nil
}

func (s *store) getAppReference(id string) *azure.Entity {
	return s.app.GetEntityReference(id, "-")
}

func (s *store) getSubscriptionReference(chatID int64, appID string) []*azure.Entity {
	chatIDstring := strconv.FormatInt(chatID, 10)
	return []*azure.Entity{
		s.subscription.GetEntityReference("chat_"+chatIDstring, appID),
		s.subscription.GetEntityReference("app_"+appID, chatIDstring),
	}
}
