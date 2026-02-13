package expiredcache

import (
	"fmt"
	"github.com/Pr3c10us/fanatix/internals/domains/queueport"
	"strconv"
	"strings"
)

func Handler(c *queueport.Context) {
	dataString := strings.Split(string(c.BytesMessage), ":")

	if len(dataString) > 1 {
		key := dataString[0]
		value := dataString[1]

		if key == c.EnvironmentVariables.RedisKeys.FixturesExternalKey {
			externalFixtureID, err := strconv.Atoi(value)
			if err != nil {
				fmt.Println("converting string to int error:", err)
				return
			}
			err = c.Services.PredictionServices.InitiateLineupVerification.Handle(int64(externalFixtureID))
			if err != nil {
				fmt.Println("verifying prediction error:", err)
				return
			}
		} else if key == c.EnvironmentVariables.RedisKeys.ScoresPredictionKey {
			externalFixtureID, err := strconv.Atoi(value)
			if err != nil {
				fmt.Println("converting string to int error:", err)
				return
			}
			err = c.Services.PredictionServices.InitiateScoresVerification.Handle(int64(externalFixtureID))
			if err != nil {
				fmt.Println("verifying prediction error:", err)
				return
			}
		}
	}

	c.Logger.Log("info", "event task complete")
}
