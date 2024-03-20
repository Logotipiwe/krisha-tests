package aps_mock

import (
	"errors"
	"github.com/Logotipiwe/krisha_model/model"
	"strconv"
)

var (
	EmulatedApsByPath = make(map[string][]*model.Ap)
	Ids               = make(map[int64]bool)
)

const pageSize = 20

func GetApsByPageAndPath(subPath string, pageNum int) (*map[string]*model.Ap, error) {
	apsOfPath, has := EmulatedApsByPath[subPath]
	if !has {
		return nil, getNonExistPathErr(subPath)
	}

	if len(apsOfPath) <= ((pageNum - 1) * pageSize) {
		m := map[string]*model.Ap{}
		return &m, nil
	}
	firstIndex := (pageNum - 1) * pageSize
	lastIndex := firstIndex + pageSize
	if lastIndex > (len(apsOfPath) - 1) {
		lastIndex = len(apsOfPath)
	}
	page := apsOfPath[firstIndex:lastIndex]
	pageMap := make(map[string]*model.Ap)
	for _, ap := range page {
		pageMap[strconv.FormatInt(ap.ID, 10)] = ap
	}
	return &pageMap, nil
}

func GetApsCount(subPath string) (int, error) {
	aps, has := EmulatedApsByPath[subPath]
	if !has {
		return 0, getNonExistPathErr(subPath)
	}
	return len(aps), nil
}

func getNonExistPathErr(subPath string) error {
	return errors.New("tried to get aps count from non-existing path: " + subPath)
}

func AddMockAp(subPath string, bean MockApBean) {
	if _, has := Ids[bean.Id]; has {
		return
	}
	_, has := EmulatedApsByPath[subPath]
	if !has {
		EmulatedApsByPath[subPath] = make([]*model.Ap, 0)
	}
	apsOfPathNew := append(EmulatedApsByPath[subPath], createMockAp(bean))
	EmulatedApsByPath[subPath] = apsOfPathNew
	Ids[bean.Id] = true
}

func ClearAps() {
	EmulatedApsByPath = make(map[string][]*model.Ap)
	Ids = make(map[int64]bool)
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
