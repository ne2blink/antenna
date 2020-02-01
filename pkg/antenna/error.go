package antenna

import "go.uber.org/zap"

type warning error

func logEvent(log *zap.SugaredLogger, msg string, err error) {
	if err != nil {
		log = log.With("err", err.Error())
		if _, ok := err.(warning); ok {
			log.Warn(msg)
		} else {
			log.Error(msg)
		}
	} else {
		log.Info(msg)
	}
}
