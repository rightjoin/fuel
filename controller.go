package fuel

type service interface {
	BeginRequest()
	EndRequest()
}

type Controller struct {
	Fixture
}

func (c Controller) BeginRequest() {

}

func (c Controller) EndRequest() {

}
