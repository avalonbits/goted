package text

import (
	"bufio"
	"io"
)

// Buffer represents the text being edited.
type Buffer struct {
	chars *chars
	lines *lines
}

func New(size int) *Buffer {
	return &Buffer{
		chars: newChars(size),
		lines: newLines(32_000),
	}
}

func (b *Buffer) Save(out io.Writer) error {
	bufOut := bufio.NewWriter(out)

	for _, text := range [][]rune{b.chars.prefix(), b.chars.suffix()} {
		for _, r := range text {
			if _, err := bufOut.WriteRune(r); err != nil {
				return err
			}
		}
	}

	return bufOut.Flush()
}

// chars is a character buffer used to store the text for the editor.
// It uses a gap buffer as its backing store.
type chars struct {
	buf    []rune
	cursor int
	curEnd int
}

// Newchars returns a *chars with the appropriate size.
func newChars(size int) *chars {
	return &chars{
		buf:    make([]rune, size),
		cursor: 0,
		curEnd: size,
	}
}

// Clear clease the gap buffer.
func (gb *chars) Clear() {
	gb.cursor = 0
	gb.curEnd = cap(gb.buf)
}

// Capacity returns the capacity of the gap buffer.
func (gb *chars) Capacity() int {
	return cap(gb.buf)
}

// Used returns how much of the capacity of the gap buffer has beend used.
func (gb *chars) Used() int {
	return gb.cursor + cap(gb.buf) - gb.curEnd
}

// Put stores a value in the gap buffer at th current position and advances the cursor.
// If there is no capacity available, returns false.
func (gb *chars) Put(val rune) bool {
	if gb.Capacity() == gb.Used() {
		return false
	}

	gb.buf[gb.cursor] = val
	gb.cursor++
	return true
}

// Delete removes the value under the cursor and retreats all values after the cursor one position.
// If there is no value to remove, returns false.
func (gb *chars) Delete() bool {
	if gb.curEnd >= cap(gb.buf) {
		return false
	}

	gb.curEnd++
	return true
}

// Backspace remove the value to the before the cursor and retreats all values starting at the cursor one position.
// If there is no value to remove, returns false.
func (gb *chars) Backspace() bool {
	if gb.cursor == 0 {
		return false
	}

	gb.cursor--
	return true
}

// Next advances the cursor count positions and returns how many positions it actually advanced.
func (gb *chars) Next(count int) int {
	target := count
	for count > 0 && gb.curEnd < cap(gb.buf) {
		gb.buf[gb.cursor] = gb.buf[gb.curEnd]

		gb.cursor++
		gb.curEnd++
		count--
	}

	return target - count
}

// Prev retreats the cursor count positions and returns how many positions it actually retreated.
func (gb *chars) Prev(count int) int {
	target := count
	for count > 0 && gb.cursor > 0 {
		gb.curEnd--
		gb.cursor--
		count--

		gb.buf[gb.curEnd] = gb.buf[gb.cursor]
	}

	return target - count
}

// Peak returns the value under the cursor.
func (gb *chars) Peek() (rune, bool) {
	if gb.curEnd == cap(gb.buf) {
		return 0, false
	}

	return gb.buf[gb.curEnd], true
}

func (gb *chars) prefix() []rune {
	return gb.buf[:gb.cursor]
}

func (gb *chars) suffix() []rune {
	return gb.buf[gb.curEnd:]
}

// lines is a line count buffer, used to track how much chars per line the teext editor has.
// It is also backed by a gap buffer.
type lines struct {
	buf    []int
	cursor int
	curEnd int
}

func newLines(size int) *lines {
	return &lines{
		buf:    make([]int, size),
		cursor: 0,
		curEnd: size,
	}
}

// Current returns the current line number.
func (l *lines) Current() int {
	return l.cursor
}

// Capacity returns the number of lines supported.
func (l *lines) Capacity() int {
	return cap(l.buf)
}

// Used returns how many lines were created.
func (l *lines) Used() int {
	return l.cursor + cap(l.buf) - l.curEnd
}

// Up moves the line pointer up.
func (l *lines) Up(count int) int {
	target := count

	for count > 0 && l.cursor > 0 {
		l.curEnd--
		l.buf[l.curEnd] = l.buf[l.cursor]
		l.cursor--

		count--
	}

	return target - count
}

// Down movs the line pointer down.
func (l *lines) Down(count int) int {
	target := count

	for count > 0 && l.curEnd < cap(l.buf) {
		l.cursor++
		l.buf[l.cursor] = l.buf[l.curEnd]
		l.curEnd++

		count--
	}

	return target - count
}

// Inc increments the character count for the line.
func (l *lines) Inc() int {
	l.buf[l.cursor]++
	return l.buf[l.cursor]
}

// Dec decrements the charactor count for the line.
func (l *lines) Dec() int {
	count := l.buf[l.cursor]
	if count > 0 {
		count--
		l.buf[l.cursor] = count
	}
	return count
}

// New adds a new line to the buffer with the capacity being (current line size) - splitSize.
// The current line size is updated to splitSize.
func (l *lines) New(splitSize int) bool {
	if l.Capacity() == l.Used() {
		return false
	}

	curSize := l.buf[l.cursor]
	l.buf[l.cursor] = splitSize
	l.cursor++
	l.buf[l.cursor] = curSize - splitSize

	return true
}
