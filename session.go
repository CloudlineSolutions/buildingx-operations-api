package buildingx

import "errors"

type Session struct {
	Partition string
	JWT       string
	Locations []Location
}

func (t *Session) Initialize(partition string) error {

	if partition == "" {
		t.Destroy()
		return errors.New("partition cannot be empty")
	}

	token, err := GetToken()
	if err != nil {
		t.Destroy()
		return errors.New("error while getting token: " + err.Error())
	}
	t.JWT = token
	t.Partition = partition

	// get locations associated with this partition
	locations, err := GetLocations(*t)
	if err != nil {
		t.Destroy()
		return errors.New("error getting locations for this partition: " + err.Error())
	}

	t.Locations = locations

	return nil
}
func (t *Session) Destroy() {

	t.Partition = ""
	t.JWT = ""
	t.Locations = nil

}
