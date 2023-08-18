local tcp = require("tcp")
local strings = require("strings")
local base64 = require("base64")

function tcp_test()
    -- ping pong game
    local conn, err = tcp.open(":12345")
    if err then error(err) end

    err = conn:write("ping")
    if err then error(err) end

    local result, err = conn:read()
    if err then error(err) end
    result = strings.trim(result, "\n")
    print(result)
    print(#result)
    print(result == "你好")

    print(base64.StdEncoding:encode_to_string("你好"))
    local s, err = base64.StdEncoding:decode_string("5L2g5aW9")
    assert(not err, err)
    print("base64:", s)
    print(result == s)


    -- if (result ~= "pong") then error("must be pong message") end

    local t = {} -- 用来存储ASCII码值的table
    for i = 1, #result do -- 遍历字符串中的每个字符
    t[i] = string.byte(result, i) -- 将ASCII码值插入table中
    end

    for i, c in ipairs(t) do
        print(string.format("idx:%d, char:%d, type:%s",i, c, type(c)))
    end
end

function udp_test()
    local conn, err = tcp.open(":12345", 1, "udp")
    if err then error(err) end

    err = conn:write("早上好")
    if err then error(err) end

    local result, err = conn:read()
    if err then error(err) end
    result = strings.trim(result, "\n")

    print(result)
end

tcp_test()
-- udp_test()
