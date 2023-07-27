print("hello")

-- 定义一个函数，用于计算两个数的和
function add(a, b)
    return a + b
end

-- 调用函数并打印结果
local result = add(5, 10)  
print("结果是:", result)

-- 定义一个条件语句
local num = 15
if num > 10 then
    print("数字大于10")
elseif num == 10 then
    print("数字等于10")
else
    print("数字小于10")
end

-- 定义一个数据容器，这里使用Lua的表（table）
local person = {
    name = "John",
    age = 30,
    city = "New York",
    isMarried = true,
    hobbies = {"reading", "swimming", "coding"}
}

-- 打印数据容器中的内容
print("姓名:", person.name)
print("年龄:", person.age)
print("城市:", person.city)
print("婚姻状态:", person.isMarried)
print("爱好:")
for i, hobby in ipairs(person.hobbies) do
    print("  -", hobby)
end


-- 创建一个空的map容器
local map = {}

-- 添加键值对到map中
map["apple"] = 10
map["banana"] = 5
map["orange"] = 8

-- 获取map中的值
local numApples = map["apple"]
local numBananas = map["banana"]
local numOranges = map["orange"]

-- 输出map中的值
print("苹果数量:", numApples)
print("香蕉数量:", numBananas)
print("橙子数量:", numOranges)

-- 遍历map容器
print("遍历map容器:")
for key, value in pairs(map) do
    print(key, ":", value)
end

local m = require("module")
-- 判断是否加载成功
if m then
    print("模块加载成功！")
    -- 调用模块中的函数和变量
    m.hi("Alice")
else
    print("模块加载失败！")
end

print("调用go hi")

print(gohi("Jim"))

local http = require("http")
local client = http.client({timeout = 10})

-- GET
local request = http.request("GET", "http://www.iqiyi.com" )
local result, err = client:do_request(request)
if err then
    error(err)
end
if not (result.code == 200) then
    error("code")
end
-- print(result.body)
print(string.sub(result.body, 0, 20))





