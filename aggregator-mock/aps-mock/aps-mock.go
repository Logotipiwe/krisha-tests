package aps_mock

import (
	"github.com/Logotipiwe/krisha_model/model"
	"strconv"
)

var (
	EmulatedAps   = make([]*model.Ap, 0)
	Ids           = make(map[int64]bool)
	EmulatedCount int
)

const pageSize = 20

func GetApsByPage(pageNum int) map[string]*model.Ap {
	allAps := EmulatedAps
	if len(allAps) <= ((pageNum - 1) * pageSize) {
		return map[string]*model.Ap{}
	}
	firstIndex := (pageNum - 1) * pageSize
	lastIndex := firstIndex + pageSize
	if lastIndex > (len(EmulatedAps) - 1) {
		lastIndex = len(EmulatedAps)
	}
	page := allAps[firstIndex:lastIndex]
	pageMap := make(map[string]*model.Ap)
	for _, ap := range page {
		pageMap[strconv.FormatInt(ap.ID, 10)] = ap
	}
	return pageMap
}

func AddMockAp(bean MockApBean) {
	if _, has := Ids[bean.Id]; has {
		return
	}
	EmulatedAps = append(EmulatedAps, createMockAp(bean))
	Ids[bean.Id] = true
	EmulatedCount++
}

func ClearAps() {
	EmulatedAps = make([]*model.Ap, 0)
	Ids = make(map[int64]bool)
	EmulatedCount = 0
}

func createMockAp(bean MockApBean) *model.Ap {
	var photos []*model.Photo
	for _, image := range bean.Images {
		photos = append(photos, &model.Photo{
			Src: image,
		})
	}
	return &model.Ap{
		ID:                      bean.Id,
		Storage:                 "",
		CommentsType:            "",
		IsCommentable:           false,
		IsCommentableByEveryone: false,
		IsOnMap:                 false,
		HasPrice:                false,
		Price:                   bean.Price,
		Photos:                  photos,
		HasPackages:             false,
		Title:                   bean.Title,
		Addresstitle:            "",
		UserType:                "",
		Square:                  0,
		Rooms:                   0,
		OwnerName:               "",
		Status:                  "",
		Map:                     nil,
	}
}
