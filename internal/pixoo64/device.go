package pixoo64

const (
	DeviceWidth  = 64
	DeviceHeight = 64
)

func PlayBuzzer(client *Client, activeTimeInCycle int, offTimeInCycle int, playTotalTime int) error {
	data := map[string]any{
		"Command":           "Device/PlayBuzzer",
		"ActiveTimeInCycle": activeTimeInCycle,
		"OffTimeInCycle":    offTimeInCycle,
		"PlayTotalTime":     playTotalTime,
	}

	_, err := client.Post(data)
	if err != nil {
		return err
	}
	return nil
}
