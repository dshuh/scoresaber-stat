package mongodb

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoConfig MongoConfig 설정 정보
type MongoConfig struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	ID       string `json:"id" yaml:"id"`
	Password string `json:"password" yaml:"password"`
	Database string `json:"database" yaml:"database"`
	PoolSize string `json:"poolsize" yaml:"poolsize"`
	IdleSize string `json:"idlesize" yaml:"idlesize"`
	Timeout  string `json:"timeout" yaml:"timeout"`
	ReadPref string `json:"readpref" yaml:"readpref"`
}

// SelectorUpdate Upsert를 위한 Selector/Update Pair
type SelectorUpdate struct {
	Selector bson.M
	Update   bson.M
}

// MongoDB 관련 기능 포함
// ReadPref : mgo.v2 library session.go 의 read preference 모드 참고 (mgo.Mode)
type MongoDB struct {
	Session  *mgo.Session
	Host     string
	ID       string
	Password string
	Database string
	PoolSize int
	Timeout  int
	ReadPref int
}

var gMongoDB *MongoDB
var gOnceMongoDB sync.Once

func GetMongoDB(params ...string) (*MongoDB, error) {
	var err error
	err = nil
	gOnceMongoDB.Do(func() {
		host := ""
		id := ""
		password := ""
		database := ""
		var poolSize int
		var timeout int
		readPref := -1

		if len(params) >= 1 {
			host = params[0]
		}

		if len(params) >= 2 {
			id = params[1]
		}

		if len(params) >= 3 {
			password = params[2]
		}

		if len(params) >= 4 {
			database = params[3]
		}

		if len(params) >= 5 {
			poolSize, err = strconv.Atoi(params[4])
			if err != nil {
				poolSize = 0
			}
		}

		if len(params) >= 6 {
			timeout, err = strconv.Atoi(params[5])
			if err != nil {
				timeout = 0
			}
		}

		if len(params) >= 7 {
			readPref, err = strconv.Atoi(params[6])
			if err != nil {
				readPref = -1
			}
		}

		gMongoDB, err = NewMongoDB(host, id, password, database, poolSize, timeout, readPref)
	})

	if err != nil {
		gMongoDB = nil
		return nil, err
	}

	if len(params) == 0 {
		if gMongoDB == nil {
			return nil, errors.New("DB is not initialized")
		}
		return gMongoDB, nil
	}

	return gMongoDB, nil
}

// NewMongoDB 세션을 생성한다.
func NewMongoDB(host string, id string, password string, database string, poolSize int, timeout int, readPref int) (mdb *MongoDB, err error) {

	mdb = new(MongoDB)

	hosts := strings.Split(host, ",")
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    hosts,
		Database: database,
		Username: id,
		Password: password,
	}

	if poolSize > 0 {
		mongoDBDialInfo.PoolLimit = poolSize
	}
	if timeout > 0 {
		mongoDBDialInfo.Timeout = time.Duration(timeout) * time.Second
	} else {
		mongoDBDialInfo.Timeout = 1 * time.Second
	}

	mdb.Session, err = mgo.DialWithInfo(mongoDBDialInfo)

	//mdb.Session, err = mgo.Dial(id + ":" + password + "@" + host + "/" + database)
	mdb.Host = host
	mdb.ID = id
	mdb.Password = password
	mdb.Database = database
	mdb.PoolSize = poolSize
	mdb.Timeout = timeout
	mdb.ReadPref = readPref

	//mdb.Session, err = mgo.Dial(mdb.Host)
	if err != nil {
		return nil, err
	}

	if readPref >= 0 {
		mdb.Session.SetMode(mgo.Mode(readPref), true)
	} else {
		mdb.Session.SetMode(mgo.Monotonic, true)
	}

	return mdb, err
}

//

// NewID is a new object id in string
func NewID() string {
	return bson.NewObjectId().Hex()
}

// Close the DB Connection
func (o *MongoDB) Close() {
	if o.Session != nil {
		o.Session.Close()
	}
}

// GetCol Get Collection
func (o *MongoDB) GetCol(colName string) (*mgo.Collection, *mgo.Session) {
	if o.Session == nil {
		var err error
		o, err = NewMongoDB(o.Host, o.ID, o.Password, o.Database, o.PoolSize, o.Timeout, o.ReadPref)
		if err != nil {
			panic(err)
		}
	}

	ses := o.Session.Copy()
	return ses.DB(o.Database).C(colName), ses
}

func (o *MongoDB) DropCollection(colName string) error {
	col, ses := o.GetCol(colName)
	defer ses.Close()

	return col.DropCollection()
}

func (o *MongoDB) GetCollections() ([]string, error) {
	if o.Session == nil {
		var err error
		o, err = NewMongoDB(o.Host, o.ID, o.Password, o.Database, o.PoolSize, o.Timeout, o.ReadPref)
		if err != nil {
			return nil, err
		}
	}

	ses := o.Session.Copy()
	return ses.DB(o.Database).CollectionNames()
}

func (o *MongoDB) DropDatabase() error {
	if o.Session == nil {
		var err error
		o, err = NewMongoDB(o.Host, o.ID, o.Password, o.Database, o.PoolSize, o.Timeout, o.ReadPref)
		if err != nil {
			return err
		}
	}

	ses := o.Session.Copy()
	return ses.DB(o.Database).DropDatabase()
}

// Remove doc match filter
func (o *MongoDB) Remove(colName string, filter interface{}) error {
	col, ses := o.GetCol(colName)
	defer ses.Close()
	err := col.Remove(filter)
	return err
}

// RemoveByID Remove ID
func (o *MongoDB) RemoveByID(colName string, id interface{}) error {
	col, ses := o.GetCol(colName)
	defer ses.Close()
	return col.RemoveId(id)
}

// RemoveAll all doc match filter
func (o *MongoDB) RemoveAll(colName string, filter interface{}) error {
	col, ses := o.GetCol(colName)
	defer ses.Close()
	_, err := col.RemoveAll(filter)
	return err
}

// Insert a document to a collection
func (o *MongoDB) Insert(colName string, d ...interface{}) error {
	col, ses := o.GetCol(colName)
	defer ses.Close()
	return col.Insert(d...)
}

// UpsertID Update Or Insert By ID
func (o *MongoDB) UpsertID(colName string, id interface{}, update interface{}) error {
	col, ses := o.GetCol(colName)
	defer ses.Close()
	_, err := col.UpsertId(id, update)
	return err
}

// Upsert Update Or Insert
func (o *MongoDB) Upsert(colName string, selector interface{}, update interface{}) error {
	col, ses := o.GetCol(colName)
	defer ses.Close()
	_, err := col.Upsert(selector, update)
	return err
}

// Update Data
func (o *MongoDB) Update(colName string, filter interface{}, update interface{}) error {
	col, ses := o.GetCol(colName)
	defer ses.Close()
	return col.Update(filter, update)
}

// UpdateID Data
func (o *MongoDB) UpdateID(colName string, id interface{}, d interface{}) error {
	col, ses := o.GetCol(colName)
	defer ses.Close()
	return col.UpdateId(id, bson.M{
		"$set": d,
	})
}

// UpdateAll Data
func (o *MongoDB) UpdateAll(colName string, filter interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	col, ses := o.GetCol(colName)
	defer ses.Close()
	return col.UpdateAll(filter, update)
}

// FindAll Data
func (o *MongoDB) FindAll(colName string, filter bson.M, d interface{}) error {
	query, ses := o.FindBy(colName, filter)
	defer ses.Close()
	return query.All(d)
}

// FindOne Data
func (o *MongoDB) FindOne(colName string, filter bson.M, d interface{}) error {
	query, ses := o.FindBy(colName, filter)
	defer ses.Close()
	return query.One(d)
}

// Count number of documents meeting the filter.
func (o *MongoDB) Count(colName string, filter bson.M) (int, error) {
	query, ses := o.FindBy(colName, filter)
	defer ses.Close()
	return query.Count()
}

// FindBy a *mgo.Query based on filter
func (o *MongoDB) FindBy(colName string, filter bson.M) (*mgo.Query, *mgo.Session) {
	col, ses := o.GetCol(colName)
	return col.Find(filter), ses
}

// FindByID a document by its id
func (o *MongoDB) FindByID(colName string, id interface{}, d interface{}) error {
	col, ses := o.GetCol(colName)
	defer ses.Close()
	return col.FindId(id).One(d)
}

// CreateIndex Create Index
func (o *MongoDB) CreateIndex(colName string, index mgo.Index) error {
	col, ses := o.GetCol(colName)
	defer ses.Close()

	err := col.EnsureIndex(index)
	return err
}

func (o *MongoDB) BulkInsert(colName string, docs ...interface{}) error {
	col, ses := o.GetCol(colName)
	defer ses.Close()

	bulk := col.Bulk()

	bulk.Insert(docs...)

	_, err := bulk.Run()
	if err != nil {
		return err
	}

	return nil
}

func (o *MongoDB) BulkUpsert(colName string, pairs ...SelectorUpdate) error {
	col, ses := o.GetCol(colName)
	defer ses.Close()

	bulk := col.Bulk()

	for _, pair := range pairs {
		bulk.Upsert(pair.Selector, pair.Update)
	}

	_, err := bulk.Run()
	if err != nil {
		return err
	}

	return nil
}

func (o *MongoDB) BulkRemove(colName string, selectors ...bson.M) error {
	data := make([]interface{}, len(selectors))
	for i := range selectors {
		data[i] = selectors[i]
	}

	col, ses := o.GetCol(colName)
	defer ses.Close()

	bulk := col.Bulk()
	bulk.Remove(data...)

	_, err := bulk.Run()
	if err != nil {
		return err
	}

	return nil
}

func (o *MongoDB) BulkRemoveAll(colName string, selectors ...bson.M) error {
	data := make([]interface{}, len(selectors))
	for i := range selectors {
		data[i] = selectors[i]
	}

	col, ses := o.GetCol(colName)
	defer ses.Close()

	bulk := col.Bulk()
	bulk.RemoveAll(data...)

	_, err := bulk.Run()
	if err != nil {
		return err
	}

	return nil
}

func (o *MongoDB) BulkUpdate(colName string, pairs ...SelectorUpdate) error {
	col, ses := o.GetCol(colName)
	defer ses.Close()

	bulk := col.Bulk()

	for _, pair := range pairs {
		bulk.Update(pair.Selector, pair.Update)
	}

	_, err := bulk.Run()
	if err != nil {
		return err
	}

	return nil
}

func (o *MongoDB) BulkUpdateAll(colName string, pairs ...SelectorUpdate) error {
	col, ses := o.GetCol(colName)
	defer ses.Close()

	bulk := col.Bulk()

	for _, pair := range pairs {
		bulk.UpdateAll(pair.Selector, pair.Update)
	}

	_, err := bulk.Run()
	if err != nil {
		return err
	}

	return nil
}
