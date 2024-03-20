package http

import (
	"aggregator_mock/aps-mock"
	"fmt"
	"github.com/Logotipiwe/krisha_model/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Controller struct {
	Router *gin.Engine
}

func InitTestController() *Controller {
	router := gin.Default()
	controller := &Controller{Router: router}

	router.GET("/a/ajax-map/*subPath", getMapData)
	router.GET("/a/ajax-map-list/*subPath", getAps)

	router.POST("/clear-aps", func(c *gin.Context) {
		println("Clearing aps...")
		aps_mock.ClearAps()
		c.Status(200)
	})
	router.POST("/create-ap", func(c *gin.Context) {
		ap := aps_mock.MockApBean{}
		c.BindJSON(&ap)
		subPath := c.Query("subPath")
		aps_mock.AddMockAp(subPath, ap)
		fmt.Printf("Added mock ap. %v\n", ap)
		c.Status(201)
	})
	router.POST("/create-n-aps", func(c *gin.Context) {
		nStr := c.Query("n")
		subPath := c.Query("subPath")
		if subPath == "" {
			c.String(400, "subPath empty")
			return
		}
		n, err := strconv.Atoi(nStr)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		aps_mock.ClearAps()
		for i := 1; i <= n; i++ {
			bean := aps_mock.MockApBean{
				Id:     int64(i),
				Title:  "title " + strconv.Itoa(i),
				Price:  int64(100 * i),
				Images: []string{strconv.Itoa(i), strconv.Itoa(i + 1)},
			}
			aps_mock.AddMockAp(subPath, bean)
			fmt.Printf("Added mock ap. %v\n", bean)
		}
		c.Status(201)
	})
	return controller
}

func getAps(c *gin.Context) {
	fmt.Println(c)
	subPath := c.Param("subPath")
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	if page < 0 {
		c.String(400, "page less than zero")
		return
	}
	if subPath == "" {
		c.String(400, "subPath empty")
		return
	}

	aps, err := aps_mock.GetApsByPageAndPath(subPath, page)
	if err != nil {
		c.String(500, err.Error())
	}
	c.JSON(200, model.ApsResult{
		HTML:         "",
		PriceHistory: nil,
		Adverts:      *aps,
		Pager:        "",
		Page:         page,
	})
}

func getMapData(c *gin.Context) {
	count, err := aps_mock.GetApsCount(c.Param("subPath"))
	if err != nil {
		c.String(400, err.Error())
		return
	}
	c.JSON(200, model.MapData{
		IsTooManyAdverts: false,
		ListURL:          "",
		MetaData:         nil,
		NbTotal:          count,
		Results:          nil,
	})
}
