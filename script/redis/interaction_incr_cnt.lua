-- 1. 判断 key 是否存在 存在则增加计数
-- 2. HINCRBY key cntKey delta eg. HINCRBY bizID:biz readCnt 1
local key = KEYS[1] -- bizID:biz
local cntKey = ARGV[1] -- read_cnt like_cnt collect_cnt
local delta = tonumber(ARGV[2]) -- 1 or -1
local exists = redis.call("EXISTS", key)
if exists == 1 then
    redis.call("HINCRBY", key, cntKey, delta)
    return 1 -- success
else
    return 0
end