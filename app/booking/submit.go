package booking

import (
	"fmt"
	"strings"
	"time"

	"github.com/Bekreth/jane_cli/app/terminal"
)

const bookingTimeYearFormat = "06.01.02T15:04"
const bookingTimeFormat = "01.02T15:04"
const bookingDateFlag = "-d"
const treatmentFlag = "-t"
const patientFlag = "-p"

func (state *bookingState) Submit() {
	if state.currentBuffer == ".." {
		state.nextState = state.rootState
		return
	}
	flags := terminal.ParseFlags(state.currentBuffer)
	state.logger.Debugf("submitting query flags: %v", flags)
	missingFlags := map[string]string{
		"-d": "",
		"-a": "",
		"-p": "",
	}

	for key := range missingFlags {
		delete(missingFlags, key)
	}
	if len(missingFlags) != 0 {
		joined := strings.Join(terminal.MapKeysString(missingFlags), ", ")
		notifcation := fmt.Sprintf("missing arguments %v", joined)
		state.writer.WriteString(notifcation)
		state.writer.NewLine()
		return
	}
	state.currentBuffer = ""
	builder := bookingBuilder{
		substate: unknown,
	}

	builder, err := state.parsePatientValue(flags[patientFlag], builder)
	if err != nil {
		state.writer.WriteString(err.Error())
		state.writer.NewLine()
		return
	}

	builder, err = state.parseTreatmentValue(flags[treatmentFlag], builder)
	if err != nil {
		state.writer.WriteString(err.Error())
		state.writer.NewLine()
		return
	}

	builder, err = state.parseDateValue(flags[bookingDateFlag], builder)
	if err != nil {
		state.writer.WriteString(err.Error())
		state.writer.NewLine()
		return
	}

	state.booking = builder
}

func (state *bookingState) parsePatientValue(
	patientName string,
	builder bookingBuilder,
) (bookingBuilder, error) {
	if patientName == "" {
		return builder, fmt.Errorf("no name provided, use the %v flag", patientFlag)
	}
	patients, err := state.fetcher.FindPatients(patientName)
	if err != nil {
		state.nextState = state.rootState
		return builder, fmt.Errorf("failed to lookup patient %v : %v", patientName, err)
	}
	builder.patients = patients
	if len(patients) == 0 {
		return builder, fmt.Errorf("no patients found for %v", patientName)
	} else if len(patients) == 1 {
		builder.targetPatient = patients[0]
	} else if len(patients) > 8 {
		return builder, fmt.Errorf("too many patients to render nicely for %v", patientName)
	}
	return builder, nil
}

func (state *bookingState) parseTreatmentValue(
	treatmentName string,
	builder bookingBuilder,
) (bookingBuilder, error) {
	if treatmentName == "" {
		return builder, fmt.Errorf("no name provided, use the %v flag", patientFlag)
	}
	treatments, err := state.fetcher.FindTreatment(treatmentName)
	if err != nil {
		state.nextState = state.rootState
		return builder, fmt.Errorf("failed to lookup treatments %v : %v", treatmentName, err)
	}
	builder.treatments = treatments
	if len(treatments) == 0 {
		return builder, fmt.Errorf("no treatment found for %v", treatmentName)
	} else if len(treatments) == 1 {
		builder.targetTreatment = treatments[0]
	} else if len(treatments) > 8 {
		return builder, fmt.Errorf(
			"too many treatments to render nicely for %v",
			treatmentName,
		)
	}
	return builder, nil
}

func (state *bookingState) parseDateValue(
	dateString string,
	builder bookingBuilder,
) (bookingBuilder, error) {

	dateValue, err := time.Parse(bookingTimeYearFormat, dateString)
	if err == nil {
		builder.appointmentDate = dateValue
		return builder, nil
	} else {
		dateValue, err = time.Parse(bookingTimeFormat, dateString)
		if err != nil {
			return builder, fmt.Errorf(
				"unable to parse date %v, please write date in the format %v or %v",
				dateString,
				bookingTimeFormat,
				bookingTimeYearFormat,
			)
		}

		now := time.Now()
		if now.Month() > dateValue.Month() {
			dateValue = dateValue.AddDate(now.Year()+1, 0, 0)
		} else {
			dateValue = dateValue.AddDate(now.Year(), 0, 0)
		}

		builder.appointmentDate = dateValue
		return builder, nil
	}
}
