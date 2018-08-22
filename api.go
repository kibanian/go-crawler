package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	//"io/ioutil"
	//mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"github.com/PuerkitoBio/goquery"
	//"encoding/json"
	//"reflect"
)

const SERVER = "http://medical-map.shop/api/"
const USERDIR = "user"
const OFFDIR = "off"
const OFFVISITORDIR = "off_visitor"


//データをサイトサーバにPOST
func postData(values url.Values, url string) bool {

	req, err := http.NewRequest("POST", url, strings.NewReader(values.Encode()))
	if err != nil {
		panic(err)
	}

	// Content-Type 設定
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")


	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	/*
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		panic(err)
	}
	*/

	return true
}

//ユーザーのデータをポストする
func postUserData(name string) bool {
	data := map[string]string{"user_id": name}
	return postData(setValues(data), SERVER+USERDIR+"/")
}

/*
//オフ-参加者データをポストする
func postOffVisitorData(data map[string]string) bool {
	return postData(setValues(data), SERVER+OFFVISITORDIR+"/")
}

*/

//オフデータをポストする
func postOffData(id string, host string) bool {
	data := map[string]string{"id": id, "host":host}
	return postData(setValues(data), SERVER+OFFDIR+"/")
}

//データをポストする
func setValues(data map[string]string) url.Values {
	values := url.Values{}
	for key, value := range data {
		values.Add(key, value)
	}

	return values
}


//データベースからユーザーのデータを取得する
func findUserData(timestamp int64) []User {

	col := getUserCollection(getDBSession())
	var data []User
	err := col.Find(bson.M{"timestamp": bson.M{"$gt": timestamp}}).Limit(10).All(&data)

	if err != nil {
		panic(err)
	}

	return data
}

/*
//データベースからオフのデータを取得する
func findOffData(timestamp int64) [string]string {
	col := getOffCollection(getDBSession())
	var data []User
	err := col.Find(bson.M{"timestamp": bson.M{"$gt": timestamp}}).Limit(10).All(&data)

	if err != nil {
		panic(err)
	}

	return data
}

//データベースからオフ-参加者のデータを取得する
func findOffVisitorData(timestamp int64) [string]string {
	return findData("off_visitor", timestamp)
}

//サイトのサーバーから最新のユーザーデータを取得
func getLatestUserData() {
	return getLatestData(SERVER+USERDIR+"/")
}

//サイトのサーバーから最新のオフデータを取得
func getLatestOffData() {
	return getLatestData(SERVER+OFFDIR+"/")
}

//サイトのサーバーから最新のオフ-参加者データを取得
func getLatestOffVisitorData() {
	return getLatestData(SERVER+OFFVISITORDIR+"/")
}

//サイトのサーバーから最新のデータを取得
func getLatestData(url string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		panic(err)
	}

	return doc

}

//JSONからデータを取り出す
func decodeJson(doc *goquery.Document) string {

	type Data struct {
		id string
	}

	var data Data
	err := json.Unmarshal(doc, &data)
	if err != nil {
		panic(err)
	}

	return data.id

}

*/

func main() {

	data := findUserData(1)

	fmt.Println(data)

	for key, item := range data {
		fmt.Println(key)
		postUserData(item.Name)
	}

	/*
	data = findOffData(1)

	fmt.Println(data)

	for key, item := range data {
		fmt.Println(key)
		postOffData(item.ID, item.Host)
	}
	*/

}