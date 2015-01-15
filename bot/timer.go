package bot

import (
	"time"
	"math/rand"
)

const (
	SPEED_INHUMAN     int = 500
	SPEED_FAST        int = 700
	SPEED_NORMAL      int = 1000
	SPEED_SLOW        int = 1500
	SPEED_REALLY_SLOW int = 2000
)

var (
	currentSpeed int 
)

func sleep( millis int, scaleDuration bool) {
	if scaleDuration {
		millis = millis * currentSpeed / SPEED_NORMAL
	}
	time.Sleep(time.Duration(millis) * time.Millisecond)
}

func sleepRangeScale(minMillis, maxMillis int, scaleDuration bool) {
	duration := rand.Intn(maxMillis-minMillis+1) + minMillis
	sleep(duration, scaleDuration)
}

func sleepRange(minMillis, maxMillis int) {
	sleepRangeScale(minMillis, maxMillis, true)
}

func PreparedFastAction() {
	sleepRange(400, 700)
}

func ActionInPage() {
	sleepRange(600, 1000)
}

func ChangePage() {
	sleepRange(900, 1500)
}

func ChangePageLong() {
	sleepRange(1000, 2000)
}

func ReadPage() {
	sleepRange(1200, 2500)
}

func PauseWhenSafe() {
	sleepRange(2000, 3000)
}

func WaitAfterLogin() {
	sleepRangeScale(6000, 7000, false)
}
