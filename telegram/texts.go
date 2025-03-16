package telegram

func (bot *Bot) InitCityNames() (err error) {
	cityNames, err := bot.mc.CityNamesGet()
	if err != nil {
		return err
	}
	for k, v := range cityNames {
		bot.cityNames.Set(k, v)
	}
	return nil
}

func (bot *Bot) CityName(key string) string {
	value, exists := bot.cityNames.Get(key)
	if exists {
		return value
	}
	return key
}
