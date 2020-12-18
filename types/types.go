package types
//自定义的一些type类型的数据

//Func 普通无参函数
type Func func()
//FuncWithParams 有参的函数, 因为golang不支持重载此处需兼容1个或多个参数所以使用可变参数
type FuncWithParams func(...interface{})