package statsdproxy

import (
	"testing"
)

func TestAnswerPing(t *testing.T) {
	answer := answerManagementQuery("ping")

	for _, v := range answer {
		if v != "pong" {
			t.Errorf("wrong answer, expected 'pong' and got %s", v)
		}
	}
}
func TestAnswerPingWithSpaces(t *testing.T) {
	answer := answerManagementQuery("ping ")

	for _, v := range answer {
		if v != "pong" {
			t.Errorf("wrong answer, expected 'pong' and got %s", v)
		}
	}
}
func TestAnswerUnknownCommand(t *testing.T) {
	answer := answerManagementQuery("lerlo")

	for _, v := range answer {
		if v != "unknown command" {
			t.Errorf("wrong answer, expected 'unknown command' and got %s", v)
		}
	}
}
func TestAnswerPingWithNewlines(t *testing.T) {
	answer := answerManagementQuery("ping\r\n")

	for _, v := range answer {
		if v != "pong" {
			t.Errorf("wrong answer, expected 'pong' and got %s", v)
		}
	}
}
