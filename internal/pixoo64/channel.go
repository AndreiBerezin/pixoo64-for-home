package pixoo64

func SetBrightness(client Client, brightness int) error {
	data := map[string]any{
		"Command":    "Channel/SetBrightness",
		"Brightness": brightness,
	}

	_, err := client.Post(data)
	if err != nil {
		return err
	}
	return nil
}

func OnOffScreen(client Client, on bool) error {
	onOff := 0
	if on {
		onOff = 1
	}
	data := map[string]any{
		"Command": "Channel/OnOffScreen",
		"OnOff":   onOff,
	}

	_, err := client.Post(data)
	if err != nil {
		return err
	}
	return nil
}
