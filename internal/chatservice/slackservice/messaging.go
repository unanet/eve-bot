package slackservice

import "fmt"

const (
	msgErrNotification           = "Something terrible has happened..."
	msgErrNotificationAssurance  = "We've received the alert and someone is looking into the error..."
	msgNotification              = "I've got some news..."
	msgDeploymentErrNotification = "I detected some deployment *errors:*"
	msgLogLinks                  = "Here are the latest logs..."
	msgResultsNotification       = "Here are your results..."
	msgReleaseNotification       = "Successfully released...."
	msgAuthLink                  = "Here is your account auth link:"
)

const (
	devOpsMonitoringChannel = "C029ZH0BQSZ"
)

func userErrMessage(user string, err error) string {
	return fmt.Sprintf("<@%s>! %s\n\n ```%s```\n\n%s", user, msgErrNotification, err, msgErrNotificationAssurance)
}

func errMessage(err error) string {
	return fmt.Sprintf("%s\n\n ```%s```\n\n%s", msgErrNotification, err.Error(), msgErrNotificationAssurance)
}

func userNotificationMessage(user, msg string) string {
	return fmt.Sprintf("<@%s>! %s\n\n ```%s```\n\n", user, msgNotification, msg)
}

func userDeploymentNotificationMessage(user, msg string) string {
	return fmt.Sprintf("<@%s>! %s\n\n ```%s```\n\n", user, msgDeploymentErrNotification, msg)
}
