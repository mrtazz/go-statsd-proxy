package statsdproxy

import (
	"testing"
)

func TestAnswerPing(t *testing.T) {
	answer := answerManagementQuery("ping")

	if answer != "pong" {
		t.Errorf("wrong answer, expected 'pong' and got %s", answer)
	}
}
func TestAnswerPingWithSpaces(t *testing.T) {
	answer := answerManagementQuery("ping ")

	if answer != "pong" {
		t.Errorf("wrong answer, expected 'pong' and got %s", answer)
	}
}
func TestAnswerUnknownCommand(t *testing.T) {
	answer := answerManagementQuery("lerlo")

	if answer != "unknown command" {
		t.Errorf("wrong answer, expected 'unknown command' and got %s", answer)
	}
}
func TestAnswerPingWithNewlines(t *testing.T) {
	answer := answerManagementQuery("ping\r\n")

	if answer != "pong" {
		t.Errorf("wrong answer, expected 'pong' and got %s", answer)
	}
}
