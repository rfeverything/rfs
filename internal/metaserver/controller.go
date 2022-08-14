package metaserver

import "context"

type Controller struct {
	ms *MetaServer
}

func NewController(ctx context.Context) *Controller {
	c := &Controller{}
	go c.tick(ctx)
	return c
}

func (c *Controller) tick(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if c.ms.Store.isleader {

			}
		}
	}
}
