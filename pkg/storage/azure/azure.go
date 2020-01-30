package azure

import "github.com/ne2blink/antenna/pkg/storage"

type store struct {
}

func (s *store) CreateApp(_ storage.App) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (s *store) UpdateApp(_ storage.App) error {
	panic("not implemented") // TODO: Implement
}

func (s *store) GetApp(ID string) (storage.App, error) {
	panic("not implemented") // TODO: Implement
}

func (s *store) DeleteApp(ID string) error {
	panic("not implemented") // TODO: Implement
}

func (s *store) ListApps() ([]storage.App, error) {
	panic("not implemented") // TODO: Implement
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
	panic("not implemented") // TODO: Implement
}

func newStore(options map[string]interface{}) (storage.Store, error) {
	return &store{}, nil
}

func init() {
	storage.Register("azure", newStore)
}
