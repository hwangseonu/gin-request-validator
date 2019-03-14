# gin-request-validator
gin-request-validator 는 gin 에서 request body 의 유효성을 검사하는 라이브러리입니다.
middleware 형식으로 검사하여 결과를 함수로 받을 수 있습니다.

## Usage
JsonRequiredMiddleware 를 사용할 url 의 Group 에 미들웨어로 등록합니다.  
만약 요청이 유효하지 않으면 상태코드 400를 반환하여 요청을 종료합니다.
```go
package main
    func main() {
    	r := gin.Default()
    	g := r.Group("/some")
        g.Use(gin_validator.JsonRequiredMiddleware(struct {
    	        Email string `json:"email"`
    	        Age   int    `json:"age" validate:"min=1 max=100"`
        }{}))
    	...
}
```
유효성 검사를 통과한 데이터를 Handler 에서 GetJsonData 함수를 통해 얻을 수 있습니다.
```go
package main

func handler(c *gin.Context) {
	req := gin_validator.GetJsonData(c).(Data)
	c.JSON(http.StatusOK, req)
}
```
직접적으로 유효성을 검사할 때는 ValidData 함수를 이용할 수 있습니다.  
유효성을 검사하여 문제가 있으면 error를 반환하며 문제가 없다면 nil을 반환합니다.  
```go
package main

type Data struct {
	Email string `json:"email"`
	Age   int    `json:"age" validate:"min=1 max=100"`
}

func main() {
	d := Data{Email: "aaa@email.com", Age: 1005}
	if err := gin_validator.ValidData(d, reflect.TypeOf(Data{})); err != nil {
		println("fail")
	} else {
		println("success")
	}
}
```
새로운 CustomValidator 함수를 등록하려면 RegisterValidator를 사용합니다.  
```go
package main

func CustomValidator(name string, data interface{}, interfaces ...interface{}) error {
	str, ok := data.(string)
	if !ok {
		return errors.New(name + "must string")
    }
	if str != "happy" {
		return errors.New("you must be happy")
	}
	return nil
}

func main() {
	gin_validator.RegisterValidator("custom", CustomValidator, "string")
	...
}
``` 