package snowflakeutil

import (
	"log"
	"time"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func Init(startTime string, machineID int64) {
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		log.Fatalf("snowflake starttime parse failed: %s \n", err)
	}
	snowflake.Epoch = st.UnixMilli()
	node, err = snowflake.NewNode(machineID)
	if err != nil {
		log.Fatalf("snowflake NewNode failed: %s \n", err)
	}
}

func GenID() int64 {
	return node.Generate().Int64()
}
