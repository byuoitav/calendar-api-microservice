package helpers

import (
	"fmt"

	"github.com/byuoitav/common/structs"

	"github.com/byuoitav/common/db"
	"github.com/byuoitav/common/log"
)

//GetCouchConfig makes request to couchdb to get scheduling-config data
func GetCouchConfig(room string) (structs.ScheduleConfig, error) {
	config, err := db.GetDB().GetScheduleConfig(room)
	if err != nil {
		log.L.Errorf("Error while trying to get Schedule configuration from the database | %v", err)
		return config, fmt.Errorf("Error while trying to get Schedule configuration from the database | %v", err)
	}
	return config, err
}
