local ERR = ngx.ERR
local log = ngx.log
local new_timer = ngx.timer.at
local pcall = pcall

local function _require(name)
  local ok, p = pcall(require, name)
  if not ok then
    error("failed to require: " .. name)
  end
  return p
end
local redis = _require("resty.redis")
local cjson = _require("cjson")

local function errlog(...)
  log(ERR, "pubsub: ", ...)
end

local function encode_record(key, value)
  return cjson.encode{key, value}
end

local function decode_record(record)
  local res = cjson.decode(record)
  return res[1], res[2]
end

local function store_record(ctx, red, key, value)
  local channel = ctx.channel
  local log_store = ctx.log_store
  local log_pos_store = ctx.log_pos_store

  local ok, err = red:watch(log_pos_store)
  if not ok then
    return nil, nil, "failed to watch position store: " .. err
  end

  local curpos, err = red:get(log_pos_store)
  if not tonumber(curpos) then
    curpos = 0
  end
  if err then
    return nil, nil, "failed to get log pos: " .. err
  end

  local ok, err = red:multi()
  if not ok then
    return nil, nil, "failed to exec multi: " .. err
  end

  local ok, err = red:hset(log_store, curpos, encode_record(key, value))
  if not ok then
    return nil, nil, "failed to queue hset: " .. err
  end

  local ok, err = red:incr(log_pos_store)
  if not ok then
    return nil, nil, "failed to queue incr log pos: " .. err
  end

  local ok, err = red:exec()
  if not ok then
    return nil, nil, "failed to exec: " .. err
  end

  -- publish key
  local n_publish, err = red:publish(channel, curpos)
  if err then
    return nil, nil, "failed to publish: " .. err
  end

  return curpos, n_publish, nil
end

local function read_to(ctx, to)
  local log_store = ctx.log_store

  local red = redis:new()
  local ok, err = red:connect(ctx.addr, ctx.port)
  if not ok then
    return "failed to connect: " .. err
  end

  for p = ctx.pos,to do
    local record, err = red:hget(log_store, p)
    if not record then
      return "failed to hget: " .. err
    end
    local key, value = decode_record(record)
    ctx:_apply(p, key, value)
  end

  -- close and store connection to pool
  local ok, err = red:set_keepalive(ctx.keepalive_idle_timeout, ctx.keepalive_pool_size)
  if not ok then
    return "failed to set keepalive: " .. err
  end

  return nil
end


local function subscribe(ctx, red)
  red:set_timeout(0)

  local ok, err = red:subscribe(ctx.channel)
  if not ok then
    return "failed to subscribe: " .. err
  end

  while true do
    local res, err = red:read_reply()
    if err then
      if err == "timeout" then
        break
      end
      return "failed to ready reply: " .. err
    end

    local err = read_to(ctx, res[3])
    if err then
      return err
    end
  end

  return nil
end

local run_subscribe
run_subscribe = function(premature, ctx)
  if premature then
    return
  end

  local red = redis:new()
  local ok, err = red:connect(ctx.addr, ctx.port)
  if not ok then
    errlog("failed to connect: " .. err)
  end
  if ok then
    local err = subscribe(ctx, red)
    if err then
      errlog(err)
    end
  end
  red:close()

  -- register the next timer
  local ok, err = new_timer(ctx.retry_interval, run_subscribe, ctx)
  if not ok then
    return nil, "failed to create timer: " .. err
  end
end


local _Class = {}

function _Class._apply(self, pos, key, value)
  local table = self.table

  if value == nil then
    table.remove(table, key)
  end
  table[key] = value
  self.pos = pos
end

function _Class.get(self, key)
  return self.table[key]
end

function _Class.set(self, key, value)
  -- connect to redis
  local red = redis:new()
  local ok, err = red:connect(self.addr, self.port)
  if not ok then
    return nil, nil, "failed to connect: " .. err
  end

  local pos, n_publish, err = store_record(self, red, key, value)
  if err then
    return nil, nil, "failed to add record: " .. err
  end

  -- close and store to connection pool
  local ok, err = red:set_keepalive(self.keepalive_idle_timeout, self.keepalive_pool_size)
  if not ok then
    return nil, nil, "failed to set keepalive: " .. err
  end

  return pos, n_publish, nil
end

function _Class.delete(self, key)
  return _Class.set(self, key, nil)
end

-- spawn the new subscriber thread
function _Class.spawn_subscriber(self)
  local ok, err = new_timer(0, run_subscribe, self)
  if not ok then
    return "failed to create timer: " .. err
  end
end

local _M = {
  _VERSION = '0.01'
}

local function parse_opts(opts)
  local channel = opts.channel
  if not channel then
    return nil, "redis_channel parameter is required"
  end

  local log_store = opts.log_store
  if not log_store then
    log_store = opts.channel .. ":log"
  end

  local log_pos_store = opts.log_pos_store
  if not log_pos_store then
    log_pos_store = opts.channel .. ":pos"
  end

  local addr = opts.addr
  if not addr then
    addr = "127.0.0.1"
  end

  local port = opts.port
  if not port then
    port = 6379
  end

  local retry_interval = opts.retry_interval
  if not retry_interval then
    retry_interval = 10
  end

  local keepalive_idle_timeout = opts.keepalive_idle_timeout
  if not keepalive_idle_timeout then
    keepalive_idle_timeout = 60
  end

  local keepalive_pool_size = opts.keepalive_pool_size
  if not keepalive_pool_size then
    keepalive_pool_size = 100
  end

  return {
    channel = channel,                               -- redis pubsub channel key (required)
    log_store = log_store,                           -- redis hash key to store logs (default: {channel}:log)
    log_pos_store = log_pos_store,                   -- redis key to store log position (default: {channel}:pos)
    addr = addr,                                     -- address to connect (default: 127.0.0.1)
    port = port,                                     -- port to connect (default: 6379)
    retry_interval = retry_interval,                 -- subuscribe retry interval in sec. (default: 10)
    keepalive_idle_timeout = keepalive_idle_timeout, -- keepalive idle timeout for redis connection in sec. (default: 60)
    keepalive_pool_size = keepalive_pool_size,       -- keepalive size for redis connection (default: 100)
  }, nil
end

function _M.new(_, opts)
  local new, err = parse_opts(opts)
  new.table = {}
  new.pos = 0
  setmetatable(new, { __index = _Class })
  return new, err
end

return _M
