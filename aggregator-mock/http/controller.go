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

	router.GET("/a/ajax-map/map/arenda/kvartiry/almaty/", getMapData)
	router.GET("/a/ajax-map-list/map/arenda/kvartiry/almaty/", getAps)

	router.POST("/clear-aps", func(c *gin.Context) {
		println("Clearing aps...")
		aps_mock.ClearAps()
		c.Status(200)
	})
	router.POST("/create-ap", func(c *gin.Context) {
		ap := aps_mock.MockApBean{}
		c.BindJSON(&ap)
		aps_mock.AddMockAp(ap)
		fmt.Printf("Added mock ap. %v\n", ap)
		c.Status(201)
	})
	router.POST("/create-n-aps", func(c *gin.Context) {
		nStr := c.Query("n")
		n, err := strconv.Atoi(nStr)
		if err != nil {
			c.AbortWithError(400, err)
			return
		}
		aps_mock.ClearAps()
		for i := 1; i <= n; i++ {
			bean := aps_mock.MockApBean{
				Id:    int64(i),
				Title: "title " + strconv.Itoa(i),
				Price: int64(100 * i),
			}
			aps_mock.AddMockAp(bean)
			fmt.Printf("Added mock ap. %v\n", bean)
		}
		c.Status(201)
	})
	return controller
}

func getAps(c *gin.Context) {
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	if page < 0 {
		c.Status(400)
		return
	}
	c.JSON(200, model.ApsResult{
		HTML:         "",
		PriceHistory: nil,
		Adverts:      aps_mock.GetApsByPage(page),
		Pager:        "",
		Page:         page,
	})
}

func getMapData(c *gin.Context) {
	c.JSON(200, model.MapData{
		IsTooManyAdverts: false,
		ListURL:          "",
		MetaData:         nil,
		NbTotal:          aps_mock.EmulatedCount,
		Results:          nil,
	})
}
