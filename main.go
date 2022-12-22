package main

import (
	"fmt"
	"github.com/FishZe/go_bilichat_core/Handler"
	"github.com/FishZe/go_bilichat_core/client"
	"log"
	"time"
)

const (
	CmdDanmuMsg                  = "DANMU_MSG"
	CmdSuperChatMessage          = "SUPER_CHAT_MESSAGE"
	CmdSuperChatMessageJpn       = "SUPER_CHAT_MESSAGE_JPN"
	CmdWatchedChange             = "WATCHED_CHANGE"
	CmdSendGift                  = "SEND_GIFT"
	CmdOnlineRankCount           = "ONLINE_RANK_COUNT"
	CmdOnlineRankV2              = "ONLINE_RANK_V2"
	CmdOnlineRankTop3            = "ONLINE_RANK_TOP3"
	CmdLikeInfoV3Click           = "LIKE_INFO_V3_CLICK"
	CmdInteractWord              = "INTERACT_WORD"
	CmdStopLiveRoomList          = "STOP_LIVE_ROOM_LIST"
	CmdLikeInfoV3Update          = "LIKE_INFO_V3_UPDATE"
	CmdHotRankChange             = "HOT_RANK_CHANGED"
	CmdNoticeMsg                 = "NOTICE_MSG"
	CmdRoomRealTimeMessageUpdate = "ROOM_REAL_TIME_MESSAGE_UPDATE"
	CmdWidgetBanner              = "WIDGET_BANNER"
	CmdHotRankChangedV2          = "HOT_RANK_CHANGED_V2"
	CmdGuardHonorThousand        = "GUARD_HONOR_THOUSAND"
	CmdLive                      = "LIVE"
	CmdRoomChange                = "ROOM_CHANGE"
	CmdRoomBlockMsg              = "ROOM_BLOCK_MSG"
	CmdFullScreenSpecialEffect   = "FULL_SCREEN_SPECIAL_EFFECT"
	CmdCommonNoticeDanmaku       = "COMMON_NOTICE_DANMAKU"
	CmdTradingScore              = "TRADING_SCORE"
	CmdPreparing                 = "PREPARING"
	CmdGuardBuy                  = "GUARD_BUY"
	CmdGiftStarProcess           = "GIFT_STAR_PROCESS"
	CmdRoomSkinMsg               = "ROOM_SKIN_MSG"
	CmdEnterEffect               = "ENTER_EFFECT"
)

type option struct {
	Cmd    string
	RoomId int
	DoFunc func(event Handler.MsgEvent)
}

type Handle struct {
	Options []option
	cmdChan chan map[string]interface{}
}

type LiveRoom struct {
	RoomId int
	client client.Client
}

func GetNewHandler() Handle {
	return Handle{cmdChan: make(chan map[string]interface{}, 1)}
}

func (live *LiveRoom) New(roomId int) {
	live.client.RoomId = roomId
}

func (handle *Handle) New(Cmd string, RoomId int, DoFunc func(event Handler.MsgEvent)) {
	handle.Options = append(handle.Options, option{Cmd: Cmd, RoomId: RoomId, DoFunc: DoFunc})
}

func (handle *Handle) Binding(room LiveRoom) {
	room.client.RoomId = room.RoomId
	if room.client.RoomId == 0 {
		log.Printf("room id is 0")
		return
	}
	room.client.BiliChat(handle.cmdChan)
}

func (handle *Handle) Run() {
	var options []Handler.Options
	for _, option := range handle.Options {
		options = append(options, Handler.Options{RoomId: []int{option.RoomId}, Cmd: option.Cmd, DoFunc: option.DoFunc})
	}
	h := Handler.Handler{Options: options, CmdChan: handle.cmdChan}
	go h.CmdHandler()
}

func main() {
	// 创建一个新的消息处理器
	h := GetNewHandler()
	live := []int{21728563, 21330414, 6568696, 25243084, 22427257, 6154037, 47867, 25660853, 1155608, 23053924, 21919321, 23111212, 2373631, 25159955, 8514027, 5561470, 25514861, 156789, 645292, 25235084, 24837306, 13238959, 24555293, 78872, 7117085, 4549518, 6760154, 50353, 21483704, 675014, 22706788, 6632844, 847674, 1315354, 51296, 22361551, 22636800, 942178, 4895312, 3139709, 14052636, 14713062, 5461024, 22696219, 701558, 454154, 23086860, 3160112, 602158, 1226214, 21692711, 452138, 11946834, 25381972, 162022, 14893, 24255585, 23999270, 12257239, 22739677, 23502452, 452450, 25815214, 22503577, 25042360, 21564611, 23186788, 26427138, 23285010, 22379387, 26173394, 22054930, 23369248, 1521498, 1552519, 24930226, 6562109, 401900, 23319833, 24978909, 22893416, 2436830, 5424, 24363408, 3080147, 21711976, 3525213, 1557982, 25290861, 23057582, 26540437, 22778610, 109860, 23170704, 22724842, 24063763, 24167371, 21109751, 26285890, 26421557, 24697117, 25629638, 22937884, 1766019, 1496449, 8725120, 24441860, 26213288, 22309191, 4726132, 7406010, 23017343, 1301450, 22384516, 23847430, 23720185, 23517719, 26032691, 21509764, 23998995, 1068885, 90952, 22469889, 6615406, 282848, 23722214, 23142425, 23158420, 10525, 675516, 25161145, 1700323, 22758221, 2807112, 22845802, 22279402, 24930492, 4942796, 25872304, 513005, 23165805, 24813281, 271675, 3415150, 25034895, 23162141, 10244, 24061115, 739818, 17305, 814331, 25839275, 21224291, 21755112, 23937823, 515759, 1254342, 7953876, 22815988, 6374209, 14846654, 419589, 14963868, 24344580, 23224539, 26279304, 26376932, 23950097, 25814163, 23287807, 22754458, 21319517, 22196167, 23472646, 22195813, 23550749, 162180, 3755384, 14795432, 704558, 24650238, 573893, 21756256, 26045551, 23414796, 22889484, 23805059, 25745599, 94770, 11574776, 1495477, 12032317, 2537902, 22500194, 22809132, 761162, 22749172, 2446210, 10354150, 25915733, 450956, 24543885, 364027, 1385589, 23942399, 2853663, 4611808, 22770386, 21613356, 22839866, 23296355, 21392465, 23093993, 24789004, 26336970, 1208203, 24188245, 4640640, 5275, 12576972, 14409049, 5169315, 22159299, 41682, 22953683, 22920508, 22940391, 23408505, 22047448, 384601, 30294, 22901614, 1603600, 24698106, 24954295, 790566, 249310, 22880700, 2206397, 22032167, 3157408, 680462, 72888, 218550, 22901927, 4758413, 23550773, 448176, 289876, 548192, 2498734, 798668, 24391749, 342415, 21292831, 25062584, 713435, 18677, 15536, 623698, 395113, 4400273, 22195814, 681838, 21622680, 24924119, 25740452, 24369033, 25350006, 24046254, 315636, 25497805, 23246736, 7387222, 24896617, 23017349, 26352293, 24483629, 7906153, 14033373, 25622556, 25386814, 931170, 851181, 22536488, 219307, 26153692, 23455498, 827805, 25144904, 6026828, 58593, 19198, 868293, 22323445, 442683, 489701, 13448, 22603245, 918717, 21896446, 394275, 972474, 25026769, 25786450, 870004, 818240, 25756784, 24488740, 23752910, 25146115, 47369, 22887899, 26028386, 22817895, 23640445, 25813575, 2098875, 22445331, 777581, 5391723, 1191799, 22748536, 3217309, 21502501, 3484029, 6990537, 236672, 8816397, 24061020, 7404017, 342099, 23042268, 8728766, 10571417, 26133110, 3725076, 4215185, 753412, 24454136, 1517135, 22802883, 22321043, 24973541, 4992290, 21564812, 4783864, 6031397, 14085407, 23095176, 22437924, 1352574, 26494647, 24663865, 162432, 25439496, 22614739, 22612907, 23900412, 145960, 23865717, 24926753, 5298167, 21402309, 23611306, 3633573, 24255210, 23458654, 24145047, 25070229, 544625, 22687755, 662441, 915110, 11827039, 21868213, 21390, 529718, 25362856, 10801, 22787192, 1302409, 1487005, 596082, 4393208, 22320823, 24135903, 24488842, 1202360, 665722, 1107323, 337374, 411318, 1145005, 24160882, 23261911, 23578004, 23307008, 1909639, 23722192, 1921106, 1322342, 23813058, 25425764, 23427118, 26566761, 22824550, 24471121, 180305, 6485126, 4480043, 23116035, 5901838, 3032130, 81004, 26278833, 1497463, 22701393, 9447200, 22620570, 24308005, 25649533, 25006396, 7688602, 22892087, 24762598, 22426817, 66858, 24531089, 1114551, 22772718, 21713184, 22605289, 23365207, 830769, 24477947, 22857429, 81414, 22578273, 3550831, 602283, 1122827, 22387734, 23459697, 21775601, 26552829, 22882574, 282208, 5469741, 24997621, 24008492, 512301, 22861834, 23789885, 23126335, 23174842, 290889, 10317, 22810269, 10360, 4495259, 23245990, 15748, 25393774, 21637980, 24160817, 24572617, 52384, 553385, 11306, 22184511, 22700048, 25986072, 21772091, 21696950, 23069409, 24967779, 593688, 21368856, 24739681, 23750078, 6302912, 15149, 13639273, 2171843, 21879412, 885794, 23805029, 22890279, 26297216, 5357912, 23610417, 437817, 6925055, 22780009, 23403462, 22728239, 22563900, 24379966, 22210372, 23173129, 22447061, 2447903, 24643640, 22551845, 2303412, 24750843, 13599376, 708397, 21359166, 26067860, 25727850, 9990672, 22571722, 22537284, 78080, 409538, 22619596, 15122677, 24928063}
	// 将直播间号为15152878礼物消息交由PrintGift处理
	// 将直播间绑定到消息处理器
	for _, v := range live {
		h.New(CmdDanmuMsg, v, PrintDanmuMsg)
		h.Binding(LiveRoom{RoomId: v})
	}
	// 开始处理消息
	h.Run()
	for {
		time.Sleep(30 * time.Second)
	}
}

func PrintDanmuMsg(event Handler.MsgEvent) {
	fmt.Printf("[%v] %v: %v\n", event.RoomId, event.DanMuMsg.Data.Sender.Name, event.DanMuMsg.Data.Content)
}
