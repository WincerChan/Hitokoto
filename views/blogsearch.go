package views

import (
	C "DIEM-API/config"
	B "DIEM-API/models/blogs"
	T "DIEM-API/tools"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func checkSearchParams(ctx *gin.Context, p *B.Params) {
	err := ctx.Bind(p)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		ctx.Abort()
	}
}

func validDateRange(r string) string {
	if r == "" {
		return fmt.Sprintf("%d~%d", 0, time.Now().Unix())
	}
	var begin, end int64
	f := func(s *int64, t time.Time, err error) {
		if err == nil {
			*s = t.Unix()
		}
	}
	lenOfTime := 10
	if len(r) <= lenOfTime {
		return ""
	}
	if r[lenOfTime] == '~' {
		beginTime, err := time.Parse("2006-01-02", r[:lenOfTime])
		f(&begin, beginTime, err)
		r = r[lenOfTime:]
	}
	if r[0] == '~' {
		endTime, err := time.Parse("2006-01-02", r[1:])
		f(&end, endTime, err)
	}
	return fmt.Sprintf("%d~%d", begin, end)
}

func execute(ctx *gin.Context) []byte {
	p := B.Params{}
	// checkSearchParams(ctx, p)
	err := B.BindStruct(ctx.Request.URL.Query(), &p)
	if err != nil {
		ctx.JSON(200, gin.H{
			"Hello world": err.Error(),
		})
		ctx.Abort()
		return []byte{}
	}
	// v, err := json.Marshal(p)
	// log.Println(string(v))
	T.CheckException(err, "decode json error")
	c := C.SearchPool.Get()
	c.WriteLine([]byte(p.Serialize()))
	// log.Println("write")
	defer C.SearchPool.Put(c)
	return c.ReadLine()
}

func BlogSearchViews(ctx *gin.Context) {
	ret := execute(ctx)
	// log.Println(ret)
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Content-Type", "application/json")
	ctx.Writer.Write(ret)
}
