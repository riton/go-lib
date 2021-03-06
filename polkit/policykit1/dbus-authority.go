/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package policykit1

import "pkg.deepin.io/lib/dbus"
import "pkg.deepin.io/lib/dbus/property"
import "reflect"
import "sync"
import "runtime"
import "fmt"
import "errors"

/*prevent compile error*/
var _ = fmt.Println
var _ = runtime.SetFinalizer
var _ = sync.NewCond
var _ = reflect.TypeOf
var _ = property.BaseObserver{}

type Authority struct {
	Path     dbus.ObjectPath
	DestName string
	core     *dbus.Object

	signals       map[<-chan *dbus.Signal]struct{}
	signalsLocker sync.Mutex

	BackendName     *dbusPropertyAuthorityBackendName
	BackendVersion  *dbusPropertyAuthorityBackendVersion
	BackendFeatures *dbusPropertyAuthorityBackendFeatures
}

func (obj *Authority) _createSignalChan() <-chan *dbus.Signal {
	obj.signalsLocker.Lock()
	ch := getBus().Signal()
	obj.signals[ch] = struct{}{}
	obj.signalsLocker.Unlock()
	return ch
}
func (obj *Authority) _deleteSignalChan(ch <-chan *dbus.Signal) {
	obj.signalsLocker.Lock()
	delete(obj.signals, ch)
	getBus().DetachSignal(ch)
	obj.signalsLocker.Unlock()
}
func DestroyAuthority(obj *Authority) {
	obj.signalsLocker.Lock()
	defer obj.signalsLocker.Unlock()
	if obj.signals == nil {
		return
	}
	for ch, _ := range obj.signals {
		getBus().DetachSignal(ch)
	}
	obj.signals = nil

	runtime.SetFinalizer(obj, nil)

	dbusRemoveMatch("type='signal',path='" + string(obj.Path) + "',interface='org.freedesktop.DBus.Properties',sender='" + obj.DestName + "',member='PropertiesChanged'")
	dbusRemoveMatch("type='signal',path='" + string(obj.Path) + "',interface='org.freedesktop.PolicyKit1.Authority',sender='" + obj.DestName + "',member='PropertiesChanged'")

	dbusRemoveMatch("type='signal',path='" + string(obj.Path) + "',interface='org.freedesktop.PolicyKit1.Authority',sender='" + obj.DestName + "',member='Changed'")

	obj.BackendName.Reset()
	obj.BackendVersion.Reset()
	obj.BackendFeatures.Reset()
}

func (obj *Authority) EnumerateActions(locale string) (action_descriptions [][]interface{}, _err error) {
	_err = obj.core.Call("org.freedesktop.PolicyKit1.Authority.EnumerateActions", 0, locale).Store(&action_descriptions)
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *Authority) CheckAuthorization(subject interface{}, action_id string, details map[string]string, flags uint32, cancellation_id string) (result []interface{}, _err error) {
	_err = obj.core.Call("org.freedesktop.PolicyKit1.Authority.CheckAuthorization", 0, subject, action_id, details, flags, cancellation_id).Store(&result)
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *Authority) CancelCheckAuthorization(cancellation_id string) (_err error) {
	_err = obj.core.Call("org.freedesktop.PolicyKit1.Authority.CancelCheckAuthorization", 0, cancellation_id).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *Authority) RegisterAuthenticationAgent(subject interface{}, locale string, object_path string) (_err error) {
	_err = obj.core.Call("org.freedesktop.PolicyKit1.Authority.RegisterAuthenticationAgent", 0, subject, locale, object_path).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *Authority) RegisterAuthenticationAgentWithOptions(subject interface{}, locale string, object_path string, options map[string]dbus.Variant) (_err error) {
	_err = obj.core.Call("org.freedesktop.PolicyKit1.Authority.RegisterAuthenticationAgentWithOptions", 0, subject, locale, object_path, options).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *Authority) UnregisterAuthenticationAgent(subject interface{}, object_path string) (_err error) {
	_err = obj.core.Call("org.freedesktop.PolicyKit1.Authority.UnregisterAuthenticationAgent", 0, subject, object_path).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *Authority) AuthenticationAgentResponse(cookie string, identity interface{}) (_err error) {
	_err = obj.core.Call("org.freedesktop.PolicyKit1.Authority.AuthenticationAgentResponse", 0, cookie, identity).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *Authority) EnumerateTemporaryAuthorizations(subject interface{}) (temporary_authorizations [][]interface{}, _err error) {
	_err = obj.core.Call("org.freedesktop.PolicyKit1.Authority.EnumerateTemporaryAuthorizations", 0, subject).Store(&temporary_authorizations)
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *Authority) RevokeTemporaryAuthorizations(subject interface{}) (_err error) {
	_err = obj.core.Call("org.freedesktop.PolicyKit1.Authority.RevokeTemporaryAuthorizations", 0, subject).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *Authority) RevokeTemporaryAuthorizationById(id string) (_err error) {
	_err = obj.core.Call("org.freedesktop.PolicyKit1.Authority.RevokeTemporaryAuthorizationById", 0, id).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *Authority) ConnectChanged(callback func()) func() {
	sigChan := obj._createSignalChan()
	go func() {
		for v := range sigChan {
			if v.Path != obj.Path || v.Name != "org.freedesktop.PolicyKit1.Authority.Changed" || 0 != len(v.Body) {
				continue
			}

			callback()
		}
	}()
	return func() {
		obj._deleteSignalChan(sigChan)
	}
}

type dbusPropertyAuthorityBackendName struct {
	*property.BaseObserver
	core *dbus.Object
}

func (this *dbusPropertyAuthorityBackendName) SetValue(notwritable interface{}) {
	fmt.Println("org.freedesktop.PolicyKit1.Authority.BackendName is not writable")
}

func (this *dbusPropertyAuthorityBackendName) Get() string {
	v, _ := this.GetValue()
	return v.(string)
}
func (this *dbusPropertyAuthorityBackendName) GetValue() (interface{} /*string*/, error) {
	var r dbus.Variant
	err := this.core.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.PolicyKit1.Authority", "BackendName").Store(&r)
	if err == nil && r.Signature().String() == "s" {
		return r.Value().(string), nil
	}
	return *new(string), err
}
func (this *dbusPropertyAuthorityBackendName) GetType() reflect.Type {
	return reflect.TypeOf((*string)(nil)).Elem()
}

type dbusPropertyAuthorityBackendVersion struct {
	*property.BaseObserver
	core *dbus.Object
}

func (this *dbusPropertyAuthorityBackendVersion) SetValue(notwritable interface{}) {
	fmt.Println("org.freedesktop.PolicyKit1.Authority.BackendVersion is not writable")
}

func (this *dbusPropertyAuthorityBackendVersion) Get() string {
	v, _ := this.GetValue()
	return v.(string)
}
func (this *dbusPropertyAuthorityBackendVersion) GetValue() (interface{} /*string*/, error) {
	var r dbus.Variant
	err := this.core.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.PolicyKit1.Authority", "BackendVersion").Store(&r)
	if err == nil && r.Signature().String() == "s" {
		return r.Value().(string), nil
	}
	return *new(string), err
}
func (this *dbusPropertyAuthorityBackendVersion) GetType() reflect.Type {
	return reflect.TypeOf((*string)(nil)).Elem()
}

type dbusPropertyAuthorityBackendFeatures struct {
	*property.BaseObserver
	core *dbus.Object
}

func (this *dbusPropertyAuthorityBackendFeatures) SetValue(notwritable interface{}) {
	fmt.Println("org.freedesktop.PolicyKit1.Authority.BackendFeatures is not writable")
}

func (this *dbusPropertyAuthorityBackendFeatures) Get() uint32 {
	v, _ := this.GetValue()
	return v.(uint32)
}
func (this *dbusPropertyAuthorityBackendFeatures) GetValue() (interface{} /*uint32*/, error) {
	var r dbus.Variant
	err := this.core.Call("org.freedesktop.DBus.Properties.Get", 0, "org.freedesktop.PolicyKit1.Authority", "BackendFeatures").Store(&r)
	if err == nil && r.Signature().String() == "u" {
		return r.Value().(uint32), nil
	}
	return *new(uint32), err
}
func (this *dbusPropertyAuthorityBackendFeatures) GetType() reflect.Type {
	return reflect.TypeOf((*uint32)(nil)).Elem()
}

func NewAuthority(destName string, path dbus.ObjectPath) (*Authority, error) {
	if !path.IsValid() {
		return nil, errors.New("The path of '" + string(path) + "' is invalid.")
	}

	core := getBus().Object(destName, path)

	obj := &Authority{Path: path, DestName: destName, core: core, signals: make(map[<-chan *dbus.Signal]struct{})}

	obj.BackendName = &dbusPropertyAuthorityBackendName{&property.BaseObserver{}, core}
	obj.BackendVersion = &dbusPropertyAuthorityBackendVersion{&property.BaseObserver{}, core}
	obj.BackendFeatures = &dbusPropertyAuthorityBackendFeatures{&property.BaseObserver{}, core}

	dbusAddMatch("type='signal',path='" + string(path) + "',interface='org.freedesktop.DBus.Properties',sender='" + destName + "',member='PropertiesChanged'")
	dbusAddMatch("type='signal',path='" + string(path) + "',interface='org.freedesktop.PolicyKit1.Authority',sender='" + destName + "',member='PropertiesChanged'")
	sigChan := obj._createSignalChan()
	go func() {
		typeString := reflect.TypeOf("")
		typeKeyValues := reflect.TypeOf(map[string]dbus.Variant{})
		typeArrayValues := reflect.TypeOf([]string{})
		for v := range sigChan {
			if v.Name == "org.freedesktop.DBus.Properties.PropertiesChanged" &&
				len(v.Body) == 3 &&
				reflect.TypeOf(v.Body[0]) == typeString &&
				reflect.TypeOf(v.Body[1]) == typeKeyValues &&
				reflect.TypeOf(v.Body[2]) == typeArrayValues &&
				v.Body[0].(string) == "org.freedesktop.PolicyKit1.Authority" {
				props := v.Body[1].(map[string]dbus.Variant)
				for key, _ := range props {
					if false {
					} else if key == "BackendName" {
						obj.BackendName.Notify()

					} else if key == "BackendVersion" {
						obj.BackendVersion.Notify()

					} else if key == "BackendFeatures" {
						obj.BackendFeatures.Notify()
					}
				}
			} else if v.Name == "org.freedesktop.PolicyKit1.Authority.PropertiesChanged" && len(v.Body) == 1 && reflect.TypeOf(v.Body[0]) == typeKeyValues {
				for key, _ := range v.Body[0].(map[string]dbus.Variant) {
					if false {
					} else if key == "BackendName" {
						obj.BackendName.Notify()

					} else if key == "BackendVersion" {
						obj.BackendVersion.Notify()

					} else if key == "BackendFeatures" {
						obj.BackendFeatures.Notify()
					}
				}
			}
		}
	}()

	dbusAddMatch("type='signal',path='" + string(obj.Path) + "',interface='org.freedesktop.PolicyKit1.Authority',sender='" + obj.DestName + "',member='Changed'")

	runtime.SetFinalizer(obj, func(_obj *Authority) { DestroyAuthority(_obj) })
	return obj, nil
}
