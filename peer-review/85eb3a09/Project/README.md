# Elevator Project - TTK4145
---

### Summary
This code is a possible solution for the issue given in the course TTK4145: Create software for controlling n elevators working in parallel across m floors. The elevators are only tested for the scenario with four floors and up to three number of elevators. This elevator system uses a pure peer-to-peer topology and is based on a process pair, consisting of a primary and a backup.

### Contents



- [Elevator Project - TTK4145](#elevator-project---ttk4145)
    - [Summary](#summary)
    - [Contents](#contents)
    - [Run program](#run-program)
    - [Features](#features)
    - [Modules](#modules)


### Run program
---
The program connects to a server using a terminal flag for the IP adress. First step is to start an elevator environment. It is either a physical elevator or a simulator from your computer. The documentation and necessary software can be found on git for this course.
Make sure the terminal runs the program from the correct directory.

1. Start elevator
-  Simulator (linux): 
```bash
$ simelevatorserv --port yourPort
```
-  Physical (linux): 
```bash
$ elevatorserver --port yourPort
```
2. Run program
In another terminal window directed to the project directory, run main.
```bash
$ cd ./Project filepath
$ go run main.go --port elevatorPort --id elevatorID --portBackup backupPort --backup 0
```

### Features
---
**Terms**

| Term                             | Explanation                                                                                                                                                                                                                 |
| -------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Available elevators**          | Elevators that are not obstructed or have motor stop, and are able to share and receive data from the network.                                                                                                              |
| **Unavailable elevators**        | Elevators that are obstructed, have motor stop, and/or are disconnected from the network.                                                                                                                                   |
| **Peer network**                 | The network in which the elevators, peers, communicate if they are "alive". In the case of a peer disconnection, the other peers will register the disconnected elevator as a "lost" peer, and reconnected as a "new" peer. |
| **Broadcasting network (bcast)** | The network where data, the world views, from each available elevator system is being broadcasted.                                                                                                                          |
| **Confirmed orders**             | Order state for when orders have been acknowledged by all available elevators and are ready to be assigned. Can be recognized with a light turned on.                                                                       |
| **Completed orders**             | Order state where orders have been called, confirmed, and served. Can be recognized by the order light being turned off.                                                                                                    |


**Peer to peer topology**

The solution uses a pure peer-to-peer topology where all elevators share data across themselves in order to compute and assign hall orders. In order to ensure correct updating of orders and request states the program uses a counter that is broadcasted between the connected elevators. An elevator can then assume that a message with a higher counter than the itself is newer.

**Crashes and disconnection**

In the case of a programcrash and peer network disconnection in a real elevator the boarded passangers should still be able to exit at their chosen floor. In other words, if the program crashes and is disconnected from the network the elevator must finish confirmed cab orders and take new cab orders. The orders assigned to the lost peer should be reassigned to *available* elevators. 

An elevator consists of a process pair, with a primary and backup. The primary runs the program with concurrent threads, communicates with other peers on the peer network, and sends "I´m alive" messages via UDP broadcasting to the backup. If the backup does not recieve a message from primary within a given time interval the backup assumes the primary is dead and will take over as primary and spawn a new backup. The primary should then finish existing cab calls and return to a normal behaviour as soon as the elevator reaches an available state. 

### Modules
---

**assigner**

The assigner module computes an order list of hall orders for each available elevator using the peers' request order list and shared hall order status. The hall request assigner, HRA, uses an algorithm to decide which elevator that should serve a request.


**communication**

The comminucation module mainly handles the messages being sent between elevators. Examples of this is counters, elevator structs, local hallorder status, acknowledgement lists. The data recieved from the broadcasting network will be processed in communicationhandler() and forwarded to the singleElevator FSM.

**config**

The config module holds the configurations for an elevator system. Constant values that define different system parameters are listed up within this module. Configurations like the number of floors, timer durations, and RequestState values are examples of the systemparameters. 

**network**

Communication between peers is broadcaasted over a network using UDP-protocol. The transmitter and reciever functions that binds elevators the their correct ports are set up in this module. This code was mainly given in the course.

**saveJSON**

As a backup all recent cab orders ar stored in a `.json` file, so in the case of programcrash or program disconnect the elevator has access to confirmed cab calls and can tend to the request taken before the crash. The saveJSON module consist of relevant functions to make the updating, saving and reading from file possible.

**singleElevator**

The singleElevator module is a complete single elevator system. Using the given configuration it has the ability to take orders and tend to them both physically with `elevio package`  and using an elevator simulator `simelevatorserv`.