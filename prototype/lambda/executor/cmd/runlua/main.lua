-- 基于github.com/vadv/gopher-lua-libs扩展的模块
local json = require("json") 
local http = require("http")
local base64 = require("base64")
local tcp = require("tcp")
local time = require("time")

-- lua模块(同目录下的另一个lua文件)
local utils = require("utils")

-- ipes_开头的是扩展的函数

function is_table_empty(tbl)
    return next(tbl) == nil
end

function check_table(expect, examined)
    for k,v in pairs(expect) do
        if examined[k] == nil then
            return false, string.format("[%s] not in keys", k)
        end

        local value = examined[k]

        if value ~= v then
            return false, string.format("headers[%s:%s] mismatch [%s:%s]", k, v, k, value)
        end
    end

    return true, ""
end

function http_probe(input_table, report_data_table)

    local target = input_table["target"]
    local timeout = target["timeout"]
    local method = target["method"]
    local url = target["url"]
    local body = target["body"]
    local headers = target["headers"]
    local basic_auth = target["basic_auth"]

    local client = http.client({
        timeout = timeout,
        insecure_ssl = true,

    })

    local request
    if body == nil or body == "" then
        request = http.request(method, url)
    else
        request = http.request(method, url, body)
    end
    
    for k, v in pairs(headers or {}) do
        request:header_set(k, v)
    end

    if basic_auth ~= nil then
        request:set_basic_auth(basic_auth["user"], basic_auth["pass"])
    end

    local result, err = client:do_request(request)
    if err then
        report_data_table["status"] = "failed"
        report_data_table["msg"] = string.format("request failed:%s", err)
        return
    end

    function check_expect()
        local expect = input_table["expect"]

        -- check code
        if not (result.code == expect["code"] ) then
            report_data_table["status"] = "failed"
            report_data_table["msg"] = string.format("code [%s] in not [%s]", result.code, expect["code"])
            return
        end
    
        -- check body
        local base64_data = expect["base64_data"]
        local data

        if base64_data == nil or base64_data == "" then
            data = expect["data"]
        else
            local decoded, err = base64.StdEncoding:decode_string(base64_data)
            assert(not err, err)
            data = decoded
        end

        if data ~= nil and result.body ~= data then
            report_data_table["status"] = "failed"
            report_data_table["msg"] = string.format("body [%s] is not [%s]", result.body, data)
            return
        end

        -- check headers
        local headers = expect["headers"]
        if headers ~= nil then
            local ok, msg = check_table(headers, result.headers)
            if not ok then
                report_data_table["status"] = "failed"
                report_data_table["msg"] = msg
                return
            end
        end

    end

    local expect = input_table["expect"]
    if  expect ~= nil and not is_table_empty(expect) then
        print("check_expect")
        check_expect()
    else
        print("no check_expect")
        -- report code, body, headers
        report_data_table["code"] = result.code
        report_data_table["headers"] = result.headers
        report_data_table["body"] = result.body
    end

    -- print("output:")
    -- utils.printTable(report_data_table)

    return
end

function tcp_probe(input_table, report_data_table)
    sock("tcp", input_table, report_data_table)
end

function udp_probe(input_table, report_data_table)
    sock("udp", input_table, report_data_table)
end

function sock(proto, input_table, report_data_table)
    print(proto)
    local target = input_table["target"]

    local url = target["url"]
    local dial_timeout = target["dial_timeout"]
    local write_timeout = target["write_timeout"]
    local read_timeout = target["read_timeout"]
    local close_timeout = target["close_timeout"]
    local data = target["data"]
    local read = target["read"]

    if dial_timeout == nil then
        dial_timeout = 5
    end

    local conn, err = tcp.open(url, dial_timeout, proto)
    if err then
        print(err)
        report_data_table["status"] = "failed"
        report_data_table["msg"] = string.format("dail failed:%s", err)
        return
    end

    if write_timeout then
        conn.writeTimeout = write_timeout
    end
    if write_timeout then
        conn.readTimeout = read_timeout
    end
    if write_timeout then
        conn.closeTimeout = close_timeout
    end

    if data then
        err = conn:write(data)
        if err then
            report_data_table["status"] = "failed"
            report_data_table["msg"] = string.format("write failed:%s", err)
            return
        end
    end

    function check_expect()
        local expect = input_table["expect"]

        -- check body
        local base64_data = expect["base64_data"]
        local data

        if base64_data == nil or base64_data == "" then
            data = expect["data"]
        else
            local decoded, err = base64.StdEncoding:decode_string(base64_data)
            assert(not err, err)
            data = decoded
        end

        if data ~= nil then
            -- read #data bytes; 注意：readall 会造成read超时。
            local result, err = conn:read(#data)
            if err then
                report_data_table["status"] = "failed"
                report_data_table["msg"] = string.format("read failed:%s", err)
                return
            end

            if result ~= data then
                report_data_table["status"] = "failed"
                report_data_table["msg"] = string.format("body [%s] is not [%s]", result.body, data)
                return
            end
        end
    end

    local expect = input_table["expect"]
    if  expect ~= nil and not is_table_empty(expect) then
        check_expect()
    end

    conn:close()
end

local process_func = {
    http = http_probe,
    tcp = tcp_probe,
    udp = udp_probe,
}

function main()
    print( "script dir: " .. ipes_script_dir())

    -- 获取输入数据
    local input_table = ipes_input()

    local probe_type = input_table["type"]
    local task_id = input_table["task_id"]
    local lambda_meta = input_table["lambda_meta"]
    local meta = input_table["meta"]

    -- call process function
    if process_func[probe_type] == nil then
        print(probe_type, "not supported!")
        return
    end

    local func = process_func[probe_type]
    if func == nil then
        print(string.format("probe type [%s] is not supported!", probe_type))
    end

    local report_data_table = {
        id = task_id,
        status = "ok",
        lambda_meta = lambda_meta,
        msg = "",
        meta = meta
    }

    -- printTable(report_data_table)

    local begin = time.unix()
    local ok, err = pcall(func, input_table, report_data_table)
    local stop = time.unix()
    
    report_data_table["time_used"] = stop - begin

    -- error handler
    if not ok then
        report_data_table["status"] = "err"
        report_data_table["msg"] = err
    end

    -- report result
    local result, err = json.encode(report_data_table)
    assert(not err, err)

    -- local err = ipes_report("kafka", "common", "ipes-test", "", result)
    -- if err ~= nil then
    --     print("report err: "..err)
    --     return
    -- end

    local err = ipes_report(result)
    if err ~= nil then
        print("report err: "..err)
        return
    end

end

main()
