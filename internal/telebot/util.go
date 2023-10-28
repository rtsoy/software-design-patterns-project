package telebot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// extractIDFromMessage extracts and returns a uint64 ID from the given message string.
func (b *Bot) extractIDFromMessage(msg string) (uint64, error) {
	parts := strings.Split(msg, " | ")

	if len(parts) < 1 {
		return 0, errors.New("invalid message")
	}

	id, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse uint from the message: %w", err)
	}

	return id, nil
}
