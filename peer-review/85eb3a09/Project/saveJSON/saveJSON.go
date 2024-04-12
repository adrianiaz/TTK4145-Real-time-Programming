package saveJSON

import (
	"Project/config"
	"Project/singleElevator/elevio"
	"encoding/json"
	"os"
	"reflect"
)

func ReadCabButtonFromFile(file_name string) ([]bool, error) {
	var cabButtonCall []bool
	values, err := os.ReadFile(file_name)
	if err != nil {
		return cabButtonCall, err
	}

	err = json.Unmarshal(values, &cabButtonCall)
	if err != nil {
		return cabButtonCall, err
	}
	return cabButtonCall, nil
}

func RestoreCabOrders(cabButtonCh chan<- elevio.ButtonEvent, file_name string) error {
	cabButtonCall, err := ReadCabButtonFromFile(file_name)
	if err != nil {
		return err
	}

	for floor := 0; floor < config.NumFloors; floor++ {
		if cabButtonCall[floor] {
			cabButtonCh <- elevio.ButtonEvent{Floor: floor, Button: elevio.ButtonType(elevio.BT_Cab)}
		}
	}
	return nil
}

func SaveCabButtonToFile(NewCabButton []bool, file_name string) error {
	if _, err := os.Stat(file_name); os.IsNotExist(err) {
		data, err := json.Marshal(NewCabButton)
		if err != nil {
			return err
		}

		err = os.WriteFile(file_name, data, 0644)
		if err != nil {
			return err
		}

		return nil
	}

	cabFromFile, err := ReadCabButtonFromFile(file_name)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(NewCabButton, cabFromFile) {
		data, err := json.Marshal(NewCabButton)
		if err != nil {
			return err
		}

		err = os.WriteFile(file_name, data, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
