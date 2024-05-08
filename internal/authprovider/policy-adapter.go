package authprovider

import (
	"encoding/csv"
	"errors"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"os"
)

type casbinFileAdapter struct {
	policyFile    *os.File
	casbinAdapter persist.Adapter
}

func newCasbinFileAdapter(policyPath string, casbinAdapter persist.Adapter) (persist.Adapter, error) {
	file, err := os.OpenFile(policyPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	if _, err := file.Seek(0, os.SEEK_END); err != nil {
		return nil, err
	}

	return &casbinFileAdapter{
		policyFile:    file,
		casbinAdapter: casbinAdapter,
	}, nil
}

func (c casbinFileAdapter) LoadPolicy(model model.Model) error {
	return c.casbinAdapter.LoadPolicy(model)
}

func (c casbinFileAdapter) SavePolicy(model model.Model) error {
	return c.casbinAdapter.SavePolicy(model)
}

// AddPolicy an addon over the main method out of the box because the main
// method does not support sec g2.
func (c casbinFileAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	// support for these policies is perfectly
	// supported by the built-in adapter
	if rule[0] == "p" || rule[0] == "g" {
		return c.casbinAdapter.AddPolicy(sec, ptype, rule)
	} else if rule[0] != "g2" {
		return errors.New("unknown policy type: p, g, g2")
	}

	if len(rule) != 3 {
		return errors.New("invalid policy format: g2, role, resource")
	}

	writer := csv.NewWriter(c.policyFile)
	writer.UseCRLF = true
	defer writer.Flush()

	// this is done so that the new line being written does
	// not merge with previous one
	if err := writer.Write([]string{""}); err != nil {
		return err
	}

	if err := writer.Write(rule); err != nil {
		return err
	}

	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}

func (c casbinFileAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return c.casbinAdapter.RemovePolicy(sec, ptype, rule)
}

func (c casbinFileAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return c.casbinAdapter.RemoveFilteredPolicy(sec, ptype, fieldIndex, fieldValues...)
}
