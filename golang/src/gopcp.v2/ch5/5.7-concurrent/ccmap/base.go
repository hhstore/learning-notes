package ccmap

const (
	DEFAULT_BUCKET_LOAD_FACTOR float64 = 0.75 // 装载因子
	DEFAULT_BUCKET_NUMBER      int     = 16   // 单个散列桶默认数量
	DEFAULT_BUCKET_MAX_SIZE    uint64  = 1000 // 单个散列桶默认最大尺寸
)

const (
	MAX_CONCURRENCY int = 65536 // 最大并发量
)
