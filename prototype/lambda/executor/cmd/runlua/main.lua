-- 基于github.com/vadv/gopher-lua-libs扩展的模块
local json = require("json") 
local http = require("http")
local base64 = require("base64")
local tcp = require("tcp")
local time = require("time")
local net = require("net")

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

function dns_check()
    local domain_list = {
        "www.baidu.com",
        "www.qq.com",
        "www.alibaba.com",
    }

    local ok = false
    for _, domain in ipairs(domain_list) do
        addr, err = net.dnslookup(domain)
        if err ~= nil then
            goto continue
        end
        if #addr == 0 then
            goto continue
        end
       
        ok = true
        print(domain .. ": ".. addr[1])

        -- gopher lua bug
        if ok then
            break
        end

        ::continue::
    end

    return ok
end

function http_probe(input_table, report_data_table)

    local target = input_table["target"]
    local if_statistics = input_table["statistics"]
    local if_target_ipport = input_table["target_ipport"]
    local timeout = target["timeout"]
    local method = target["method"]
    local url = target["url"]
    local body = target["body"]
    local headers = target["headers"]
    local basic_auth = target["basic_auth"]

    if timeout == nil then
        timeout = 5
    end

    local http_config = {
        timeout = timeout,
        insecure_ssl = true, 
    }
    
    local sip_dial
    if if_target_ipport then
        local err
        sip_dial, err = http.custom_dial("target_ipport", timeout)
        if err ~= nil then
            print("custom dial err:", err)
        else
            http_config["dial_context"] = sip_dial
        end
    end

    local client = http.client(http_config)

    local request
    if body == nil or body == "" then
        print("-----1 "..method..url)
        request = http.request(method, url)
    else
        request = http.request(method, url, body)
    end

    if if_statistics then
        request, statistics = http.attach_statistics(request)
    end
    
    for k, v in pairs(headers or {}) do
        if string.lower(k) == "host" then
            request:set_host(v)
        else
            request:header_set(k, v)
        end
    end

    if basic_auth ~= nil then
        request:set_basic_auth(basic_auth["user"], basic_auth["pass"])
    end

    local result, err = client:do_request(request)
    
    -- try get even when err occur
    if if_statistics then
        report_data_table["statistics"] = http.get_statistics(statistics)
    end
    if if_target_ipport and sip_dial then
        report_data_table["target_ipport"] = http.get_source_ipport(sip_dial)
    end
    
    -- err from do_request
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
            report_data_table["msg"] = string.format("code [%s] in not [%s], body:[%s]", result.code, expect["code"], result.body)
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
        check_expect()
    else
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

function dns_lookup(input_table, report_data_table)
    local target = input_table["target"]

    local url = target["url"]
    local timeout = target["timeout"]

    if timeout == nil then
        timeout = 3 
    end

    local addr, err = net.dnslookup(url, timeout)
    if err then
        report_data_table["status"] = "failed"
        report_data_table["msg"] = string.format("dns lookup failed:%s", err)
        return
    end

    report_data_table["data"] = {
        domain = url,
        addr = addr,
    }

end

function ping(input_table, report_data_table)
    local target = input_table["target"]
    local url = target["url"]

    local count = 3
    local params = input_table["params"]
    if params ~= nil then
        local num_str = params["count"]
        if num_str ~= nil then
            local num = tonumber(num_str)
            if num ~= nil then
                count = num
            end
        end
    end

    if count > 3 then
        count =3
    end

    local stats, err = net.ping(url, count)
    if err then
        report_data_table["status"] = "failed"
        report_data_table["msg"] = string.format("dns lookup failed:%s", err)
        return
    end

    report_data_table["data"] = stats
    if stats["pkt_recv"] < 1 then
        report_data_table["status"] = "failed"
    end

end

local process_func = {
    http = http_probe,
    tcp = tcp_probe,
    udp = udp_probe,
    dns = dns_lookup,
    ping = ping,
}

function main()
    print( "script dir: " .. ipes_script_dir())

    -- 获取输入数据
    local input_table = ipes_input()

    local probe_type = input_table["type"]
    local task_id = input_table["task_id"]
    local lambda_meta = input_table["lambda_meta"]
    local meta = input_table["meta"]
    local if_check_dns = input_table["dns_check"]
    local retry = input_table["retry"]

    if if_check_dns == nil then
        if_check_dns = true
    end
    
    if retry == nil then
        retry = 3
    end

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

    local dns_ok
    local dns_begin
    local dns_stop
    if if_check_dns == false then
        dns_ok = true
    else
        print("check dns")
        dns_begin = time.unix()
        dns_ok = dns_check()
        dns_stop = time.unix()
    end

    -- printTable(report_data_table)

    if dns_ok then

        local begin = time.unix()
        local ok = false
        for i = 1, retry do 
            ok, err = pcall(func, input_table, report_data_table)
            local status_ok = report_data_table["status"]
            if ok and status_ok == "ok" then
                break
            end
            print("retry: ", i)
            time.sleep(1)
        end
        local stop = time.unix()
        
        report_data_table["time_used"] = stop - begin
        -- error handler
        if not ok then
            report_data_table["status"] = "err"
            report_data_table["msg"] = err
        end
    else
        report_data_table["time_used"] = dns_stop - dns_begin
        -- dns is faulty, skip this device
        report_data_table["status"] = "err"
        report_data_table["msg"] = "faulty dns"
    end


    -- report result
    local result, err = json.encode(report_data_table)
    assert(not err, err)

    local err = ipes_report(result)
    if err ~= nil then
        print("report err: "..err)
        return
    end

end

main()
