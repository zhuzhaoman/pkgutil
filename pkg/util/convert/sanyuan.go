package convert

/*
Go不支持运算符重载，因此需要先将 a<b 在函数外转换为 bool 条件
Go不支持泛型，只能用 interface{} 模拟
返回的类型安全需要用户自己保证，.(type) 的类型必须匹配
interface{} 是运行时泛型，性能没有编译时泛型高
*/
func IfSanyuan(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
