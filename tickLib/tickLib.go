package tickLib

/*
tickLib is a library for tick.
tick可以指定时间间隔执行某个函数，也可以指定时间点执行某个函数。
*/

// 指定某个时间间隔执行某个函数
// interval: 间隔时间，单位为秒
// f: 要执行的函数
// args: 要执行的函数的参数
// 返回值: 无
func SetInterval(interval int, f func(args ...interface{}), args ...interface{}) {

}
