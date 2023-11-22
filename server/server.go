package server

/**
* created by mengqi on 2023/11/15
 */

import "context"

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
