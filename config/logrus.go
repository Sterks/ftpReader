package config

import (
	"github.com/sirupsen/logrus"
)

//LogrusConf Тип логирования
type LogrusConf struct {
	logger *logrus.Logger
}

//New Ссылка на объект
func New() *LogrusConf {
	return &LogrusConf{}
}
