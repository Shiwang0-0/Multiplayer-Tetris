package protocol

import "time"

const VOTE_COUNT_DOWN_TIME_IN_SECOND = 5
const CLEAR_ROW_TIME_IN_MILISECOND = 75
const FALL_TIME_IN_MILISECOND = 800

var VoteCountdownDuration = time.Duration(VOTE_COUNT_DOWN_TIME_IN_SECOND) * time.Second

var ClearRowDuration = time.Duration(CLEAR_ROW_TIME_IN_MILISECOND) * time.Millisecond

var FallDuration = time.Duration(FALL_TIME_IN_MILISECOND) * time.Millisecond
