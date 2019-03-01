package fuel

type serviceComposite interface {
	BeginRequest()
	EndRequest()
}

type Service struct {
	Fixture
}

func (s Service) BeginRequest() {

}

func (s Service) EndRequest() {

}
