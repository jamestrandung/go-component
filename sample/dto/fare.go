package dto

type Fare struct {
	Amount float64
}

type RunningFare struct {
	fare *Fare
}

func MakeRunningFare() RunningFare {
	return RunningFare{
		fare: &Fare{
			Amount: 0,
		},
	}
}

func (f RunningFare) GetRunningFare() *Fare {
	return f.fare
}
