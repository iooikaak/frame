package redis

import "github.com/go-redis/redis/v8"

type Client = redis.Client
type Cmder = redis.Cmder
type Cmdable = redis.Cmdable
type ScanIterator = redis.ScanIterator
type Pipeline = redis.Pipeline
type PubSub = redis.PubSub
type Pipeliner = redis.Pipeliner
type Subscription = redis.Subscription
type Message = redis.Message
type Pong = redis.Pong

type Cmd = redis.Cmd
type BoolCmd = redis.BoolCmd
type BoolSliceCmd = redis.BoolSliceCmd
type StringCmd = redis.StringCmd
type StringSliceCmd = redis.StringSliceCmd
type StringStringMapCmd = redis.StringStringMapCmd
type StringIntMapCmd = redis.StringIntMapCmd
type StringStructMapCmd = redis.StringStructMapCmd
type IntCmd = redis.IntCmd
type IntSliceCmd = redis.IntSliceCmd
type SliceCmd = redis.SliceCmd
type ZSliceCmd = redis.ZSliceCmd
type FloatCmd = redis.FloatCmd
type FloatSliceCmd = redis.FloatSliceCmd
type ZWithKeyCmd = redis.ZWithKeyCmd
type XStreamSliceCmd = redis.XStreamSliceCmd
type XPendingExtCmd = redis.XPendingExtCmd
type XMessageSliceCmd = redis.XMessageSliceCmd
type XPendingCmd = redis.XPendingCmd
type XInfoStreamCmd = redis.XInfoStreamCmd
type XInfoGroupsCmd = redis.XInfoGroupsCmd
type XInfoConsumersCmd = redis.XInfoConsumersCmd
type TimeCmd = redis.TimeCmd
type StatefulCmdable = redis.StatefulCmdable
type StatusCmd = redis.StatusCmd
type SlowLogCmd = redis.SlowLogCmd
type ScanCmd = redis.ScanCmd
type GeoPosCmd = redis.GeoPosCmd
type GeoLocationCmd = redis.GeoLocationCmd
type DurationCmd = redis.DurationCmd
type ClusterSlotsCmd = redis.ClusterSlotsCmd
type CommandsInfoCmd = redis.CommandsInfoCmd
type XStream = redis.XStream
type XMessage = redis.XMessage
type Z = redis.Z

type Script = redis.Script
type RingOptions = redis.RingOptions
type Options = redis.Options
type ClusterOptions = redis.ClusterOptions
type ClusterClient = redis.ClusterClient
