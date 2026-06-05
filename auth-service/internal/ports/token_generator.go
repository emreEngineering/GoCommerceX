package ports

import "context"

type TokenGenerator interface {
	Generate(ctx context.Context, userID string, email string) (string, error)
}
