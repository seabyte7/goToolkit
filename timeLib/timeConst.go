package timeLib

const (
	SecondMillisecondRatio     = 1000
	MinuteSecondRatio          = 60
	HourMinuteRatio            = 60
	MillisecondNanosecondRatio = 1000000
	DayHourRatio               = 24

	SecondNanosecondRatio = SecondMillisecondRatio * MillisecondNanosecondRatio

	MinuteMillisecondRatio = MinuteSecondRatio * SecondMillisecondRatio
	MinuteNanosecondRatio  = MinuteMillisecondRatio * MillisecondNanosecondRatio

	HourSecondRatio      = HourMinuteRatio * MinuteSecondRatio
	HourMillisecondRatio = HourSecondRatio * SecondMillisecondRatio

	DayMinuteRatio      = DayHourRatio * HourMinuteRatio
	DaySecondRatio      = DayMinuteRatio * MinuteSecondRatio
	DayMillisecondRatio = DaySecondRatio * SecondMillisecondRatio
)
