package ipfsapp

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

var admin adminConn

func init() {
	admin = adminConn{}
	registerAnonymousFunc("welcome", func(name interface{}) string {
		return "Hi " + name.(string) + "! Welcome to use IPFS"
	})
}

type apiRequest struct {
	Op       string
	Token    string
	Nounce   int
	Data     string
	IsCrypto bool
}

type apiReply struct {
	Status   string
	Data     string
	IsCrypto bool
}

//StartAPIServer 开启RESTful API服务
func StartAPIServer(c chan error) {
	router := gin.Default()

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/login/:info", login)
	router.GET("/welcome", func(c *gin.Context) {
		//firstname := c.DefaultQuery("firstname", "Guest")
		//lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")
		c.String(http.StatusOK, "Welcomo to use IPFS APP")
	})
	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
	router.POST("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})
	router.GET("/Get", getting)
	router.POST("/Post", posting)
	router.PUT("/Put", putting)
	router.DELETE("/Delete", deleting)
	router.PATCH("/Patch", patching)
	router.HEAD("/Head", head)
	router.OPTIONS("/Options", options)
	//router.Run(":8080")
	c <- http.ListenAndServe(":8080", router)
}

//Get 请求处理，解析成apiRequest后根据op找到要调用的数据，传入参数Data将返回值返回给用户
func getting(c *gin.Context) {
	reqstr := c.DefaultQuery("request", "{\"op\":\"welcome\",\"data\":\"bigdata\"}")
	//jsonstr, err := RsaDecrypt([]byte(reqstr))
	var err error
	jsonstr := []byte(reqstr)
	req := &apiRequest{}
	if err == nil {
		err = json.Unmarshal(jsonstr, req)
	}
	if err == nil {
		fun := anonymousMap[req.Op]
		c.String(http.StatusOK, fun(req.Data))
	}
}

func posting(c *gin.Context)  {}
func putting(c *gin.Context)  {}
func deleting(c *gin.Context) {}
func patching(c *gin.Context) {}
func head(c *gin.Context)     {}
func options(c *gin.Context)  {}
func login(c *gin.Context) {
	name := c.Param("info")
	c.String(http.StatusOK, "Hello %s", name)
}

func getPeersIP(c *gin.Context) {
	peersip := make(map[string]string)
	jsonstr, _ := json.Marshal(peersip)
	c.String(http.StatusOK, string(jsonstr))
}

func getIndexTree(veriferInfo, isCrypto string) string {
	//veriferInfo := c.Param("verifer")
	//isCrypto := c.DefaultQuery("crypto", "true")
	if isCrypto == "true" {
		jsonstr, err := RsaDecrypt([]byte(veriferInfo), privateKey)
		req := &apiRequest{}
		err = json.Unmarshal([]byte(jsonstr), req)
		if err == nil {
			return admin.getIndexTree(req.Data)
		}
	} else {
		req := &apiRequest{}
		err = json.Unmarshal([]byte(veriferInfo), req)
		return admin.getIndexTree(req.Data)
	}
	return ""
}

func updateIndexTree(c *gin.Context) {}

func welcome(name string) string {
	return "Hi " + name + "! Welcome to use IPFS"
}
