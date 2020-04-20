package orderhandler

import (
	"fmt"

	. "../config"
	hw "../hardware"
)

func getOrder(floor int, buttonType Button, elevator Elev) bool {
	return elevator.Queue[floor][buttonType]
}

func setOrder(floor int, buttonType Button, elevator Elev, set bool) {
	elevator.Queue[floor][buttonType] = set
}

// en goroutine som mottar knappetrykk og states og sørger for at riktig heis går til riktig plass.
// trenger vel ikke returnere noe??

func OrderHandler(order chan Keypress, nrElev int, completedOrderCh chan int, updateLightsCh chan [NumElevators]Elev,
	newOrderCh chan Keypress, elevatorCh chan Elev, updateQueueCh chan [NumElevators]Elev, updateSyncCh chan Elev,
	orderUpdateCh chan Keypress) {

	var (
		elevators      [NumElevators]Elev
		completedOrder Keypress
	)
	completedOrder.DesignatedElevator = nrElev
	elevators[nrElev] = <-elevatorCh
	updateSyncCh <- elevators[nrElev]
	for {
		select {
		case orderLocal := <-order:
			fmt.Println("Orderlocal")
			if !elevators[nrElev].Online && orderLocal.Btn == BtnInside {
				setOrder(orderLocal.Floor, BtnInside, elevators[nrElev], true)
				updateLightsCh <- elevators
				go func() { newOrderCh <- orderLocal }()
			} else if !elevators[nrElev].Online && orderLocal.Btn != BtnInside {
				continue
			} else {
				if orderLocal.Floor == elevators[nrElev].Floor && elevators[nrElev].State != Moving {
					newOrderCh <- orderLocal
				} else {
					if !existingOrder(orderLocal, elevators, nrElev) {
						fmt.Println("new order", orderLocal.Floor+1)
						var sums [NumElevators]int
						for elevator := 0; elevator < NumElevators; elevator++ {
							sums[elevator] = costofElev(orderLocal, elevators[elevator], nrElev)
							if elevator != 0 {
								if sums[elevator] < sums[orderLocal.DesignatedElevator] {
									orderLocal.DesignatedElevator = elevator
								}
							} else {
								orderLocal.DesignatedElevator = 1
							}
						}
						fmt.Println("Cost of elevators: ", sums)
						orderUpdateCh <- orderLocal
					}
				}
			}

		case completedFloor := <-completedOrderCh:
			completedOrder.Finished = true
			completedOrder.Floor = completedFloor
			var button Button
			for btn := BtnUp; btn < NumButtons; btn++ {
				if elevators[nrElev].Queue[completedFloor][btn] {
					button = btn
					completedOrder.Btn = button
				}
				for elevator := 0; elevator < NumElevators; elevator++ {
					if btn != BtnInside || elevator == nrElev {
						elevators[elevator].Queue[completedFloor][btn] = false
					}
				}
			}
			if elevators[nrElev].Online {
				orderUpdateCh <- completedOrder
			}
			updateLightsCh <- elevators

		case newElevator := <-elevatorCh:
			newQueue := elevators[nrElev].Queue
			if elevators[nrElev].State == Undefined && newElevator.State != Undefined {
				elevators[nrElev].Online = true
			}
			elevators[nrElev] = newElevator
			elevators[nrElev].Queue = newQueue
			if elevators[nrElev].Online {
				updateSyncCh <- elevators[nrElev]
			}
		case tempElevList := <-updateQueueCh:
			newOrder := false
			for elevator := 0; elevator < NumElevators; elevator++ {
				if nrElev == elevator {
					continue
				}
				if elevators[elevator].Queue != tempElevList[elevator].Queue {
					newOrder = true
				}
				elevators[elevator] = tempElevList[elevator]
			}
			for floor := 0; floor < NumFloors; floor++ {
				for button := BtnUp; button < NumButtons; button++ {
					if !elevators[nrElev].Queue[floor][button] && tempElevList[nrElev].Queue[floor][button] {
						elevators[nrElev].Queue[floor][button] = true
						order1 := Keypress{Floor: floor, Btn: button, DesignatedElevator: nrElev, Finished: false}
						go func() { newOrderCh <- order1 }()
						newOrder = true
					} else if elevators[nrElev].Queue[floor][button] && !tempElevList[nrElev].Queue[floor][button] {
						elevators[nrElev].Queue[floor][button] = false
						order1 := Keypress{Floor: floor, Btn: button, DesignatedElevator: nrElev, Finished: true}
						go func() { newOrderCh <- order1 }()
						newOrder = true
					}
				}
			}
			if newOrder {
				updateLightsCh <- elevators
			}
		}
	}

}

func SetLights(updateLightsCh <-chan [NumElevators]Elev, nrElev int) {
	//var orders [NumElevators]bool

	for {
		elevators := <-updateLightsCh
		for floor := 0; floor < NumFloors; floor++ {
			for button := BtnUp; button < NumButtons; button++ {
				for elevator := 0; elevator < NumElevators; elevator++ {
					//orders[elevator] = false
					if elevator != nrElev && (button == BtnInside || button == BtnDown || button == BtnUp) {
						continue
					}
					if elevators[elevator].Queue[floor][button] {
						hw.SetButtonLamp(button, floor, 1)
						//orders[elevator] = true
					} else {
						hw.SetButtonLamp(button, floor, 0)
					}
				}
			}
		}
	}
}

func existingOrder(order Keypress, elevators [NumElevators]Elev, nrElev int) bool {
	if elevators[nrElev].Queue[order.Floor][BtnInside] && order.Btn == BtnInside {
		return true
	}
	for elev := 0; elev < NumElevators; elev++ {
		if elevators[nrElev].Queue[order.Floor][order.Btn] {
			return true
		}
	}
	return false
}
