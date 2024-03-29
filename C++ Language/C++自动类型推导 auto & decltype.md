### 🐋 自动类型推导(~~编译期~~) auto & decltype
#### 🍎 使用 auto(只能对变量)以及使用decltype(对表达式)进行类型推导
* *auto 变量标识符 = 字面值常量或者表达式结果或者函数返回值;*  
* `auto`不能推导函数参数和数组的类型
```
int func(auto x,auto y){};//非法
auto arr[]={0,};//非法
```
* *decltype(任意复杂的表达式) 变量标识符;*  
* 关键字`auto`可以简化变量类型声明，将冗长的声明工作交给编译器，而`decltype`是为了解决`auto`只能对变量进行类型推导的问题
* `decltype`推导的类型：  
    * `decltype(exp) var;//exp为左值或者带圆括号()，则var为引用类型`
    * 其他情况，var类型和exp类型一致：函数返回值、类成员、普通变量
#### 🍎 扩展与示例
* 🍒 `auto` 变量类型推导：
```
auto a=5;
auto b=5;
int& a0=a;
vector<int> vec={1,2,3,4,5};
auto a1=new int(10);
auto a2=vec.begin();//类型 vector<int>::iterator
```
* 🍒 `decltype` 类型推导：
```
decltype(a+b) c;//int
decltype(a0) c1;//int&
decltype((a+b)) c2;//int&
decltype(a1) c3;//int*
```
* 🍒 函数类型推导：
```
//尾返回类型
template<typename T1>
auto func(T x,T y) -> decltype(x+y)
{
 return x+y;
}
//不使用尾返回类型，C++17之前不支持
template<typename T1>
auto func(T x,T y)
{
 return x+y;
}
```
* 🍒 `decltype(auto)`用于函数参数转发时，自动类型推导，可以避免写冗长的函数返回类型
```
std::string& func0()
{
   //do something;
}

decltype(auto) func()
{
   return func0();
}
```
