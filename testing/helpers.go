package testing

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"euromoby.com/core/db"
	"euromoby.com/core/logger"
)

const execTimeout = 2 * time.Second

func SetWorkingDir() {
	_, filename, _, _ := runtime.Caller(1) //nolint:dogsled
	dir := path.Join(path.Dir(filename), "../..")
	if err := os.Chdir(dir); err != nil {
		panic(err)
	}
}

func CleanupDatabase(db db.Conn, tables []string) {
	ctx, cancel := context.WithTimeout(context.Background(), execTimeout)
	defer cancel()

	for _, table := range tables {
		_, err := db.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s", table))
		if err != nil {
			logger.Fatal(err)
		}
	}
}
