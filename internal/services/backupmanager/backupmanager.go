package backupmanager

import "context"

type Store interface {
	UploadFile(ctx context.Context, filePath string, alias string) error
	DownloadFile(string) (error, []byte)
}

type Backupmanager struct {
	store Store
}

func New(store Store) *Backupmanager {
	return &Backupmanager{
		store: store,
	}
}

func (b *Backupmanager) BackupFile(ctx context.Context, filePath string, alias string) error {
	err := b.store.UploadFile(ctx, "file path", "alias...")
	if err != nil {
		return err
	}

	return nil
}

func (b *Backupmanager) ListBackups() {

}

func (b *Backupmanager) ListCompleteBackups() {

}
