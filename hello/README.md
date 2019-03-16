# gin-validator

gin-validator는 gin에서 request body의 유효성을 검사하는 라이브러리입니다.  

## Usage

구조체를 정의할 때 각 필드에 validate태그를 이용하여 유효성 검사에 대한 정보를 기입합니다.  
binding 태그를 이용하여 꼭 필요한 필드를 나타낼 수 있습니다.
```go
package main

type Data struct {
	Email string `json:"email" validate:"email" binding:"required"`
	Age   int    `json:"age" validate:"min=1 max=100"`
}
```

JsonRequiredMiddleware를 url의 미들웨어로 등록하여 유효성검사를 할 수 있습니다.
만약 요청이 유효하지 않으면 상태코드 400을 반환하며 요청을 종료합니다.
```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-request-validator"
)

func main() {
	e := gin.Default()
	g := e.Group("/awesome")

	//인자값으로 정의된 구조체의 빈 인스턴스를 넘겨줍니다
	g.Use(gin_validator.JsonRequiredMiddleware(Data{}))
	...
}
```

유효성 검사를 통과한 데이터는 Handler에서 GetJsonData함수를 통해 얻을 수 있습니다.
```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-request-validator"
	"net/http"
)

func Handler(c *gin.Context) {
	req := gin_validator.GetJsonData(c).(Data)
	c.JSON(http.StatusOK, req)
}
```

## 제공되는 validator

- min : 구조체 필드가 최소값보다 크거나 같은지 확인합니다. (ex: validate:"min=1")
- max : 구조체 필드가 최대값보다 작거나 같은지 확인합니다. (ex: validate:"max=100")
- email : 구조체 필드가 이메일 형식을 지키는지 확인합니다. (ex: validate:"email")
- notblank : 구조체 필드가 비어있지 않는지 확인합니다. (ex: validate:"notblank")
- pattern : 구조체 필드가 이메일 형식을 지키는지 확인합니다. (ex: validate:"pattern=^[A-Z]{3}$")

### 예시
```go
type Data struct {
	Name string `json:"name" validate:"pattern=^[A-Z]{3}$" binding:"required"`
	Nickname string `json:"nickname" validate:"notblank" binding:"required"`
	Email string `json:"email" validate:"email" binding:"required"`
	Age   int    `json:"age" validate:"min=1 max=100"`
}
```