package buildingx

import "errors"

type Session struct {
	IsInitialized bool
	Partition     string
	JWT           string
}

// Initialize valides the API credentials and gets an array of all available locations
func (t *Session) Initialize(partition string) error {

	if partition == "" {
		t.Invalidate()
		return errors.New("partition cannot be empty")
	}

	token, err := GetToken()
	if err != nil {
		t.Invalidate()
		return errors.New("error while getting token: " + err.Error())
	}

	t.IsInitialized = true
	t.JWT = token
	t.Partition = partition

	return nil
}

// Invalidate resets all session properties
func (t *Session) Invalidate() {

	t.IsInitialized = false
	t.Partition = ""
	t.JWT = ""

}
