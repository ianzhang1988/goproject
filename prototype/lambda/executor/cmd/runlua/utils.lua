function printTable(tbl, indent)
    indent = indent or 0

    for key, value in pairs(tbl) do
        if type(value) == "table" then
            print(string.rep("  ", indent) .. key .. ":")
            printTable(value, indent + 1)
        else
            print(string.rep("  ", indent) .. key .. ": " .. tostring(value))
        end
    end
end

function copyTable(tbl)

    copy = {}
    for key, value in pairs(tbl) do
        copy[key] = value
    end
    return copy
end
return {
    printTable = printTable,
    copyTable = copyTable
}