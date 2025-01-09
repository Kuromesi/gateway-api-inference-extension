// Package scheduling implements request scheduling algorithms.
package scheduling

import (
	"fmt"
	"math/rand"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"inference.networking.x-k8s.io/gateway-api-inference-extension/pkg/ext-proc/backend"
	klog "k8s.io/klog/v2"
)

const (
	// TODO(https://github.com/kubernetes-sigs/gateway-api-inference-extension/issues/16) Make this configurable.
	defaultKvCacheThreshold = 0.8
	// TODO(https://github.com/kubernetes-sigs/gateway-api-inference-extension/issues/16) Make this configurable.
	defaultQueueThresholdCritical = 5
	// TODO(https://github.com/kubernetes-sigs/gateway-api-inference-extension/issues/16) Make this configurable.
	// the threshold for queued requests to be considered low below which we can prioritize LoRA affinity.
	// The value of 50 is arrived heuristicically based on experiments.
	defaultQueueingThresholdLoRA = 50
)

var (
	defaultFilter = &filterChainImpl{
		name:          "critical request",
		filter:        toFilter(criticalRequestPredicate),
		nextOnSuccess: lowLatencyFilter,
		nextOnFailure: sheddableRequestFilter,
	}

	// queueLoRAAndKVCacheFilter applied least queue -> low cost lora ->  least KV Cache filter
	queueLoRAAndKVCacheFilter = &filterChainImpl{
		name:   "least queuing",
		filter: leastQueuingFilterFunc,
		nextOnSuccessOrFailure: &filterChainImpl{
			name:   "low cost LoRA",
			filter: toFilter(lowLoRACostPredicate),
			nextOnSuccessOrFailure: &filterChainImpl{
				name:   "least KV cache percent",
				filter: leastKVCacheFilterFunc,
			},
		},
	}

	// queueAndKVCacheFilter applies least queue followed by least KV Cache filter
	queueAndKVCacheFilter = &filterChainImpl{
		name:   "least queuing",
		filter: leastQueuingFilterFunc,
		nextOnSuccessOrFailure: &filterChainImpl{
			name:   "least KV cache percent",
			filter: leastKVCacheFilterFunc,
		},
	}

	lowLatencyFilter = &filterChainImpl{
		name:   "low queueing filter",
		filter: toFilter((lowQueueingPodPredicate(defaultQueueingThresholdLoRA))),
		nextOnSuccess: &filterChainImpl{
			name:          "affinity LoRA",
			filter:        toFilter(loRAAffinityPredicate),
			nextOnSuccess: queueAndKVCacheFilter,
			nextOnFailure: &filterChainImpl{
				name:                   "can accept LoRA Adapter",
				filter:                 toFilter(canAcceptNewLoraPredicate),
				nextOnSuccessOrFailure: queueAndKVCacheFilter,
			},
		},
		nextOnFailure: queueLoRAAndKVCacheFilter,
	}

	sheddableRequestFilter = &filterChainImpl{
		// When there is at least one model server that's not queuing requests, and still has KV
		// cache below a certain threshold, we consider this model server has capacity to handle
		// a sheddable request without impacting critical requests.
		name:          "has capacity for sheddable requests",
		filter:        toFilter(noQueueAndLessThanKVCacheThresholdPredicate(defaultQueueThresholdCritical, defaultKvCacheThreshold)),
		nextOnSuccess: queueLoRAAndKVCacheFilter,
		// If all pods are queuing or running above the KVCache threshold, we drop the sheddable
		// request to make room for critical requests.
		nextOnFailure: &filterChainImpl{
			name: "drop request",
			filter: func(req *LLMRequest, pods []*backend.PodMetrics) ([]*backend.PodMetrics, error) {
				klog.Infof("Dropping request %v", req)
				return []*backend.PodMetrics{}, status.Errorf(
					codes.ResourceExhausted, "dropping request due to limited backend resources")
			},
		},
	}
)

func NewScheduler(pmp PodMetricsProvider, opts ...SchedulerOption) *Scheduler {
	s := &Scheduler{
		podMetricsProvider: pmp,
		filter:             defaultFilter,
	}

	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithOrchestrator(orchestrator FilterOrchestrator) SchedulerOption {
	return func(s *Scheduler) {
		s.filterOrchestrator = orchestrator
	}
}

type SchedulerOption func(*Scheduler)

type Scheduler struct {
	podMetricsProvider PodMetricsProvider
	filter             FilterChain
	filterOrchestrator FilterOrchestrator
}

// PodMetricsProvider is an interface to provide set of pods in the backend and information such as
// metrics.
type PodMetricsProvider interface {
	AllPodMetrics() []*backend.PodMetrics
}

// Schedule finds the target pod based on metrics and the requested lora adapter.
func (s *Scheduler) Schedule(req *LLMRequest) (targetPod backend.Pod, err error) {
	klog.V(3).Infof("request: %v; metrics: %+v", req, s.podMetricsProvider.AllPodMetrics())
	pods, err := s.filterOrchestrator.Orchestrate().Filter(req, s.podMetricsProvider.AllPodMetrics())
	if err != nil || len(pods) == 0 {
		return backend.Pod{}, fmt.Errorf(
			"failed to apply filter, resulted %v pods, this should never happen: %w", len(pods), err)
	}
	klog.V(3).Infof("Going to randomly select a pod from the candidates: %+v", pods)
	i := rand.Intn(len(pods))
	return pods[i].Pod, nil
}
