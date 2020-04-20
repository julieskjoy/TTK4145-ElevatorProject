# Order handler module

This is a module that keeps track of the orders coming and assigning them to an elevator based on each elevators 'cost'. 

The cost of an elevator is calculated based on:
- The differnece between the floor that the order has come at and the floor that the elevator is on. 
- The number of orders that the elevator already has.
- The direction that the elevator is heading compared to the order that has come in to the system.
- If the elevator is moving but doesn't have any other orders. 
