package telebot

import (
	"fmt"
	"strconv"

	"github.com/rtsoy/software-design-patterns-project/internal/model"
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
	somethingWentWrongMessage     = "😬 Что-то пошло не так! Попробуйте ещё раз!"
	cardAttachmentSuccessMessage  = "😉 Ваша карта была успешно привязана!"
	insertFuelAmountMessage       = "⤵️ Сколько литров топлива вы бы хотели залить?"
	noAttachedCreditCards         = "\U0001FAE4 У вас нет привязанных кредитных карт." //nolint
	chooseCreditCard              = "👀 Выберите кредитную карту для оплаты."           //nolint
	startRefueling                = "❓ Начать заправку?"
	stopRefueling                 = "❗️ Нажмите кнопку, по окончании заправки."

	enterCardNumberMessage = "🔢 Введите номер карты, расположен на ее лицевой стороне и обычно состоит из 16 цифр.\n" +
		"Например, 1234 5678 9012 3456."
	enterCardExpDateMessage = "⌛️ Введите дату истечения карты, дата под номером банковской карты.\n" +
		"Например, 12/24."
	enterCardHolderMessage = "🤵 Введите имя и фамилию держателя карты, расположено на лицевой стороне карты.\n" +
		"Например, Ivan Ivanov."
	enterCardCVVMessage = "🤫 Введите трёхзначный код расположенный на задней части карты.\n" +
		"Например, 123."
)

func (b *Bot) getCashPaymentMessage(total uint64) string {
	format := "🚶‍♂️ Пройдите, пожалуйста, на кассу и осуществите оплату (%d₸), указав номер вашей заправочной колонки."

	return fmt.Sprintf(format, total)
}

func (b *Bot) getStopMarkup() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	row := markup.Row(markup.Text("️🛑 Закончить заправку!"))

	markup.Reply(row)

	return markup
}

func (b *Bot) getStartMarkup() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	row := markup.Row(markup.Text("▶️ Начать!"))

	markup.Reply(row)

	return markup
}

func (b *Bot) getCreditCardsMarkup(cards []model.CreditCard) *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}
	rows := make([]tele.Row, len(cards)+1)

	for idx, card := range cards {
		number := strconv.FormatUint(card.Number, 10)

		text := fmt.Sprintf("%d. %s .... .... %s (%d/%d)",
			idx+1, number[:4], number[12:], card.ExpirationDate.Month(), card.ExpirationDate.Year())

		rows[idx] = markup.Row(markup.Text(text))
	}

	rows[len(cards)] = markup.Row(markup.Text("⛔️ Отмена!"))

	markup.Reply(rows...)

	return markup
}

func (b *Bot) getDoneOrCancelMarkup() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	rows := []tele.Row{
		markup.Row(markup.Text("☑️ Готово!")),
		markup.Row(markup.Text("⛔️ Отмена!")),
	}

	markup.Reply(rows...)

	return markup
}

func (b *Bot) getCancelMarkup() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	row := markup.Row(markup.Text("⛔️ Отмена!"))

	markup.Reply(row)

	return markup
}

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
	rows := make([]tele.Row, len(pumps)+1)

	for idx, pump := range pumps {
		availability := "✅"
		if !pump.IsAvailable {
			availability = "❌"
		}

		text := fmt.Sprintf("%d | 🛢 %s - %d₸ %s", pump.ID, pump.FuelType, pump.Price, availability)
		rows[idx] = markup.Row(markup.Text(text))
	}

	rows[len(pumps)] = markup.Row(markup.Text("⛔️ Отмена!"))

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
	rows := make([]tele.Row, len(stations)+1)

	for idx, station := range stations {
		text := fmt.Sprintf("%d | ⛽️ %s (%s)", station.ID, station.Company, station.Address)
		rows[idx] = markup.Row(markup.Text(text))
	}

	rows[len(stations)] = markup.Row(markup.Text("💼 Привязать карту."))

	markup.Reply(rows...)

	return markup, nil
}
