/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

package sound

// #cgo pkg-config: glib-2.0 libcanberra
// #include <stdlib.h>
// #include "player.h"
import "C"

import (
	"fmt"
	"os/exec"
	"unsafe"
)

const (
	PlayerLibcanberra = "libcanberra"
	PlayerMPV         = "mpv"
	PlayerPaplay      = "paplay"
)

// PlayThemeSound: play sound theme event
// @theme : the special sound theme
// @event : the special theme event id
// @device: the special backend card, default: the current sound card
// @driver: the special backend driver, as: 'pulse','alsa','null' ...
//          default: 'pulse'
// @player: the special sound play method, as: 'libcanberra','mpv','paplay'
//          if empty or other value, just as 'libcanberra'
//
// return: the error message
func PlayThemeSound(theme, event, device, driver, player string) error {
	switch player {
	case PlayerPaplay, PlayerMPV:
		file, err := findThemeFile(theme, event)
		if err != nil {
			return err
		}
		_, err = exec.Command(player, []string{file}...).Output()
		return err
	}

	cTheme := C.CString(theme)
	defer C.free(unsafe.Pointer(cTheme))
	cEvent := C.CString(event)
	defer C.free(unsafe.Pointer(cEvent))
	cDevice := C.CString(device)
	defer C.free(unsafe.Pointer(cDevice))
	cDriver := C.CString(driver)
	defer C.free(unsafe.Pointer(cDriver))

	ret := C.canberra_play_system_sound(cTheme, cEvent,
		cDevice, cDriver)
	if ret != 0 {
		msg := C.GoString(C.ca_strerror(ret))
		return fmt.Errorf("Play failed: %s", msg)
	}
	return nil
}

// PlaySoundFile: play sound file
// @file  : the file which needed to play
// @device: the special backend card, default: the current sound card
// @driver: the special backend driver, as: 'pulse','alsa','null' ...
//          default: 'pulse'
// @player: the special sound play method, as: 'libcanberra','mpv','paplay'
//          if empty or other value, just as 'libcanberra'
//
// return: the error message
func PlaySoundFile(file, device, driver, player string) error {
	switch player {
	case PlayerPaplay, PlayerMPV:
		_, err := exec.Command(player, []string{file}...).Output()
		return err
	}

	cFile := C.CString(file)
	defer C.free(unsafe.Pointer(cFile))
	cDevice := C.CString(device)
	defer C.free(unsafe.Pointer(cDevice))
	cDriver := C.CString(driver)
	defer C.free(unsafe.Pointer(cDriver))

	ret := C.canberra_play_sound_file(cFile, cDevice, cDriver)
	if ret != 0 {
		msg := C.GoString(C.ca_strerror(ret))
		return fmt.Errorf("Play failed: %s", msg)
	}
	return nil
}
