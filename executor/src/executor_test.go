package executor

import (
	"strings"
	"testing"

	"github.com/maratona-run-time/Maratona-Runtime/model"
	"github.com/maratona-run-time/Maratona-Runtime/utils"
)

func TestExecuteOK(t *testing.T) {
	logger := utils.InitDummyLogger()
	status := Execute("./tests/ok.sh", "../../verdict/src/tests/sum/inputs", 1.0, 1, logger)
	if status[0].Status != "OK" {
		t.Errorf("Expected status OK but got %s", status[0].Status)
	}
	if strings.EqualFold(status[0].Message, "3") {
		t.Errorf("Expected output 3 but got %s", status[0].Message)
	}
}

func TestExecuteTLE(t *testing.T) {
	logger := utils.InitDummyLogger()
	status := Execute("./tests/tle.sh", "../../verdict/src/tests/sum/inputs", 1.0, 1, logger)
	if status[3].Status != model.TIME_LIMIT_EXCEEDED {
		t.Errorf("Expected status %s but got %s", model.TIME_LIMIT_EXCEEDED, status[3].Status)
	}
}

func TestExecuteRTE(t *testing.T) {
	logger := utils.InitDummyLogger()
	status := Execute("./tests/rte.sh", "../../verdict/src/tests/sum/inputs", 1.0, 1, logger)
	if status[0].Status != model.RUNTIME_ERROR {
		t.Errorf("Expected status %s but got %s", model.RUNTIME_ERROR, status[0].Status)
	}
}

func TestExecuteMLE(t *testing.T) {
	logger := utils.InitDummyLogger()
	status := Execute("./tests/mle.sh", "../../verdict/src/tests/sum/inputs", 1.0, 1, logger)
	if status[0].Status != "MLE" {
		t.Errorf("Expected status MLE but got %s", status[0].Status)
	}
}
