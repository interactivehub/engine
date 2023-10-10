package wheel

import (
	"log"
	"time"
)

type WheelRoundAutoer interface {
	Start() WheelRoundAutoer
	OnRoundStart(func(r *WheelRound) error) WheelRoundAutoer
	OnRoundEnd(func(r *WheelRound) error) WheelRoundAutoer
	OnStatusChange(func(r *WheelRound) error) WheelRoundAutoer
}

type WheelRoundAuto struct {
	*WheelRound
	alreadyStarted  bool
	openToSpinTimer *time.Timer
	spinToEndTimer  *time.Timer
	onStatusChange  func(r *WheelRound) error
	onRoundStart    func(r *WheelRound) error
	onRoundEnd      func(r *WheelRound) error
}

func (r *WheelRoundAuto) Start() WheelRoundAutoer {
	if r.alreadyStarted {
		return r
	}

	r.alreadyStarted = true

	go func() {
		for {
			if r.onRoundStart != nil {
				r.onRoundStart(r.WheelRound)
			}

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

				r.openToSpinTimer.Stop()
				r.spinToEndTimer.Stop()

				if r.onStatusChange != nil {
					r.onStatusChange(r.WheelRound)
				}

				if r.onRoundEnd != nil {
					r.onRoundEnd(r.WheelRound)
				}

				return
			}
		}
	}()

	return r
}

func (r *WheelRoundAuto) OnStatusChange(onStatusChange func(r *WheelRound) error) WheelRoundAutoer {
	r.onStatusChange = onStatusChange

	return r
}

func (r *WheelRoundAuto) OnRoundStart(onRoundStart func(r *WheelRound) error) WheelRoundAutoer {
	r.onRoundStart = onRoundStart

	return r
}

func (r *WheelRoundAuto) OnRoundEnd(onRoundEnd func(r *WheelRound) error) WheelRoundAutoer {
	r.onRoundEnd = onRoundEnd

	return r
}

func (r *WheelRound) Auto() WheelRoundAutoer {
	openToSpinTimer := time.NewTimer(r.OpenDuration)
	spinToEndTimer := time.NewTimer(r.OpenDuration + r.SpinDuration)

	auto := &WheelRoundAuto{
		WheelRound:      r,
		openToSpinTimer: openToSpinTimer,
		spinToEndTimer:  spinToEndTimer,
	}

	return auto
}
