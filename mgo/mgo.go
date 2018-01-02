package mgo

import (
	"time"

	"github.com/influx6/mgodataset/config"
	mgo "gopkg.in/mgo.v2"
)

// mongoServer defines a mongo connection manager that builds
// allows usage of a giving configuration to generate new mongo
// sessions and database instances.
type mongoServer struct {
	config.MongoDBConf
	master *mgo.Session
}

// New returns a new session and database from the giving configuration.
func (m *mongoServer) New(isread bool) (*mgo.Database, *mgo.Session, error) {
	if m.master != nil && isread {
		copy := m.master.Copy()
		return copy.DB(m.MongoDBConf.DB), copy, nil
	}

	if m.master != nil && !isread {
		clone := m.master.Clone()
		return clone.DB(m.MongoDBConf.DB), clone, nil
	}

	ses, err := getSession(m.MongoDBConf)
	if err != nil {
		return nil, nil, err
	}

	m.master = ses

	if isread {
		copy := m.master.Copy()
		return copy.DB(m.MongoDBConf.DB), copy, nil
	}

	clone := m.master.Copy()
	return clone.DB(m.MongoDBConf.DB), clone, nil
}

// getSession attempts to retrieve the giving session for the given config.
func getSession(config config.MongoDBConf) (*mgo.Session, error) {
	info := mgo.DialInfo{
		Addrs:    []string{config.Host},
		Timeout:  60 * time.Second,
		Database: config.AuthDB,
		Username: config.User,
		Password: config.Password,
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	ses, err := mgo.DialWithInfo(&info)
	if err != nil {
		return nil, err
	}

	ses.SetMode(mgo.Monotonic, true)

	return ses, nil
}
