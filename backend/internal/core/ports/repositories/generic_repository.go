package ports_repositories

import "context"

// GenericRepository เป็น Interface พอร์ตกลางสำหรับทำ CRUD
type GenericRepository[Tentity any, Tfilter any] interface {
	Add(ctx context.Context, entity *Tentity) error
	AddRange(ctx context.Context, entities []*Tentity) error
	Get(ctx context.Context, filters *Tfilter) (*Tentity, error)
	GetAll(ctx context.Context, filters *Tfilter) ([]*Tentity, error)
	Update(ctx context.Context, entity *Tentity) error
	UpdateRange(ctx context.Context, entities []*Tentity) error
	Delete(ctx context.Context, entity *Tentity) error
	DeleteRange(ctx context.Context, entities []*Tentity) error
}
