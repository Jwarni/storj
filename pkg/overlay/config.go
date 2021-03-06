// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package overlay

import (
	"context"
	"strings"
	"time"

	"github.com/zeebo/errs"
	"go.uber.org/zap"
	monkit "gopkg.in/spacemonkeygo/monkit.v2"

	"storj.io/storj/pkg/kademlia"
	"storj.io/storj/pkg/pb"
	"storj.io/storj/pkg/provider"
	"storj.io/storj/pkg/statdb"
	"storj.io/storj/pkg/storj"
	"storj.io/storj/pkg/utils"
	"storj.io/storj/storage"
	"storj.io/storj/storage/boltdb"
	"storj.io/storj/storage/redis"
)

var (
	mon = monkit.Package()
	// Error represents an overlay error
	Error = errs.Class("overlay error")
)

// Config is a configuration struct for everything you need to start the
// Overlay cache responsibility.
type Config struct {
	DatabaseURL     string        `help:"the database connection string to use" default:"bolt://$CONFDIR/overlay.db"`
	RefreshInterval time.Duration `help:"the interval at which the cache refreshes itself in seconds" default:"1s"`
}

// LookupConfig is a configuration struct for querying the overlay cache with one or more node IDs
type LookupConfig struct {
	NodeIDsString string `help:"one or more string-encoded node IDs, delimited by Delimiter"`
	Delimiter     string `help:"delimiter used for parsing node IDs" default:","`
}

// CtxKey used for assigning cache and server
type CtxKey int

const (
	ctxKeyOverlay CtxKey = iota
	ctxKeyOverlayServer
)

// Run implements the provider.Responsibility interface. Run assumes a
// Kademlia responsibility has been started before this one.
func (c Config) Run(ctx context.Context, server *provider.Provider) (
	err error) {
	defer mon.Task()(&ctx)(&err)

	kad := kademlia.LoadFromContext(ctx)
	if kad == nil {
		return Error.New("programmer error: kademlia responsibility unstarted")
	}

	sdb, ok := ctx.Value("masterdb").(interface {
		StatDB() statdb.DB
	})
	if !ok {
		return Error.New("unable to get master db instance")
	}

	driver, source, err := utils.SplitDBURL(c.DatabaseURL)
	if err != nil {
		return Error.Wrap(err)
	}

	var db storage.KeyValueStore

	switch driver {
	case "bolt":
		db, err = boltdb.New(source, OverlayBucket)
		if err != nil {
			return err
		}
		zap.S().Info("Starting overlay cache with BoltDB")
	case "redis":
		db, err = redis.NewClientFrom(c.DatabaseURL)
		if err != nil {
			return err
		}
		zap.S().Info("Starting overlay cache with Redis")
	default:
		return Error.New("database scheme not supported: %s", driver)
	}

	cache := NewOverlayCache(db, kad, sdb.StatDB())
	srv := NewServer(zap.L(), cache, kad)
	pb.RegisterOverlayServer(server.GRPC(), srv)

	ctx2 := context.WithValue(ctx, ctxKeyOverlay, cache)
	ctx2 = context.WithValue(ctx2, ctxKeyOverlayServer, srv)
	return server.Run(ctx2)
}

// LoadFromContext gives access to the cache from the context, or returns nil
func LoadFromContext(ctx context.Context) *Cache {
	if v, ok := ctx.Value(ctxKeyOverlay).(*Cache); ok {
		return v
	}
	return nil
}

// LoadServerFromContext gives access to the overlay server from the context, or returns nil
func LoadServerFromContext(ctx context.Context) *Server {
	if v, ok := ctx.Value(ctxKeyOverlayServer).(*Server); ok {
		return v
	}
	return nil
}

// ParseIDs converts the base58check encoded node ID strings from the config into node IDs
func (c LookupConfig) ParseIDs() (ids storj.NodeIDList, err error) {
	var idErrs []error
	idStrs := strings.Split(c.NodeIDsString, c.Delimiter)
	for _, s := range idStrs {
		id, err := storj.NodeIDFromString(s)
		if err != nil {
			idErrs = append(idErrs, err)
			continue
		}
		ids = append(ids, id)
	}
	if err := utils.CombineErrors(idErrs...); err != nil {
		return nil, err
	}
	return ids, nil
}
