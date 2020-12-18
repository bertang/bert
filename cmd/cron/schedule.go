package cron

type Callback func(...interface{})
type JobFunc func()


type Schedule struct {
	ID            uint8
	JobName       string
	MethodNoParam JobFunc
	Method        Callback
	Params        []interface{} //多个参数 都改为
}

func (s *Schedule) Run() {
	if s.Method == nil && s.MethodNoParam == nil {
		panic("Invalid Parameters")
	}

	if s.MethodNoParam != nil {
		s.Method()
	} else if s.Method != nil {
		s.Method(s.Params...)
	}

}