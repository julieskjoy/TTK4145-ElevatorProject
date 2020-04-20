# Elevator State Machine module

This is the finite state machine of the system that has 3 states and 3 events. 

The elevator state machine takes local orders from the queue and executes them. It handles the motor direction, door timeout, door light as well as engine error timer. 

States: 
- **`IDLE`**: Elevator is stationary (standing still) at a floor with closed doors, waiting to get an order.
- **`MOVING`**: Elevator is either moving between floors or going past a floor. 
- **`DOOROPEN`**: Elevator is at a floor with the door open. 

Events: 
- **`New Order`**: There is a new order added to the queue. 
- **`Foor Reached`**: The elevator has reached a new floor. 
- **`Door Timeout`**: The door has been open for a set period of time and is timed out; the doors should close. 
