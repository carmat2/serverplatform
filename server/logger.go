package server

import (
	"fmt"
	"os"
	"time"
	"context"
	"encoding/base64"

	"github.com/mattermost/logr/v2"
	"github.com/mattermost/logr/v2/formatters"
	"github.com/mattermost/logr/v2/targets"
	
	"golang.org/x/oauth2"
  	"golang.org/x/oauth2/google"
  	"google.golang.org/api/gmail/v1"
  	"google.golang.org/api/option"
)

var lgr *logr.Logr
var logger logr.Sugar

// createLoggerCommon initializes the logr logger and creates the custome email alerts target.
func createLoggerCommon() {
	lgr,_ = logr.New()
	logger = lgr.NewLogger().Sugar()

	formatterAlert := &formatters.Plain{Delim: " | "}
	filterAlert := &logr.StdFilter{Lvl: logr.Fatal}

	targetEmail := NewLogTargetEmail()
	if targetEmail != nil {
		tAlert := targets.NewWriterTarget(targetEmail)
		lgr.AddTarget(tAlert, "EmailAlert", filterAlert, formatterAlert, 1000)
	} else {
		// TODO
	}
}

// CreateLoggerTargetTesting creates the logger target for unit tests.
// Log output is sent to os.Stdout from the Trace level, adding the stack trace from the Error level.
// Email alerts will log from the Fatal level.
func CreateLoggerTargetTesting() {
	createLoggerCommon()

	formatterStdOut := &formatters.Plain{DisableTimestamp: true, Delim: " | "}
	filterStdOut := &logr.StdFilter{Lvl: logr.Trace, Stacktrace: logr.Error}
	tStdOut := targets.NewWriterTarget(os.Stdout)
	lgr.AddTarget(tStdOut, "StdOut", filterStdOut, formatterStdOut, 1000)
}

// CreateLoggerTargetTesting creates the logger targets for application run.
// Log output is sent to ./logs/serverplatform.log from the Info level, adding the stack trace from the Error level.
// Email alerts will log from the Fatal level.
func CreateLoggerTarget() {
	createLoggerCommon()

	formatterLogFile := &formatters.JSON{}
	filtersLogFile := &logr.StdFilter{Lvl: logr.Info, Stacktrace: logr.Error}
	// max file size 10MB, keep log files for 30 days, keep up to 5 old backup log files, no file compression
	opts := targets.FileOptions{
		Filename:   "./logs/serverplatform.log",
		MaxSize:    10,
		MaxAge:     30,
		MaxBackups: 5,
		Compress:   false,
	}
	tLogFile := targets.NewFileTarget(opts)
	lgr.AddTarget(tLogFile, "LogFile", filtersLogFile, formatterLogFile, 1000)
}

// ShutdownLogger ensures targets are drained before application exit.
func ShutdownLogger() {
	lgr.Shutdown()
}

//logTargetEmail defines a log target to send alert emails
type logTargetEmail struct {
	targets.Writer
	gmailService *gmail.Service
}

// NewLogTargetEmail creates a new alert targwt instance
func NewLogTargetEmail() *logTargetEmail {
	t := &logTargetEmail{}

	// TODO - do not call Init directly, it should be called by logr
	err := t.Init()
	if err != nil {
		return nil
	}
	return t
}

// Called once to initialize target.
// Creates and intializes the gmail service
func (t *logTargetEmail) Init() error {
	config := oauth2.Config{
    	ClientID:     "hidden",
    	ClientSecret: "hidden",
    	Endpoint:     google.Endpoint,
    	RedirectURL:  "http://localhost",
  	}

  	token := oauth2.Token{
    	AccessToken:  "hidden",
    	RefreshToken: "hidden",
    	TokenType:    "Bearer",
    	Expiry:       time.Now(),
  	}

  	var tokenSource = config.TokenSource(context.Background(), &token)

	var err error
  	t.gmailService, err = gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
  	if err != nil {
		logger.Errorf("Failed to initialize the gmail service %v", err)
		return err
 	 }
  	
    logger.Info("Gmail service is initialized")
  	return nil
}

// Write will always be called by a single internal Logr goroutine, so no locking needed.
// Sends an alert email
func (t *logTargetEmail) Write(p []byte) (int, error) {
	var message gmail.Message
  
	toEmail := "hidden"
	srcEmail := "hidden"
	srcDisplayName := "ServerPlatform v1"
	fromMail := fmt.Sprintf("From: %s <%s> \r\n", srcDisplayName, srcEmail)
	emailTo := "To: " + toEmail + "\r\n"
	subject := "Subject: " + "Alert from ServerPlatform v1" + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(fromMail + emailTo + subject + mime + "\n" + string(p))
  
	message.Raw = base64.URLEncoding.EncodeToString(msg)
	
	_, err := t.gmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
	  logger.Errorf("Email logger target error sending email %s", err.Error())
	  return 0, err
	}

	return len(message.Raw), nil
}
  
// Called once to cleanup/free resources for target.
func (t *logTargetEmail) Shutdown() error {
	return nil
}