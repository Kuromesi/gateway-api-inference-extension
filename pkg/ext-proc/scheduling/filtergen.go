package scheduling

const (
	FilterCriticalRequestName  = "critical_request"
	FilterLeastQueuingName     = "least_queuing"
	FilterLowCostLoraName      = "low_cost_lora"
	FilterLowLatencyName       = "low_latency"
	FilterAffinityLoraName     = "affinity_lora"
	FilterSheddableRequestName = "sheddable_request"
	FilterLeastKvCacheName     = "least_kv_cache"
	FilterDropRequestName      = "drop_request"
	FilterCanAcceptNewLoraName = "can_accept_new_lora"
)

const (
	TopKByWaitingQueueSize    = "waiting_queue_size"
	TopKByKVCacheUsagePercent = "kv_cache_usage_percent"
)

var filterMap = map[string]FilterGen{
	FilterLowLatencyName:       FilterLowLatency,
	FilterCriticalRequestName:  FilterCriticalRequest,
	FilterLeastQueuingName:     FilterLeastQueuing,
	FilterCanAcceptNewLoraName: FilterCanAcceptNewLora,
	FilterSheddableRequestName: FilterSheddableRequest,
	FilterDropRequestName:      FilterDropRequest,
	FilterAffinityLoraName:     FilterAffinityLora,
	FilterLowCostLoraName:      FilterLowCostLora,
	FilterLeastKvCacheName:     FilterLeastKvCache,
}

// FilterGen generate a filter from a filter option
type FilterGen interface {
	Name() string
	Get(*FilterOption) filter
	Validate(*FilterOption) error
}

type FilterOption struct {
	KvCacheThreshold *float64 `json:"kvCacheThreshold,omitempty"`

	QueueThresholdCritical *int `json:"queueThresholdCritical,omitempty"`
	QueueingThresholdLoRA  *int `json:"queueingThresholdLoRA,omitempty"`
}

type filterGenImpl struct {
	name      string
	getter    func(*FilterOption) filter
	validator func(*FilterOption) error
}

var _ FilterGen = &filterGenImpl{}

func (fg *filterGenImpl) Name() string {
	return fg.name
}

func (fg *filterGenImpl) Get(fo *FilterOption) filter {
	return fg.getter(fo)
}

func (fg *filterGenImpl) Validate(fo *FilterOption) error {
	return fg.validator(fo)
}

var (
	FilterCriticalRequest FilterGen = &filterGenImpl{
		name: FilterCriticalRequestName,
		getter: func(fo *FilterOption) filter {
			return toFilter(criticalRequestPredicate)
		},
		validator: func(fo *FilterOption) error { return nil },
	}

	FilterLeastQueuing FilterGen = &filterGenImpl{
		name: FilterLeastQueuingName,
		getter: func(fo *FilterOption) filter {
			return leastQueuingFilterFunc
		},
		validator: func(fo *FilterOption) error { return nil },
	}

	FilterLowCostLora FilterGen = &filterGenImpl{
		name: FilterLowCostLoraName,
		getter: func(fo *FilterOption) filter {
			return toFilter(lowLoRACostPredicate)
		},
		validator: func(fo *FilterOption) error { return nil },
	}

	FilterLowLatency FilterGen = &filterGenImpl{
		name: FilterLowLatencyName,
		getter: func(fo *FilterOption) filter {
			return toFilter(lowQueueingPodPredicate)
		},
		validator: func(fo *FilterOption) error { return nil },
	}

	FilterAffinityLora FilterGen = &filterGenImpl{
		name: FilterAffinityLoraName,
		getter: func(fo *FilterOption) filter {
			return toFilter(loRAAffinityPredicate)
		},
		validator: func(fo *FilterOption) error { return nil },
	}

	FilterSheddableRequest FilterGen = &filterGenImpl{
		name: FilterSheddableRequestName,
		getter: func(opt *FilterOption) filter {
			qtc, kct := queueThresholdCritical, kvCacheThreshold
			if opt != nil {
				if opt.KvCacheThreshold != nil {
					kct = *opt.KvCacheThreshold
				}
				if opt.QueueThresholdCritical != nil {
					qtc = *opt.QueueThresholdCritical
				}
			}
			return toFilter(noQueueAndLessThanKVCacheThresholdPredicate(qtc, kct))
		},
		validator: func(fo *FilterOption) error { return nil },
	}

	FilterLeastKvCache FilterGen = &filterGenImpl{
		name: FilterLeastKvCacheName,
		getter: func(fo *FilterOption) filter {
			return leastKVCacheFilterFunc
		},
		validator: func(fo *FilterOption) error { return nil },
	}

	FilterDropRequest FilterGen = &filterGenImpl{
		name: FilterDropRequestName,
		getter: func(fo *FilterOption) filter {
			return dropRequestFilterFunc
		},
		validator: func(fo *FilterOption) error { return nil },
	}

	FilterCanAcceptNewLora FilterGen = &filterGenImpl{
		name: FilterCanAcceptNewLoraName,
		getter: func(fo *FilterOption) filter {
			return toFilter(canAcceptNewLoraPredicate)
		},
		validator: func(fo *FilterOption) error { return nil },
	}
)
