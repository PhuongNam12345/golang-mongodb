package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)


type Person struct {
Fullname string `json:"fullname"`
Email string `json:"email"`
Password string `json:"password"`
}

func CORSMiddleware() gin.HandlerFunc {
return func(c *gin.Context) {

c.Header("Access-Control-Allow-Origin", "*")
c.Header("Access-Control-Allow-Credentials", "true")
c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,accept, origin, Cache-Control, X-Requested-With")
c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")
if c.Request.Method == "OPTIONS" {
c.AbortWithStatus(204)
return
}

c.Next()
}
}

var collection, collection_admin *mongo.Collection
func init() {
// Thiết lập thông tin kết nối MongoDB
clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
client, err := mongo.Connect(context.TODO(), clientOptions)
if err != nil {
log.Fatal(err)
}
err = client.Ping(context.TODO(), nil)
if err != nil {
log.Fatal(err)
}
// Lấy collection trong MongoDB
collection = client.Database("DB-test").Collection("customer")
collection_admin = client.Database("DB-test").Collection("admin")

}
func main(){// Dữ liệu để thêm vào MongoDB
gin.SetMode(gin.ReleaseMode)
r := gin.New()
// setup gin middleware
r.Use(gin.Recovery())
r.Use(CORSMiddleware())
// Register
r.POST("/register", func(c *gin.Context) {
var person Person
if err := c.BindJSON(&person); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
return
}
// Lưu trữ dữ liệu người dùng vào MongoDB
_, err := collection_admin.InsertOne(context.TODO(), person)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}
c.JSON(http.StatusOK, gin.H{"message": "Đăng ký thành công"})
})
// Login
r.POST("/api/login", func(c *gin.Context) {
var loginData struct {
Email string `bson:"email" binding:"required"`
Password string `bson:"password" `
}

if err := c.ShouldBindJSON(&loginData); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}
err := collection_admin.FindOne(context.TODO(), bson.M{"email":loginData.Email, "password":
loginData.Password}).Decode(&loginData)
if err != nil { 
c.JSON(http.StatusNotFound, gin.H{"message": "Đăng nhập không thành công"})
return
}else{
c.JSON(http.StatusOK, gin.H{"message": "Đăng nhập thành công"})
}
// if isValidUser(loginData.Email, loginData.Password) {
// c.JSON(http.StatusOK, gin.H{"message": "Đăng nhập thành công"})
// } else {
// c.JSON(http.StatusUnauthorized, gin.H{"error": "Thông tin đăng nhập không chính xác"})
// }
// Thêm khách hàng
})
r.POST("/addcustomer", func(c *gin.Context) {
	var customerData struct {
		Fullname string `bson:"fullname"`
		Email string `bson:"email" `	
		Phone string `bson:"phone" `
		Address string `bson:"address" `
		}
	if err := c.ShouldBindJSON(&customerData); err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error" : "lỗi"})
	return
	}
	// Lưu trữ dữ liệu khách hàng vào MongoDB
	_, err := collection.InsertOne(context.TODO(), customerData)
	if err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Thêm thành công"})
	})
// hiển thị khách hàng
	// Truy vấn dữ liệu từ collection

// 	http.HandleFunc("/users", GetUsers)
// 	log.Fatal(http.ListenAndServe(":5000", nil));

// func GetUsers(w http.ResponseWriter, r *http.Request) {
// 	type User struct {
// 		Name  string
// 		Email string
// 	}	

// collection = client.Database("DB-test").Collection("test-collection")

// 	cursor, err := collection.Find(context.TODO(), nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer cursor.Close(context.TODO())

// 	var users []User
// 	for cursor.Next(context.TODO()) {
// 		var user User
// 		err := cursor.Decode(&user)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		users = append(users, user)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(users)
r.GET("/showcustomer", func(c *gin.Context) {

	type Customer struct {
		ID string `bson:"_id"`
		Fullname string `bson:"fullname"`
		Email string `bson:"email" `	
		Phone string `bson:"phone" `
		Address string `bson:"address" `
		}
	var people []Customer

	// Truy vấn dữ liệu từ MongoDB
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	// defer cursor.Close(context.TODO())

	// Lấy dữ liệu từ cursor và đưa vào slice people
	for cursor.Next(context.TODO()) {
		var customer Customer
		if err := cursor.Decode(&customer); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		people = append(people, customer)
	}
	c.JSON(http.StatusOK, people)
})
	// cur, err := collection.Find(context.TODO(), nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// var users []User
	// for cur.Next(context.TODO()) {
	// 	var user User
	// 	err := cur.Decode(&user)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	users = append(users, user) 
	// }
	// fmt.Println("Danh sách người dùng:")
	// for _, user := range users {
	// 	fmt.Println("fullname:", user.Name, "email:", user.Email)
	// }
r.Run(":5000")
}
// type Cusomer struct {
// 	Fullname  string `json:"fullname"`
// 	Email   string    `json:"email"`
// }
// func isValidUser(email, password string) bool {
// // Đây là nơi bạn thực hiện kiểm tra thông tin đăng nhập từ cơ sở dữ liệu
// // Trong ví dụ này, chỉ kiểm tra xem username và password có giống nhau không.
// return email == "users" && password == "pass"
// }

// Kiểm tra thông tin đăng nhập từ cơ sở dữ liệu
// Điều này thường liên quan đến truy vấn cơ sở dữ liệu để kiểm tra thông tin người dùng.

// Nếu thông tin đúng, trả về thành công
// if isValidUser(loginData.Email, loginData.Password) {
// c.JSON(http.StatusOK, gin.H{"message": "Đăng nhập thành côngs"})

// else {
// c.JSON(http.StatusUnauthorized, gin.H{"error": "Thông tin đăng nhập không chính xác"})
// }
// r.Run(":5000")
// }

// func isValidUser(email, password string) bool {

// // Đây là nơi bạn thực hiện kiểm tra thông tin đăng nhập từ cơ sở dữ liệu
// // Trong ví dụ này, chỉ kiểm tra xem username và password có giống nhau không.
// return email == "users" && password == "pass"

// }