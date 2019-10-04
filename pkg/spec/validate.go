package spec

import (
	"fmt"

	"github.com/deviceplane/deviceplane/pkg/validation"
	"gopkg.in/yaml.v2"
)

var (
	validators = map[string][]func(interface{}) error{
		"cap_add":          []func(interface{}) error{validation.ValidateStringArray},
		"cap_drop":         []func(interface{}) error{validation.ValidateStringArray},
		"command":          []func(interface{}) error{validation.ValidateStringOrStringArray},
		"cpuset":           []func(interface{}) error{validation.ValidateString},
		"cpu_shares":       []func(interface{}) error{validation.ValidateStringOrInteger},
		"cpu_quota":        []func(interface{}) error{validation.ValidateStringOrInteger},
		"dns":              []func(interface{}) error{validation.ValidateStringOrStringArray},
		"dns_opt":          []func(interface{}) error{validation.ValidateStringOrStringArray},
		"dns_search":       []func(interface{}) error{validation.ValidateStringOrStringArray},
		"domainname":       []func(interface{}) error{validation.ValidateString},
		"entrypoint":       []func(interface{}) error{validation.ValidateStringOrStringArray},
		"environment":      []func(interface{}) error{validation.ValidateArrayOrObject},
		"extra_hosts":      []func(interface{}) error{validation.ValidateArrayOrObject},
		"group_add":        []func(interface{}) error{validation.ValidateStringIntegerArray},
		"image":            []func(interface{}) error{validation.ValidateString},
		"hostname":         []func(interface{}) error{validation.ValidateString},
		"ipc":              []func(interface{}) error{validation.ValidateString},
		"labels":           []func(interface{}) error{validation.ValidateArrayOrObject},
		"mem_limit":        []func(interface{}) error{validation.ValidateStringOrInteger},
		"mem_reservation":  []func(interface{}) error{validation.ValidateStringOrInteger},
		"memswap_limit":    []func(interface{}) error{validation.ValidateStringOrInteger},
		"network_mode":     []func(interface{}) error{validation.ValidateString},
		"oom_kill_disable": []func(interface{}) error{validation.ValidateBoolean},
		"oom_score_adj":    []func(interface{}) error{validation.ValidateInteger},
		"pid":              []func(interface{}) error{validation.ValidateString},
		"ports":            []func(interface{}) error{validation.ValidateStringIntegerArray},
		"privileged":       []func(interface{}) error{validation.ValidateBoolean},
		"read_only":        []func(interface{}) error{validation.ValidateBoolean},
		"restart":          []func(interface{}) error{validation.ValidateString},
		"security_opt":     []func(interface{}) error{validation.ValidateStringArray},
		"shm_size":         []func(interface{}) error{validation.ValidateStringOrInteger},
		"stop_signal":      []func(interface{}) error{validation.ValidateString},
		"user":             []func(interface{}) error{validation.ValidateString},
		"uts":              []func(interface{}) error{validation.ValidateString},
		"volumes":          []func(interface{}) error{validation.ValidateStringArray},
		"working_dir":      []func(interface{}) error{validation.ValidateString},
	}
)

func Validate(c []byte) error {
	var m map[string]interface{}
	if err := yaml.Unmarshal(c, &m); err != nil {
		return err
	}

	for serviceName := range m {
		if len(serviceName) > 100 {
			return fmt.Errorf("service name '%s' is longer than 100 characters", serviceName)
		}
	}

	for serviceName, service := range m {
		service, ok := service.(map[interface{}]interface{})
		if !ok {
			return fmt.Errorf("service '%s' is not an object", serviceName)
		}

		for key := range service {
			typedKey, ok := key.(string)
			if !ok {
				return fmt.Errorf("service '%s': invalid key '%v'", serviceName, key)
			}
			if _, ok = validators[typedKey]; !ok {
				return fmt.Errorf("service '%s': invalid key '%s'", serviceName, typedKey)
			}
		}

		for key, validators := range validators {
			value, ok := service[key]
			if !ok {
				continue
			}
			for _, validator := range validators {
				if err := validator(value); err != nil {
					return fmt.Errorf("service '%s', key '%s': %v", serviceName, key, err)
				}
			}
		}
	}

	return nil
}
