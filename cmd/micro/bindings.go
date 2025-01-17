package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"

	"github.com/flynn/json5"
	"github.com/zyedidia/tcell"
)

var bindingsStr map[string]string
var bindings map[Key][]func(*View, bool) bool
var mouseBindings map[Key][]func(*View, bool, *tcell.EventMouse) bool
var helpBinding string
var kmenuBinding string

var mouseBindingActions = map[string]func(*View, bool, *tcell.EventMouse) bool{
	"MousePress":       (*View).MousePress,
	"MouseMultiCursor": (*View).MouseMultiCursor,
}

var bindingActions = map[string]func(*View, bool) bool{
	"CursorUp":               (*View).CursorUp,
	"CursorDown":             (*View).CursorDown,
	"CursorPageUp":           (*View).CursorPageUp,
	"CursorPageDown":         (*View).CursorPageDown,
	"CursorLeft":             (*View).CursorLeft,
	"CursorRight":            (*View).CursorRight,
	"CursorStart":            (*View).CursorStart,
	"CursorEnd":              (*View).CursorEnd,
	"SelectToStart":          (*View).SelectToStart,
	"SelectToEnd":            (*View).SelectToEnd,
	"SelectUp":               (*View).SelectUp,
	"SelectDown":             (*View).SelectDown,
	"SelectLeft":             (*View).SelectLeft,
	"SelectRight":            (*View).SelectRight,
	"WordRight":              (*View).WordRight,
	"WordLeft":               (*View).WordLeft,
	"SelectWordRight":        (*View).SelectWordRight,
	"SelectWordLeft":         (*View).SelectWordLeft,
	"DeleteWordRight":        (*View).DeleteWordRight,
	"DeleteWordLeft":         (*View).DeleteWordLeft,
	"SelectLine":             (*View).SelectLine,
	"SelectToStartOfLine":    (*View).SelectToStartOfLine,
	"SelectToEndOfLine":      (*View).SelectToEndOfLine,
	"ParagraphPrevious":      (*View).ParagraphPrevious,
	"ParagraphNext":          (*View).ParagraphNext,
	"InsertNewline":          (*View).InsertNewline,
	"InsertSpace":            (*View).InsertSpace,
	"Backspace":              (*View).Backspace,
	"Delete":                 (*View).Delete,
	"InsertTab":              (*View).InsertTab,
	"Save":                   (*View).Save,
	"SaveAll":                (*View).SaveAll,
	"SaveAs":                 (*View).SaveAs,
	"Find":                   (*View).Find,
	"FindNext":               (*View).FindNext,
	"FindPrevious":           (*View).FindPrevious,
	"Center":                 (*View).Center,
	"Undo":                   (*View).Undo,
	"Redo":                   (*View).Redo,
	"Copy":                   (*View).Copy,
	"Cut":                    (*View).Cut,
	"CutLine":                (*View).CutLine,
	"DuplicateLine":          (*View).DuplicateLine,
	"DeleteLine":             (*View).DeleteLine,
	"MoveLinesUp":            (*View).MoveLinesUp,
	"MoveLinesDown":          (*View).MoveLinesDown,
	"IndentSelection":        (*View).IndentSelection,
	"OutdentSelection":       (*View).OutdentSelection,
	"OutdentLine":            (*View).OutdentLine,
	"Paste":                  (*View).Paste,
	"PastePrimary":           (*View).PastePrimary,
	"SelectAll":              (*View).SelectAll,
	"OpenFile":               (*View).OpenFile,
	"Start":                  (*View).Start,
	"End":                    (*View).End,
	"PageUp":                 (*View).PageUp,
	"PageDown":               (*View).PageDown,
	"SelectPageUp":           (*View).SelectPageUp,
	"SelectPageDown":         (*View).SelectPageDown,
	"HalfPageUp":             (*View).HalfPageUp,
	"HalfPageDown":           (*View).HalfPageDown,
	"StartOfLine":            (*View).StartOfLine,
	"EndOfLine":              (*View).EndOfLine,
	"ToggleHelp":             (*View).ToggleHelp,
	"ToggleKeyMenu":          (*View).ToggleKeyMenu,
	"ToggleRuler":            (*View).ToggleRuler,
	"JumpLine":               (*View).JumpLine,
	"ClearStatus":            (*View).ClearStatus,
	"ShellMode":              (*View).ShellMode,
	"CommandMode":            (*View).CommandMode,
	"ToggleOverwriteMode":    (*View).ToggleOverwriteMode,
	"Escape":                 (*View).Escape,
	"Quit":                   (*View).Quit,
	"QuitAll":                (*View).QuitAll,
	"AddTab":                 (*View).AddTab,
	"PreviousTab":            (*View).PreviousTab,
	"NextTab":                (*View).NextTab,
	"NextSplit":              (*View).NextSplit,
	"PreviousSplit":          (*View).PreviousSplit,
	"Unsplit":                (*View).Unsplit,
	"VSplit":                 (*View).VSplitBinding,
	"HSplit":                 (*View).HSplitBinding,
	"ToggleMacro":            (*View).ToggleMacro,
	"PlayMacro":              (*View).PlayMacro,
	"Suspend":                (*View).Suspend,
	"ScrollUp":               (*View).ScrollUpAction,
	"ScrollDown":             (*View).ScrollDownAction,
	"SpawnMultiCursor":       (*View).SpawnMultiCursor,
	"SpawnMultiCursorSelect": (*View).SpawnMultiCursorSelect,
	"RemoveMultiCursor":      (*View).RemoveMultiCursor,
	"RemoveAllMultiCursors":  (*View).RemoveAllMultiCursors,
	"SkipMultiCursor":        (*View).SkipMultiCursor,
	"JumpToMatchingBrace":    (*View).JumpToMatchingBrace,

	// This was changed to InsertNewline but I don't want to break backwards compatibility
	"InsertEnter": (*View).InsertNewline,
}

var bindingMouse = map[string]tcell.ButtonMask{
	"MouseLeft":       tcell.Button1,
	"MouseMiddle":     tcell.Button2,
	"MouseRight":      tcell.Button3,
	"MouseWheelUp":    tcell.WheelUp,
	"MouseWheelDown":  tcell.WheelDown,
	"MouseWheelLeft":  tcell.WheelLeft,
	"MouseWheelRight": tcell.WheelRight,
}

var bindingKeys = map[string]tcell.Key{
	"Up":             tcell.KeyUp,
	"Down":           tcell.KeyDown,
	"Right":          tcell.KeyRight,
	"Left":           tcell.KeyLeft,
	"UpLeft":         tcell.KeyUpLeft,
	"UpRight":        tcell.KeyUpRight,
	"DownLeft":       tcell.KeyDownLeft,
	"DownRight":      tcell.KeyDownRight,
	"Center":         tcell.KeyCenter,
	"PageUp":         tcell.KeyPgUp,
	"PageDown":       tcell.KeyPgDn,
	"Home":           tcell.KeyHome,
	"End":            tcell.KeyEnd,
	"Insert":         tcell.KeyInsert,
	"Delete":         tcell.KeyDelete,
	"Help":           tcell.KeyHelp,
	"Exit":           tcell.KeyExit,
	"Clear":          tcell.KeyClear,
	"Cancel":         tcell.KeyCancel,
	"Print":          tcell.KeyPrint,
	"Pause":          tcell.KeyPause,
	"Backtab":        tcell.KeyBacktab,
	"F1":             tcell.KeyF1,
	"F2":             tcell.KeyF2,
	"F3":             tcell.KeyF3,
	"F4":             tcell.KeyF4,
	"F5":             tcell.KeyF5,
	"F6":             tcell.KeyF6,
	"F7":             tcell.KeyF7,
	"F8":             tcell.KeyF8,
	"F9":             tcell.KeyF9,
	"F10":            tcell.KeyF10,
	"F11":            tcell.KeyF11,
	"F12":            tcell.KeyF12,
	"F13":            tcell.KeyF13,
	"F14":            tcell.KeyF14,
	"F15":            tcell.KeyF15,
	"F16":            tcell.KeyF16,
	"F17":            tcell.KeyF17,
	"F18":            tcell.KeyF18,
	"F19":            tcell.KeyF19,
	"F20":            tcell.KeyF20,
	"F21":            tcell.KeyF21,
	"F22":            tcell.KeyF22,
	"F23":            tcell.KeyF23,
	"F24":            tcell.KeyF24,
	"F25":            tcell.KeyF25,
	"F26":            tcell.KeyF26,
	"F27":            tcell.KeyF27,
	"F28":            tcell.KeyF28,
	"F29":            tcell.KeyF29,
	"F30":            tcell.KeyF30,
	"F31":            tcell.KeyF31,
	"F32":            tcell.KeyF32,
	"F33":            tcell.KeyF33,
	"F34":            tcell.KeyF34,
	"F35":            tcell.KeyF35,
	"F36":            tcell.KeyF36,
	"F37":            tcell.KeyF37,
	"F38":            tcell.KeyF38,
	"F39":            tcell.KeyF39,
	"F40":            tcell.KeyF40,
	"F41":            tcell.KeyF41,
	"F42":            tcell.KeyF42,
	"F43":            tcell.KeyF43,
	"F44":            tcell.KeyF44,
	"F45":            tcell.KeyF45,
	"F46":            tcell.KeyF46,
	"F47":            tcell.KeyF47,
	"F48":            tcell.KeyF48,
	"F49":            tcell.KeyF49,
	"F50":            tcell.KeyF50,
	"F51":            tcell.KeyF51,
	"F52":            tcell.KeyF52,
	"F53":            tcell.KeyF53,
	"F54":            tcell.KeyF54,
	"F55":            tcell.KeyF55,
	"F56":            tcell.KeyF56,
	"F57":            tcell.KeyF57,
	"F58":            tcell.KeyF58,
	"F59":            tcell.KeyF59,
	"F60":            tcell.KeyF60,
	"F61":            tcell.KeyF61,
	"F62":            tcell.KeyF62,
	"F63":            tcell.KeyF63,
	"F64":            tcell.KeyF64,
	"CtrlSpace":      tcell.KeyCtrlSpace,
	"CtrlA":          tcell.KeyCtrlA,
	"CtrlB":          tcell.KeyCtrlB,
	"CtrlC":          tcell.KeyCtrlC,
	"CtrlD":          tcell.KeyCtrlD,
	"CtrlE":          tcell.KeyCtrlE,
	"CtrlF":          tcell.KeyCtrlF,
	"CtrlG":          tcell.KeyCtrlG,
	"CtrlH":          tcell.KeyCtrlH,
	"CtrlI":          tcell.KeyCtrlI,
	"CtrlJ":          tcell.KeyCtrlJ,
	"CtrlK":          tcell.KeyCtrlK,
	"CtrlL":          tcell.KeyCtrlL,
	"CtrlM":          tcell.KeyCtrlM,
	"CtrlN":          tcell.KeyCtrlN,
	"CtrlO":          tcell.KeyCtrlO,
	"CtrlP":          tcell.KeyCtrlP,
	"CtrlQ":          tcell.KeyCtrlQ,
	"CtrlR":          tcell.KeyCtrlR,
	"CtrlS":          tcell.KeyCtrlS,
	"CtrlT":          tcell.KeyCtrlT,
	"CtrlU":          tcell.KeyCtrlU,
	"CtrlV":          tcell.KeyCtrlV,
	"CtrlW":          tcell.KeyCtrlW,
	"CtrlX":          tcell.KeyCtrlX,
	"CtrlY":          tcell.KeyCtrlY,
	"CtrlZ":          tcell.KeyCtrlZ,
	"CtrlLeftSq":     tcell.KeyCtrlLeftSq,
	"CtrlBackslash":  tcell.KeyCtrlBackslash,
	"CtrlRightSq":    tcell.KeyCtrlRightSq,
	"CtrlCarat":      tcell.KeyCtrlCarat,
	"CtrlUnderscore": tcell.KeyCtrlUnderscore,
	"CtrlPageUp":     tcell.KeyCtrlPgUp,
	"CtrlPageDown":   tcell.KeyCtrlPgDn,
	"Tab":            tcell.KeyTab,
	"Esc":            tcell.KeyEsc,
	"Escape":         tcell.KeyEscape,
	"Enter":          tcell.KeyEnter,
	"Backspace":      tcell.KeyBackspace2,
	"OldBackspace":   tcell.KeyBackspace,

	// I renamed these keys to PageUp and PageDown but I don't want to break someone's keybindings
	"PgUp":   tcell.KeyPgUp,
	"PgDown": tcell.KeyPgDn,
}

// The Key struct holds the data for a keypress (keycode + modifiers)
type Key struct {
	keyCode   tcell.Key
	modifiers tcell.ModMask
	buttons   tcell.ButtonMask
	r         rune
	escape    string
}

// InitBindings initializes the keybindings for micro
func InitBindings() {
	bindings = make(map[Key][]func(*View, bool) bool)
	bindingsStr = make(map[string]string)
	mouseBindings = make(map[Key][]func(*View, bool, *tcell.EventMouse) bool)

	var parsed map[string]string
	defaults := DefaultBindings()

	filename := configDir + "/bindings.json"
	if _, e := os.Stat(filename); e == nil {
		input, err := ioutil.ReadFile(filename)
		if err != nil {
			TermMessage("Error reading bindings.json file: " + err.Error())
			return
		}

		err = json5.Unmarshal(input, &parsed)
		if err != nil {
			TermMessage("Error reading bindings.json:", err.Error())
		}
	}

	parseBindings(defaults)
	parseBindings(parsed)
}

func parseBindings(userBindings map[string]string) {
	for k, v := range userBindings {
		BindKey(k, v)
	}
}

// findKey will find binding Key 'b' using string 'k'
func findKey(k string) (b Key, ok bool) {
	modifiers := tcell.ModNone

	// First, we'll strip off all the modifiers in the name and add them to the
	// ModMask
modSearch:
	for {
		switch {
		case strings.HasPrefix(k, "-"):
			// We optionally support dashes between modifiers
			k = k[1:]
		case strings.HasPrefix(k, "Ctrl") && k != "CtrlH":
			// CtrlH technically does not have a 'Ctrl' modifier because it is really backspace
			k = k[4:]
			modifiers |= tcell.ModCtrl
		case strings.HasPrefix(k, "Alt"):
			k = k[3:]
			modifiers |= tcell.ModAlt
		case strings.HasPrefix(k, "Shift"):
			k = k[5:]
			modifiers |= tcell.ModShift
		case strings.HasPrefix(k, "\x1b"):
			return Key{
				keyCode:   -1,
				modifiers: modifiers,
				buttons:   -1,
				r:         0,
				escape:    k,
			}, true
		default:
			break modSearch
		}
	}

	if len(k) == 0 {
		return Key{buttons: -1}, false
	}

	// Control is handled specially, since some character codes in bindingKeys
	// are different when Control is depressed. We should check for Control keys
	// first.
	if modifiers&tcell.ModCtrl != 0 {
		// see if the key is in bindingKeys with the Ctrl prefix.
		k = string(unicode.ToUpper(rune(k[0]))) + k[1:]
		if code, ok := bindingKeys["Ctrl"+k]; ok {
			// It is, we're done.
			return Key{
				keyCode:   code,
				modifiers: modifiers,
				buttons:   -1,
				r:         0,
			}, true
		}
	}

	// See if we can find the key in bindingKeys
	if code, ok := bindingKeys[k]; ok {
		return Key{
			keyCode:   code,
			modifiers: modifiers,
			buttons:   -1,
			r:         0,
		}, true
	}

	// See if we can find the key in bindingMouse
	if code, ok := bindingMouse[k]; ok {
		return Key{
			modifiers: modifiers,
			buttons:   code,
			r:         0,
		}, true
	}

	// If we were given one character, then we've got a rune.
	if len(k) == 1 {
		return Key{
			keyCode:   tcell.KeyRune,
			modifiers: modifiers,
			buttons:   -1,
			r:         rune(k[0]),
		}, true
	}

	// We don't know what happened.
	return Key{buttons: -1}, false
}

// findAction will find 'action' using string 'v'
func findAction(v string) (action func(*View, bool) bool) {
	action, ok := bindingActions[v]
	if !ok {
		// If the user seems to be binding a function that doesn't exist
		// We hope that it's a lua function that exists and bind it to that
		action = LuaFunctionBinding(v)
	}
	return action
}

func findMouseAction(v string) func(*View, bool, *tcell.EventMouse) bool {
	action, ok := mouseBindingActions[v]
	if !ok {
		// If the user seems to be binding a function that doesn't exist
		// We hope that it's a lua function that exists and bind it to that
		action = LuaFunctionMouseBinding(v)
	}
	return action
}

// TryBindKey tries to bind a key by writing to configDir/bindings.json
// This function is unused for now
func TryBindKey(k, v string) {
	filename := configDir + "/bindings.json"
	if _, e := os.Stat(filename); e == nil {
		input, err := ioutil.ReadFile(filename)
		if err != nil {
			TermMessage("Error reading bindings.json file: " + err.Error())
			return
		}

		conflict := -1
		lines := strings.Split(string(input), "\n")
		for i, l := range lines {
			parts := strings.Split(l, ":")
			if len(parts) >= 2 {
				if strings.Contains(parts[0], k) {
					conflict = i
					TermMessage("Warning: Keybinding conflict:", k, " has been overwritten")
				}
			}
		}

		binding := fmt.Sprintf("    \"%s\": \"%s\",", k, v)
		if conflict == -1 {
			lines = append([]string{lines[0], binding}, lines[conflict:]...)
		} else {
			lines = append(append(lines[:conflict], binding), lines[conflict+1:]...)
		}
		txt := strings.Join(lines, "\n")
		err = ioutil.WriteFile(filename, []byte(txt), 0644)
		if err != nil {
			TermMessage("Error")
		}
	}
}

// BindKey takes a key and an action and binds the two together
func BindKey(k, v string) {
	key, ok := findKey(k)
	if !ok {
		TermMessage("Unknown keybinding: " + k)
		return
	}
	if v == "ToggleHelp" {
		helpBinding = k
	}
	if v == "ToggleKeyMenu" {
		kmenuBinding = k
	}
	if helpBinding == k && v != "ToggleHelp" {
		helpBinding = ""
	}
	if kmenuBinding == k && v != "ToggleKeyMenu" {
		kmenuBinding = ""
	}

	actionNames := strings.Split(v, ",")
	if actionNames[0] == "UnbindKey" {
		delete(bindings, key)
		delete(mouseBindings, key)
		delete(bindingsStr, k)
		if len(actionNames) == 1 {
			return
		}
		actionNames = append(actionNames[:0], actionNames[1:]...)
	}
	actions := make([]func(*View, bool) bool, 0, len(actionNames))
	mouseActions := make([]func(*View, bool, *tcell.EventMouse) bool, 0, len(actionNames))
	for _, actionName := range actionNames {
		if strings.HasPrefix(actionName, "Mouse") {
			mouseActions = append(mouseActions, findMouseAction(actionName))
		} else if strings.HasPrefix(actionName, "command:") {
			cmd := strings.SplitN(actionName, ":", 2)[1]
			actions = append(actions, CommandAction(cmd))
		} else if strings.HasPrefix(actionName, "command-edit:") {
			cmd := strings.SplitN(actionName, ":", 2)[1]
			actions = append(actions, CommandEditAction(cmd))
		} else {
			actions = append(actions, findAction(actionName))
		}
	}

	if len(actions) > 0 {
		// Can't have a binding be both mouse and normal
		delete(mouseBindings, key)
		bindings[key] = actions
		bindingsStr[k] = v
	} else if len(mouseActions) > 0 {
		// Can't have a binding be both mouse and normal
		delete(bindings, key)
		mouseBindings[key] = mouseActions
	}
}

// DefaultBindings returns a map containing micro's default keybindings
func DefaultBindings() map[string]string {
	return map[string]string{
		"Up":             "CursorUp",
		"Down":           "CursorDown",
		"Right":          "CursorRight",
		"Left":           "CursorLeft",
		"ShiftUp":        "SelectUp",
		"ShiftDown":      "SelectDown",
		"ShiftLeft":      "SelectLeft",
		"ShiftRight":     "SelectRight",
		"AltLeft":        "WordLeft",
		"AltShiftRight":  "SelectWordRight",
		"AltShiftLeft":   "SelectWordLeft",
		"CtrlUp":         "MoveLinesUp",
		"CtrlDown":       "MoveLinesDown",
		"CtrlLeft":       "StartOfLine",
		"CtrlRight":      "EndOfLine",
		"CtrlShiftLeft":  "SelectToStartOfLine",
		"ShiftHome":      "SelectToStartOfLine",
		"CtrlShiftRight": "SelectToEndOfLine",
		"ShiftEnd":       "SelectToEndOfLine",
		"AltRight":       "WordRight",
		"CtrlShiftUp":    "SelectToStart",
		"CtrlShiftDown":  "SelectToEnd",
		"Alt-{":          "ParagraphPrevious",
		"Alt-}":          "ParagraphNext",
		"Enter":          "InsertNewline",
		"CtrlH":          "Backspace",
		"Backspace":      "Backspace",
		"Alt-CtrlH":      "DeleteWordLeft",
		"Alt-Backspace":  "DeleteWordLeft",
		"Tab":            "IndentSelection,InsertTab",
		"Backtab":        "OutdentSelection,OutdentLine",
		"CtrlO":          "OpenFile",
		"CtrlS":          "Save",
		"CtrlF":          "Find",
		"CtrlN":          "FindNext",
		"CtrlP":          "FindPrevious",
		"CtrlZ":          "Undo",
		"CtrlShiftZ":     "Redo",
		"CtrlC":          "Copy",
		"CtrlX":          "Cut",
		"CtrlK":          "CutLine",
		"CtrlShiftD":     "DuplicateLine",
		"CtrlV":          "Paste",
		"CtrlA":          "SelectAll",
		"CtrlT":          "AddTab",
		"Alt,":           "PreviousTab",
		"Alt.":           "NextTab",
		"Home":           "StartOfLine",
		"End":            "EndOfLine",
		"CtrlHome":       "CursorStart",
		"CtrlEnd":        "CursorEnd",
		"PageUp":         "CursorPageUp",
		"PageDown":       "CursorPageDown",
		"CtrlPageUp":     "PreviousTab",
		"CtrlPageDown":   "NextTab",
		"CtrlG":          "ToggleHelp",
		"Alt-g":          "ToggleKeyMenu",
		"CtrlR":          "ToggleRuler",
		"CtrlL":          "JumpLine",
		"Delete":         "Delete",
		"CtrlB":          "ShellMode",
		"CtrlQ":          "Quit",
		"CtrlE":          "CommandMode",
		"CtrlW":          "NextSplit",
		"CtrlU":          "ToggleMacro",
		"CtrlJ":          "PlayMacro",
		"Insert":         "ToggleOverwriteMode",

		// Emacs-style keybindings
		"Alt-f": "WordRight",
		"Alt-b": "WordLeft",
		"Alt-a": "StartOfLine",
		"Alt-e": "EndOfLine",
		// "Alt-p": "CursorUp",
		// "Alt-n": "CursorDown",

		// Integration with file managers
		"F2":  "Save",
		"F3":  "Find",
		"F4":  "Quit",
		"F7":  "Find",
		"F10": "Quit",
		"Esc": "Escape",

		// Mouse bindings
		"MouseWheelUp":   "ScrollUp",
		"MouseWheelDown": "ScrollDown",
		"MouseLeft":      "MousePress",
		"MouseMiddle":    "PastePrimary",
		"Ctrl-MouseLeft": "MouseMultiCursor",

		"Alt-n": "SpawnMultiCursor",
		"Alt-m": "SpawnMultiCursorSelect",
		"Alt-p": "RemoveMultiCursor",
		"Alt-c": "RemoveAllMultiCursors",
		"Alt-x": "SkipMultiCursor",
	}
}
