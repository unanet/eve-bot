package slackservice

import "fmt"

const (
	msgErrNotification           = "Something terrible has happened..."
	msgErrNotificationAssurance  = "I've dispatched a semi-competent team of monkeys to battle the issue..."
	msgNotification              = "I've got some news..."
	msgDeploymentErrNotification = "I detected some deployment *errors:*"
	msgLogLinks                  = "here are the latest logs..."
)

const (
	// https://clearview.slack.com/archives/CUK5MSMPU
	devOpsMonitoringChannel = "CUK5MSMPU"
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
