package container

import (
	"github.com/melodywen/supports/exceptions"
	"github.com/melodywen/supports/str"
	"reflect"
	"strings"
	"sync"
)

// outTypeInfo
//  @Description:
type outTypeInfo struct {
	isPtr      bool // 是否为指针类型
	originType reflect.Type
	newType    reflect.Type
}

// InTypeInfo
//  @Description:
type inTypeInfo struct {
	originType reflect.Type
}

// buildInfo
//  @Description:
type buildInfo struct {
	identifier        string        // 标识符
	isConstructMethod bool          //是否为构造方法
	concreteType      reflect.Type  // concrete 反射type
	concreteValue     reflect.Value // concrete 反射 value
	numOut            int           // concrete 出参个数
	numIn             int           // 入参个数
	outType           []outTypeInfo // 出参的类型
	inType            []inTypeInfo  // 入参的类型
}

// Build
//  @Description:实例化给定类型的具体实例
//  @receiver container
//  @param concrete
//  @param parameters
//  @param stack
//  @return object
func (container *Container) Build(concrete any, parameters []any, stack *containerStack) (object any) {
	var params []reflect.Value

	info := cacheBuildInfo(concrete)

	// 是否是构造函数压入栈
	if info.isConstructMethod {
		stack.Stack = append(stack.Stack, info.identifier)
	}

	for _, parameter := range parameters {
		params = append(params, reflect.ValueOf(parameter))
	}

	params = container.getParameters(params, info, stack)

	// 调用函数
	resultList := info.concreteValue.Call(params)
	// 然后进行克隆反射
	response := []any{}
	for m := 0; m < info.numOut; m++ {
		if info.outType[m].isPtr {
			returnNew := reflect.New(info.outType[m].newType).Elem() //创建对象 //获取源实际类型(否则为指针类型)
			returnValue := resultList[m]                             //源数据值
			returnValue = returnValue.Elem()                         //源数据实际值（否则为指针）
			returnNew.Set(returnValue)                               //设置数据
			returnNew = returnNew.Addr()                             //创建对象的地址（否则返回值）
			response = append(response, returnNew.Interface())       //返回地址
		} else {
			returnNew := reflect.New(info.outType[m].newType).Elem() //创建对象
			returnValue := resultList[m]                             //源数据值
			returnNew.Set(returnValue)                               //设置数据
			response = append(response, returnNew.Interface())       //返回
		}
	}
	if info.isConstructMethod {
		stack.Stack = stack.Stack[:len(stack.Stack)-1]
	}
	if len(response) == 1 {
		return response[0]
	}
	return response
}

var buildInfoOfConcreteCache sync.Map // map[reflect.Type]buildInfo

// cacheBuildInfo
//  @Description:
//  @param concrete
//  @return *buildInfo
func cacheBuildInfo(concrete any) *buildInfo {
	identifier := AbstractToString(concrete)
	identifier = str.AfterLast(identifier, "/")
	// 如果是匿名函数 则 跳过缓存
	if strings.Count(identifier, ".") >= 3 {
		return getBuildInfo(concrete)
	}
	t := reflect.TypeOf(concrete)
	if item, ok := buildInfoOfConcreteCache.Load(t); ok {
		return item.(*buildInfo)
	}
	item, _ := buildInfoOfConcreteCache.LoadOrStore(t, getBuildInfo(concrete))
	return item.(*buildInfo)
}

// getBuildInfo
//  @Description: 获取构建信息
//  @param concrete
//  @return *buildInfo
func getBuildInfo(concrete any) *buildInfo {
	info := &buildInfo{
		identifier:        AbstractToString(concrete),
		isConstructMethod: false,
		concreteType:      reflect.TypeOf(concrete),
		concreteValue:     reflect.ValueOf(concrete),
		numOut:            0,
		numIn:             0,
		outType:           []outTypeInfo{},
		inType:            []inTypeInfo{},
	}
	if info.concreteType.Kind() != reflect.Func {
		panic(exceptions.NewInvalidParamErrorWithData("concrete type must be a function", errorOther(map[string]any{
			"concrete":   concrete,
			"identifier": AbstractToString(concrete),
		})))
	}
	info.isConstructMethod = isConstructMethod(info.identifier)
	info.numOut = info.concreteValue.Type().NumOut()
	info.numIn = info.concreteValue.Type().NumIn()
	for m := 0; m < info.numOut; m++ {
		outInfo := outTypeInfo{
			isPtr:      false,
			originType: info.concreteValue.Type().Out(m),
			newType:    info.concreteValue.Type().Out(m),
		}
		if outInfo.originType.Kind() == reflect.Ptr {
			outInfo.isPtr = true
			outInfo.newType = outInfo.newType.Elem()
		}
		info.outType = append(info.outType, outInfo)
	}
	for m := 0; m < info.numIn; m++ {
		info.inType = append(info.inType, inTypeInfo{originType: info.concreteValue.Type().In(m)})
	}
	return info
}

// getParameters
//  @Description: 获取构造参数。
//  @receiver container
//  @param currentParams
//  @param info
//  @param stack
//  @return finalParam
func (container *Container) getParameters(currentParams []reflect.Value, info *buildInfo, stack *containerStack) (finalParam []reflect.Value) {
	if info.numIn == len(currentParams) {
		return currentParams
	}
	finalParam = make([]reflect.Value, info.numIn)
	// 如果不是则所有的参数通过反射获取
	for i := 0; i < info.numIn; i++ {
		finalParam[i] = container.resolveDependencies(info.inType[i].originType, stack)
	}
	return finalParam
}

var stackIdentifier string
var containerIdentifier string

func init() {
	stackIdentifier = AbstractToString(&containerStack{})
	containerIdentifier = AbstractToString(&Container{})
}

// resolveDependencies
//  @Description: 反射获取依赖实参
//  @receiver container
//  @param parameterType
//  @param stack
//  @return reflect.Value
func (container *Container) resolveDependencies(parameterType reflect.Type, stack *containerStack) reflect.Value {
	// 对类型进行判定
	var identifier string
	switch parameterType.Kind() {
	case reflect.Interface, reflect.Struct, reflect.Ptr:
		identifier = cachedAbstractTypeToString(parameterType)
	default:
		panic(exceptions.NewInvalidParamErrorWithData(
			"can not auto load param,because param type can not suppose ,please connect admin",
			errorOther(map[string]any{
				"identifier": AbstractToString(identifier),
				"type":       parameterType.Kind().String(),
			})))
	}
	var object interface{}
	if container.Bound(identifier) {
		object = container.makeWithBuildStack(identifier, nil, stack)
	} else if identifier == stackIdentifier {
		// 如果是特殊是的字段————构建的堆栈信息
		object = stack
	} else if identifier == containerIdentifier {
		object = container
	} else {
		panic(exceptions.NewInvalidParamErrorWithData(
			"can not auto load param，because this interface is not registered,please connect admin",
			errorOther(map[string]any{
				identifier: identifier,
			})))
	}
	return reflect.ValueOf(object)
}
