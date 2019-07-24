package tests

import (
	"errors"
	"testing"

	"github.com/rightjoin/fuel"
	baloo "gopkg.in/h2non/baloo.v3"
)

type FaultyService struct {
	fuel.Service
	getIt fuel.GET `route:"get-it/{input}"`
}

func (s *FaultyService) GetIt(input string) error {
	if input == "0" {
		return nil
	} else if input == "1" {
		return errors.New("time is now")
	} else {
		return fuel.Fault{
			HTTPCode: 414,
			ErrorNum: 100,
			Message:  "There is a fault",
		}
	}
}

func TestFaulty(t *testing.T) {
	server := fuel.NewServer()
	server.AddService(&FaultyService{})
	url, _ := server.RunTestInstance()

	var web = baloo.New(url)

	web.Get("/faulty/get-it/0").
		Expect(t).
		Status(200).
		Done()

	web.Get("/faulty/get-it/1").
		Expect(t).
		Status(417).
		JSON(map[string]interface{}{"error_num": 9999, "message": "An error occurred", "inner": "time is now"}).
		Done()

	web.Get("/faulty/get-it/2").
		Expect(t).
		Status(414).
		JSON(map[string]interface{}{"error_num": 100, "message": "There is a fault"}).
		Done()

}
