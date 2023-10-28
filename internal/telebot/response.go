package telebot

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

// Messages displayed to the user.
var (
	startMessage                  = "⛽️ Привет! Выбери необходимую заправку, чтобы забронировать очередь!"
	errorMessage                  = "😬 Что-то пошло не так. Попробуйте еще раз позже."
	unknownCommandMessage         = "⛔️ Такой команды не существует!"
	chooseFuelPumpMessage         = "🚘 Выбери необходимую колонку!"
	fuelPumpIsNotAvailableMessage = "🚫 Колонка занята! Мы уведовим вас, когда она освободится!"
	choosePaymentMethodMessage    = "💸 Пожалуйста выберите метод отплаты!"
)

// getPaymentMarkup returns a markup for selecting payment methods.
func (b *Bot) getPaymentMarkup() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	rows := []tele.Row{
		markup.Row(markup.Text("💵 Наличный рассчет")),
		markup.Row(markup.Text("💳 Безналичный рассчет")),
		markup.Row(markup.Text("⛔️ Отмена!")),
	}

	markup.Reply(rows...)

	return markup
}

// getFuelPumpsMarkup returns a markup for available fuel pumps at a station.
func (b *Bot) getFuelPumpsMarkup(stationID uint64) (*tele.ReplyMarkup, error) {
	pumps, err := b.repo.FuelPumpRepository.GetAll(stationID)
	if err != nil {
		return nil, err
	}

	markup := &tele.ReplyMarkup{}
	rows := make([]tele.Row, len(pumps))

	for idx, pump := range pumps {
		availability := "✅"
		if !pump.IsAvailable {
			availability = "❌"
		}

		text := fmt.Sprintf("%d | 🛢 %s - %d₸ %s", pump.ID, pump.FuelType, pump.Price, availability)
		rows[idx] = markup.Row(markup.Text(text))
	}

	markup.Reply(rows...)

	return markup, nil
}

// getStationsMarkup returns a markup for selecting gas stations.
func (b *Bot) getStationsMarkup() (*tele.ReplyMarkup, error) {
	stations, err := b.repo.GasStationRepository.GetAll()
	if err != nil {
		return nil, err
	}

	markup := &tele.ReplyMarkup{}
	rows := make([]tele.Row, len(stations))

	for idx, station := range stations {
		text := fmt.Sprintf("%d | ⛽️ %s (%s)", station.ID, station.Company, station.Address)
		rows[idx] = markup.Row(markup.Text(text))
	}

	markup.Reply(rows...)

	return markup, nil
}
