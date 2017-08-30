package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gopkg.in/couchbase/gocb.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/redis.v4"
	"reflect"
	"sync"
	"testing"
	"time"
)

type Person struct {
	NAME  string
	PHONE string
}

type Men struct {
	Persons []Person
}

const (
	MongoURL = "172.17.0.3:27017"
)

func mongoTest(c chan int) {
	session, err := mgo.Dial(MongoURL) //连接数据库
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	db := session.DB("godb")     //数据库名称
	collection := db.C("person") //如果该集合已经存在的话，则直接返回

	//*****集合中元素数目********
	countNum, err := collection.Count()
	if err != nil {
		panic(err)
	}
	fmt.Println("Things objects count: ", countNum)

	//*******插入元素*******
	//temp := &Person{PHONE: "18811577546", NAME: "zhangzheHero"}
	//一次可以插入多个对象 插入两个Person对象

	//err = collection.Insert(&Person{"Ale", "+55 53 8116 9639"}, temp)
	if err != nil {
		panic(err)
	}

	//*****查询单条数据*******
	result := Person{}
	err = collection.Find(bson.M{"name": "Ale"}).One(&result)
	fmt.Println("One:", result.NAME, result.PHONE)

	//*****查询多条数据*******
	var personAll Men //存放结果
	iter := collection.Find(nil).Iter()
	for iter.Next(&result) {
		fmt.Printf("Result: %v, %v\n", result.NAME, result.PHONE)
		personAll.Persons = append(personAll.Persons, result)
	}
	//*******更新数据*********
	err = collection.Update(bson.M{"name": "ccc"}, bson.M{"$set": bson.M{"name": "ddd"}})
	err = collection.Update(bson.M{"name": "Ale"}, bson.M{"$set": bson.M{"phone": "12345678"}})
	err = collection.Update(bson.M{"name": "aaa"}, bson.M{"phone": "1245", "name": "bbb"})
	//*******删除单条匹配数据****
	err = collection.Remove(bson.M{"name": "zhangzheHero"})
	//*******删除所有匹配数据******
	_, err = collection.RemoveAll(bson.M{"name": "Ale"})
	if err != nil {
		fmt.Println("Delete failure!")
	}
	err = collection.Find(bson.M{"name": "Ale"}).One(&result)
	if err == nil {
		fmt.Println("One:", result.NAME, result.PHONE)
	} else {
		fmt.Println("No person named Ale")
	}
	close(c)
}

type Product struct {
	Code  string
	Price uint
}

func sqlite3Delete(db *gorm.DB, query, args string, limit int, o interface{}) {
	tp := reflect.TypeOf(o)
	tn := tp.Name()
	if tn == "Product" {
		products := make([]Product, 100)
		db.Where(query, args).Limit(limit).Find(&Product{})
		for _, r := range products {
			//fmt.Println(r.Price)
			db.Delete(&r)
		}
	}
}

func sqlite3Test() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	//db.Create(&Product{Code: "L1678", Price: 124})

	// Read
	product := &Product{}
	db.Where("code = ?", "L1678").First(product) // find product with code L1678
	fmt.Println(product.Price)
	products := make([]Product, 100)
	db.Model(Product{}).Find(&products)
	fmt.Println(len(products))
	for _, p := range products {
		//fmt.Printf("%v -- %v\n", p.Code, p.Price)
		fmt.Printf("%5s : %d\n", p.Code, p.Price)
	}
	// Update - update products  who's price=18 to 2000
	db.Model(Product{}).Where("price = ?", 18).Update("Price", 10080)
	db.Table("products").Where("id IN (?)", []int{28}).Updates(map[string]interface{}{"price": 1809})
	//db.Save(&product)
	db.Where("id = ?", 1).First(product)
	fmt.Println(product.Price)
	// Delete - delete product
	//db.Delete()
	//db.Where("code = ?", "L1678").Limit(2).Find(products)
	//db.Where("code = ?", "L1678").Limit(1).Delete(Product{})
	/**
	query := "code = ?"
	args := "L1678"
	limit := 1
	sqlite3Delete(db, query, args, limit, &Product{})
	**/
	//db1.Delete(&Product{}
}

func mysqlTest() {
	ms, err := gorm.Open("mysql", "vrit:1234@tcp(172.17.0.5:3306)/golang?charset=utf8&parseTime=True&loc=Local")
	fmt.Println(err)
	defer ms.Close()
	//Create table
	ms.AutoMigrate(&Product{})
	//insert one record
	product := &Product{"sustc", 1234}
	ms.Create(product)
	p := &Product{}
	//find one record
	ms.Where("code = ?", "sustc").First(p)
	fmt.Println(p.Price)
	//update one record
	ms.Model(&Product{}).Where("code = ?", "sustc").Update("price", 10086)
	ms.Where("code = ?", "sustc").First(p)
	fmt.Println(p.Price)
	//delet any count records
	//ms.Where("code = ?", "sustc").Limit(2).Delete(&Product{})
}

func postgresTest() {
	ps, err := gorm.Open("postgres", "host=172.17.0.4 user=vrit dbname=test sslmode=disable password=1234")
	fmt.Println(err)
	defer ps.Close()
	ps.AutoMigrate(&Product{})
	//insert one record
	//product := &Product{"sustc", 129}
	ps.Create(&Product{"sustc", 129})
	ps.Create(&Product{"sutc", 19})
	ps.Create(&Product{"suc", 29})
	p := &Product{}
	//find one record
	ps.Where("code = ?", "sustc").First(p)
	fmt.Println(p.Price)
	//update one record
	//ps.Model(&Product{}).Where("code = ?", "sustc").Update("price", 100)
	ps.Where("code = ?", "sustc").First(p)
	fmt.Println(p.Price)
	//delet any count records
	products := []Product{}
	ps.Where("price < ?", 500).Limit(2).Find(&products)
	fmt.Println(len(products))
	for _, pp := range products {
		fmt.Println(pp.Price)
		//这里使用primary key作为每条记录的唯一标识符，确保删除记录数不会超过指定大小
		ps.Where("code = ?", pp.Code).Delete(&pp)
	}
}

// 创建 redis 客户端
func createClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "172.17.0.6:6379",
		Password: "",
		DB:       0,
		PoolSize: 5,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	return client
}

// String 操作
func stringOperation(client *redis.Client) {
	// 第三个参数是过期时间, 如果是0, 则表示没有过期时间.
	err := client.Set("name", "xys", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("name").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("name", val)

	// 这里设置过期时间.
	err = client.Set("age", "20", 1*time.Second).Err()
	if err != nil {
		panic(err)
	}

	client.Incr("age") // 自增
	client.Incr("age") // 自增
	client.Decr("age") // 自减

	val, err = client.Get("age").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("age", val) // age 的值为21

	// 因为 key "age" 的过期时间是一秒钟, 因此当一秒后, 此 key 会自动被删除了.
	time.Sleep(1 * time.Second)
	val, err = client.Get("age").Result()
	if err != nil {
		// 因为 key "age" 已经过期了, 因此会有一个 redis: nil 的错误.
		fmt.Printf("error: %v\n", err)
	}
	fmt.Println("age", val)
}

// list 操作
func listOperation(client *redis.Client) {
	client.RPush("fruit", "apple")               //在名称为 fruit 的list尾添加一个值为value的元素
	client.LPush("fruit", "banana")              //在名称为 fruit 的list头添加一个值为value的 元素
	length, err := client.LLen("fruit").Result() //返回名称为 fruit 的list的长度
	if err != nil {
		panic(err)
	}
	fmt.Println("length: ", length) // 长度为2

	value, err := client.LPop("fruit").Result() //返回并删除名称为 fruit 的list中的首元素
	if err != nil {
		panic(err)
	}
	fmt.Println("fruit: ", value)

	value, err = client.RPop("fruit").Result() // 返回并删除名称为 fruit 的list中的尾元素
	if err != nil {
		panic(err)
	}
	fmt.Println("fruit: ", value)
}

// set 操作
func setOperation(client *redis.Client) {
	client.SAdd("blacklist", "Obama")     // 向 blacklist 中添加元素
	client.SAdd("blacklist", "Hillary")   // 再次添加
	client.SAdd("blacklist", "the Elder") // 添加新元素

	client.SAdd("whitelist", "the Elder") // 向 whitelist 添加元素

	// 判断元素是否在集合中
	isMember, err := client.SIsMember("blacklist", "Bush").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Is Bush in blacklist: ", isMember)

	// 求交集, 即既在黑名单中, 又在白名单中的元素
	names, err := client.SInter("blacklist", "whitelist").Result()
	if err != nil {
		panic(err)
	}
	// 获取到的元素是 "the Elder"
	fmt.Println("Inter result: ", names)

	// 获取指定集合的所有元素
	all, err := client.SMembers("blacklist").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("All member: ", all)
}

// hash 操作
func hashOperation(client *redis.Client) {
	client.HSet("user_xys", "name", "xys") // 向名称为 user_xys 的 hash 中添加元素 name
	client.HSet("user_xys", "age", "18")   // 向名称为 user_xys 的 hash 中添加元素 age

	// 批量地向名称为 user_test 的 hash 中添加元素 name 和 age
	client.HMSet("user_test", map[string]string{"name": "test", "age": "20"})
	// 批量获取名为 user_test 的 hash 中的指定字段的值.
	fields, err := client.HMGet("user_test", "name", "age").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("fields in user_test: ", fields)

	// 获取名为 user_xys 的 hash 中的字段个数
	length, err := client.HLen("user_xys").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("field count in user_xys: ", length) // 字段个数为2

	// 删除名为 user_test 的 age 字段
	client.HDel("user_test", "age")
	age, err := client.HGet("user_test", "age").Result()
	if err != nil {
		fmt.Printf("Get user_test age error: %v\n", err)
	} else {
		fmt.Println("user_test age is: ", age) // 字段个数为2
	}
}

// redis.v4 的连接池管理
func connectPool(client *redis.Client) {
	wg := sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < 100; j++ {
				client.Set(fmt.Sprintf("name%d", j), fmt.Sprintf("xys%d", j), 0).Err()
				client.Get(fmt.Sprintf("name%d", j)).Result()
			}

			fmt.Printf("PoolStats, TotalConns: %d, FreeConns: %d\n", client.PoolStats().TotalConns, client.PoolStats().FreeConns)
		}()
	}

	wg.Wait()
}

func redisTest() {
	client := createClient()
	defer client.Close()

	stringOperation(client)
	listOperation(client)
	setOperation(client)
	hashOperation(client)

	connectPool(client)
}

type User struct {
	Id        string   `json:"uid"`
	Email     string   `json:"email"`
	Interests []string `json:"interests"`
}

func couchdbTest() {
	cluster, _ := gocb.Connect("couchbase://localhost")
	bucket, _ := cluster.OpenBucket("default", "")

	bucket.Manager("", "").CreatePrimaryIndex("", true, false)

	bucket.Upsert("u:kingarthur",
		User{
			Id:        "kingarthur",
			Email:     "kingarthur@couchbase.com",
			Interests: []string{"Holy Grail", "African Swallows"},
		}, 0)

	// Get the value back
	var inUser User
	bucket.Get("u:kingarthur", &inUser)
	fmt.Printf("User: %v\n", inUser)

	// Use query
	query := gocb.NewN1qlQuery("SELECT * FROM default WHERE $1 IN interests")
	rows, _ := bucket.ExecuteN1qlQuery(query, []interface{}{"African Swallows"})
	var row interface{}
	for rows.Next(&row) {
		fmt.Printf("Row: %v", row)
	}
}

//TestDatabase 测试各个数据库
func TestDatabase(t *testing.T) {
	fmt.Println("*****MongoTest start...********")
	c1 := make(chan int)
	go mongoTest(c1)
	<-c1
	fmt.Println("*****sqlite3Test start...******")
	sqlite3Test()
	fmt.Println("*****postgresTest start...******")
	postgresTest()
	fmt.Println("*****mysqlTest start...******")
	mysqlTest()
	fmt.Println("*****redisTest start...******")
	redisTest()
}
