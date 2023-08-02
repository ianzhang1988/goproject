local json = require("json") 
local http = require("http")
-- 上面是注册到lua里面的模块
local utils = require("utils")
-- 这个是同目录下的另一个lua文件

-- ipes_开头的是我写
local input_table = ipes_input()

-- print("input:")
-- local target = utils.printTable(input_table)

local target = input_table["target"]

local client = http.client({timeout = 10})

-- GET
local ok, report_data = pcall( function ()
    local request = http.request("GET", target["url"] )
    local result, err = client:do_request(request)
    if err then
        error(err)
    end

    report_data_table = {
        status = "ok",
        msg = ""
    }

    function helper()
        if not (result.code == target["code"] ) then
            report_data_table["status"] = "failed"
            report_data_table["msg"] = string.format("code [%s] in not [%s]", result.code, target["code"])
            return
        end
    
        if result.body ~= target["body"] then
            report_data_table["status"] = "failed"
            report_data_table["msg"] = string.format("body [%s] in not [%s]", result.body, target["body"])
            return
        end
    end

    helper()

    -- print("output:")
    -- utils.printTable(report_data_table)

    local result, err = json.encode(report_data_table)
    if err then
        error(err)
    end

    return result
end)


if not ok then
    report_data = "err: "..report_data
end

local err = ipes_report(report_data)
if err ~= nil then
    print("report err: "..err)
end