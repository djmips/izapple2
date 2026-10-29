package izapple2

import (
	"testing"
)

// TestNMI checks that CommandNMI interrupts the program through the NMI
// vector. On the Apple II+ the ROM sends it to $03fb, where a jump to an
// endless loop at $0300 is left.
func TestNMI(t *testing.T) {
	a, err := CreateAppleFromModel("2plus", map[string]string{"s6": "empty"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	a.Init()
	cyclesToText(t, a, "APPLE ][", 10_000_000)

	for i, b := range []uint8{0x4c, 0x00, 0x03} { // JMP $0300
		a.mmu.Poke(0x0300+uint16(i), b)
		a.mmu.Poke(0x03fb+uint16(i), b)
	}
	_, _, _, _, sp, pc := a.GetRegisters()
	if pc >= 0x0300 && pc <= 0x0302 {
		t.Fatalf("the program was already at the loop before the NMI")
	}

	a.SendCommand(CommandNMI)
	a.RunCycles(1_000)

	if pc := a.GetPC(); pc < 0x0300 || pc > 0x0302 {
		t.Fatalf("PC is $%04x after the NMI, it should be in the loop at $0300", pc)
	}
	pushed := a.Peek(0x0100 + uint16(sp-2))
	if pushed&0x10 != 0 || pushed&0x20 == 0 {
		t.Errorf("the status pushed by the NMI is $%02x, it must have B clear and bit 5 set", pushed)
	}
}
