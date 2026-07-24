package platform

import (
	"context"
	"errors"
)

func (p *Platform) Migrate(ctx context.Context) error {
	return errors.New("database migrations not implemented")
}
