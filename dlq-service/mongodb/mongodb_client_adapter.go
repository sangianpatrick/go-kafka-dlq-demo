package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ClientAdapter is a concrete struct of mongodb client adapter.
type ClientAdapter struct {
	c *mongo.Client
}

// NewClientAdapter is a constructor.
func NewClientAdapter(mongoClient *mongo.Client) Client {
	return &ClientAdapter{
		c: mongoClient,
	}
}

// Connect initializes the Client by starting background monitoring goroutines.
// If the Client was created using the NewClient function, this method must be called before a Client can be used.
//
// Connect starts background goroutines to monitor the state of the deployment and does not do any I/O in the main
// goroutine. The Client.Ping method can be used to verify that the connection was created successfully.
func (c *ClientAdapter) Connect(ctx context.Context) (err error) {
	err = c.c.Connect(ctx)
	return
}

// Database returns a handle for a database with the given name configured with the given DatabaseOptions.
func (c *ClientAdapter) Database(name string, opts ...*options.DatabaseOptions) (db Database) {
	mongoDatabase := c.c.Database(name, opts...)
	db = &DatabaseAdapter{db: mongoDatabase}
	return
}

// Disconnect closes sockets to the topology referenced by this Client. It will
// shut down any monitoring goroutines, close the idle connection pool, and will
// wait until all the in use connections have been returned to the connection
// pool and closed before returning. If the context expires via cancellation,
// deadline, or timeout before the in use connections have returned, the in use
// connections will be closed, resulting in the failure of any in flight read
// or write operations. If this method returns with no errors, all connections
// associated with this Client have been closed.
func (c *ClientAdapter) Disconnect(ctx context.Context) (err error) {
	err = c.c.Disconnect(ctx)
	return
}
