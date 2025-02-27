package driver

import (
	"encoding/binary"
	"encoding/json"
	"sync"
	"time"

	"github.com/i4de/rulex/common"
	"github.com/i4de/rulex/glogger"
	"github.com/i4de/rulex/typex"
	modbus "github.com/wwhai/gomodbus"
)

/*
*
* Modbus RTU
*
 */
type modBusRtuDriver struct {
	state      typex.DriverState
	handler    *modbus.RTUClientHandler
	client     modbus.Client
	RuleEngine typex.RuleX
	Registers  []common.RegisterRW
	device     *typex.Device
	lock       sync.Mutex
}

func NewModBusRtuDriver(
	d *typex.Device,
	e typex.RuleX,
	Registers []common.RegisterRW,
	handler *modbus.RTUClientHandler,
	client modbus.Client) typex.XExternalDriver {
	return &modBusRtuDriver{
		state:      typex.DRIVER_UP,
		device:     d,
		RuleEngine: e,
		client:     client,
		handler:    handler,
		Registers:  Registers,
		lock:       sync.Mutex{},
	}

}
func (d *modBusRtuDriver) Test() error {
	return nil
}

func (d *modBusRtuDriver) Init(map[string]string) error {
	return nil
}

func (d *modBusRtuDriver) Work() error {
	return nil
}

func (d *modBusRtuDriver) State() typex.DriverState {
	return d.state
}

func (d *modBusRtuDriver) Read(cmd []byte, data []byte) (int, error) {
	dataMap := map[string]common.RegisterRW{}
	for _, r := range d.Registers {
		d.handler.SlaveId = r.SlaverId
		if r.Function == common.READ_COIL {
			d.lock.Lock()
			results, err := d.client.ReadCoils(r.Address, r.Quantity)
			d.lock.Unlock()
			if err != nil {
				glogger.GLogger.Error(err)
			}
			value := common.RegisterRW{
				Tag:      r.Tag,
				Function: r.Function,
				SlaverId: r.SlaverId,
				Address:  r.Address,
				Quantity: r.Quantity,
				Value:    covertEmptyHex(results),
			}
			dataMap[r.Tag] = value
		}
		if r.Function == common.READ_DISCRETE_INPUT {
			d.lock.Lock()
			results, err := d.client.ReadDiscreteInputs(r.Address, r.Quantity)
			d.lock.Unlock()
			if err != nil {
				glogger.GLogger.Error(err)
			}
			value := common.RegisterRW{
				Tag:      r.Tag,
				Function: r.Function,
				SlaverId: r.SlaverId,
				Address:  r.Address,
				Quantity: r.Quantity,
				Value:    covertEmptyHex(results),
			}
			dataMap[r.Tag] = value

		}
		if r.Function == common.READ_HOLDING_REGISTERS {
			d.lock.Lock()
			results, err := d.client.ReadHoldingRegisters(r.Address, r.Quantity)
			d.lock.Unlock()
			if err != nil {
				glogger.GLogger.Error(err)
			}
			value := common.RegisterRW{
				Tag:      r.Tag,
				Function: r.Function,
				SlaverId: r.SlaverId,
				Address:  r.Address,
				Quantity: r.Quantity,
				Value:    covertEmptyHex(results),
			}
			dataMap[r.Tag] = value
		}
		if r.Function == common.READ_INPUT_REGISTERS {
			d.lock.Lock()
			results, err := d.client.ReadInputRegisters(r.Address, r.Quantity)
			d.lock.Unlock()
			if err != nil {
				glogger.GLogger.Error(err)
			}
			value := common.RegisterRW{
				Tag:      r.Tag,
				Function: r.Function,
				SlaverId: r.SlaverId,
				Address:  r.Address,
				Quantity: r.Quantity,
				Value:    covertEmptyHex(results),
			}
			dataMap[r.Tag] = value
		}

		// 设置一个间隔时间防止低级CPU黏包等
		// TODO 未来通过参数形式传递
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
	bytes, _ := json.Marshal(dataMap)
	copy(data, bytes)
	return len(bytes), nil

}

func (d *modBusRtuDriver) Write(cmd []byte, data []byte) (int, error) {
	dataMap := []common.RegisterRW{}
	if err := json.Unmarshal(data, &dataMap); err != nil {
		return 0, err
	}
	for _, r := range dataMap {
		if r.Function == common.WRITE_SINGLE_COIL {
			d.lock.Lock()
			_, err := d.client.WriteSingleCoil(r.Address, binary.BigEndian.Uint16([]byte(r.Value)[0:2]))
			d.lock.Unlock()
			if err != nil {
				return 0, err
			}
		}
		if r.Function == common.WRITE_MULTIPLE_COILS {
			d.lock.Lock()
			_, err := d.client.WriteMultipleCoils(r.Address, r.Quantity, []byte(r.Value))
			d.lock.Unlock()
			if err != nil {
				return 0, err
			}
		}
		if r.Function == common.WRITE_SINGLE_HOLDING_REGISTER {
			d.lock.Lock()
			_, err := d.client.WriteSingleRegister(r.Address, binary.BigEndian.Uint16([]byte(r.Value)[0:2]))
			d.lock.Unlock()
			if err != nil {
				return 0, err
			}
		}
		if r.Function == common.WRITE_MULTIPLE_HOLDING_REGISTERS {
			d.lock.Lock()
			_, err := d.client.WriteMultipleRegisters(r.Address, r.Quantity, []byte(r.Value))
			d.lock.Unlock()
			if err != nil {
				return 0, err
			}
		}
	}
	return 0, nil
}

func (d *modBusRtuDriver) DriverDetail() typex.DriverDetail {
	return typex.DriverDetail{
		Name:        "Generic ModBus RTU Driver",
		Type:        "UART",
		Description: "Generic ModBus RTU Driver",
	}
}

func (d *modBusRtuDriver) Stop() error {
	d.handler.Close()
	d = nil
	return nil
}
