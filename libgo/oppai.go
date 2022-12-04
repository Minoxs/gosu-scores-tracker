package libgo

// #define OPPAI_IMPLEMENTATION
// #include "../libc/oppai.c"
import "C"
import "fmt"

func GetPP() {
	const stdinString string = "-"
	ez := C.ezpp_new()
	C.ezpp_set_mods(ez, C.MODS_HR)
	C.ezpp_set_accuracy(ez, 37, 4)
	C.ezpp_set_nmiss(ez, 0)
	C.ezpp_set_combo(ez, 476)
	C.ezpp_set_aim_stars(ez, 3.06299)
	C.ezpp_set_speed_stars(ez, 2.21616)
	C.ezpp(ez, C.CString(stdinString))
	fmt.Println(C.ezpp_pp(ez))
	fmt.Println(C.ezpp_accuracy_percent(ez))
	fmt.Println(C.ezpp_aim_stars(ez))
	fmt.Println(C.ezpp_speed_stars(ez))
	defer C.ezpp_free(ez)
}
