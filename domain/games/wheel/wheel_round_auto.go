package wheel

import (
	"log"
	"time"
)

type onRoundEnder interface {
	OnRoundEnd(func(r *WheelRound) error) starterOnStatusChanger
}

type onStatusChanger interface {
	OnStatusChange(func(r *WheelRound) error) starterOnRoundEnder
}

type starter interface {
	Start() onRoundEnderOnStatusChanger
}

type starterOnRoundEnder interface {
	starter
	onRoundEnder
}

type starterOnStatusChanger interface {
	starter
	onStatusChanger
}

type onRoundEnderOnStatusChanger interface {
	onRoundEnder
	onStatusChanger
}

type AutoWheelRounder interface {
	starter
	onStatusChanger
	onRoundEnder
}

type WheelRoundAuto struct {
	*WheelRound
	openToSpinTimer *time.Timer
	spinToEndTimer  *time.Timer
	onStatusChange  func(r *WheelRound) error
	onRoundEnd      func(r *WheelRound) error
}

func (r *WheelRoundAuto) Start() onRoundEnderOnStatusChanger {
	go func() {
		for {
			select {
			case <-r.openToSpinTimer.C:
				_, err := r.Roll()
				if err != nil {
					log.Println(err)
					continue
				}

				if r.onStatusChange != nil {
					r.onStatusChange(r.WheelRound)
				}

			case <-r.spinToEndTimer.C:
				err := r.EndRound()
				if err != nil {
					log.Println(err)
					continue
				}

				if r.onStatusChange != nil {
					r.onStatusChange(r.WheelRound)
				}

				if r.onRoundEnd != nil {
					r.onRoundEnd(r.WheelRound)
				}
			}
		}
	}()

	return r
}

func (r *WheelRoundAuto) OnStatusChange(onStatusChange func(r *WheelRound) error) starterOnRoundEnder {
	r.onStatusChange = onStatusChange

	return r
}

func (r *WheelRoundAuto) OnRoundEnd(onRoundEnd func(r *WheelRound) error) starterOnStatusChanger {
	r.onRoundEnd = onRoundEnd

	return r
}

func (r *WheelRound) Auto() AutoWheelRounder {
	openToSpinTimer := time.NewTimer(r.OpenDuration)
	spinToEndTimer := time.NewTimer(r.OpenDuration + r.SpinDuration)

	auto := &WheelRoundAuto{
		WheelRound:      r,
		openToSpinTimer: openToSpinTimer,
		spinToEndTimer:  spinToEndTimer,
	}

	return auto
}
