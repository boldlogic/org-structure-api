package repository

import (
	"github.com/boldlogic/packages/shutdown"
	"go.uber.org/zap"
)

func (r *Repo) logWrapper(f string, err error) {
	if err == nil {
		r.logger.Debug("успех", zap.String("func", f))
		return
	}

	if shutdown.IsExceeded(err) {
		return
	}
	r.logger.Error("ошибка выполнения", zap.String("func", f), zap.Error(err))

}
