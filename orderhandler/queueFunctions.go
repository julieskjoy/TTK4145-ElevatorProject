package orderhandler

import (
	. "../config"
)

func numberOfOrders(elev Elev) int {
	orders := 0
	for floor := 0; floor < NumFloors; floor++ {
		ordersAtFloor := 0
		for button := 0; button < NumButtons; button++ {
			if elev.Queue[floor][button] {
				ordersAtFloor = 1
			}
		}
		orders += ordersAtFloor
	}
	return orders
}

func costofElev(order Keypress, elev Elev, nrElev int) int {
	if !elev.Online {
		return 1000
	}
	sum := order.Floor - elev.Floor
	if sum < 0 {
		sum = -sum
	}
	if elev.Dir == DirStop && sum == 0 {
		return 0
	} else if elev.Dir == DirDown {
		if order.Floor > elev.Floor {
			sum += 3
		}
	} else if elev.Dir == DirUp {
		if order.Floor < elev.Floor {
			sum += 3
		}
	} else if elev.State == Moving && sum == 0 {
		sum += 5
	}
	return sum
}
