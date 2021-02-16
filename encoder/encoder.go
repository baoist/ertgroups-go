package encoder

import (
	"strings"
	"regexp"
	"fmt"
)

type Delimiters struct {
	ColDelimRe string `json:"col_delimiter_regex"`
	ColExtraRe string `json:"col_extra_character_regex"`
	RowDelimRe string `json:"row_delimiter"`
}

type Group struct {
	Players []string `json:"groups"`
	Csv string `json:"csv"`
	Delimiters Delimiters `json:"delimiters"`
	ErtCode string `json:"ert_code"`
}

var _byte_to_6bit_char = [...]string {
	"a", "b", "c", "d", "e", "f", "g", "h",
	"i", "j", "k", "l", "m", "n", "o", "p",
	"q", "r", "s", "t", "u", "v", "w", "x",
	"y", "z", "A", "B", "C", "D", "E", "F",
	"G", "H", "I", "J", "K", "L", "M", "N",
	"O", "P", "Q", "R", "S", "T", "U", "V",
	"W", "X", "Y", "Z", "0", "1", "2", "3",
	"4", "5", "6", "7", "8", "9", "(", ")",
}

func encodeForPrint(str []byte) string {
	strArr := str
	var i = 0
	var buffer []string

	for i <= len(str) {
		vals := strArr[i:i+3]
		x1, x2, x3 := int(vals[0]), int(vals[1]), int(vals[2])

		i = i + 3

		cache := x1 + x2 * 256 + x3 * 65536

		byte1 := cache % 64
		cache = (cache - byte1) / 64
		byte2 := cache % 64
		cache = (cache - byte2) / 64
		byte3 := cache % 64
		byte4 := (cache - byte3) / 64

		buffer = append(buffer, _byte_to_6bit_char[byte1], _byte_to_6bit_char[byte2], _byte_to_6bit_char[byte3], _byte_to_6bit_char[byte4])
	}

	return strings.Join(buffer[:len(buffer)], "")
}

func (group *Group) assignPlayersToGroup() {
	players_csv := group.Csv

	row_delim_re := regexp.MustCompile(group.Delimiters.RowDelimRe)
	col_delim_re := regexp.MustCompile(group.Delimiters.ColDelimRe)
	col_extra_re := regexp.MustCompile(group.Delimiters.ColExtraRe)

	players_arr := make([]string, 40, 40)
	split_groups := make([]string, 8, 8)
	split_groups = row_delim_re.Split(players_csv, -1)[0:8]

	for group_idx, group := range split_groups {
		var tmp_row = make([]string, 5, 5)
	  copy(tmp_row[0:5], col_delim_re.Split(group, -1))

		for player_idx, player := range tmp_row {
			position := (group_idx * 5) + player_idx
			players_arr[position] = strings.Title(strings.ToLower(string(col_extra_re.ReplaceAll([]byte(player), []byte("")))))
		}
	}

	group.Players = players_arr
}

func (group *Group) format() string {
	header := "EXRTRGR0"
	tmp_player_arr := make([]string, 40, 40)

	for player_idx, player := range group.Players {
		formatted_player := fmt.Sprintf("[%d]=\"%s\"", player_idx + 1, player)
		tmp_player_arr[player_idx] = formatted_player
	}

	formatted_table_str := fmt.Sprintf("0,{%s}", strings.Join(tmp_player_arr[:], ","))
	encoded_str := encodeForPrint([]byte(formatted_table_str))

	return fmt.Sprintf("%s%s", header, encoded_str)
}

func (group *Group) Format() {
	group.ErtCode = group.format()
}

func CreateGroup(players_csv string, delimiters Delimiters) Group {
	var group *Group = new(Group)
	group.Delimiters = delimiters
	group.Csv = players_csv
	group.assignPlayersToGroup()

	return *group
}
