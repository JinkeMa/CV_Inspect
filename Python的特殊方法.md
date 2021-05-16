### 🐍 Python的特殊方法
Python的特殊方法是指类中以双下划线`__`开头和结尾的成员函数，这些都是保留函数
### 🐍 扩展与示例
### 🦅
```
class myClass(Object):
  def __init__(self,name):
     self.name = name
     
  def __call__(self):#直接将对象名当做函数名使用时执行
     print("__call__:".format(self.name))
     
  def __str__(self):#print(对象)时执行
     return "__str__:{}".format(self.name) 

def main():
  myObj = myClass("Hello")
  print(myObj)#__str__:Hello
  myObj()#__call__:Hello
  
if __name__='__main__':
  main()
```
结果：  
```
__str__:Hello  
myObj()__call__:Hello
```
