package main

import (
	"strings"
	"fmt"
	"time"
	//"reflect"
	//"os"
	"github.com/PuerkitoBio/goquery"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)


type Url struct {
	Domain	string
	Dir_1		string
	ID			string
}

type OffUrl struct {
	Url
	Dir_2	string
}

type UserUrl struct {
	Url
}

type PrefUrl struct {
	Url
}


//ユーザーアカウント
var users_g []string

//最初に検索するユーザー
var f_user string

//最初に検索するオフ
var f_off string


//配列の要素を逆順にならべる
func arrayReverse(items []string) []string {
	if len(items) == 0 {
		return items
	}
	return append(arrayReverse(items[1:]), items[0])
}

//指定した値が配列の要素に存在するか
func inArray(data string, items []string) bool {
	for _, item := range items {
		if item == data { return true }
	}
	return false
}

func getIDs(doc *goquery.Document, selector string, url_type string, ids []string) []string {

	//セレクタでオフのタイトルに貼ってあるaタグを絞り込み
	doc.Find(selector).Each(func(_ int, s *goquery.Selection) {
			url, _ := s.Attr("href")
			
			var tmp []string
			tmp = make([]string, 0, 3)
			tmp = strings.Split(url, "/")
			if (len(tmp) == 3) && (tmp[1] == url_type) {
				ids = append(ids, tmp[2])
			}
	})

	return ids
}

//idが最大値のユーザーを取得
func getMaxIDUser(col *mgo.Collection) string {

	if f_user != "" { return f_user }

	user := new(User)
	users, err := col.Count()
	if err != nil { panic(err) }
	col.Find(nil).Skip(users - 10).Limit(1).One(&user)

	f_user = user.Name
	return user.Name
}

//idが最大値のオフを取得
func getIDOff(col *mgo.Collection, num int) string {

	off := new(Off)
	col.Find(nil).Skip(num).Limit(1).One(&off)
	f_off = off.ID
	return off.ID
}

func getOffCount(col *mgo.Collection) int {
	offs, err := col.Count()
	if err != nil { panic(err) }
	return offs
}

func parseHtml(url string, second time.Duration) *goquery.Document {
	//待機
	time.Sleep(second * time.Second)

	fmt.Println(url)

	doc, err := goquery.NewDocument(url)
	if err != nil { panic(err) }

	return doc
}

func getOffIDsFromTop(doc *goquery.Document) bool {

	off_ids := make([]string, 0, 30)
	off_ids = getIDs(doc, "tr > td > a", "off", off_ids)
	fmt.Println(off_ids)

	if off_ids == nil { return false }

	for _, off_id := range off_ids {
		getUserIDsFromOff(parseHtml(getUrl("off", off_id), 5), off_id)
	}

	return true
}


func getOffIDsFromUser(doc *goquery.Document, user_id string) bool {

	off_ids := make([]string, 0, 30)
	off_ids = getIDs(doc, "tr > td > a", "off", off_ids)
	fmt.Println(off_ids)

	if off_ids == nil { return false }

	off_col := getOffCollection(getDBSession())

	for _, off_id := range off_ids {

		//オフのIDを
		o := new(Off)
		err := off_col.Find(bson.M{"id": off_id, "host": user_id}).One(&o)
		if err != nil {
			err = off_col.Insert(&Off{ID: off_id, Host: user_id, Timestamp: time.Now().Unix()})
		} else {
			fmt.Println("id:"+off_id+" host:"+user_id+" is already registed.")
			continue
		}

		getUserIDsFromOff(parseHtml(getUrl("off", off_id), 5), off_id)
	}

	return true
}

//テニスオフのページをスクレイピングして、参加者のユーザーIDを取得
func getUserIDsFromOff(doc *goquery.Document, off_id string) bool {

	user_ids := make([]string, 0, 10)
	user_ids = getIDs(doc, "div > table > tbody > tr > td > a", "profile", user_ids)

	//user_ids配列の一番目が主催者になってしまうので、主催者のidしかとれていない場合はfalse
	if len(user_ids) < 2 { return false }

	fmt.Println(arrayReverse(user_ids))

	user_col := getUserCollection(getDBSession())
	off_visitor_col := getOffVisitorCollection(getDBSession())
	//user_ids配列の一番目が主催者になってしまうので、配列を逆順にしておく
	for _, user_id := range arrayReverse(user_ids) {

		u := new(User)
		err := user_col.Find(bson.M{"name": user_id}).One(&u)
		if err == nil {
			fmt.Println(user_id+" is already registed.")
		} else {
			//データベースにアカウント情報を登録
			err = user_col.Insert(&User{Name: user_id, Timestamp: time.Now().Unix()})
			if err != nil { panic(err) }
		}
		//データベースに参加したオフのIDと参加者のIDを紐づけて保存
		ov := new(OffVisitor)
		err = off_visitor_col.Find(bson.M{"offid": off_id, "userid": user_id}).One(&ov)
		if err != nil {
			err = off_visitor_col.Insert(&OffVisitor{OffID: off_id, UserID: user_id, Timestamp: time.Now().Unix()})
		} else {
			fmt.Println("off_id:"+off_id+" user_id:"+user_id+" is already registed.")
		}


		fmt.Println(user_id)

		//アカウントの開催しているオフのIDを取得
		getOffIDsFromUser(parseHtml(getUrl("user", user_id), 7), user_id)
	}

	return true
}

//アクセスするURLを設定
func getUrl(url_type string, first_id string) string {

	var access_url string
	if url_type == "off" {
		url := OffUrl{}
		url.Domain = TENNISOFF
		url.Dir_1 = "off"
		url.ID = first_id
		url.Dir_2 = "member"

		access_url = url.Domain+"/"+url.Dir_1+"/"+url.ID+"/"+url.Dir_2
	} else if url_type == "user" {
		url := UserUrl{}
		url.Domain = TENNISOFF
		url.Dir_1 = "profile"
		url.ID = first_id

		access_url = url.Domain+"/"+url.Dir_1+"/"+url.ID
	} else if url_type == "pref" {
		url := PrefUrl{}
		url.Domain = TENNISOFF
		url.Dir_1 = "pref"
		url.ID = first_id
		access_url = url.Domain+"/"+url.Dir_1+"/"+url.ID
	}

	return access_url

}


func main() {

	/*
	for i:=getOffCount(getOffCollection(getDBSession())); i > 0; i-- {
		off_id := getIDOff(getOffCollection(getDBSession()), i)
		getUserIDsFromOff(parseHtml(getUrl("off", off_id), 5), off_id)
	}
	*/
	//getUserIDsFromOff(parseHtml(getUrl("off", "1307855"), 5), "1307855")
	getOffIDsFromTop(parseHtml(getUrl("pref", "12"), 5))
}
