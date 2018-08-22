package main

import (
	mgo "gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

const TENNISOFF = "somesite"
const DBNAME = "db"
const DBPASS = "mongodb://localhost/"+DBNAME

type User struct {
	Name			string
	Timestamp	int64
}

type OffVisitor struct {
	OffID			string
	UserID		string
	Timestamp	int64
}

type Off struct {
	ID				string
	Host			string
	Timestamp	int64
}

//MongDB
var session *mgo.Session

//MongoDBに接続
//シングルトンパターン
func getDBSession() *mgo.Session {

	if session != nil { return session }

	db_session, err := mgo.Dial(DBPASS)
	if err != nil {
		panic(err)
	}

	//グローバル変数に代入
	session = db_session

	return session
}

//userコレクションに接続
func getUserCollection(session *mgo.Session) *mgo.Collection {
	return session.DB(DBNAME).C("user")
}

//offコレクションに接続
func getOffCollection(session *mgo.Session) *mgo.Collection {
	return session.DB(DBNAME).C("off")
}

//off_visitorコレクションに接続
func getOffVisitorCollection(session *mgo.Session) * mgo.Collection {
	return session.DB(DBNAME).C("off_visitor")
}
