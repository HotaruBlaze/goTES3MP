local GoTES3MPUtils = {}

local charset = {}  do -- [0-9a-zA-Z]
    for c = 48, 57  do table.insert(charset, string.char(c)) end
    for c = 65, 90  do table.insert(charset, string.char(c)) end
    for c = 97, 122 do table.insert(charset, string.char(c)) end
end

GoTES3MPUtils.randomString = function(length)
    if not length or length <= 0 then return '' end
    math.randomseed(os.clock()^5)
    return GoTES3MPUtils.randomString(length - 1) .. charset[math.random(1, #charset)]
end

return GoTES3MPUtils